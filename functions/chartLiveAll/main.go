package main

import (
	"fmt"

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

	projBuilder := dal.Projection(dal.Name("timestamp_sort"), names...)
	res, err := dal.List("live", dal.Name("timestamp").GreaterThanEqual(dal.Value(request.From)), projBuilder)
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
		d := make(map[float64][]float64, 0)
		for _, v := range dbData {
			d[v["timestamp_sort"]] = append(d[v["timestamp_sort"]], v[request.Sensor])
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

		response.Data = result
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}

	resultChart := make([]map[string]float64, 0)
	maxValues := make(map[string]float64, 0)
	d := make(map[float64]map[string][]float64, 0)

	for _, v := range dbData {
		date := v["timestamp_sort"]
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

	response.Data = resultDataMulti{Chart: resultChart, Max: maxValues}

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func main() {
	lambda.Start(Handler)
}
