package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

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
	request := new(vm.ChartHourAllRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	dbRawData := make(map[string][]map[string]float64, 0)
	for _, sid := range request.SensorAll {
		res, err := dal.QueryMultiple("chart_day",
			dal.Condition{
				"sensor": {
					ComparisonOperator: aws.String("EQ"),
					AttributeValueList: []*dal.AttributeValue{
						{
							S: aws.String(sid),
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
			dal.Projection(dal.Name("date"), dal.Name("value")),
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

		dbRawData[sid] = dbData
	}

	if len(dbRawData) <= 1 {
		result := make([]*resultData, 0)
		for _, v := range dbRawData[request.Sensor] {
			result = append(result, &resultData{
				Date:  v["date"],
				Value: v["value"],
			})
		}

		result = qsort(result)
		response.Data = result
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}

	resultRaw := make(map[float64]map[string]float64, 0)
	resultRawIndex := make([]float64, 0)
	result := make([]map[string]float64, 0)

	for sid, sd := range dbRawData {
		for _, v := range sd {
			d := v["date"]
			val := v["value"]
			_, ok := resultRaw[d]
			if !ok {
				resultRaw[d] = make(map[string]float64, 0)
				resultRawIndex = append(resultRawIndex, d)
			}

			resultRaw[d][sid] = val
		}
	}

	sort.Float64s(resultRawIndex)
	sort.Sort(sort.Reverse(sort.Float64Slice(resultRawIndex)))
	maxValues := make(map[string]float64, 0)

	for _, ri := range resultRawIndex {
		d := resultRaw[ri]
		rd := make(map[string]float64, 0)
		rd["date"] = ri

		for _, s := range request.SensorAll {
			sd, ok := d[s]
			if ok {
				if s == "BATTERY_VOLTAGE" {
					rd[s] = sd
				} else {
					rd[s] = math.Round(sd)
				}
				mv, okmv := maxValues[s]
				if sd > mv {
					maxValues[s] = sd
				} else if !okmv {
					maxValues[s] = 0
				}
			} else {
				// rd[s] = 0
				continue
			}
		}

		result = append(result, rd)
	}

	result = fillDataMulti(result, request.SensorAll)
	response.Data = resultDataMulti{Chart: result, Max: maxValues}
	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
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
