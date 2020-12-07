package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/appsmonkey/core.server.functions/dal"
	dala "github.com/appsmonkey/core.server.functions/dal/access"
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

	var usersTable = "users"
	if value, ok := os.LookupEnv("dynamodb_table_users"); ok {
		usersTable = value
	}

	sid, st, _, err := dala.GetTempUser(request.Email)
	if err != nil {
		fmt.Println("email register", request)
		// We do not have a temp user
		// insert data into the DB as is
		dal.Insert(usersTable, request)

		return req, nil
	}

	request.SocialID = sid
	request.SocialType = st

	fmt.Println("social register", request)
	dal.Insert(usersTable, request)

	return req, nil
}

func main() {
	lambda.Start(Handler)
}
