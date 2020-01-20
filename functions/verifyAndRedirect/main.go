package main

import (
	"fmt"
	"net/http"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/avct/uasurfer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
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

	fmt.Println("User agent :::", req.Headers["User-Agent"])
	ua := uasurfer.Parse(req.Headers["User-Agent"])

	if request.Type == "verify" {
		// create verification URL
		verificationURL := "https://cityos.auth.us-east-1.amazoncognito.com/confirmUser?client_id=" + request.ClientID + "&user_name=" + request.UserName + "&response_type=code" + "&confirmation_code=" + request.ConfirmationCode
		fmt.Println("Verification URL:", verificationURL)

		_, err := http.Get(verificationURL)

		if err != nil {
			fmt.Println("Verification error: ", err)
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
		}
	}

	fmt.Println("Fetch User ::: ", request.UserName, request.CognitoID)
	res, err := dal.Get("users", map[string]*dal.AttributeValue{
		"cognito_id": {
			S: aws.String(request.CognitoID),
		},
		"email": {
			S: aws.String(request.UserName),
		},
	})

	if err != nil {
		fmt.Println("User missing error: ", err)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	user := new(m.User)
	res.Unmarshal(&user)

	fmt.Println("USER :::", user)
	if user.Attributes["cognito:user_status"] != "CONFIRMED" {
		fmt.Println("User not confirmed, verification failed.")
		response.AddError(&es.Error{Message: err.Error(), Data: "User not confirmed, verification failed or not attempted."})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 403, Headers: response.Headers()}, nil
	}

	response.Data = user

	headers := response.Headers()

	if ua.OS.Name.String() == "OSAndroid" {
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: headers}, nil
	} else if ua.OS.Name.String() == "OSiOS" || (ua.OS.Name.String() == "OSUnknown" && ua.DeviceType.String() == "DeviceUnknown") {
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: headers}, nil
	} else {
		fmt.Println("Default response ::: ", ua.OS.Name.String(), ua.OS.Platform.String(), ua.DeviceType.String(), ua.Browser.Name.String())

		route := "complete-registration"
		if request.Type == "verify" {
			route = "complete-registration"
		} else if request.Type == "pwreset" {
			route = "reset-password"
		} else if request.Type == "info" {
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: headers}, nil
		}

		headers["Location"] = "https://dev.cityos.io/" + route + "?username=" + user.Email + "&token=" + user.Token + "&id" + user.CognitoID
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 303, Headers: headers}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
