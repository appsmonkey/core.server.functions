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
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
// City list fetches minimal city preview, it does not require user to be authenticated
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := new(vm.CityListResponse)
	response.Init()

	var citiesTable = "cities"
	if value, ok := os.LookupEnv("dynamodb_table_cities"); ok {
		citiesTable = value
	}

	// fetch all cities with minimal data
	dbRes, err := dal.ListNoFilter(citiesTable, dal.Projection(dal.Name("city_id"), dal.Name("name"), dal.Name("country"), dal.Name("timestamp")))
	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Code: 0, Message: err.Error(), Data: ""})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	dbData := make([]m.City, 0)
	err = dbRes.Unmarshal(&dbData)

	if err != nil {
		fmt.Println(err)
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	rd := make([]*vm.CityGetDataMinimal, 0)

	for _, c := range dbData {
		data := vm.CityGetDataMinimal{
			CityID:    c.CityID,
			Country:   c.Country,
			Timestamp: c.Timestamp,
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
