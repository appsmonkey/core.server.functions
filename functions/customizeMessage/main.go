package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(event events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {
	userSub := event.Request.UserAttributes["sub"]

	fmt.Println("USER SUB ::: ", userSub)

	if event.TriggerSource == "CustomMessage_SignUp" {
		event.Response.EmailSubject = "Welcome to CityOS, please click the following link to verify your email"
		event.Response.EmailMessage = fmt.Sprintf(
			`Please click the link below to verify your email address.\n\n
			https://links.cityos.io/auth/validate?client_id=%s&user_name=%s&confirmation_code=%s&type=verify&cog_id=%s`,
			event.CallerContext.ClientID, event.UserName, event.Request.CodeParameter, userSub,
		)
	} else if event.TriggerSource == "CustomMessage_ForgotPassword" {
		event.Response.EmailSubject = "Password reset requested"
		event.Response.EmailMessage = fmt.Sprintf(
			`Password reset request, if this was you please go to the link below.\n\n
			 https://links.cityos.io/auth/validate?client_id=%s&user_name=%s&confirmation_code=%s&type=pwreset&cog_id=%s`,
			event.CallerContext.ClientID, event.UserName, event.Request.CodeParameter, userSub,
		)
	}

	return event, nil
}

func main() {
	lambda.Start(Handler)
}
