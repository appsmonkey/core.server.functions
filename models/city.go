package models

// City can contain many Zones, devices are grouped by cities,
type City struct {
	CityID    string   `json:"city_id"`
	Name      string   `json:"name"`
	Country   string   `json:"country"` // maybe we want to group them by countries later on ??
	Zones     []string `json:"zones"`
	Timestamp float64  `json:"timestamp"`
}
