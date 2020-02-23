package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(event events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {
	userSub := event.Request.UserAttributes["sub"]

	if event.TriggerSource == "CustomMessage_SignUp" {
		event.Response.EmailSubject = "Welcome to CityOS, please click the following link to verify your email"
		event.Response.EmailMessage = "Please click the link below to verify your email address."
		event.Response.EmailMessage += "<br/>"
		event.Response.EmailMessage += "<br/>"
		link := "https://links.cityos.io/auth/validate?client_id=" + event.CallerContext.ClientID
		link += "&user_name=" + event.UserName
		link += "&confirmation_code=" + event.Request.CodeParameter
		link += "&type=verify&cog_id=" + userSub.(string)
		event.Response.EmailMessage += fmt.Sprintf(`<a href="%s">Verify email</a>`, link)
		event.Response.EmailMessage += "<br/> If you can not open the link above please go to: <br/>" + link

	} else if event.TriggerSource == "CustomMessage_ForgotPassword" {
		event.Response.EmailSubject = "Password reset requested"
		event.Response.EmailMessage = "Password reset request, if this was you please go to the link below."
		event.Response.EmailMessage += "<br/>"
		event.Response.EmailMessage += "<br/>"
		link := "https://links.cityos.io/auth/validate?client_id=" + event.CallerContext.ClientID
		link += "&user_name=" + event.UserName
		link += "&confirmation_code=" + event.Request.CodeParameter
		link += "&type=pwreset&cog_id=" + userSub.(string)
		event.Response.EmailMessage += fmt.Sprintf(`<a href="%s">Reset password</a>`, link)
		event.Response.EmailMessage += "<br/> If you can not open the link above please go to: <br/>" + link
	}

	return event, nil
}

func main() {
	lambda.Start(Handler)
}
