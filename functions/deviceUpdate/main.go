package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	envType = os.Getenv("ENV_TYPE")
)

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req interface{}) error {
	input, ok := req.(map[string]interface{})
	if !ok {
		err := errors.New("incorrect data received. input has incorrect format")
		fmt.Println(err)
		return err
	}

	state, ok := input["state"].(map[string]interface{})
	if !ok {
		err := errors.New("incorrect data received. 'state' field is missing")
		fmt.Println(err)
		return err
	}

	desired, ok := state["desired"].(map[string]interface{})
	if !ok {
		err := errors.New("incorrect data received. 'desired' field is missing")
		fmt.Println(err)
		return err
	}

	type data struct {
		Token        string
		DeviceID     string
		DeviceType   string
		Measurements []interface{}
	}

	deviceData := data{
		Token:        desired["token"].(string),
		DeviceID:     desired["device_id"].(string),
		DeviceType:   desired["device_type"].(string),
		Measurements: desired["measurements"].([]interface{}),
	}

	dbRes, err := dal.List("devices", dal.Name("token").Equal(dal.Value(deviceData.Token)), dal.Projection(dal.Name("token"), dal.Name("device_id"), dal.Name("meta"), dal.Name("map_meta"), dal.Name("active"), dal.Name("measurements"), dal.Name("cognito_id")))
	dbData := make([]m.Device, 0)
	err = dbRes.Unmarshal(&dbData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if len(dbData) == 0 {
		return err
	}

	device := dbData[0]
	device.Active = true
	device.DeviceID = deviceData.DeviceID
	if len(device.MapMeta) == 0 {
		device.MapMeta = make(map[string]m.MapMeta, 0)
	}
	if len(device.Measurements) == 0 {
		device.Measurements = make(map[string]interface{}, 0)
	}

	for _, v := range deviceData.Measurements {
		data := v.(map[string]interface{})
		var mk string
		var mv float64
		for i, j := range data {
			mk = i
			mv = j.(float64)
			break
		}

		value := fmt.Sprintf("%f", mv)
		level := m.Level(mk, mv)
		mm := m.MapMeta{
			Name:        device.Meta.Name,
			Level:       level,
			Coordinates: device.Meta.Coordinates,
			Value:       value,
			Icon:        m.LevelIconMap[level],
			Measurement: m.MeasureMapName[mk],
			Unit:        m.MeasureMapUnit[mk],
		}

		// Update map meta for the sensor
		device.MapMeta[mk] = mm

		// Update the meassurement for the sensor
		device.Measurements[mk] = mv
	}

	mapData, err := dynamodbattribute.MarshalMap(device.MapMeta)
	if err != nil {
		fmt.Println(err)
		return err
	}

	measureData, err := dynamodbattribute.MarshalMap(device.Measurements)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dal.Update("devices", "set active = :a, device_id = :d, map_meta = :mm, measurements = :m",
		map[string]*dal.AttributeValue{
			"token": {
				S: aws.String(deviceData.Token),
			},
			"cognito_id": {
				S: aws.String(device.CognitoID),
			},
		}, map[string]*dal.AttributeValue{
			":a": {
				BOOL: aws.Bool(device.Active),
			},
			":d": {
				S: aws.String(deviceData.DeviceID),
			},
			":mm": {
				M: mapData,
			},
			":m": {
				M: measureData,
			},
		})

	return nil
}

func main() {
	lambda.Start(Handler)
}
