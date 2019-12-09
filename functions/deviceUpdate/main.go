package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	sl "github.com/aws/aws-sdk-go/service/lambda"
)

var lambdaClient *sl.Lambda

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req interface{}) error {
	input, ok := req.(map[string]interface{})
	fmt.Println("DEVICE_UPDATE: ", input)
	if !ok {
		err := errors.New("incorrect data received. input has incorrect format")
		fmt.Println(err)
		return err
	}

	timestamp := input["timestamp"].(float64)

	state, ok := input["state"].(map[string]interface{})

	if !ok {
		err := errors.New("incorrect data received. 'state' field is missing")
		fmt.Println(err)
		return err
	}

	// reported, ok := state["reported"].(map[string]interface{})
	// if !ok {
	// 	err := errors.New("incorrect data received. 'reported' field is missing")
	// 	fmt.Println(err)
	// 	return err
	// }

	type data struct {
		Token        string
		DeviceID     string
		DeviceType   string
		Measurements map[string]interface{}
	}

	type schemaData struct {
		Data    m.Schema
		Version string
	}

	deviceData := data{
		Token:        input["token"].(string),
		DeviceID:     input["token"].(string),
		DeviceType:   "BOXY",
		Measurements: state["reported"].(map[string]interface{}),
	}

	dbRes, err := dal.Get("devices", map[string]*dal.AttributeValue{
		"token": {
			S: aws.String(deviceData.Token),
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	device := new(m.Device)
	err = dbRes.Unmarshal(device)

	if err != nil {
		fmt.Println(err)
		return err
	}

	device.Active = true
	device.Token = deviceData.Token
	device.DeviceID = deviceData.DeviceID
	device.Timestamp = timestamp
	if len(device.MapMeta) == 0 {
		device.MapMeta = make(map[string]m.MapMeta, 0)
	}

	if len(device.Measurements) == 0 {
		device.Measurements = make(map[string]interface{}, 0)
	}

	// TODO: we should use schema to refer to measurement units
	schemaRes, err := dal.Get("schema", map[string]*dal.AttributeValue{
		"version": {
			S: aws.String("1"),
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	schema := new(schemaData)
	err = schemaRes.Unmarshal(schema)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for k, v := range deviceData.Measurements {
		var mk string = k
		mv, _ := strconv.ParseFloat(v.(string), 64)

		fieldData := schema.Data[mk]
		level := fieldData.Result(mv)
		// level := m.Level(mk, mv)

		fmt.Println("FIELD DATA", fieldData.Name, fieldData.Unit, fieldData.Result(mv))
		mm := m.MapMeta{
			Level:       level,
			Value:       mv,
			Measurement: fieldData.Name,
			Unit:        fieldData.Unit,
			// Measurement: m.MeasureMapName[mk],
			// Unit:        m.MeasureMapUnit[mk],
		}

		// Update map meta for the sensor
		device.MapMeta[mk] = mm

		// Update the meassurement for the sensor
		device.Measurements[mk] = mv
		fmt.Println("TEST TEST", mm)

	}

	// Since zone_id is an index, we need to have some value in it
	if len(device.ZoneID) == 0 {
		device.ZoneID = "none"
	}

	// Since cognito_id is an index, we need to have some value in it
	if len(device.CognitoID) == 0 {
		device.CognitoID = h.CognitoIDZeroValue
	}

	err = dal.Insert("devices", device)
	if err == nil && len(device.ZoneID) > 0 && device.ZoneID != "none" {
		payload := fmt.Sprintf(`{ "zone_id": "%v" }`, device.ZoneID)

		invOut, err := lambdaClient.Invoke(&sl.InvokeInput{FunctionName: aws.String("CityOSZoneUpdate"), Payload: []byte(payload)})
		if err != nil {
			fmt.Println("invOut", invOut)
			fmt.Println("err", err)
		}
	}

	// { // Run the calculations
	// 	payload := fmt.Sprintf(`{ "value": "%v" }`, device.Measurements["7"])

	// 	invOut, err := lambdaClient.Invoke(&sl.InvokeInput{FunctionName: aws.String("CityOSTest"), Payload: []byte(payload)})
	// 	if err != nil {
	// 		fmt.Println("invOut", invOut)
	// 		fmt.Println("err", err)
	// 	}
	// }

	dal.Insert("live", device.ToLiveData())

	return nil
}

func main() {
	lambda.Start(Handler)
}

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	lambdaClient = sl.New(sess, &aws.Config{Region: aws.String("eu-west-1")})
}
