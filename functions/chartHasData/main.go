package main

import (
	"encoding/json"
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

var mapping map[string]string
var mappingAll map[string]string

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.ChartHasDataRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	m := mapping
	if !request.Device {
		m = mappingAll
	}

	table, ok := m[request.Chart]
	if !ok {
		response.AddError(&es.ErrMissingChart)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	var qry dal.ConditionBuilder
	if request.Device {
		qry = dal.Name("hash").Equal(dal.Value(request.Token + ":" + request.Sensor)).And(dal.Name("date").GreaterThanEqual(dal.Value(request.From)))
	} else {
		qry = dal.Name("sensor").Equal(dal.Value(request.Sensor)).And(dal.Name("date").GreaterThanEqual(dal.Value(request.From)))
	}

	response.Data = dal.HasItemsWithFilter(table, qry)

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func main() {
	mappingAll = make(map[string]string, 0)
	mappingAll["live"] = "live"
	mappingAll["day"] = "chart_hour"
	mappingAll["week"] = "chart_six"
	mappingAll["month"] = "chart_day"

	mapping = make(map[string]string, 0)
	mapping["live"] = "live"
	mapping["day"] = "chart_device_hour"
	mapping["week"] = "chart_device_six"
	mapping["month"] = "chart_device_day"

	lambda.Start(Handler)
}

func printJson(in interface{}) {
	b, _ := json.Marshal(in)
	fmt.Println(string(b))
}
