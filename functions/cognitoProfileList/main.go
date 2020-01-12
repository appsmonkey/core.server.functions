package main

import (
	"fmt"
	"log"
	"os"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	"github.com/appsmonkey/core.server.functions/integration/cognito"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/joho/godotenv"
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID, email := CognitoData(req.RequestContext.Authorizer)
	response := new(vm.CognitoProfileListResponse)
	response.Init()

	res, err := dal.Get("users", map[string]*dal.AttributeValue{
		"cognito_id": {
			S: aws.String(cognitoID),
		},
		"email": {
			S: aws.String(email),
		},
	})
	if err != nil {
		response.AddError(&es.Error{Code: 0, Message: err.Error()})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	model := m.User{SocialID: "none", SocialType: "none"}
	err = res.Unmarshal(&model)
	if err != nil {
		response.AddError(&es.Error{Code: 0, Message: err.Error()})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	usrGroups, err := cog.ListGroupsForUser(email)
	if err != nil {
		fmt.Println("Fetch user groups error ::: ", err)
		errData := es.ErrProfileMissingEmail
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	response.Groups = usrGroups
	response.Data = model.Profile

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func init() {
	fmt.Println("INIT TRIGGERED")
	if os.Getenv("ENV") == "local" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("error loading .env: %v\n", err)
		}
	}

	cog = cognito.NewCognito()
}

// CognitoData for user
func CognitoData(in map[string]interface{}) (string, string) {
	data := in["claims"].(map[string]interface{})

	return data["sub"].(string), data["email"].(string)
}

func main() {
	lambda.Start(Handler)
}
