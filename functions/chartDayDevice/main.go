package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
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
	request := new(vm.ChartHourDeviceRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	var dbData []map[string]interface{}
	for _, s := range request.SensorAll {
		res, err := dal.QueryMultiple("chart_device_day",
			dal.Condition{
				"hash": {
					ComparisonOperator: aws.String("EQ"),
					AttributeValueList: []*dal.AttributeValue{
						{
							S: aws.String(fmt.Sprintf("%v<->%v", request.Token, s)),
						},
					},
				},
				"date": {
					ComparisonOperator: aws.String("GT"),
					AttributeValueList: []*dal.AttributeValue{
						{
							N: aws.String(request.From),
						},
					},
				},
			},
			dal.Projection(dal.Name("hash"), dal.Name("date"), dal.Name("value")),
			true, true)

		if err != nil {
			response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
			fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
		}

		var tmpData []map[string]interface{}
		err = res.Unmarshal(&tmpData)
		if err != nil {
			response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
			fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
		}

		dbData = append(dbData, tmpData...)
	}

	type schemaData struct {
		Version   string   `json:"version"`
		Data      m.Schema `json:"data"`
		Heartbeat int      `json:"heartbeat"`
	}

	schemaRes, err := dal.Get("schema", map[string]*dal.AttributeValue{
		"version": {
			S: aws.String("1"),
		},
	})
	if err != nil {
		fmt.Println("Error fetching schema from db", err)
		response.AddError(&es.Error{Message: err.Error(), Data: "could not fetch schema data from the DB"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	schema := new(schemaData)
	err = schemaRes.Unmarshal(schema)
	if err != nil {
		fmt.Println("Error unmarshaling schema ::. ", err)
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	if len(request.SensorAll) <= 1 {
		result := make([]*resultData, 0)
		for _, v := range dbData {
			val := v["value"].(float64)
			if request.Sensor != "BATTERY_VOLTAGE" {
				val = math.Round(val)
			}
			result = append(result, &resultData{
				Date:  v["date"].(float64),
				Value: val,
			})
		}

		result = qsort(result)
		result = fillDataOffline(result, schema.Heartbeat)
		result = qsort(result)
		response.Data = result

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}

	resultChart := make([]map[string]float64, 0)
	maxValues := make(map[string]float64, 0)
	d := make(map[float64]map[string][]float64, 0)

	for _, v := range dbData {
		date := v["date"].(float64)
		_, ok := d[date]
		if !ok {
			d[date] = make(map[string][]float64, 0)
		}
		for _, s := range request.SensorAll {
			splitHash := strings.Split(v["hash"].(string), "<->")
			if len(splitHash) > 1 && splitHash[1] == s {

				d[date][splitHash[1]] = make([]float64, 0)
				d[date][s] = append(d[date][s], v["value"].(float64))
			}

		}
	}

	fmt.Println("D ::: ", d)

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
					if s == "BATTERY_VOLTAGE" {
						rd[s] = av / float64(len(values))
					} else {
						rd[s] = math.Round(av / float64(len(values)))
					}
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
	resultChart = fillDataMulti(resultChart, request.SensorAll)
	resultChart = fillDataMultiOffline(resultChart, schema.Heartbeat)
	resultChart = qsortMulti(resultChart)
	// resultChart = smoothMulti(resultChart)

	response.Data = resultDataMulti{Chart: resultChart, Max: maxValues}
	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// fills device offline periods for multiple sensors
func fillDataMultiOffline(data []map[string]float64, heartbeat int) []map[string]float64 {
	// if no data return
	if len(data) < 1 {
		return data
	}

	var interval float64 = 60 * 30
	var onlineTime float64 = 60 * 2 * float64(heartbeat) // heartbeat is stated in minutes
	latest := data[0]["date"]
	diff := float64(time.Now().Unix()) - latest

	if diff > interval {
		// device is int artif. online mode, add data
		for i := diff; i > interval; i -= interval {
			dataToFill := make(map[string]float64, 0)
			for k, v := range data[0] {
				dataToFill[k] = v
			}
			dataToFill["date"] = dataToFill["date"] + interval

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

				for j := 1; j <= timesToAdd; j++ {
					dataToFill := make(map[string]float64, 0)
					for k, v := range data[k+1] {
						dataToFill[k] = v
					}
					dataToFill["date"] = dataToFill["date"] + (interval * float64(j))

					// insert data on the needed index
					data = append(data[:k], append([]map[string]float64{dataToFill}, data[k:]...)...)
					k++
				}
			}
		}
	}
	return data
}

// fills device offline periods for single sensor
func fillDataOffline(data []*resultData, heartbeat int) []*resultData {
	// if no data return
	if len(data) < 1 {
		return data
	}

	var interval float64 = 60 * 30
	var onlineTime float64 = 60 * 2 * float64(heartbeat) // heartbeat is stated in minutes
	latest := data[0].Date
	diff := float64(time.Now().Unix()) - latest

	if diff > interval {
		// device is int artif. online mode, add data
		for i := diff; i > interval; i -= interval {
			dataToFill := *data[0]
			dataToFill.Date = dataToFill.Date + interval

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
					timesToAdd = maxTimesToAdd
				}

				fmt.Println("TIMES TO ADD", timesToAdd, maxTimesToAdd)
				for j := 1; j <= timesToAdd; j++ {
					dataToFill := *data[k+1]
					dataToFill.Date = dataToFill.Date + (interval * float64(j))

					// insert data on the needed index
					data = append(data[:k], append([]*resultData{&dataToFill}, data[k:]...)...)
					k++
				}
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
