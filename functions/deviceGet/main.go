package main

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	defaultDevice "github.com/appsmonkey/core.server.functions/functions/deviceGet/general"
	m "github.com/appsmonkey/core.server.functions/models"
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID := CognitoData(req.RequestContext.Authorizer)
	request := new(vm.DeviceGetRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	if len(request.Token) == 0 {
		response.Data = defaultDevice.Get()
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}

	res, err := dal.Get("devices", map[string]*dal.AttributeValue{
		"token": {
			S: aws.String(request.Token),
		},
	})
	if err != nil {
		errData := es.ErrDeviceNotFound
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	model := m.Device{}
	err = res.Unmarshal(&model)
	if err != nil {
		errData := es.ErrDeviceNotFound
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	data := vm.DeviceGetData{
		DeviceID:  model.Token,
		Name:      model.Meta.Name,
		Active:    model.Active,
		Mine:      model.CognitoID == cognitoID,
		Model:     model.Meta.Model,
		Indoor:    model.Meta.Indoor,
		Location:  model.Meta.Coordinates,
		MapMeta:   model.MapMeta,
		Latest:    model.Measurements,
		Timestamp: model.Timestamp,
	}
	response.Data = data

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// CognitoData for user
func CognitoData(in map[string]interface{}) string {
	data, ok := in["claims"].(map[string]interface{})

	if !ok {
		return h.CognitoIDZeroValue
	}

	return data["sub"].(string)
}

func main() {
	lambda.Start(Handler)
}
