package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	bg "github.com/appsmonkey/core.server.functions/tools/guid"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"

	// Loading the sarajevo map

	_ "github.com/appsmonkey/core.server.functions/tools/zones/sarajevo"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.CityAddRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	city := m.City{}
	city.CityID = request.CityID
	if len(city.CityID) == 0 {
		city.CityID = bg.New()
	}

	// if this is true we are updating existing city
	if len(city.CityID) > 0 {
		city.CityID = request.CityID
	}

	city.Country = request.Country
	city.Name = request.Name

	// Zones can be empty and updated later, or initally set
	if request.Zones != nil && len(request.Zones) > 0 {
		// Check if each zones exists and add to city, if any zone missing break
		for _, z := range request.Zones {
			zoneRes, err := dal.List("zones", dal.Name("zone_id").Equal(dal.Value(z)), dal.Projection(dal.Name("zone_id"), dal.Name("data")))

			if err != nil {

			}

		}
	} else {
		city.Zones = make([]string, 0)
	}

	response.Data = vm.CityAddData{Token: city.CityID}

	// insert data into the DB
	dal.Insert("cities", city)

	// Log and return result
	fmt.Println("Wrote item:  ", city)

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// CognitoData for user
func CognitoData(in map[string]interface{}) string {
	data := in["claims"].(map[string]interface{})

	return data["sub"].(string)
}

func main() {
	// zone := z.ZoneByPoint(&z.Point{Lat: 43.8444278, Lng: 18.408692})
	// fmt.Println(zone.Title)
	lambda.Start(Handler)
}
