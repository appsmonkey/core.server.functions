package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	bg "github.com/appsmonkey/core.server.functions/tools/guid"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"

	// Loading the sarajevo map
	z "github.com/appsmonkey/core.server.functions/tools/zones"
	_ "github.com/appsmonkey/core.server.functions/tools/zones/sarajevo"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.DeviceAddRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	device := m.Device{}
	device.Token = bg.New()
	device.DeviceID = ""
	device.CognitoID = CognitoData(req.RequestContext.Authorizer)
	device.Meta = request.Metadata
	device.Active = false
	device.ZoneID = "none"

	// If coordinates are set, then find the zone it belongs to
	if !device.Meta.Coordinates.IsEmpty() {
		if zone := z.ZoneByPoint(&z.Point{Lat: device.Meta.Coordinates.Lat, Lng: device.Meta.Coordinates.Lng}); zone != nil {
			device.ZoneID = zone.Title
		}
	}

	response.Data = vm.DeviceAddData{Token: device.Token}

	// insert data into the DB
	dal.Insert("devices", device)

	// Log and return result
	fmt.Println("Wrote item:  ", device)

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
