package models

import (
	"github.com/appsmonkey/core.server.functions/dal"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Device is the base model representing the state of the DB of a single device
type Device struct {
	Token        string                 `json:"token"`
	DeviceID     string                 `json:"device_id"`
	CognitoID    string                 `json:"cognito_id"`
	Meta         Metadata               `json:"meta"`
	MapMeta      map[string]MapMeta     `json:"map_meta,omitempty"`
	Active       bool                   `json:"active"`
	Measurements map[string]interface{} `json:"measurements,omitempty"`
}

// Metadata holds all the meda around the device
type Metadata struct {
	Name        string   `json:"name"`
	Model       string   `json:"model"`
	Coordinates Location `json:"coordinates"`
	Indoor      bool     `json:"indoor"`
}

// Location coordinates
type Location struct {
	Long string `json:"long"`
	Lat  string `json:"lat"`
}

// MapMeta holds the calculated data used to dispay on the map
type MapMeta struct {
	Name        string   `json:"name"`
	Level       string   `json:"level"`
	Coordinates Location `json:"coordinates"`
	Value       string   `json:"value"`
	Icon        string   `json:"icon"`
	Measurement string   `json:"measurement"`
	Unit        string   `json:"unit"`
}

// ToAttributeMap data
func (m MapMeta) ToAttributeMap() map[string]*dynamodb.AttributeValue {
	return map[string]*dal.AttributeValue{
		"name": {
			S: aws.String(m.Name),
		},
		"level": {
			S: aws.String(m.Level),
		},
		"coordinates": &dynamodb.AttributeValue{
			M: map[string]*dal.AttributeValue{
				"long": {
					S: aws.String(m.Coordinates.Long),
				},
				"lat": {
					S: aws.String(m.Coordinates.Lat),
				},
			},
		},
		"icon": {
			S: aws.String(m.Icon),
		},
		"measurement": {
			S: aws.String(m.Measurement),
		},
		"unit": {
			S: aws.String(m.Unit),
		},
		"value": {
			N: aws.String(m.Value),
		},
	}
}

// MapMetaAttributes data
func (d Device) MapMetaAttributes() map[string]*dynamodb.AttributeValue {
	result := make(map[string]*dynamodb.AttributeValue, 0)
	for k, v := range d.MapMeta {
		result[k] = &dynamodb.AttributeValue{
			M: v.ToAttributeMap(),
		}
	}

	return result
}
