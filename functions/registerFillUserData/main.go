package main

import (
	"fmt"
	"os"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	"github.com/appsmonkey/core.server.functions/integration/cognito"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.RegisterFillUserDataRequest)
	response := request.Validate(req.Body)

	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	profileData, err := dynamodbattribute.MarshalMap(request.UserProfile)
	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Message: err.Error(), Code: 0, Data: "marshaling error"})
		response.Code = es.StatusProfileUpdateError
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	var usersTable = "users"
	if value, ok := os.LookupEnv("dynamodb_table_users"); ok {
		usersTable = value
	}

	res, err := dal.Get(usersTable, map[string]*dal.AttributeValue{
		"cognito_id": {
			S: aws.String(request.CognitoID),
		},
		"email": {
			S: aws.String(request.UserName),
		},
	})

	if err != nil {
		fmt.Println("User missing error: ", err)
		response.AddError(&es.Error{Message: err.Error(), Data: "User not found"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	user := new(m.User)
	res.Unmarshal(&user)

	if user.Token != request.Token {
		fmt.Println("Unauthorized request", user.Token, request.Token)
		response.AddError(&es.Error{Message: "Token invalid, unauthorized", Data: "Unauthorized request"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	_, err = cog.SetUserPassword(request.UserName, request.Password, true)
	fmt.Println("Password set ::: ", request.Password, err)

	if err != nil {
		fmt.Println("Set password error")
		response.AddError(&es.Error{Message: err.Error(), Data: "Set password error"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	err = dal.Update(usersTable, "set profile = :p",
		map[string]*dal.AttributeValue{
			"cognito_id": {
				S: aws.String(request.CognitoID),
			},
			"email": {
				S: aws.String(request.UserName),
			},
		}, map[string]*dal.AttributeValue{
			":p": {
				M: profileData,
			},
		})

	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Message: err.Error(), Code: 0, Data: "db error"})
		response.Code = es.StatusProfileUpdateError
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	//login user after successful profile update
	loginRes, err := cog.SignIn(request.UserName, request.Password)
	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Message: err.Error(), Code: 0, Data: "login error, profile update is succesful"})
		response.Code = es.StatusProfileUpdateError
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	response.Data = loginRes

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func init() {
	cog = cognito.NewCognito()
}

func main() {
	lambda.Start(Handler)
}
