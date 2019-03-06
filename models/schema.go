package models

// DefaultSensor for sensors that don't need calcujlations
const DefaultSensor = "7"

// DefaultValue of the default sensor
const DefaultValue = 0

// IsMeasuredSensor returns if the sensor needs to calcuilate its state or use defaults.
// `true = needs calculation`
func IsMeasuredSensor(sensorID string) bool {
	if sensorID == "7" || sensorID == "8" {
		return true
	}

	return false
}

// MeasureMapName sensor data mapping with names
var MeasureMapName = map[string]string{
	"1":  "Temperature",
	"2":  "Humidity",
	"3":  "Temperature Feel",
	"4":  "Pressure",
	"5":  "Altitude",
	"6":  "PM 1",
	"7":  "PM 2.5",
	"8":  "PM 10",
	"9":  "API Range",
	"10": "PM 2.5 Range",
	"11": "PM 10 Range",
	"12": "Light Lux",
	"13": "Eco 2",
	"14": "TVOC",
	"15": "Soil Temperature",
	"16": "Soil Moisture",
	"17": "Unix Time",
	"18": "Water Level",
	"19": "Motion",
}

// MeasureMapUnit sensor data mapping with Unit
var MeasureMapUnit = map[string]string{
	"1":  "℃",
	"2":  "?",
	"3":  "℃",
	"4":  "Pa",
	"5":  "m",
	"6":  "μg/m³",
	"7":  "μg/m³",
	"8":  "μg/m³",
	"9":  "?",
	"10": "μg/m³",
	"11": "μg/m³",
	"12": "?",
	"13": "?",
	"14": "?",
	"15": "℃",
	"16": "?",
	"17": "ms",
	"18": "m",
	"19": "?",
}

// Level of the sensor
func Level(sensor string, value float64) string {
	if IsMeasuredSensor(sensor) {
		if sensor == "7" {
			return PM25Level(value)
		}
		if sensor == "8" {
			return PM10Level(value)
		}
	}

	return "Great"
}

// PM25Level get the level
func PM25Level(in float64) string {
	if in >= 0 && in <= 12 {
		return "Great"
	}
	if in >= 13 && in <= 35 {
		return "OK"
	}
	if in >= 36 && in <= 55 {
		return "Sensitive beware"
	}
	if in >= 56 && in <= 150 {
		return "Unhealthy"
	}
	if in >= 152 && in <= 250 {
		return "Very Unhealthy"
	}
	if in > 250 {
		return "hazardous"
	}

	return "Unknown"
}

// PM10Level get the level
func PM10Level(in float64) string {
	if in >= 0 && in <= 54 {
		return "Great"
	}
	if in >= 55 && in <= 154 {
		return "OK"
	}
	if in >= 155 && in <= 254 {
		return "Sensitive beware"
	}
	if in >= 255 && in <= 354 {
		return "Unhealthy"
	}
	if in >= 355 && in <= 424 {
		return "Very Unhealthy"
	}
	if in > 425 {
		return "Hazardous"
	}

	return "Unknown"
}

// LevelIconMap icons
var LevelIconMap = map[string]string{
	"Great":            "images/icon-great.png",
	"OK":               "images/icon-ok.png",
	"Sensitive beware": "images/icon-sensitive.png",
	"Unhealthy":        "images/icon-unhealthy.png",
	"Very Unhealthy":   "images/icon-very-unhealthy.png",
	"Hazardous":        "images/icon-hazardous.png",
	"Unknown":          "images/icon-unknown.png",
}
