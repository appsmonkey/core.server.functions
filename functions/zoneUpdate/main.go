package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	s "github.com/appsmonkey/core.server.functions/models/schema"
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
	city := input["city_id"].(string)

	if city == "none" || len(city) == 0 {
		city = "Sarajevo"
	}

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
		sensors := s.ExtractVersion("1")
		for sk := range sensors {
			ld, ln := s.SensorReading("1", sk, 0)
			// zd := m.Zone{
			// 	ZoneID:   zoneID,
			// 	SensorID: sk,
			// 	CityID:   city,
			// 	Data: m.ZoneMeta{
			// 		SensorID:    sk,
			// 		Name:        zoneID,
			// 		Level:       ln,
			// 		Value:       0,
			// 		Measurement: ld.Name,
			// 		Unit:        ld.Unit,
			// 	},
			// }

			// dal.Insert("zones", zd)

			zoneD := m.ZoneMeta{
				SensorID:    sk,
				Name:        zoneID,
				Level:       ln,
				Value:       0,
				Measurement: ld.Name,
				Unit:        ld.Unit,
			}

			mZoneD, err := dynamodbattribute.MarshalMap(zoneD)

			if err != nil {
				fmt.Println("ZONE DATA MARSHAL ERROR")
			}

			err = dal.Update("zones", "set data = :d",
				map[string]*dal.AttributeValue{
					"zone_id": {
						S: aws.String(zoneID),
					},
					"sensor_id": {
						S: aws.String(sk),
					},
				}, map[string]*dal.AttributeValue{
					":d": {
						M: mZoneD,
					},
				})
		}

		return nil
	}

	data := make(map[string]float64, 0)
	datak := make(map[string]float64, 0)

	for _, d := range dbData {
		for mmk, mmv := range d.MapMeta {
			data[mmk] += mmv.Value
			datak[mmk]++
		}
	}

	for rk, rv := range data {
		val := rv / datak[rk]

		schema := s.ExtractVersion("1")
		fieldData := schema[rk]
		val = fieldData.ConvertRawValue(val)

		ld, ln := s.SensorReading("1", rk, val)
		// zd := m.Zone{
		// 	ZoneID:   zoneID,
		// 	SensorID: rk,
		// 	CityID:   city,
		// 	Data: m.ZoneMeta{
		// 		SensorID:    rk,
		// 		Name:        zoneID,
		// 		Level:       ln,
		// 		Value:       val,
		// 		Measurement: ld.Name,
		// 		Unit:        ld.Unit,
		// 	},
		// }

		// dal.Insert("zones", zd)

		zoneD := m.ZoneMeta{
			SensorID:    rk,
			Name:        zoneID,
			Level:       ln,
			Value:       val,
			Measurement: ld.Name,
			Unit:        ld.Unit,
		}

		mZoneD, err := dynamodbattribute.MarshalMap(zoneD)

		if err != nil {
			fmt.Println("ZONE DATA MARSHAL ERROR")
		}

		err = dal.Update("zones", "set data = :d",
			map[string]*dal.AttributeValue{
				"zone_id": {
					S: aws.String(zoneID),
				},
				"sensor_id": {
					S: aws.String(rk),
				},
			}, map[string]*dal.AttributeValue{
				":d": {
					M: mZoneD,
				},
			})
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
