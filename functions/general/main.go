package main

import (
	"fmt"
	"time"

	"os"

	"github.com/appsmonkey/core.server.functions/dal"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	apiKey = os.Getenv("API_KEY")
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log body and pass to the DAO
	fmt.Printf("Received body: %v\n", req)

	request := new(vm.GeneralRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500}, nil
	}

	request.Date = time.Now().Unix()

	var mainTable = "main"
	if value, ok := os.LookupEnv("dynamodb_table_main"); ok {
		mainTable = value
	}

	// insert data into the DB
	dal.Insert(mainTable, request)

	// Log and return result
	fmt.Println("Wrote item:  ", request)
	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
