package main

import (
	"fmt"
	"strings"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	"github.com/appsmonkey/core.server.functions/integration/cognito"
	m "github.com/appsmonkey/core.server.functions/models"
	"github.com/appsmonkey/core.server.functions/tools/defaultDevice"
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID := h.CognitoIDZeroValue
	authHdr := header("AccessToken", req.Headers)
	if len(authHdr) > 0 {
		c, _, err := cog.ValidateToken(authHdr)
		if err != nil {
			fmt.Println(err)
		} else {
			cognitoID = c
		}
	}

	response := new(vm.DeviceListResponse)
	response.Init()

	if cognitoID == h.CognitoIDZeroValue {
		rd := make([]*vm.DeviceGetData, 0)

		// Add the default device on the top
		dd := defaultDevice.Get()
		rd = append(rd, &dd)

		response.Data = rd
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}

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

	rd := make([]*vm.DeviceGetData, 0)

	// Add the default device on the top
	dd := defaultDevice.Get()
	rd = append(rd, &dd)

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

func header(hdr string, in map[string]string) string {
	result, ok := in[hdr]
	if !ok {
		lwr := strings.ToLower(hdr)
		result = in[lwr]
	}

	return result
}

func init() {
	cog = cognito.NewCognito()
}

func main() {
	lambda.Start(Handler)
}
