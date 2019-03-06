package main

import (
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
	cognitoID, email := CognitoData(req.RequestContext.Authorizer)
	response := new(vm.CognitoProfileListResponse)
	response.Init()

	res, err := dal.Get("users", map[string]*dal.AttributeValue{
		"cognito_id": {
			S: aws.String(cognitoID),
		},
		"email": {
			S: aws.String(email),
		},
	})
	if err != nil {
		response.AddError(&es.Error{Code: 0, Message: err.Error()})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	model := m.User{}
	err = res.Unmarshal(&model)
	if err != nil {
		response.AddError(&es.Error{Code: 0, Message: err.Error()})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	response.Data = model

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// CognitoData for user
func CognitoData(in map[string]interface{}) (string, string) {
	data := in["claims"].(map[string]interface{})

	return data["sub"].(string), data["email"].(string)
}

func main() {
	lambda.Start(Handler)
}
