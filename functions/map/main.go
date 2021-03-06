package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	"github.com/appsmonkey/core.server.functions/integration/cognito"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	s "github.com/appsmonkey/core.server.functions/models/schema"
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	"github.com/joho/godotenv"
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.MapRequest)
	response := request.Validate(req.QueryStringParameters)
	cognitoID := h.CognitoIDZeroValue
	authHdr := header("AccessToken", req.Headers)
	if len(authHdr) > 0 {
		c, _, isExpired, err := cog.ValidateToken(authHdr)
		if err != nil {
			fmt.Println(err)
			if isExpired {
				return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 401, Headers: response.Headers()}, nil
			}
		} else {
			cognitoID = c
		}
	}

	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	userGroups, err := cog.ListGroupsForUserFromID(cognitoID)

	isAdmin := false
	if err == nil {
		for _, g := range userGroups.Groups {
			if g.GroupName != nil && (*g.GroupName == "AdminGroup" || *g.GroupName == "SuperAdminGroup") {
				isAdmin = true
			}
		}
	}

	fmt.Println("Is Admin ::: ", isAdmin)

	// Get the polygon data
	type zoneResult struct {
		ZoneID   string       `json:"zone_id"`
		ZoneName string       `json:"zone_name"`
		Data     []m.ZoneMeta `json:"data"`
	}

	zoneMap := make(map[string]zoneResult, 0)
	zoneData := make([]zoneResult, 0)
	var zonesTable = "zones"
	if value, ok := os.LookupEnv("dynamodb_table_zones"); ok {
		zonesTable = value
	}
	for _, z := range request.Zone {
		fmt.Println("Fetching zone data for sensor: " + z)

		zoneRes, err := dal.GetFromIndex(zonesTable, "sensor_id-index", dal.Condition{
			"sensor_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dal.AttributeValue{
					{
						S: aws.String(z),
					},
				},
			},
		})

		zd := make([]m.Zone, 0)
		err = zoneRes.Unmarshal(&zd)

		if err != nil {
			fmt.Println(err)
			response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
		}

		for _, zz := range zd {
			_, ok := zoneMap[zz.ZoneID]
			if ok {
				t := zoneMap[zz.ZoneID]
				t.Data = append(t.Data, zz.Data)
				zoneMap[zz.ZoneID] = t
			} else {
				zoneMap[zz.ZoneID] = zoneResult{ZoneID: zz.ZoneID, Data: []m.ZoneMeta{zz.Data}, ZoneName: zz.ZoneName}
			}
		}
	}

	var devicesTable = "devices"
	if value, ok := os.LookupEnv("dynamodb_table_devices"); ok {
		devicesTable = value
	}

	dbRes, err := dal.ListNoFilter(devicesTable, dal.Projection(dal.Name("token"), dal.Name("device_id"), dal.Name("meta"), dal.Name("map_meta"), dal.Name("active"), dal.Name("city"), dal.Name("measurements"), dal.Name("cognito_id"), dal.Name("timestamp"), dal.Name("zone_id")))

	dbData := make([]m.Device, 0)
	err = dbRes.Unmarshal(&dbData)
	if err != nil {
		response.AddError(&es.Error{Code: 0, Message: err.Error(), Data: ""})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	data := make([]vm.DeviceGetData, 0)
	for _, tz := range zoneMap {
		var hasDevice = false
		data := make(map[string]float64, 0)
		datak := make(map[string]float64, 0)

		for _, d := range dbData {
			mine := d.CognitoID != h.CognitoIDZeroValue && cognitoID != h.CognitoIDZeroValue && d.CognitoID == cognitoID
			if (!d.Active) || d.Meta.Coordinates.IsEmpty() {
				continue
			} else if request.Filter == "mine" && !mine {
				continue
			} else if request.Filter == "indoor" && d.Meta.Indoor == false {
				continue
			} else if request.Filter == "outdoor" && d.Meta.Indoor {
				continue
			}

			if tz.ZoneID == d.ZoneID {
				hasDevice = true

				tz.Data = make([]m.ZoneMeta, 0)

				// take sensors in query in consideration
				for mmk, mmv := range d.MapMeta {
					for _, z := range request.Zone {
						if z == mmk {
							// Don't take in account if device has no sensor
							data[mmk] += mmv.Value
							datak[mmk]++
						}
					}
				}

				for rk, rv := range data {
					val := rv / datak[rk]

					schema := s.ExtractVersion("1")
					fieldData := schema[rk]
					val = fieldData.ConvertRawValue(val)

					ld, ln := s.SensorReading("1", rk, val)
					Data := m.ZoneMeta{
						SensorID:    rk,
						Name:        tz.ZoneID,
						Level:       ln,
						Value:       val,
						Measurement: ld.Name,
						Unit:        ld.Unit,
					}

					tz.Data = append(tz.Data, Data)
				}
			}
		}

		if !hasDevice {
			for index, zs := range tz.Data {
				zs.Value = -1
				zs.Level = "No device"
				tz.Data[index] = zs
			}
		} else {
			for _, z := range request.Zone {
				hasSensor := false
				for _, zsVal := range tz.Data {
					if zsVal.SensorID == z {
						hasSensor = true
					}
				}
				if !hasSensor {
					ld, _ := s.SensorReading("1", z, -1)
					Data := m.ZoneMeta{
						SensorID:    z,
						Name:        tz.ZoneID,
						Level:       "No device",
						Value:       -1,
						Measurement: ld.Name,
						Unit:        ld.Unit,
					}

					tz.Data = append(tz.Data, Data)
				}

			}
		}

		zoneData = append(zoneData, tz)
	}

	for _, d := range dbData {
		mine := d.CognitoID != h.CognitoIDZeroValue && cognitoID != h.CognitoIDZeroValue && d.CognitoID == cognitoID
		if (!mine && !d.Active && !isAdmin) || d.Meta.Coordinates.IsEmpty() {
			continue
		}

		md := make(map[string]m.MapMeta, 0)
		if len(request.Sensor) > 0 {
			for _, rs := range request.Sensor {
				md[rs] = d.MapMeta[rs]
			}
		}

		dData := vm.DeviceGetData{
			DeviceID:  d.Token,
			Name:      d.Meta.Name,
			Active:    d.Active,
			Mine:      mine,
			City:      d.City,
			Model:     d.Meta.Model,
			Indoor:    d.Meta.Indoor,
			Location:  d.Meta.Coordinates,
			MapMeta:   md,
			Timestamp: d.Timestamp,
			ZoneID:    d.ZoneID,
		}

		data = append(data, dData)
	}

	type MapResponseData struct {
		Zones   []zoneResult       `json:"zones"`
		Devices []vm.DeviceGetData `json:"devices"`
	}

	resData := MapResponseData{
		Zones:   zoneData,
		Devices: data,
	}
	response.Data = resData

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func init() {
	if os.Getenv("ENV") == "local" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("error loading .env: %v\n", err)
		}
	}

	cog = cognito.NewCognito()
}

