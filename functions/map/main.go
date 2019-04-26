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
	vm "github.com/appsmonkey/core.server.functions/viewModels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	h "github.com/appsmonkey/core.server.functions/tools/helper"
	"github.com/joho/godotenv"
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cognitoID := h.CognitoIDZeroValue
	authHdr := header("AccessToken", req.Headers)
	if len(authHdr) > 0 {
		c, _, err := cog.ValidateToken(authHdr)
		if err != nil {
			fmt.Println(err)
		} else {
			cognitoID = c
		}
	}

	request := new(vm.MapRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	// Get the polygon data
	type zoneResult struct {
		ZoneID string       `json:"zone_id"`
		Data   []m.ZoneMeta `json:"data"`
	}
	zoneMap := make(map[string]zoneResult, 0)
	zoneData := make([]zoneResult, 0)
	for _, z := range request.Zone {
		zoneRes, err := dal.List("zones", dal.Name("sensor_id").Equal(dal.Value(z)), dal.Projection(dal.Name("zone_id"), dal.Name("data")))
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
				zoneMap[zz.ZoneID] = zoneResult{ZoneID: zz.ZoneID, Data: []m.ZoneMeta{zz.Data}}
			}
		}
	}

	for _, tz := range zoneMap {
		zoneData = append(zoneData, tz)
	}

	dbRes, err := dal.ListNoFilter("devices", dal.Projection(dal.Name("token"), dal.Name("device_id"), dal.Name("meta"), dal.Name("map_meta"), dal.Name("active"), dal.Name("measurements"), dal.Name("cognito_id"), dal.Name("timestamp")))
	dbData := make([]m.Device, 0)
	err = dbRes.Unmarshal(&dbData)
	if err != nil {
		response.AddError(&es.Error{Code: 0, Message: err.Error(), Data: ""})
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	data := make([]vm.DeviceGetData, 0)

	for _, d := range dbData {
		mine := d.CognitoID != h.CognitoIDZeroValue && cognitoID != h.CognitoIDZeroValue && d.CognitoID == cognitoID
		if (!mine && !d.Active) || d.Meta.Coordinates.IsEmpty() {
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
			Model:     d.Meta.Model,
			Indoor:    d.Meta.Indoor,
			Location:  d.Meta.Coordinates,
			MapMeta:   md,
			Timestamp: d.Timestamp,
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
