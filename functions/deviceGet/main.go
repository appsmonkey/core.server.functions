package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	"github.com/appsmonkey/core.server.functions/integration/cognito"
	m "github.com/appsmonkey/core.server.functions/models"
	defaultDevice "github.com/appsmonkey/core.server.functions/tools/defaultDevice"
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
	request := new(vm.DeviceGetRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	if len(authHdr) > 0 {
		c, _, isExpired, err := cog.ValidateToken(authHdr)
		if err != nil {
			fmt.Println(err)
			if isExpired {
				return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 401, Headers: response.Headers()}, nil
			}
		} else {
			cognitoID = c
		}
	}

	if len(request.Token) == 0 {
		response.Data = defaultDevice.Get("Sarajevo")
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}

	var devicesTable = "devices"
	if value, ok := os.LookupEnv("dynamodb_table_devices"); ok {
		devicesTable = value
	}

	res, err := dal.Get(devicesTable, map[string]*dal.AttributeValue{
		"token": {
			S: aws.String(request.Token),
		},
	})
	if err != nil {
		errData := es.ErrDeviceNotFound
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	model := m.Device{}
	err = res.Unmarshal(&model)
	if err != nil {
		errData := es.ErrDeviceNotFound
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	fmt.Println(model.CognitoID, cognitoID, "print creds")

	data := vm.DeviceGetData{
		DeviceID:  model.Token,
		Name:      model.Meta.Name,
		Active:    model.Active,
		Mine:      model.CognitoID == cognitoID, // FIX ME - Mine is always false even
		Model:     model.Meta.Model,
		Indoor:    model.Meta.Indoor,
		Location:  model.Meta.Coordinates,
		MapMeta:   model.MapMeta,
		Latest:    model.Measurements,
		Timestamp: model.Timestamp,
	}

	// if ID missing then there is no device
	if data.DeviceID == "" {
		errData := es.ErrDeviceNotFound
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	response.Data = data

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func header(hdr string, in map[string]string) string {
	result, ok := in[hdr]
	if !ok {
		lwr := strings.ToLower(hdr)
		result = in[lwr]
	}

	return result
}

// CognitoData for user
func CognitoData(in map[string]interface{}) string {
	data, ok := in["claims"].(map[string]interface{})

	if !ok {
		return h.CognitoIDZeroValue
	}

	return data["sub"].(string)
}

func init() {
	cog = cognito.NewCognito()
}

func main() {
	lambda.Start(Handler)
}
