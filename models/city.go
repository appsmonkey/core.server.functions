package models

// City can contain many Zones, devices are grouped by cities,
type City struct {
	CityID    string  `json:"city_id"`
	Country   string  `json:"country"` // maybe we want to group them by countries later on ??
	Timestamp float64 `json:"timestamp"`
}
