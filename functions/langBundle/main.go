package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/appsmonkey/core.server.functions/integration/cognito"
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
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

func init() {
	if os.Getenv("ENV") == "local" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("error loading .env: %v\n", err)
		}
	}

	cog = cognito.NewCognito()
}

func local() {
	data, _ := json.Marshal(vm.ForgotPasswordStartRequest{
		Email: os.Getenv("USER_EMAIL"),
	})

	resp, err := Handler(context.Background(), events.APIGatewayProxyRequest{
		Body: string(data),
	})

	if err != nil {
		fmt.Printf("unhandled error! \nError: %v\n", err)
	} else {
		j, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(j))
	}
}

func main() {
	if os.Getenv("ENV") == "local" {
		local()
		return
	}

	lambda.Start(Handler)
}
