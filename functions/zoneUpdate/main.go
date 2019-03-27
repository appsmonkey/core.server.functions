package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req interface{}) error {
	input, ok := req.(map[string]interface{})
	if !ok {
		err := errors.New("incorrect data received. input has incorrect format")
		fmt.Println(err)
		return err
	}

	zoneID := input["zone_id"].(string)

	res, err := dal.GetFromIndex("devices", "ZoneID-index", dal.Condition{
		"zone_id": {
			ComparisonOperator: aws.String("EQ"),
			AttributeValueList: []*dal.AttributeValue{
				{
					S: aws.String(zoneID),
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	dbData := make([]m.Device, 0)
	err = res.Unmarshal(&dbData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if len(dbData) == 0 {
		return nil
	}

	data := make(map[string]float64, 0)
	datak := make(map[string]float64, 0)

	for _, d := range dbData {
		for mmk, mmv := range d.MapMeta {
			f, _ := strconv.ParseFloat(mmv.Value, 64)
			data[mmk] += f
			datak[mmk]++
		}
	}

	for rk, rv := range data {
		val := rv / datak[rk]
		vals := fmt.Sprintf("%f", val)

		level := m.Level(rk, val)
		ti := m.Zone{
			ZoneID:   zoneID,
			SensorID: rk,
			Data: m.ZoneMeta{
				Name:        zoneID,
				Level:       level,
				Value:       vals,
				Measurement: m.MeasureMapName[rk],
				Unit:        m.MeasureMapUnit[rk],
			},
		}

		dal.Insert("zones", ti)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
