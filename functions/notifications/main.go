package main

import (
	"fmt"
	"math"

	ua "github.com/appsmonkey/core.server.functions/integration/airship"
	s "github.com/appsmonkey/core.server.functions/models/schema"
	"github.com/appsmonkey/core.server.functions/tools/defaultDevice"
	"github.com/aws/aws-lambda-go/lambda"
)

type sensor struct {
	name    string
	display string
	value   float64
	level   string
}

// Value representing a composite of the name and value
func (s *sensor) Value() string {
	return fmt.Sprintf("%v = %v", s.display, math.Round(s.value))
}

// Channel in UA
func (s *sensor) Channel() ua.ChanelType {
	switch s.level {
	case "Sensitive beware":
		return ua.Sensitive
	case "Unhealthy":
		return ua.Unhealthy
	case "Very Unhealthy":
		return ua.VeryUnhealthy
	case "Hazardous":
		return ua.Hazardous
	default:
		return ua.Good
	}
}

// Handler will handle our request comming from the API gateway
func Handler() error {
	data := defaultDevice.Get()
	schemaDefault := s.ExtractVersion("")
	// pm25, pm25Sensor := schemaDefault.ExtractData("PM 2.5")
	// pm10, pm10Sensor := schemaDefault.ExtractData("PM 10")

	pm25, pm25Sensor := schemaDefault.ExtractData("AIR_PM2P5")
	pm10, pm10Sensor := schemaDefault.ExtractData("AIR_PM10")

	fmt.Println(pm10, pm25)

	sens25 := sensor{name: pm25Sensor, display: "PM 2.5", value: data.Latest[pm25Sensor].(float64)}
	sens10 := sensor{name: pm10Sensor, display: "PM 10", value: data.Latest[pm10Sensor].(float64)}

	sens25.level = pm25.Result(sens25.value)
	sens10.level = pm10.Result(sens10.value)

	small := smaller(&sens25, &sens10)

	fmt.Println("CHECK PRINT: ", small, small.Value())
	ua.New().Send(small.Value(), small.Channel())

	return nil
}

func main() {
	lambda.Start(Handler)
}

func smaller(a *sensor, b *sensor) *sensor {
	if s.LevelOrder(a.level) <= s.LevelOrder(b.level) {
		return a
	}

	return b
}
