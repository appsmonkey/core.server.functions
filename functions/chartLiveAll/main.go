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

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.ChartLiveAllRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	res, err := dal.List("live", dal.Name("timestamp").GreaterThanEqual(dal.Value(request.From)), dal.Projection(dal.Name(request.Sensor), dal.Name("timestamp_sort")))
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

func main() {
	lambda.Start(Handler)
}
