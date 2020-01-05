package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(event events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessageResponse, error) {
	fmt.Println("CUSTOM MESSAGE REQ EVET: ", event)
	fmt.Println("CUSTOM MESSAGE REQ EVET CONTEXT: ", event.CallerContext.ClientID)

	if event.TriggerSource == "CustomMessage_SignUp" {

		fmt.Println("SIGNUP EVENT TRIGGERED")
		event.Response.EmailMessage = "Welcome to CityOS, please click the following link to verify your email, this is a custom message"
		event.Response.EmailMessage = fmt.Sprintf(`Please click the link below to verify your email address. https://apigway.com/auth/validate?client_id=%s&user_name=%s&confirmation_code=%s`, event.CallerContext.ClientID, event.UserName, event.Request.CodeParameter)

		return event.Response, nil
	}

	return event.Response, nil
}

func main() {
	lambda.Start(Handler)
}
