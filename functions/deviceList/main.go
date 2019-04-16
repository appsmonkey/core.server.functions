package main

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
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

	dbData := make([]m.Device, 0)
	err = res.Unmarshal(&dbData)
	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}
	// if len(dbData) == 0 {
	// 	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	// }

	// dbRes, err := dal.List("devices", dal.Name("cognito_id").Equal(dal.Value(cognitoID)), dal.Projection(dal.Name("token"), dal.Name("device_id"), dal.Name("meta"), dal.Name("map_meta"), dal.Name("active"), dal.Name("measurements")))
	// dbData := make([]m.Device, 0)
	// err = dbRes.Unmarshal(&dbData)
	// if err != nil {
	// 	fmt.Println(err)
	// 	response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
	// 	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	// }
	// if len(dbData) == 0 {
	// 	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	// }

	rd := make([]*vm.DeviceGetData, 0)

	for _, d := range dbData {
		data := vm.DeviceGetData{
			DeviceID:  d.Token,
			Name:      d.Meta.Name,
			Active:    d.Active,
			Mine:      d.CognitoID == cognitoID,
			Model:     d.Meta.Model,
			Indoor:    d.Meta.Indoor,
			Location:  d.Meta.Coordinates,
			MapMeta:   d.MapMeta,
			Timestamp: d.Timestamp,
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
