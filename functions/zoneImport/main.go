package main

import (
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) error {

	return nil
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
