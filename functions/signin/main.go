package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	dal "github.com/appsmonkey/core.server.functions/dal/access"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	"github.com/appsmonkey/core.server.functions/integration/cognito"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"net/http"
)

var (
	cog        *cognito.Cognito
	httpClient = &http.Client{}
)

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.SigninRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 403, Headers: response.Headers()}, nil
	}

	// 1. Check if we have social login data, if so then validate the token first
	fmt.Println("REQUEST_SOCIAL: ", request)
	if request.Social.HasData() {
		if request.Social.Type == "G" {
			data, err := cog.Google(request.Social.ID, request.Social.Token, request.Email, httpClient)
			if err != nil {
				errData := es.ErrRegistrationSignInError
				errData.Data = err.Error()
				response.Errors = append(response.Errors, errData)

				return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 403, Headers: response.Headers()}, nil
			}

			response.Data = data
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
		} else if request.Social.Type == "FB" {
			data, err := cog.Facebook(request.Social.ID, request.Social.Token, request.Email)
			if err != nil {
				errData := es.ErrRegistrationSignInError
				errData.Data = err.Error()
				response.Errors = append(response.Errors, errData)

				return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 403, Headers: response.Headers()}, nil
			}

			response.Data = data
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
		} else {
			errData := es.ErrRegistrationSignInError
			errData.Data = "Unrecognized signin method"
			response.Errors = append(response.Errors, errData)
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
		}

	} else {
		// do not allow social user to login with username and password
		_, email, _, _, suc, err := dal.CheckSocial(request.Password)
		if email == request.Email && !suc && err == nil {
			errData := es.ErrRegistrationSignInError
			errData.Data = err.Error()
			response.Errors = append(response.Errors, errData)

			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 403, Headers: response.Headers()}, nil
		}
	}

	data, err := cog.SignIn(request.Email, request.Password)
	if err != nil {
		errData := es.ErrRegistrationSignInError
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 403, Headers: response.Headers()}, nil
	}

	response.Data = data
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
	data, _ := json.Marshal(vm.SigninRequest{
		Email:    os.Getenv("USER_EMAIL"),
		Password: os.Getenv("USER_PASS"),
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
