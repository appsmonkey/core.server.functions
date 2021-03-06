package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/appsmonkey/core.server.functions/dal/access"
	"github.com/aws/aws-lambda-go/lambda"
)

// Hour based data
type Hour struct {
	Date   int64   `json:"timestamp"`
	Value  float64 `json:"value"`
	Token  string  `json:"token"`
	Sensor string  `json:"sensor"`
	City   string  `json:"city"`
}

// Save the data into the DB
func (h *Hour) Save(last *int64) {
	table := "chart_device_six"
	data := make(map[string]interface{}, 0)
	if len(h.Token) == 0 {
		table = "chart_six"
		data["sensor"] = h.Sensor
		data["date"] = h.Date
		data["value"] = h.Value
		data["city"] = h.City
	} else {
		data["hash"] = fmt.Sprintf("%v<->%v", h.Token, h.Sensor)
		data["date"] = h.Date
		data["value"] = h.Value
	}

	// fmt.Println("DATA :::", h.Date, h.Sensor, h.Value, h.Token)
	// fmt.Println("DATA :::", data)

	err := access.SaveHourChart(table, &data)
	errString := ""
	if err != nil {
		errString = err.Error()
		fmt.Printf("Could not Save data [table: %v || err: %v || data: %v]\n", table, errString, h)
	} else if h.Date > *last {
		*last = h.Date
	}
}

type empty struct{}

// Handler will handle our request comming from the API gateway
func Handler() error {
	last := int64(0)
	from, ok := access.State("six_Last", "time_stamp").(float64)
	if !ok {
		fmt.Println("NOT FLOAT64")
		from = float64(0)
	}

	data := access.ChartInput(from, "chart_six_input")
	n := len(data)
	sem := make(chan empty, n) // Using semaphore for efficiency

	for _, key := range data {
		go func(key access.ChartHourData) {
			h := processKey(key)
			if h != nil && h.Date != 0 {
				h.Save(&last)
			}
			sem <- empty{}
		}(key)
	}

	// wait for goroutines to finish
	for i := 0; i < n; i++ {
		<-sem
	}

	// Now that everything is updated we will save the new state
	if last != 0 {
		access.SaveState("six_Last", "time_stamp", last)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}

func processKey(key access.ChartHourData) *Hour {
	res := formulateKey(key)

	if key.Count > 0 && key.Value > 0 {
		res.Value = key.Value / float64(key.Count)
	}

	res.City = key.City

	return res
}

func formulateKey(key access.ChartHourData) *Hour {
	d := strings.Split(key.Hash, "<->")
	res := new(Hour)
	n := len(d)
	switch n {
	case 3:
		t, err := strconv.ParseInt(d[1], 10, 64)
		if err != nil {
			return nil
		}
		res.Date = t
		res.Sensor = d[2]
	case 4:
		t, err := strconv.ParseInt(d[1], 10, 64)
		if err != nil {
			return nil
		}
		res.Date = t
		res.Token = d[2]
		res.Sensor = d[3]
	default:
		fmt.Println("FORMAT KEY ERROR! LEN :: ", n, "DATA ::", d)
	}

	return res
}
