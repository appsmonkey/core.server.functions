package models

import (
	"encoding/json"
	"math"
)

// DefaultSensor for sensors that don't need calcujlations
const DefaultSensor = "AIR_PM2P5"

// DefaultValue of the default sensor
const DefaultValue = 0

// IsMeasuredSensor returns if the sensor needs to calcuilate its state or use defaults.
// `true = needs calculation`
func IsMeasuredSensor(sensorID string) bool {
	if sensorID == "AIR_PM10" || sensorID == "AIR_PM2P5" {
		return true
	}

	return false
}

// MeasureMapName sensor data mapping with names
var MeasureMapName = map[string]string{
	"AIR_TEMPERATURE":      "Temperature",
	"AIR_HUMIDITY":         "Air Humidity",
	"AIR_TEMPERATURE_FEEL": "Temperature Feel",
	"AIR_PRESSURE":         "Pressure",
	"AIR_ALTITUDE":         "Altitude",
	"AIR_PM1":              "PM 1",
	"AIR_PM2P5":            "PM 2.5",
	"AIR_PM10":             "PM 10",
	"AIR_AQI_RANGE":        "AQI Range",
	"AIR_PM2P5_RANGE":      "PM 2.5 Range",
	"AIR_PM10_RANGE":       "PM 10 Range",
	"LIGHT_INTENSITY":      "Light Lux",
	"AIR_ECO2":             "Eco 2",
	"AIR_TVOC":             "TVOC",
	"AIR_CO2":              "CO2",
	"SOIL_TEMPERATURE":     "Soil Temperature",
	"SOIL_MOISTURE":        "Soil Moisture",
	"TIME_UNIXTIME":        "Unix Time",
	"WATER_LEVEL_SWITCH":   "Water Level",
	"MOTION":               "Motion",
	"DISTANCE":             "Distance",
	"AIR_VOC":              "VOC",
}

// MeasureMapUnit sensor data mapping with Unit
var MeasureMapUnit = map[string]string{
	"AIR_TEMPERATURE":      "℃",
	"AIR_HUMIDITY":         "%",
	"AIR_TEMPERATURE_FEEL": "℃",
	"AIR_PRESSURE":         "Pa",
	"AIR_ALTITUDE":         "m",
	"AIR_PM1":              "μg/m³",
	"AIR_PM2P5":            "μg/m³",
	"AIR_PM10":             "μg/m³",
	"AIR_AQI_RANGE":        "",
	"AIR_PM2P5_RANGE":      "",
	"AIR_PM10_RANGE":       "",
	"LIGHT_INTENSITY":      "℃",
	"AIR_ECO2":             "℃",
	"AIR_CO2":              "ppm",
	"AIR_TVOC":             "℃",
	"SOIL_TEMPERATURE":     "℃",
	"SOIL_MOISTURE":        "%",
	"TIME_UNIXTIME":        "ms",
	"WATER_LEVEL_SWITCH":   "m",
	"MOTION":               "",
	"DISTANCE":             "m",
	"AIR_VOC":              "℃",
}

// Level of the sensor
func Level(sensor string, value float64) string {
	if IsMeasuredSensor(sensor) {
		if sensor == "AIR_PM2P5" {
			return PM25Level(value)
		}
		if sensor == "AIR_PM10" {
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

// Schema definition
type Schema map[string]*SchemaData

// Marshal the schema
func (s Schema) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// MarshalSchema the schema
func MarshalSchema() string {
	return sch.Marshal()
}

// SensorReading will return the calculated data for the provided data on teh sensor
func SensorReading(sensor string, value float64) (*SchemaData, string) {
	sc, ok := sch[sensor]
	if ok {
		return sc, sc.Result(value)
	}

	return nil, "NOPE"
}

// SchemaData definition of a calculation
type SchemaData struct {
	Name           string            `json:"name"`
	Unit           string            `json:"unit"`
	CalcSteps      []*SchemaCalcStep `json:"steps"`
	DefaultValue   string            `json:"default"`
	ParseCondition string            `json:"parse_condition"`
}

// Result will return the calculated result
func (s *SchemaData) Result(v float64) string {
	if s.CalcSteps == nil || len(s.CalcSteps) == 0 {
		return s.DefaultValue
	}

	for _, cs := range s.CalcSteps {
		if cs.IsMe(v) {
			return cs.Result
		}
	}

	return s.DefaultValue
}

// ConvertRawValue will return parsed raw value form device
func (s *SchemaData) ConvertRawValue(v float64) interface{} {
	if len(s.ParseCondition) > 0 {
		switch s.ParseCondition {
		case "round":
			{
				return math.Round(v)
			}
		// add support for more parsing logic if needed
		default:
			{
				return v
			}
		}
	}

	return v
}

// SchemaCalcStep to check the data
type SchemaCalcStep struct {
	From   float64 `json:"from"`
	To     float64 `json:"to"`
	Result string  `json:"result"`
}

// IsMe returns `true` if the provided value is between its defined bounds
func (s *SchemaCalcStep) IsMe(v float64) bool {
	if v >= s.From && v <= s.To {
		return true
	}

	return false
}

var sch Schema

func init() {
	sch = make(map[string]*SchemaData, 0)
	sch["1"] = &SchemaData{
		Name:         "Temperature",
		Unit:         "℃",
		DefaultValue: "OK_DEFAULT",
		CalcSteps: []*SchemaCalcStep{
			&SchemaCalcStep{
				From:   -12,
				To:     0,
				Result: "Great",
			},
			&SchemaCalcStep{
				From:   0,
				To:     24,
				Result: "OK",
			},
		},
	}
}
