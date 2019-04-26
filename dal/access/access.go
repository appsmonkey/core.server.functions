package access

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	"github.com/aws/aws-sdk-go/aws"
)

func init() {
}

// Increment the key inside the table
func Increment(in *IncrementInput) error {
	return dal.Update(in.Table, in.Expression(), in.Key(), in.Data())
}

// State returns the chart state
func State(state, key string) interface{} {
	gr, err := dal.Get("chart_state", map[string]*dal.AttributeValue{
		"name": {S: aws.String(state)},
	})
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	var res map[string]interface{}
	err = gr.Unmarshal(&res)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	data, ok := res[key]
	if !ok {
		return ""
	}

	return data
}

// SaveState will update the provided state
func SaveState(state, key string, value interface{}) error {
	data := map[string]interface{}{"name": state, key: value}

	return dal.Insert("chart_state", data)
}

// ChartHourInput will retrieve the hourly chart data from a specific point in time.
// `from` should be a timestamp in the past
func ChartHourInput(from interface{}) []ChartHourData {
	res, err := dal.List("chart_hour_input", dal.Name("time_stamp").GreaterThan(dal.Value(from)), dal.Projection(dal.Name("hash"), dal.Name("data_count"), dal.Name("data_value"), dal.Name("time_stamp")))
	if err != nil {
		fmt.Println(err.Error())
		return []ChartHourData{}
	}

	var data []ChartHourData
	err = res.Unmarshal(&data)
	if err != nil {
		fmt.Println(err.Error())
		return []ChartHourData{}
	}

	return data
}

// SaveHourChart will save the hourly chart data
func SaveHourChart(table string, data interface{}) error {
	return dal.Insert(table, data)
}
