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
	fmt.Println("VERIFICATION REQUEST: ", req.QueryStringParameters)
	request := new(vm.VerifyRedirectRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	// redired URL - default
	redirectURL := "https://dev.cityos.io"

	switch request.ClientID {
	// - Android
	case "km0afsc8ua4f0bc56brcn7t90":
		redirectURL = "https://dev.cityos.io" // - TODO: change URL for Android
	// - IOS
	case "70mq6uphtmmorkjt74ei0rj5fr":
		redirectURL = "https://dev.cityos.io" // - TODO: change URL for IOS
	}

	// create verification URL
	verificationURL := "https://cityos.auth.us-east-1.amazoncognito.com/confirmUser?client_id=" + request.ClientID + "&user_name=" + request.UserName + "&response_type=code" + "&confirmation_code=" + request.ConfirmationCode + "&redirect_uri=" + redirectURL
	fmt.Println("VERIFICATION URL:", verificationURL)

	_, err := http.Get(verificationURL)
	if err != nil {
		fmt.Println("Verification error: ", err)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	headers := response.Headers()
	headers["Location"] = redirectURL
	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 302, Headers: headers}, nil
}

func main() {
	lambda.Start(Handler)
}
