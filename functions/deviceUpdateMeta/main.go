package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"

	// Loading the sarajevo map
	z "github.com/appsmonkey/core.server.functions/tools/zones"
	_ "github.com/appsmonkey/core.server.functions/tools/zones/sarajevo"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.DeviceUpdateMetaRequest)
	resData := vm.DeviceUpdateMetaData{Success: false}
	response := request.Validate(req.Body)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		response.Data = resData
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	dbRes, err := dal.Get("devices", map[string]*dal.AttributeValue{
		"token": {
			S: aws.String(request.Token),
		},
	})
	if err != nil {
		fmt.Println(err)
		response.Data = resData
		response.AddError(&es.Error{Message: err.Error(), Data: "could not find the provided device"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	device := m.Device{}
	err = dbRes.Unmarshal(&device)
	if err != nil {
		fmt.Println(err)
		response.Data = resData
		response.AddError(&es.Error{Message: err.Error(), Data: "could not find the provided device"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	if len(device.CognitoID) == 0 {
		// TODO: Add so that only the admin user can do this
		// right now we are assigning the dvice to loged in user if it was not assigned before
		device.CognitoID = CognitoData(req.RequestContext.Authorizer)
	} else if device.CognitoID != CognitoData(req.RequestContext.Authorizer) {
		response.Data = resData
		response.AddError(&es.Error{Message: "", Data: "device does not belong to you"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	// If coordinates are set, then find the zone it belongs to
	if !request.Coordinates.IsEmpty() {
		var lat, lng float64
		if s, err := strconv.ParseFloat(device.Meta.Coordinates.Lat, 64); err == nil {
			lat = s
		}
		if s, err := strconv.ParseFloat(device.Meta.Coordinates.Lng, 64); err == nil {
			lng = s
		}

		if lng > 0 && lat > 0 {
			zone := z.ZoneByPoint(&z.Point{Lat: lat, Lng: lng})
			device.ZoneID = zone.Title
			device.Meta.Coordinates = request.Coordinates
		}
	}

	device.Meta.Name = request.Name
	device.Meta.Model = request.Model
	device.Meta.Indoor = request.Indoor

	resData.Success = true
	response.Data = resData

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
	lambda.Start(Handler)
}
