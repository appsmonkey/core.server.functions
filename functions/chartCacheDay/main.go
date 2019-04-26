package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/appsmonkey/core.server.functions/dal/access"
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
	timestampStr := fmt.Sprintf("%f", timestamp)
	measurements := input["measurements"].([]interface{})

	for _, m := range measurements {
		sensor, value := sensorData(m)
		dev, gen := calculateHash(timestamp, token, sensor)

		// Set counter and value for the device specific value
		strValue := fmt.Sprintf("%f", value)
		access.Increment(incrementData(dev, timestampStr, "data_count", "1", "data_value", strValue))
		access.Increment(incrementData(gen, timestampStr, "data_count", "1", "data_value", strValue))
	}

	return nil
}

func main() {
	seconds = fmt.Sprint(time.Now().Add(time.Second * 2592000 * 3).Unix()) // One Month * 3 in seconds
	lambda.Start(Handler)
}

func incrementData(hash, timestamp, key1, value1, key2, value2 string) *access.IncrementInput {
	return &access.IncrementInput{
		Table:     "chart_day_input",
		KeyName:   "hash",
		KeyValue:  hash,
		TTL:       seconds,
		Timestamp: timestamp,
		Columns: []access.IncrementItem{
			{
				Column: key1,
				Value:  value1,
				Type:   "number",
			},
			{
				Column: key2,
				Value:  value2,
				Type:   "number",
			},
		},
	}
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

	devToken = fmt.Sprintf("day:%v:%v:%v", t.Unix(), token, sensor)
	generalToken = fmt.Sprintf("day:%v:%v", t.Unix(), sensor)

	return
}

func formulateTimestamp(in int64) time.Time {
	then := time.Unix(in, 0)
	return time.Date(then.Year(), then.Month(), then.Day(), 0, 0, 0, 0, time.UTC)
}
