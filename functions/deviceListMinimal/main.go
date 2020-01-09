package main

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
	"github.com/appsmonkey/core.server.functions/tools/defaultDevice"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID := CognitoData(req.RequestContext.Authorizer)
	response := new(vm.DeviceListResponse)
	response.Init()

	res, err := dal.GetFromIndex("devices", "CognitoID-index", dal.Condition{
		"cognito_id": {
			ComparisonOperator: aws.String("EQ"),
			AttributeValueList: []*dal.AttributeValue{
				{
					S: aws.String(cognitoID),
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	City := req.QueryStringParameters["city"]
	if len(City) < 1 {
		// Sarajevo is the default city
		City = "Sarajevo"
	}

	dbData := make([]m.Device, 0)
	err = res.Unmarshal(&dbData)
	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	rd := make([]*vm.DeviceGetDataMinimal, 0)

	// Add the default device on the top
	dd := defaultDevice.GetMinimal(City)
	rd = append(rd, &dd)

	for _, d := range dbData {
		data := vm.DeviceGetDataMinimal{
			DeviceID: d.Token,
			Name:     d.Meta.Name,
			Active:   d.Active,
			Model:    d.Meta.Model,
			Indoor:   d.Meta.Indoor,
		}
		rd = append(rd, &data)
	}

	response.Data = rd

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
