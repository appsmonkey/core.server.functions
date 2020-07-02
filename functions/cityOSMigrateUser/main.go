package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/appsmonkey/core.server.functions/integration/cognito"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	cog        *cognito.Cognito
	httpClient = &http.Client{}
)

// Handler will handle our request comming from the API gateway
func Handler(event events.CognitoEventUserPoolsMigrateUser) (events.CognitoEventUserPoolsMigrateUser, error) {

	if event.TriggerSource == "UserMigration_Authentication" {
		// user migrates during login flow

		// authenticate the user with your existing user directory service
		user, err := cog.SignIn(event.UserName, event.CognitoEventUserPoolsMigrateUserRequest.Password)

		if err != nil {
			fmt.Println("Login error")
		}

		if user != nil {
			event.CognitoEventUserPoolsMigrateUserResponse.UserAttributes["email"] = *user.UserData.User.Username
			event.CognitoEventUserPoolsMigrateUserResponse.UserAttributes["email_verified"] = "true"

			event.CognitoEventUserPoolsMigrateUserResponse.FinalUserStatus = "CONFIRMED"
			event.CognitoEventUserPoolsMigrateUserResponse.MessageAction = "SUPPRESS"

			return event, nil
		} else {
			// Return error to Amazon Cognito
			fmt.Println("Bad password, auth failed")
			return event, errors.New("Bad password, auth failed")
		}
	} else if event.TriggerSource == "UserMigration_ForgotPassword" {
		// user migrates during forgot passowrd flow

		// Lookup the user in your existing user directory service
		user, err := cog.Profile(event.UserName)

		if err != nil {
			// return to Amazon cognito
			fmt.Println("FP-flow, User not found, returning to Amazon Cognito")
		}

		if user != nil {
			event.CognitoEventUserPoolsMigrateUserResponse.UserAttributes["email"] = *user.Username

			// required to enable password-reset code to be sent to user
			event.CognitoEventUserPoolsMigrateUserResponse.UserAttributes["email_verified"] = "true"
			event.CognitoEventUserPoolsMigrateUserResponse.MessageAction = "SUPPRESS"
			return event, nil
		} else {
			// Return error to Amazon Cognito
			return event, errors.New("FP-flow, User not found, returning to Amazon Cognito")
		}
	} else {
		// Return error to Amazon Cognito
		return event, errors.New("Bad TriggerSource " + event.TriggerSource)
	}
}

func main() {
	lambda.Start(Handler)
}
