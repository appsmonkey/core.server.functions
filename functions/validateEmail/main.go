package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/appsmonkey/core.server.functions/integration/cognito"
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
	request := new(vm.ValidateEmailRequest)
	response := request.Validate(req.Body)

	type resData struct {
		Exists    bool `json:"exists"`
		Confirmed bool `json:"confirmed"`
	}

	res := resData{Exists: false, Confirmed: false}

	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		response.Data = res
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	data, err := cog.Profile(request.Email)
	if err != nil {
		response.Data = res
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	res.Exists = true
	res.Confirmed = data.UserStatus != nil && *data.UserStatus == "CONFIRMED"

	response.Data = res
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
	data, _ := json.Marshal(vm.ValidateEmailRequest{
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
		// a, b := mmm.SensorReading("23", "7", 200)
		// fmt.Println(a.Name)
		// fmt.Println(a.Unit)
		// fmt.Println(b)
		// fmt.Println()
		// fmt.Println(mmm.MarshalSchema())

		return
	}

	lambda.Start(Handler)
}
