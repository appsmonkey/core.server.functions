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
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log body and pass to the DAO
	fmt.Println("Received body: ", req.Body)

	request := new(vm.RegisterRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500}, nil
	}

	var usersTable = "users"
	if value, ok := os.LookupEnv("dynamodb_table_users"); ok {
		usersTable = value
	}

	// insert data into the DB
	dal.Insert(usersTable, request)

	// Log and return result
	fmt.Println("Wrote item:  ", request)
	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
