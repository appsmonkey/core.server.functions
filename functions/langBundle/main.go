package main

import (
	"context"
	"os"

	"github.com/appsmonkey/core.server.functions/integration/cognito"
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.LangBundleRequest)
	response := request.Validate(req.QueryStringParameters)

	response.Data = h.GetLang(request.Language)
	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func main() {
	if os.Getenv("ENV") == "local" {
		return
	}

	lambda.Start(Handler)
}
