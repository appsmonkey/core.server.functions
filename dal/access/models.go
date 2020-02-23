package access

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	"github.com/aws/aws-sdk-go/aws"
)

// IncrementInput for incrementing a value in a tables column
type IncrementInput struct {
	Table        string
	KeyName      string
	KeyValue     string
	SortKeyName  string
	SortKeyValue string
	TTL          string
	City         string
	Timestamp    string
	Columns      []IncrementItem
}

// IncrementItem a single column to increment
type IncrementItem struct {
	Column string
	Value  string
	Type   string
}

// Key for dynamoDb statement
func (i *IncrementInput) Key() map[string]*dal.AttributeValue {
	keyMap := make(map[string]*dal.AttributeValue, 0)
	keyMap[i.KeyName] = &dal.AttributeValue{S: aws.String(i.KeyValue)}

	if len(i.SortKeyName) > 0 {
		keyMap[i.SortKeyName] = &dal.AttributeValue{N: aws.String(i.SortKeyValue)}
	}

	return keyMap
}

// Data for the increment statement
func (i *IncrementInput) Data() map[string]*dal.AttributeValue {
	keyMap := make(map[string]*dal.AttributeValue, 0)

	for _, c := range i.Columns {
		if c.Type == "number" {
			keyMap[":"+c.Column] = &dal.AttributeValue{N: aws.String(c.Value)}
		} else {
			keyMap[":"+c.Column] = &dal.AttributeValue{S: aws.String(c.Value)}
		}
	}

	if len(i.TTL) > 0 {
		keyMap[":ttl"] = &dal.AttributeValue{N: aws.String(i.TTL)}
	}
	if len(i.Timestamp) > 0 {
		keyMap[":ts"] = &dal.AttributeValue{N: aws.String(i.Timestamp)}
	}
	if len(i.City) > 0 {
		keyMap[":ci"] = &dal.AttributeValue{S: aws.String(i.City)}
	}

	keyMap[":ks"] = &dal.AttributeValue{S: aws.String("true")}

	return keyMap
}

// Expression used to increment the data
func (i *IncrementInput) Expression() string {
	expr := ""

	if len(i.TTL) > 0 && len(i.Timestamp) > 0 {
		expr = "SET time_to_live = :ttl, time_stamp = :ts"
	} else if len(i.TTL) > 0 {
		expr = "SET time_to_live = :ttl"
	} else if len(i.Timestamp) > 0 {
		expr = "SET time_stamp = :ts"
	}

	if len(i.City) > 0 {
		expr += ",city = :ci"
	}

	expr += ",take = :ks"

	expr += " ADD"
	length := len(i.Columns)
	for ind, c := range i.Columns {
		comma := ", "
		if length-1 == ind%length {
			comma = ""
		}

		expr = fmt.Sprintf("%v %v :%v%v", expr, c.Column, c.Column, comma)
	}

	return expr
}

// ChartHourData representing onw row in the input table
type ChartHourData struct {
	Hash      string  `json:"hash"`
	Count     float64 `json:"data_count"`
	Value     float64 `json:"data_value"`
	Timestamp int64   `json:"time_stamp"`
	City      string  `json:"city"`
}

// HourChart represents one calculated row in the hourly chart
type HourChart struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
	Token     string  `json:"token"`
	Sensor    string  `json:"sensor"`
}
