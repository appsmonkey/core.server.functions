package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
)

type resultData struct {
	Date  float64 `json:"date"`
	Value float64 `json:"value"`
}

type resultDataMulti struct {
	Chart []map[string]float64 `json:"chart"`
	Max   map[string]float64   `json:"max"`
}

type resultDataAll struct {
	Date  float64            `json:"date"`
	Value map[string]float64 `json:"value"`
}

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.ChartLiveAllRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	names := make([]dal.NameBuilder, 0)
	for _, s := range request.SensorAll {
		names = append(names, dal.Name(s))
	}
	names = append(names, dal.Name("indoor"))
	names = append(names, dal.Name("token"))

	projBuilder := dal.Projection(dal.Name("timestamp_sort"), names...)
	res, err := dal.List("chart_all_minute", dal.Name("timestamp").GreaterThanEqual(dal.Value(request.From)), projBuilder)
	if err != nil {
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	var dbDataForFilter []map[string]interface{}
	var dbData []map[string]float64
	err = res.Unmarshal(&dbDataForFilter)
	if err != nil {
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	for _, v := range dbDataForFilter {
		if v["indoor"] == false || v["indoor"] == "false" {
			r := make(map[string]float64, 0)
			for ka, va := range v {
				if ka != "indoor" && ka != "token" {
					r[ka] = va.(float64)
				}
			}

			dbData = append(dbData, r)
		}
	}

	if len(request.SensorAll) <= 1 {
		result := make([]*resultData, 0)
		d := make(map[float64][]float64, 0)
		for _, v := range dbData {

			merged := false
			for k := range d {
				iDate := time.Unix(int64(k), 0)
				jDate := time.Unix(int64(v["timestamp"]), 0)
				year, month, day, hour, min, sec := diff(iDate, jDate)
				if year == 0 && month == 0 && day == 0 && hour == 0 && min == 0 && sec > 0 {
					merged = true
					d[k] = append(d[v["timestamp"]], v[request.Sensor])
					break
				}
			}

			if !merged {
				d[v["timestamp"]] = append(d[v["timestamp"]], v[request.Sensor])
			}
		}

		for k, v := range d {
			rd := new(resultData)
			rd.Date = k
			if len(v) > 0 {
				var av float64
				for _, sv := range v {
					av += sv
				}

				if av == 0 {
					rd.Value = 0
				} else {
					rd.Value = av / float64(len(v))
				}
			}

			result = append(result, rd)
		}

		// sort data according to timestamp
		result = qsort(result)
		// result = smooth(result)

		response.Data = result
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}

	resultChart := make([]map[string]float64, 0)
	maxValues := make(map[string]float64, 0)
	d := make(map[float64]map[string][]float64, 0)

	for _, v := range dbData {
		date := v["timestamp"]
		_, ok := d[date]
		if !ok {
			d[date] = make(map[string][]float64, 0)
		}
		for _, s := range request.SensorAll {
			_, ok := d[date][s]
			if !ok {
				d[date][s] = make([]float64, 0)
			}
			d[date][s] = append(d[date][s], v[s])
		}
	}

	for k, v := range d {
		rd := make(map[string]float64, 0)
		rd["date"] = k
		for _, s := range request.SensorAll {
			values := v[s]
			if len(values) > 0 {
				var av float64
				for _, sv := range values {
					av += sv
				}

				if av == 0 {
					rd[s] = 0
				} else {
					rd[s] = av / float64(len(values))
				}

				mv, okmv := maxValues[s]
				if rd[s] > mv {
					maxValues[s] = rd[s]
				} else if !okmv {
					maxValues[s] = 0
				}
			}
		}

		resultChart = append(resultChart, rd)
	}

	resultChart = qsortMulti(resultChart)
	// resultChart = smoothMulti(resultChart)

	response.Data = resultDataMulti{Chart: resultChart, Max: maxValues}

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
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

func smoothPoints(it, jt time.Time, iv, jv float64) []*resultData {
	year, month, day, hour, min, _ := diff(it, jt)
	minutes := float64(year*525600 + month*43800 + day*1440 + hour*60 + min)
	res := make([]*resultData, 0)

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

func main() {
	lambda.Start(Handler)
}
