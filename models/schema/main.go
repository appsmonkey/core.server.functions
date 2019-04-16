package schema

import (
	"encoding/json"
)

// Default schema (used if no version was loaded)
var defaultData Schema

// All the loaded version of the schema
var versions map[string]Schema

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
	Name         string      `json:"name"`
	Unit         string      `json:"unit"`
	CalcSteps    []*CalcStep `json:"steps"`
	DefaultValue string      `json:"default"`
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

func init() {
	versions = make(map[string]Schema, 0)
}
