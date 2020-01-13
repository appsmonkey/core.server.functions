package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID := CognitoData(req.RequestContext.Authorizer)
	request := new(vm.CityAddRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	fmt.Println("Cognito ID ::: ", cognitoID)

	// FIXME: Only admin can CUD cities - check user cognito user pool when set up
	// type resToUser struct {
	// 	Success bool   `json:"success"`
	// 	Message string `json:"message"`
	// }

	// r := resToUser{Success: true, Message: ""}
	// if cognitoID == h.CognitoIDZeroValue {
	// 	r.Success = false
	// 	r.Message = "no permissions to add new city"

	// 	response.Data = r
	// 	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	// }

	city := m.City{}
	city.CityID = request.CityID

	// if this is true we are updating existing city
	if len(city.CityID) > 0 {
		city.CityID = request.CityID

		// check if city exists
		existingCity, err := dal.Get("cities", map[string]*dal.AttributeValue{
			"city_id": {
				S: aws.String(request.CityID),
			},
		})

		if err == nil {
			var c m.City
			err = existingCity.Unmarshal(&c)
			fmt.Println("City already exists ::: ", c, err)
			errData := es.ErrCityAlreadyExists
			response.Errors = append(response.Errors, errData)
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 503, Headers: response.Headers()}, nil
		}
	}

	city.Country = request.Country
	response.Data = vm.CityAddData{CityID: city.CityID}

	// insert data into the DB
	dal.Insert("cities", city)

	// Log and return result
	fmt.Println("Wrote item:  ", city)

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// CognitoData for user
func CognitoData(in map[string]interface{}) string {
	data, ok := in["claims"].(map[string]interface{})

	if !ok {
		return h.CognitoIDZeroValue
	}

	return data["sub"].(string)
}

func main() {
	lambda.Start(Handler)
}
