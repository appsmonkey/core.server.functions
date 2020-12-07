package main

import (
	"fmt"
	"os"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID, email := CognitoData(req.RequestContext.Authorizer)
	request := new(vm.CognitoProfileUpdateRequest)
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

	err = dal.Update(usersTable, "set profile = :p",
		map[string]*dal.AttributeValue{
			"cognito_id": {
				S: aws.String(cognitoID),
			},
			"email": {
				S: aws.String(email),
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

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// CognitoData for user
func CognitoData(in map[string]interface{}) (string, string) {
	data := in["claims"].(map[string]interface{})

	return data["sub"].(string), data["email"].(string)
}

func main() {
	lambda.Start(Handler)
}
