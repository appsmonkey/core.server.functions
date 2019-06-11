package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"context"

	"github.com/appsmonkey/core.server.functions/integration/cognito"
)

// type sensor struct {
// 	name    string
// 	display string
// 	value   float64
// 	level   string
// }

// // Value representing a composite of the name and value
// func (s *sensor) Value() string {
// 	return fmt.Sprintf("%v = %v", s.display, s.value)
// }

// // Channel in UA
// func (s *sensor) Channel() ua.ChanelType {
// 	switch s.level {
// 	case "Sensitive beware":
// 		return ua.Sensitive
// 	case "Unhealthy":
// 		return ua.Unhealthy
// 	case "Very Unhealthy":
// 		return ua.VeryUnhealthy
// 	case "Hazardous":
// 		return ua.Hazardous
// 	default:
// 		return ua.Good
// 	}
// }

// // Handler will handle our request comming from the API gateway
// func Handler(ctx context.Context, req interface{}) (int64, error) {
// 	in := req.(map[string]interface{})
// 	lvl := in["level"].(float64)

// 	schemaDefault := s.ExtractVersion("")
// 	pm10, pm10Sensor := schemaDefault.ExtractData("PM 10")
// 	sens10 := sensor{name: pm10Sensor, display: "PM 10", value: lvl, level: pm10.Result(lvl)}

// 	ua.New().Send(sens10.Value(), sens10.Channel())

// 	return 0, nil
// }

// Handler will handle our request comming from the API gateway
// func Handler(ctx context.Context, req interface{}) (int64, error) {
// 	in := req.(map[string]interface{})
// 	return dal.Count(in["table"].(string)), nil
// }

// // Handler will handle our request comming from the API gateway
// func Handler(ctx context.Context, req interface{}) (string, error) {
// 	in := req.(map[string]interface{})
// 	action := in["key"].(string)
// 	params := in["params"].([]interface{})

// 	switch action {
// 	case "FLUSHDB":
// 		r := redis.FlushDB()
// 		return "", r
// 	case "KEYS":
// 		if len(params) > 0 {
// 			k := params[0].(string)
// 			r, err := redis.Keys(k)
// 			rr, err := json.Marshal(r)
// 			return string(rr), err
// 		} else {
// 			return "Missing Parameters", nil
// 		}
// 	case "DEL":
// 		if len(params) > 0 {
// 			k := params[0].(string)
// 			err := redis.Del(k)
// 			return "Executed", err
// 		} else {
// 			return "Missing Parameters", nil
// 		}
// 	case "HGET.I":
// 		if len(params) > 1 {
// 			h := params[0].(string)
// 			k := params[1].(string)
// 			i, err := redis.GetIntHash(h, k)
// 			return fmt.Sprintf("%v", i), err
// 		} else {
// 			return "Missing Parameters", nil
// 		}
// 	case "HGET.F":
// 		if len(params) > 1 {
// 			h := params[0].(string)
// 			k := params[1].(string)
// 			i, err := redis.GetFloatHash(h, k)
// 			return fmt.Sprintf("%v", i), err
// 		} else {
// 			return "Missing Parameters", nil
// 		}
// 	}

// 	return "No Command Found", nil
// }

var (
	cog *cognito.Cognito
)

func Handler(ctx context.Context, req interface{}) error {
	in := req.(map[string]interface{})
	email := in["email"].(string)
	token := in["token"].(string)

	_, err := cog.SignInTest(email, token)
	return err
}

func main() {
	lambda.Start(Handler)
}

func init() {
	cog = cognito.NewCognito()
}

// func sensorData(v interface{}) (sensor string, value float64) {
// 	data := v.(map[string]interface{})
// 	for i, j := range data {
// 		sensor = i
// 		value = j.(float64)
// 		break
// 	}

// 	return
// }

// func calculateHash(timestamp float64, token, sensor string) (devToken, generalToken string) {
// 	t := formulateTimestamp(int64(timestamp))

// 	devToken = fmt.Sprintf("hour:%v:%v:%v", t.Unix(), token, sensor)
// 	generalToken = fmt.Sprintf("hour:%v:%v", t.Unix(), sensor)

// 	return
// }

// func formulateTimestamp(in int64) time.Time {
// 	then := time.Unix(in, 0)
// 	return time.Date(then.Year(), then.Month(), then.Day(), then.Hour(), 0, 0, 0, time.UTC)
// }

// func Handler2(ctx context.Context, s3Event events.S3Event) {
// 	for _, record := range s3Event.Records {
// 		s3 := record.S3
// 		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
// 	}
// }
