package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
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
	fmt.Println("SIGNUP_REQUEST: ", req.Body)
	request := new(vm.SignupRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	// Register User in Cognito
	signupData, err := cog.SignUpWithVerif(request.Email, request.Password, request.Gender, request.FirstName, request.LastName)
	if err != nil {
		errData := es.ErrRegistrationCognitoSignupError
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	// Now save it into our DB
	cogReq := new(vm.CognitoRegisterRequest)
	cogResponse := cogReq.ValidateCognitoWithVerif(signupData)
	fmt.Println("CHECK", cogReq)
	if cogResponse.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	// insert data into the DB
	if os.Getenv("ENV") != "local" {
		dal.Insert("users", cogReq)
	}

	response.Data = signupData
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
	data, _ := json.Marshal(vm.SignupRequest{
		Email:     os.Getenv("USER_EMAIL"),
		Password:  os.Getenv("USER_PASS"),
		Gender:    os.Getenv("USER_GENDER"),
		FirstName: os.Getenv("USER_FIRSTNAME"),
		LastName:  os.Getenv("USER_LASTNAME"),
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
