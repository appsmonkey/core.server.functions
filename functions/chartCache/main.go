// Chart aggregation every 5 minutes
package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/appsmonkey/core.server.functions/dal/access"
	"github.com/aws/aws-lambda-go/lambda"
)

var seconds string
var timeSteps map[int]int

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req interface{}) error {
	fmt.Println("INPUT: ", req)

	input, ok := req.(map[string]interface{})
	if !ok {
		err := errors.New("incorrect data received. input has incorrect format")
		fmt.Println(err)
		return err
	}

	token := input["token"].(string)
	timestamp := input["timestamp"].(float64)
	timestampStr := fmt.Sprintf("%f", timestamp)
	measurements := input["reported"].(map[string]interface{})

	for k, m := range measurements {
		sensor := k
		value, _ := strconv.ParseFloat(m.(string), 64)
		// value := strconv.ParseFloat(v.(string), 64)
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
	timeSteps = formulateTimeSteps(5)
	lambda.Start(Handler)
}

func incrementData(hash, timestamp, key1, value1, key2, value2 string) *access.IncrementInput {
	return &access.IncrementInput{
		Table:     "chart_hour_input",
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

func calculateHash(timestamp float64, token, sensor string) (devToken, generalToken string) {
	t := formulateTimestamp(int64(timestamp))

	devToken = fmt.Sprintf("hour:%v:%v:%v", t.Unix(), token, sensor)
	generalToken = fmt.Sprintf("hour:%v:%v", t.Unix(), sensor)

	return
}

func formulateTimestamp(in int64) time.Time {
	then := time.Unix(in, 0)
	minute := timeSteps[then.Minute()]
	return time.Date(then.Year(), then.Month(), then.Day(), then.Hour(), minute, 0, 0, time.UTC)
}

func formulateTimeSteps(step int) map[int]int {
	res := make(map[int]int, 0)

	s := 0
	d := -1
	for h := 0; h < 60; h++ {
		if s < step {
			if d < 0 {
				p := 0
				if h > 0 {
					p = res[h-1]
				}
				d = int(math.Round(float64(h+step+p) / 2))
			}
			res[h] = d
		}

		s++
		if s == step {
			s = 0
			d = -1
		}
	}

	return res
}
