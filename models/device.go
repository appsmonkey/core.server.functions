package models

import (
	"time"
)

// Device is the base model representing the state of the DB of a single device
type Device struct {
	Token        string                 `json:"token"`
	DeviceID     string                 `json:"device_id"`
	CognitoID    string                 `json:"cognito_id,omitempty"`
	ZoneID       string                 `json:"zone_id"`
	Meta         Metadata               `json:"meta"`
	MapMeta      map[string]MapMeta     `json:"map_meta,omitempty"`
	Active       bool                   `json:"active"`
	Measurements map[string]interface{} `json:"measurements,omitempty"`
	Timestamp    float64                `json:"timestamp"`
	// TODO: add thing name attr.
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
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// IsEmpty indicates if the coordinates are set or not.
// Returns `false` if coordinates are set
func (l Location) IsEmpty() bool {
	if l.Lat == 0 && l.Lng == 0 {
		return true
	}

	return false
}

// MapMeta holds the calculated data used to dispay on the map
type MapMeta struct {
	Level       string  `json:"level"`
	Value       float64 `json:"value"`
	Measurement string  `json:"measurement"`
	Unit        string  `json:"unit"`
}

// ToLiveData will convert the data into live data needsd for the live table
func (d *Device) ToLiveData() map[string]interface{} {
	data := make(map[string]interface{}, 0)
	then := time.Now()
	data["token"] = d.Token
	data["timestamp"] = then.Unix()
	data["timestamp_sort"] = formulateTimestamp(then.Unix()).Unix()
	data["ttl"] = then.Add(time.Hour * 24 * 3).Unix()

	for k, v := range d.Measurements {
		data[k] = v
	}

	return data
}

func formulateTimestamp(in int64) time.Time {
	then := time.Unix(in, 0)
	return time.Date(then.Year(), then.Month(), then.Day(), then.Hour(), then.Minute(), then.Second(), 0, time.UTC)
}
