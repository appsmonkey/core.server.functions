package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

type resultData struct {
	Date  float64 `json:"date"`
	Value float64 `json:"value"`
}

type resultDataMulti struct {
	Chart []map[string]float64 `json:"chart"`
	Max   map[string]float64   `json:"max"`
}

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.ChartLiveDeviceRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	names := make([]dal.NameBuilder, 0)
	for _, s := range request.SensorAll {
		names = append(names, dal.Name(s))
	}

	projBuilder := dal.Projection(dal.Name("timestamp"), names...)
	res, err := dal.QueryMultiple("live",
		dal.Condition{
			"token": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dal.AttributeValue{
					{
						S: aws.String(request.Token),
					},
				},
			},
			"timestamp": {
				ComparisonOperator: aws.String("GT"),
				AttributeValueList: []*dal.AttributeValue{
					{
						N: aws.String(request.From),
					},
				},
			},
		},
		projBuilder,
		true, true)

	if err != nil {
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	var dbData []map[string]float64
	err = res.Unmarshal(&dbData)
	if err != nil {
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	if len(request.SensorAll) <= 1 {
		result := make([]*resultData, 0)
		for _, v := range dbData {
			val := v[request.Sensor]
			if request.Sensor != "BATTERY_VOLTAGE" {
				val = math.Round(val)
			}
			result = append(result, &resultData{
				Date:  v["timestamp"],
				Value: val,
			})
		}

		result = qsort(result)
		result = fillDataOffline(result)
		// result = smooth(result)
		response.Data = result
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}

	resultChart := make([]map[string]float64, 0)
	maxValues := make(map[string]float64, 0)

	for _, v := range dbData {
		rd := make(map[string]float64, 0)
		for _, s := range request.SensorAll {
			rd["date"] = v["timestamp"]
			if s == "BATTERY_VOLTRAGE" {
				rd[s] = v[s]
			} else {
				rd[s] = math.Round(v[s])
			}

			mv, okmv := maxValues[s]
			if rd[s] > mv {
				maxValues[s] = rd[s]
			} else if !okmv {
				maxValues[s] = 0
			}
		}

		resultChart = append(resultChart, rd)
	}

	resultChart = qsortMulti(resultChart)
	resultChart = fillDataMulti(resultChart, request.SensorAll)
	resultChart = fillDataMultiOffline(resultChart)
	// resultChart = smoothMulti(resultChart)

	response.Data = resultDataMulti{Chart: resultChart, Max: maxValues}

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// fills device offline periods for single sensor
func fillDataOffline(data []*resultData) []*resultData {
	// if no data return
	if len(data) < 1 {
		return data
	}

	var interval float64 = 60
	var onlineTime float64 = 60 * 120
	latest := data[0].Date
	diff := float64(time.Now().Unix()) - latest

	if diff > interval {
		// device is int artif. online mode, add data
		for i := diff; i > interval; i -= interval {
			dataToFill := *data[0]
			dataToFill.Date = dataToFill.Date + 60

			// stop filling after online period is exceeded
			if dataToFill.Date >= latest+onlineTime {
				break
			}

			// prepend data
			data = append([]*resultData{&dataToFill}, data...)
		}
	}

	if len(data) > 2 {
		// data point difference in sec
		for k := 0; k < len(data)-1; k++ {
			diff := data[k].Date - data[k+1].Date

			if diff > interval {
				timesToAdd := int(diff) / int(interval)
				maxTimesToAdd := int(onlineTime) / int(interval)

				// if exceeds onlineTime fill only max online time
				if timesToAdd > maxTimesToAdd {
					fmt.Println("MAX TIMES TO ADD EXCEEDED, SETTING LIMIT")
					timesToAdd = maxTimesToAdd
				}

				fmt.Println("TIMES TO ADD", timesToAdd, maxTimesToAdd)
				for j := 1; j <= timesToAdd; j++ {
					dataToFill := *data[k+1]
					dataToFill.Date = data[k].Date - (interval * float64(j))

					// insert data on the needed index
					data = append(data[:k], append([]*resultData{&dataToFill}, data[k:]...)...)
					k++
				}
				k++
			}
		}
	}
	return data
}

// fills device offline periods for multiple sensors
func fillDataMultiOffline(data []map[string]float64) []map[string]float64 {
	// if no data return
	if len(data) < 1 {
		return data
	}

	var interval float64 = 60
	var onlineTime float64 = 60 * 120
	latest := data[0]["date"]
	diff := float64(time.Now().Unix()) - latest

	if diff > interval {
		// device is int artif. online mode, add data
		for i := diff; i > interval; i -= interval {
			dataToFill := data[0]
			dataToFill["date"] = dataToFill["date"] + 60

			// stop filling after online period is exceeded
			if dataToFill["date"] >= latest+onlineTime {
				break
			}

			// prepend data
			data = append([]map[string]float64{dataToFill}, data...)
		}
	}

	if len(data) > 2 {
		// data point difference in sec
		for k := 0; k < len(data)-1; k++ {
			diff := data[k]["date"] - data[k+1]["date"]

			if diff > interval {
				timesToAdd := int(diff) / int(interval)
				maxTimesToAdd := int(onlineTime) / int(interval)

				// if exceeds onlineTime fill only max online time
				if timesToAdd > maxTimesToAdd {
					timesToAdd = maxTimesToAdd
				}

				fmt.Println("TIMES TO ADD", timesToAdd, maxTimesToAdd)
				for j := 1; j <= timesToAdd; j++ {
					dataToFill := data[k+1]
					dataToFill["date"] = data[k]["date"] - (interval * float64(j))

					// insert data on the needed index
					data = append(data[:k], append([]map[string]float64{dataToFill}, data[k:]...)...)
					k++
				}
				k++
			}
		}
	}
	return data
}

// data has to be sorted asc. by date in order for this to work
func fillDataMulti(data []map[string]float64, sensors []string) []map[string]float64 {

	for k, v := range data {
		if k == len(data)-1 {
			continue
		}

		for _, vs := range sensors {
			_, ok := v[vs]

			if !ok {
				for i := k + 1; i < len(data); i++ {
					val, ok := data[i][vs]
					if ok {
						v[vs] = val
						break
					}
				}
			}
		}
	}
	return data
}

// qsort is a quicksort implmentation for sorting chart data
func qsortMulti(a []map[string]float64) []map[string]float64 {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	// Pick a pivot
	pivotIndex := rand.Int() % len(a)

	// Move the pivot to the right
	a[pivotIndex], a[right] = a[right], a[pivotIndex]

	// Pile elements smaller than the pivot on the left
	for i := range a {
		if a[i]["date"] > a[right]["date"] {
			a[i], a[left] = a[left], a[i]
			left++
		}
	}

	// Place the pivot after the last smaller element
	a[left], a[right] = a[right], a[left]

	// Go down the rabbit hole
	qsortMulti(a[:left])
	qsortMulti(a[left+1:])

	return a
}

// qsort is a quicksort implmentation for sorting chart data
func qsort(a []*resultData) []*resultData {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	// Pick a pivot
	pivotIndex := rand.Int() % len(a)

	// Move the pivot to the right
	a[pivotIndex], a[right] = a[right], a[pivotIndex]

	// Pile elements smaller than the pivot on the left
	for i := range a {
		if a[i].Date > a[right].Date {
			a[i], a[left] = a[left], a[i]
			left++
		}
	}

	// Place the pivot after the last smaller element
	a[left], a[right] = a[right], a[left]

	// Go down the rabbit hole
	qsort(a[:left])
	qsort(a[left+1:])

	return a
}

func main() {
	lambda.Start(Handler)
}

func smooth(in []*resultData) []*resultData {
	result := make([]*resultData, 0)
	lenIN := len(in)

	for i := 0; i < lenIN; i++ {
		result = append(result, in[i])

		j := i + 1
		if j == lenIN {
			break
		}

		iDate := time.Unix(int64(in[i].Date), 0)
		jDate := time.Unix(int64(in[j].Date), 0)
		iValue := in[i].Value
		jValue := in[j].Value

		res := smoothPoints(iDate, jDate, iValue, jValue)
		for _, dp := range res {
			result = append(result, dp)
		}
	}

	return result
}

func smoothPoints(it, jt time.Time, iv, jv float64) []*resultData {
	year, month, day, hour, min, _ := diff(it, jt)
	minutes := float64(year*525600 + month*43800 + day*1440 + hour*60 + min)
	res := make([]*resultData, 0)

	// for large differences we can just add a simple curve
	i := 0
	if year > 0 || month > 0 || day > 0 {
		mod := float64(10)
		v := iv
		t := it
		for {
			i++
			// Get the time for the new data point (substract 10%)
			m := time.Duration(minutes / mod)
			t = t.Add(time.Minute * m * -1)

			// Get the value for the new data point (substract 10%) of the difference between the two points
			if iv > jv {
				v -= (iv - jv) / mod
			} else if iv < jv {
				v += (jv - iv) / mod
			}

			// if we overshot, stop
			if t.Before(jt) {
				break
			}

			res = append(res, &resultData{
				Date:  float64(t.Unix()),
				Value: v,
			})
		}

		return res
	}

	if hour > 0 {
		mod := float64(5)
		v := iv
		t := it
		for {
			i++
			// Get the time for the new data point (substract 10%)
			m := time.Duration(minutes / mod)
			t = t.Add(time.Minute * m * -1)

			// Get the value for the new data point (substract 10%) of the difference between the two points
			if iv > jv {
				v -= (iv - jv) / mod
			} else if iv < jv {
				v += (jv - iv) / mod
			}

			// if we overshot, stop
			if t.Before(jt) {
				break
			}

			res = append(res, &resultData{
				Date:  float64(t.Unix()),
				Value: v,
			})
		}

		return res
	}

	// if we have three minutes missing, just do nothing
	if min <= 3 {
		return res
	}

	// we have a minutes chart so we need to figure out the amount of data point
	// to put between the two existing points
	// we base it on the minimum value jump in our dataset
	if min > 3 {
		mod := float64(min)
		v := iv
		t := it
		for {
			i++
			// Get the time for the new data point (substract 10%)
			m := time.Duration(minutes / mod)
			t = t.Add(time.Minute * m * -1)

			// Get the value for the new data point (substract 10%) of the difference between the two points
			if iv > jv {
				v -= (iv - jv) / mod
			} else if iv < jv {
				v += (jv - iv) / mod
			}

			// if we overshot, stop
			if t.Before(jt) {
				break
			}

			res = append(res, &resultData{
				Date:  float64(t.Unix()),
				Value: v,
			})
		}

		return res
	}

	return res
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func smoothMulti(in []map[string]float64) []map[string]float64 {
	result := make([]map[string]float64, 0)
	lenIN := len(in)

	for i := 0; i < lenIN; i++ {
		result = append(result, in[i])

		j := i + 1
		if j == lenIN {
			break
		}

		iDate := time.Unix(int64(in[i]["date"]), 0)
		jDate := time.Unix(int64(in[j]["date"]), 0)
		iValue := in[i]["value"]
		jValue := in[j]["value"]

		res := smoothPointsMulti(iDate, jDate, iValue, jValue)
		for _, dp := range res {
			result = append(result, dp)
		}
	}

	return result
}

func smoothPointsMulti(it, jt time.Time, iv, jv float64) []map[string]float64 {
	year, month, day, hour, min, _ := diff(it, jt)
	minutes := float64(year*525600 + month*43800 + day*1440 + hour*60 + min)
	res := make([]map[string]float64, 0)

	// for large differences we can just add a simple curve
	if year > 0 || month > 0 || day > 0 {
		mod := float64(10)
		v := iv
		t := it
		for {
			// Get the time for the new data point (substract 10%)
			m := time.Duration(minutes / mod)
			t = t.Add(time.Minute * m * -1)

			// Get the value for the new data point (substract 10%) of the difference between the two points
			if iv > jv {
				v -= (iv - jv) / mod
			} else if iv < jv {
				v += (jv - iv) / mod
			}

			// if we overshot, stop
			if t.Before(jt) {
				break
			}

			res = append(res, map[string]float64{
				"date":  float64(t.Unix()),
				"value": v,
			})
		}

		return res
	}

	if hour > 0 {
		mod := float64(5)
		v := iv
		t := it
		for {
			// Get the time for the new data point (substract 10%)
			m := time.Duration(minutes / mod)
			t = t.Add(time.Minute * m * -1)

			// Get the value for the new data point (substract 10%) of the difference between the two points
			if iv > jv {
				v -= (iv - jv) / mod
			} else if iv < jv {
				v += (jv - iv) / mod
			}

			// if we overshot, stop
			if t.Before(jt) {
				break
			}

			res = append(res, map[string]float64{
				"date":  float64(t.Unix()),
				"value": v,
			})
		}

		return res
	}

	// if we have three minutes missing, just do nothing
	if min <= 3 {
		return res
	}

	// we have a minutes chart so we need to figure out the amount of data point
	// to put between the two existing points
	// we base it on the minimum value jump in our dataset
	if min > 3 {
		mod := float64(min)
		v := iv
		t := it
		for {
			// Get the time for the new data point (substract 10%)
			m := time.Duration(minutes / mod)
			t = t.Add(time.Minute * m * -1)

			// Get the value for the new data point (substract 10%) of the difference between the two points
			if iv > jv {
				v -= (iv - jv) / mod
			} else if iv < jv {
				v += (jv - iv) / mod
			}

			// if we overshot, stop
			if t.Before(jt) {
				break
			}

			res = append(res, map[string]float64{
				"date":  float64(t.Unix()),
				"value": v,
			})
		}

		return res
	}

	return res
}
