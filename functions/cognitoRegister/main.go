package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/appsmonkey/core.server.functions/dal"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
)

var (
	apiKey1 = os.Getenv("API_KEY")
)

// Handler will handle our request comming from the API gateway
func Handler(req events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	request := new(vm.CognitoRegisterRequest)
	response := request.Validate(&req)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v\n", response.Errors, response.RequestID)

		return req, nil
	}

	// insert data into the DB
	dal.Insert("users", request)

	// // Log and return result
	fmt.Println("Wrote item:  ", request)

	return req, nil
}

func main() {
	lambda.Start(Handler)
}
