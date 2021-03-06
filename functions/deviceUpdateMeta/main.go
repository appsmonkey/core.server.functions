package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	"github.com/appsmonkey/core.server.functions/integration/cognito"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"

	h "github.com/appsmonkey/core.server.functions/tools/helper"

	// Loading the sarajevo map
	z "github.com/appsmonkey/core.server.functions/tools/zones"
	_ "github.com/appsmonkey/core.server.functions/tools/zones/sarajevo"

	ss "github.com/aws/aws-sdk-go/aws/session"
	sl "github.com/aws/aws-sdk-go/service/lambda"
)

var lambdaClient *sl.Lambda
var (
	cog *cognito.Cognito
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

	var devicesTable = "devices"
	if value, ok := os.LookupEnv("dynamodb_table_devices"); ok {
		devicesTable = value
	}

	dbRes, err := dal.Get(devicesTable, map[string]*dal.AttributeValue{
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

	userGroups, err := cog.ListGroupsForUserFromID(CognitoData(req.RequestContext.Authorizer))

	isAdmin := false
	if err == nil {
		for _, g := range userGroups.Groups {
			if g.GroupName != nil && (*g.GroupName == "AdminGroup" || *g.GroupName == "SuperAdminGroup") {
				isAdmin = true
			}
		}
	}
	fmt.Println("Is Admin ::: ", isAdmin)

	if h.IsCognitoIDEmpty(device.CognitoID) {
		// TODO: Add so that only the admin user can do this
		// right now we are assigning the dvice to loged in user if it was not assigned before
		device.CognitoID = CognitoData(req.RequestContext.Authorizer)
	} else if !isAdmin && device.CognitoID != CognitoData(req.RequestContext.Authorizer) {
		response.Data = resData
		response.AddError(&es.Error{Message: "", Data: "device does not belong to you"})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	device.City = request.City

	oldZone := device.ZoneID
	// If coordinates are set, then find the zone it belongs to
	if !request.Coordinates.IsEmpty() {
		zone := z.ZoneByPoint(&z.Point{Lat: request.Coordinates.Lat, Lng: request.Coordinates.Lng})
		fmt.Println("ZONE NOT FOUND", zone)
		if zone != nil {
			device.ZoneID = "Sarajevo" + "@" + zone.Title
			// device.ZoneID = device.City + "@" + zone.Title

			device.City = "Sarajevo"
			device.Meta.Coordinates = request.Coordinates
		} else {
			device.Meta.Coordinates = request.Coordinates
			device.City = "Unknown"
			device.ZoneID = "none"
		}
	}

	fmt.Println("OLD_ZONE :: ", oldZone)
	fmt.Println("NEW_ZONE :::", device.ZoneID)

	device.Meta.Name = request.Name
	device.Meta.Model = request.Model
	device.Meta.Indoor = request.Indoor

	resData.Success = true
	response.Data = resData

	// insert data into the DB
	dal.Insert(devicesTable, device)

	// Update the old zone data
	if oldZone != "none" {
		payload := fmt.Sprintf(`{ "zone_id": "%v", "city_id": "%v"  }`, oldZone, device.City)

		invOut, err := lambdaClient.Invoke(&sl.InvokeInput{FunctionName: aws.String("CityOS-zoneUpdate-1H3L31K60T4LW"), Payload: []byte(payload)})
		if err != nil {
			fmt.Println("invOut", invOut)
			fmt.Println("err", err)
		}
	}

	// Update the new zone data if different
	if oldZone != device.ZoneID {
		payload := fmt.Sprintf(`{ "zone_id": "%v", "city_id": "%v" }`, device.ZoneID, device.City)

		invOut, err := lambdaClient.Invoke(&sl.InvokeInput{FunctionName: aws.String("CityOS-zoneUpdate-1H3L31K60T4LW"), Payload: []byte(payload)})
		if err != nil {
			fmt.Println("invOut", invOut)
			fmt.Println("err", err)
		}
	}

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

func init() {
	sess := ss.Must(ss.NewSessionWithOptions(ss.Options{
		SharedConfigState: ss.SharedConfigEnable,
	}))

	lambdaClient = sl.New(sess, &aws.Config{Region: aws.String("us-east-1")})
	cog = cognito.NewCognito()
}
