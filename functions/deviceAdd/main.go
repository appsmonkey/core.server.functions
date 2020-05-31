package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
	bg "github.com/appsmonkey/core.server.functions/tools/guid"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-sdk-go/aws"

	// Loading the sarajevo map
	z "github.com/appsmonkey/core.server.functions/tools/zones"
	_ "github.com/appsmonkey/core.server.functions/tools/zones/sarajevo"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Request ADD_DEVICE: ", req.Body)
	request := new(vm.DeviceAddRequest)
	response := request.Validate(req.Body)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	existingDevice := m.Device{}
	res, err := dal.Get("devices", map[string]*dal.AttributeValue{
		"token": {
			S: aws.String(request.Token),
		},
	})
	if err != nil {
		fmt.Println("Existing device not found ::: ", request.Token)
	} else {
		err = res.Unmarshal(&existingDevice)
		if err != nil {
			fmt.Println(err)
			response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
		}
	}

	device := m.Device{}

	if len(existingDevice.Token) > 0 {
		device = existingDevice
	}

	device.Token = request.Token
	if len(device.Token) == 0 {
		device.Token = bg.New()
	}
	device.DeviceID = device.Token
	// FIXME: If cognito id is bad ??
	device.CognitoID = CognitoData(req.RequestContext.Authorizer)
	device.Meta = request.Metadata
	device.ZoneID = "none"

	// We can add manually or we can check with lat lon
	// if len(request.City) > 0 {
	// 	device.City = h.MapCity(h.TransformCityString(request.City))
	// } else {
	// 	device.City = "Sarajevo" // default value is Sarajevo
	// }

	// If coordinates are set, then find the zone it belongs to
	if !device.Meta.Coordinates.IsEmpty() {
		if zone := z.ZoneByPoint(&z.Point{Lat: device.Meta.Coordinates.Lat, Lng: device.Meta.Coordinates.Lng}); zone != nil {
			device.ZoneID = "Sarajevo" + "@" + zone.Title
			device.City = "Sarajevo"
			// device.ZoneID = device.City + "@" + zone.Title
		} else {
			device.City = "Unknown"
			device.ZoneID = "none"
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
