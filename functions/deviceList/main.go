package main

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID := CognitoData(req.RequestContext.Authorizer)
	response := new(vm.DeviceListResponse)
	response.Init()

	dbRes, err := dal.List("devices", dal.Name("cognito_id").Equal(dal.Value(cognitoID)), dal.Projection(dal.Name("token"), dal.Name("device_id"), dal.Name("meta"), dal.Name("map_meta"), dal.Name("active"), dal.Name("measurements")))
	dbData := make([]m.Device, 0)
	err = dbRes.Unmarshal(&dbData)
	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}
	if len(dbData) == 0 {
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	response.Data = dbData

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// CognitoData for user
func CognitoData(in map[string]interface{}) string {
	data := in["claims"].(map[string]interface{})

	return data["sub"].(string)
}

func main() {
	lambda.Start(Handler)
}
