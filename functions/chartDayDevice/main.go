package main

import (
	"fmt"
	"math/rand"
	"strings"

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
			true)

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

	if len(request.SensorAll) <= 1 {
		result := make([]*resultData, 0)
		for _, v := range dbData {
			result = append(result, &resultData{
				Date:  v["date"].(float64),
				Value: v["value"].(float64),
			})
		}

		result = qsort(result)
		response.Data = result

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}

	resultChart := make([]map[string]float64, 0)
	maxValues := make(map[string]float64, 0)

	for _, v := range dbData {
		rd := make(map[string]float64, 0)
		for _, s := range request.SensorAll {
			splitHash := strings.Split(v["hash"].(string), "<->")

			if len(splitHash) > 1 && splitHash[1] == s {
				rd["date"] = v["date"].(float64)
				rd[s] = v["value"].(float64)

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
	// resultChart = smoothMulti(resultChart)

	response.Data = resultDataMulti{Chart: resultChart, Max: maxValues}
	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
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