func local() {
	data, _ := json.Marshal(vm.MapRequest{
		Sensor: []string{os.Getenv("SENSOR")},
	})

	resp, err := Handler(events.APIGatewayProxyRequest{
		Headers: map[string]string{"Authorization": `eyJraWQiOiI4NzVObDBXY0dwVVhqaUVNWmsxXC9rUEtoWUkxTlZFa0gxc1p4OW5jT05IWT0iLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI0OWEzMDAyYS1kYjljLTQyZWUtODRmNC03YzdlYjI0NDAyMWIiLCJldmVudF9pZCI6IjE0MjllNTU1LTU1NTItMTFlOS04MmE2LWE3ZWIzMWY2OGI1MCIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE1NTQyMTQ2NjIsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS13ZXN0LTEuYW1hem9uYXdzLmNvbVwvZXUtd2VzdC0xXzBndDZzNVRBUSIsImV4cCI6MTU1NDIxODI2MiwiaWF0IjoxNTU0MjE0NjYyLCJqdGkiOiJjMDBiY2MyOS03OTliLTQ4MDgtODg1MS1hODI1NmU4ZWE1MWEiLCJjbGllbnRfaWQiOiI2a25ubG4zMXRqZmJ0OWU3amFwOHZjNTU5MyIsInVzZXJuYW1lIjoiYWxiaW4uZGlkaWNAZ21haWwuY29tIn0.XzCMs5OczUXgmkJcSgOwTHqdcigpN3ZK7idFdQSlmmqvirHunCUI19-kEZDV8e2WfYeNW4IbkM5E3afcGCHchyHXvLYKeDY8KAVfcLEpJhXizixAld4XrFu4Jy0GjmFznAB5eitWKb7i7CvH1sHi933fMT3piHZptTg0ZF4M-q_KS1OOUrvaovzCGbfZaDZLjtXYCzop20h3KNtzlzJkg2avoa21wdyHlSHBFsUG66xbLYryy7a42PkQST-GX6BcnKHoYGWDwe3FN68M2tYB7ofVbUPboW9iee3pNbtzlTVrwxdU8-QiXxZULdBN6KRyx3yPOSnAgWFDFUIBrHwOwQ`},
		Body:    string(data),
	})

	if err != nil {
		fmt.Printf("unhandled error! \nError: %v\n", err)
	} else {
		j, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(j))
	}
}

func header(hdr string, in map[string]string) string {
	result, ok := in[hdr]
	if !ok {
		lwr := strings.ToLower(hdr)
		result = in[lwr]
	}

	return result
}

func main() {
	if os.Getenv("ENV") == "local" {
		local()
		return
	}

	lambda.Start(Handler)
}
