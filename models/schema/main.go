package schema

import (
	"encoding/json"
	"fmt"
	"math"
)

// Default schema (used if no version was loaded)
var defaultData Schema

// All the loaded version of the schema
var versions map[string]Schema

// levels maped in order
var levels map[string]int

// Schema definition
type Schema map[string]*Data

// Marshal the schema
func (s Schema) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// AddSensor will add the sensor data into the schema
func (s Schema) AddSensor(sensor string, data *Data) {
	s[sensor] = data
}

// ExtractData reurn the sensor based on it's name
func (s Schema) ExtractData(name string) (*Data, string) {
	for k, v := range s {
		if k == name {
			return v, k
		}
	}

	return nil, ""
}

// AddVersion will add the schema as the provided verfsion
func AddVersion(version string, s Schema) {
	versions[version] = s
}

// ExtractVersion will return the selected version or the default version
func ExtractVersion(version string) Schema {
	d, ok := versions[version]
	if ok {
		return d
	}

	return defaultData
}

// MarshalSchema the schema
func MarshalSchema(version string) string {
	return ExtractVersion(version).Marshal()
}

// SensorReading will return the calculated data for the provided data on teh sensor
func SensorReading(version, sensor string, value float64) (*Data, string) {
	v := ExtractVersion(version)
	sc, ok := v[sensor]
	if ok {
		return sc, sc.Result(value)
	}

	return nil, ""
}

// Data definition of a calculation
type Data struct {
	Name           string      `json:"name"`
	Unit           string      `json:"unit"`
	CalcSteps      []*CalcStep `json:"steps"`
	DefaultValue   string      `json:"default"`
	ParseCondition string      `json:"parse_condition"`
}

// ConvertRawValue will return parsed raw value form device
func (s *Data) ConvertRawValue(v float64) float64 {
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

// Result will return the calculated result
func (s *Data) Result(v float64) string {
	if len(s.CalcSteps) == 0 {
		return s.DefaultValue
	}

	for _, cs := range s.CalcSteps {
		if cs.IsMe(v) {
			return cs.Result
		}
	}

	return s.DefaultValue
}

// CalcStep to check the data
type CalcStep struct {
	From   float64 `json:"from"`
	To     float64 `json:"to"`
	Result string  `json:"result"`
}

// IsMe returns `true` if the provided value is between its defined bounds
func (s *CalcStep) IsMe(v float64) bool {
	if v >= s.From && v <= s.To {
		return true
	}

	return false
}

// LevelOrder will return the weight of the level
func LevelOrder(level string) int {
	fmt.Println("LEVEL ORDER", levels)
	lvl, ok := levels[level]
	if !ok {
		return -1
	}

	return lvl
}

func init() {
	versions = make(map[string]Schema, 0)
	levels = make(map[string]int, 0)
	levels["No devices"] = -1
	levels["Great"] = 0
	levels["OK"] = 1
	levels["Sensitive beware"] = 2
	levels["Unhealthy"] = 3
	levels["Very Unhealthy"] = 4
	levels["Hazardous"] = 5
}
