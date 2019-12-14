package main

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID := CognitoData(req.RequestContext.Authorizer)
	fmt.Println("DEVICE_DEL_REQ: ", req.Body)
	request := new(vm.DeviceDelRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
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

	type resToUser struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	r := resToUser{Success: true, Message: ""}

	if cognitoID != model.CognitoID {
		r.Success = false
		r.Message = "this device does not belong to you"

		response.Data = r
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	err = dal.Delete("devices", map[string]*dal.AttributeValue{
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

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// CognitoData for user
func CognitoData(in map[string]interface{}) string {
	data := in["claims"].(map[string]interface{})

	return data["sub"].(string)
}

func main() {
	lambda.Start(Handler)
}
