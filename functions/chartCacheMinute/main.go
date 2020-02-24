// Chart aggregation per day
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"
	"github.com/appsmonkey/core.server.functions/tools/defaultDevice"

	m "github.com/appsmonkey/core.server.functions/models"
	"github.com/aws/aws-lambda-go/lambda"
)

var seconds string

type empty struct{}

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req interface{}) error {

	// fetch all cities with minimal data
	dbRes, err := dal.ListNoFilter("cities", dal.Projection(dal.Name("city_id"), dal.Name("name"), dal.Name("country"), dal.Name("timestamp")))
	if err != nil {
		fmt.Println("Error fetching cities")
	}

	cities := make([]m.City, 0)
	err = dbRes.Unmarshal(&cities)

	if err != nil {
		fmt.Println("Unmarshaling error ::: ", err)
	}

	// n := len(cities)
	// sem := make(chan empty, n) // Using semaphore for efficiency

	// for _, key := range cities {
	// 	go func(key m.City) {

	// 		dd := defaultDevice.Get(key.CityID)
	// 		then := time.Now()
	// 		data := make(map[string]interface{}, 0)

	// 		data["token"] = dd.City
	// 		data["timestamp"] = then.Unix()
	// 		data["timestamp_sort"] = formulateTimestamp(then.Unix()).Unix()
	// 		data["ttl"] = then.Add(time.Hour * 6).Unix()
	// 		data["indoor"] = false

	// 		for k, v := range dd.Latest {
	// 			data[k] = v
	// 		}

	// 		err = dal.Insert("chart_all_minute", data)
	// 		if err != nil {
	// 			fmt.Println("Couldn't insert data into table")
	// 		}
	// 		sem <- empty{}
	// 	}(key)
	// } - use this wehn list table stops using scan calls

	// dd := defaultDevice.Get("Sarajevo")
	// then := time.Now()
	// data := make(map[string]interface{}, 0)

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

	err = dal.Insert("chart_all_minute", data)
	if err != nil {
		fmt.Println("Couldn't insert data into table")
		return err
	}

	// wait for goroutines to finish
	// for i := 0; i < n; i++ {
	// 	<-sem
	// }

	return nil
}

func main() {
	seconds = fmt.Sprint(time.Now().Add(time.Second * 691200 * 1).Unix()) // One Month * 3 in seconds
	lambda.Start(Handler)
}

func formulateTimestamp(in int64) time.Time {
	then := time.Unix(in, 0)
	return time.Date(then.Year(), then.Month(), then.Day(), 0, 0, 0, 0, time.UTC)
}
