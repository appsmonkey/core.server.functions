package main

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID := CognitoData(req.RequestContext.Authorizer)
	request := new(vm.DeviceGetRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	fmt.Print("cid", cognitoID)
	fmt.Print("token", request.Token)

	res, err := dal.Get("devices", map[string]*dal.AttributeValue{
		"token": {
			S: aws.String(request.Token),
		},
		"cognito_id": {
			S: aws.String(cognitoID),
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "1: " + err.Error(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	model := m.Device{}
	err = res.Unmarshal(&model)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "2: " + err.Error(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	// remove the user-id (do not expose it)
	model.CognitoID = ""
	response.Data = model

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
