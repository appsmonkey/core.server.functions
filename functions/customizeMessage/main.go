package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(event events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {

	if event.TriggerSource == "CustomMessage_SignUp" {
		event.Response.EmailMessage = "Welcome to CityOS, please click the following link to verify your email, this is a custom message"
		event.Response.EmailMessage = fmt.Sprintf(`Please click the link below to verify your email address. https://links.cityos.io/auth/validate?client_id=%s&user_name=%s&confirmation_code=%stype="verify"`, event.CallerContext.ClientID, event.UserName, event.Request.CodeParameter)
	} else if event.TriggerSource == "CustomMessage_ForgotPassword" {
		event.Response.EmailMessage = "Password reset requested"
		event.Response.EmailMessage = fmt.Sprintf(`Please click the link below to verify your email address. https://links.cityos.io/auth/validate?client_id=%s&user_name=%s&confirmation_code=%stype="pwreset"`, event.CallerContext.ClientID, event.UserName, event.Request.CodeParameter)
	}

	return event, nil
}

func main() {
	lambda.Start(Handler)
}
