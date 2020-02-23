package access

import (
	"errors"
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
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
	// res, err := dal.List("chart_hour_input", dal.Name("time_stamp").GreaterThan(dal.Value(from)), dal.Projection(dal.Name("hash"), dal.Name("data_count"), dal.Name("data_value"), dal.Name("city"), dal.Name("time_stamp")), true)
	res, err := dal.GetFromIndex("chart_hour_input", "take-timestamp-index",
		dal.Condition{
			"take": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dal.AttributeValue{
					{
						S: aws.String("true"),
					},
				},
			},
			"times_tamp": {
				ComparisonOperator: aws.String("GT"),
				AttributeValueList: []*dal.AttributeValue{
					{
						N: aws.String(fmt.Sprintf("%v", from)),
					},
				},
			},
		})

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

// ChartInput will retrieve the specified chart data from a specific point in time.
// `table` should be a table name from which to get data from in the chart based schema
// `from` should be a timestamp in the past
func ChartInput(from interface{}, table string) []ChartHourData {
	// res, err := dal.List(table, dal.Name("time_stamp").GreaterThan(dal.Value(from)), dal.Projection(dal.Name("hash"), dal.Name("data_count"), dal.Name("data_value"), dal.Name("city"), dal.Name("time_stamp")), true)
	res, err := dal.GetFromIndex(table, "take-timestamp-index",
		dal.Condition{
			"take": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dal.AttributeValue{
					{
						S: aws.String("true"),
					},
				},
			},
			"times_tamp": {
				ComparisonOperator: aws.String("GT"),
				AttributeValueList: []*dal.AttributeValue{
					{
						N: aws.String(fmt.Sprintf("%v", from)),
					},
				},
			},
		})

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

// CheckSocial will return data from the DB for the selected social data
func CheckSocial(id string) (sub, email, sid, st string, success bool, err error) {
	gr, err := dal.GetFromIndex("users", "UsersSocialID", dal.Condition{
		"social_id": {
			ComparisonOperator: aws.String("EQ"),
			AttributeValueList: []*dal.AttributeValue{
				{
					S: aws.String(id),
				},
			},
		},
	})
	if err != nil {
		fmt.Println("CheckSocial Error", err.Error())
		return "", "", "", "", false, err
	}

	res := make([]m.User, 0)
	err = gr.Unmarshal(&res)
	if err != nil {
		fmt.Println("CheckSocial Error [could not unmarshal]", err.Error())
		return "", "", "", "", false, err
	}

	if len(res) == 0 {
		fmt.Println("CheckSocial Error [No data received]")
		fmt.Println("ID passed: " + id)
		return "", "", "", "", false, errors.New("no data received")
	}

	usr := res[0]
	sub = usr.CognitoID
	email = usr.Email
	sid = usr.SocialID
	st = usr.SocialType
	err = nil

	return
}

// AddTempUser will save a temp user to be used by cognito later on
func AddTempUser(email, socialID, socialType string) error {
	data := make(map[string]string, 0)
	data["cognito_id"] = "TEMP" // We will store all temp users under one key
	data["email"] = email
	data["social_id"] = socialID
	data["social_type"] = socialType

	return dal.Insert("users", data)
}

// AddUser will save a user to the DB
func AddUser(data interface{}) error {
	return dal.Insert("users", data)
}

// GetTempUser will retrieve the temp user data, if any
func GetTempUser(email string) (socialID, socialType string, success bool, err error) {
	gr, err := dal.Get("users", map[string]*dal.AttributeValue{
		"cognito_id": {S: aws.String("TEMP")},
		"email":      {S: aws.String(email)},
	})
	if err != nil {
		fmt.Println(err.Error())
		return "", "", false, err
	}

	var res map[string]string
	err = gr.Unmarshal(&res)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", false, err
	}

	socialID, ok := res["social_id"]
	if !ok {
		return "", "", true, errors.New("missing social_id")
	}

	socialType, ok = res["social_type"]
	if !ok {
		return "", "", true, errors.New("missing social_type")
	}

	return
}
