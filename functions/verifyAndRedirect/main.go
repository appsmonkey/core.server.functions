package main

import (
	"fmt"
	"net/http"

	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("VERIFICATION REQUEST: ", req.Body)
	request := new(vm.VerifyRedirectRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	// create verification URL
	verificationURL := "https://cityos.auth.us-east-1.amazoncognito.com/confirm?client_id" + request.ClientID + "&user_name=" + request.UserName + "&confirmation_code=" + request.ConfirmationCode + "&redirect_uri=https://dev.cityos.io"
	http.Get(verificationURL)

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 302, Headers: response.Headers()}, nil
}

func main() {
	lambda.Start(Handler)
}
