// Chart aggregation per day
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"

	"github.com/appsmonkey/core.server.functions/tools/defaultDevice"
	"github.com/aws/aws-lambda-go/lambda"
)

var seconds string

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req interface{}) error {
	dd := defaultDevice.Get("Sarajevo")
	then := time.Now()
	data := make(map[string]interface{}, 0)

	data["token"] = dd.City
	data["timestamp"] = then.Unix()
	data["timestamp_sort"] = formulateTimestamp(then.Unix()).Unix()
	data["ttl"] = then.Add(time.Hour * 6).Unix()
	data["indoor"] = false

	for k, v := range dd.Latest {
		data[k] = v
	}

	err := dal.Insert("chart_all_minute", data)
	if err != nil {
		fmt.Println("Couldn't insert data into table")
		return err
	}
	return nil
}

func main() {
	seconds = fmt.Sprint(time.Now().Add(time.Second * 2592000 * 3).Unix()) // One Month * 3 in seconds
	lambda.Start(Handler)
}

func formulateTimestamp(in int64) time.Time {
	then := time.Unix(in, 0)
	return time.Date(then.Year(), then.Month(), then.Day(), 0, 0, 0, 0, time.UTC)
}
