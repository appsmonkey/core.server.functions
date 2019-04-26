package models

// Zone is the base model representing the state of the DB of a single zone
type Zone struct {
	ZoneID   string   `json:"zone_id"`
	SensorID string   `json:"sensor_id"`
	Data     ZoneMeta `json:"data,omitempty"`
}

// ZoneMeta holds the calculated data used to dispay on the map
type ZoneMeta struct {
	SensorID    string  `json:"sensor_id"`
	Name        string  `json:"name"`
	Level       string  `json:"level"`
	Value       float64 `json:"value"`
	Measurement string  `json:"measurement"`
	Unit        string  `json:"unit"`
}
