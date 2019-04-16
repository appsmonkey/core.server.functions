package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/appsmonkey/core.server.functions/integration/redis"
	"github.com/aws/aws-lambda-go/lambda"
)

var seconds string

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req interface{}) error {
	input, ok := req.(map[string]interface{})
	if !ok {
		err := errors.New("incorrect data received. input has incorrect format")
		fmt.Println(err)
		return err
	}

	token := input["token"].(string)
	timestamp := input["timestamp"].(float64)
	measurements := input["measurements"].([]interface{})

	for _, m := range measurements {
		sensor, value := sensorData(m)
		dev, gen := calculateHash(timestamp, token, sensor)

		// Set counter and value for the device specific value
		strValue := fmt.Sprintf("%f", value)
		redis.IncrementByHash(dev, "counter", "1")
		redis.IncrementByHashFloat(dev, "value", strValue)
		redis.Expire(dev, seconds)

		// Set counter and value for the overal calculation
		redis.IncrementByHash(gen, "counter", "1")
		redis.IncrementByHashFloat(gen, "value", strValue)
		redis.Expire(gen, seconds)
	}

	return nil
}

func main() {
	seconds = fmt.Sprint(2592000 * 3) // One Month * 3 in seconds
	lambda.Start(Handler)
}

func sensorData(v interface{}) (sensor string, value float64) {
	data := v.(map[string]interface{})
	for i, j := range data {
		sensor = i
		value = j.(float64)
		break
	}

	return
}

func calculateHash(timestamp float64, token, sensor string) (devToken, generalToken string) {
	t := formulateTimestamp(int64(timestamp))

	devToken = fmt.Sprintf("hour:%v:%v:%v", t.Unix(), token, sensor)
	generalToken = fmt.Sprintf("hour:%v:%v", t.Unix(), sensor)

	return
}

func formulateTimestamp(in int64) time.Time {
	then := time.Unix(in, 0)
	return time.Date(then.Year(), then.Month(), then.Day(), then.Hour(), 0, 0, 0, time.UTC)
}
