package main

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.MapRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	// Get the polygon data
	zoneRes, err := dal.List("zones", dal.Name("sensor_id").Equal(dal.Value(request.Sensor)), dal.Projection(dal.Name("zone_id"), dal.Name("sensor_id"), dal.Name("data")))
	zoneData := make([]m.Zone, 0)
	err = zoneRes.Unmarshal(&zoneData)
	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	dbRes, err := dal.ListNoFilter("devices", dal.Projection(dal.Name("token"), dal.Name("device_id"), dal.Name("meta"), dal.Name("map_meta"), dal.Name("active"), dal.Name("measurements")))
	dbData := make([]m.Device, 0)
	err = dbRes.Unmarshal(&dbData)
	if err != nil {
		response.AddError(&es.Error{Code: 0, Message: err.Error(), Data: ""})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	data := make([]m.MapMeta, 0)

	for _, d := range dbData {
		if !d.Active {
			continue
		}

		data = append(data, d.MapMeta[request.Sensor])
	}

	resData := vm.MapResponseData{
		Zones:   zoneData,
		Devices: data,
	}
	response.Data = resData

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func main() {
	lambda.Start(Handler)
}
