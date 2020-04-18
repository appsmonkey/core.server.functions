package defaultDevice

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	s "github.com/appsmonkey/core.server.functions/models/schema"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-sdk-go/aws"
)

// GetMinimal default device data with minimal data
func GetMinimal(city string) (result vm.DeviceGetDataMinimal) {
	if len(city) > 0 {
		result.Name = city + " Air"
	} else {
		result.Name = "Sarajevo Air"
	}

	result.DeviceID = ""
	result.Model = "BOXY"
	result.Indoor = false
	result.Active = true
	result.DefaultDevice = true

	return
}

// GetFrom will return the default device from the specific time
func GetFrom(from int64, city string) (result vm.DeviceGetData) {
	if len(city) > 0 {
		result.Name = city + " Air"
	} else {
		result.Name = "Sarajevo Air"
	}

	result.DefaultDevice = true
	result.Active = true
	result.DeviceID = ""
	result.Indoor = false
	result.Location = m.Location{}
	result.MapMeta = make(map[string]m.MapMeta, 0)
	result.Latest = make(map[string]interface{}, 0)
	result.Mine = true
	result.Model = "BOXY"
	result.Timestamp = float64(time.Now().Unix())
	result.ActiveCount = 0
	result.City = city

	dbDevicesData := make([]m.Device, 0)
	dRes, err := dal.ListNoProjection("devices", dal.Name("active").Equal(dal.Value(true)), true)

	if err != nil {
		fmt.Println("could not retirieve devices data")
	}

	err = dRes.Unmarshal(&dbDevicesData)
	if err != nil {
		fmt.Println("Could not unmarshal devices from db", err)
	}

	measurementsData := make(map[string][]float64, 0)
	for _, v := range dbDevicesData {
		// if indoor or diff. city skip
		if v.Meta.Indoor || v.City != city {
			continue
		}

		// count of devices included in the city avg.
		result.ActiveCount++

		// sensors to be ignored for city avg.
		toIgnore := map[string]bool{
			"timestamp":          true,
			"WATER_LEVEL_SWITCH": true,
			"SOIL_MOISTURE":      true,
			"LIGHT_INTENSITY":    true,
			"token":              true,
			"BATTERY_VOLTAGE":    true,
			"BATTERY_PERCENTAGE": true,
			"MOTION":             true,
			"DEVICE_TEMPERATURE": true,
			"timestamp_sort":     true,
			"ttl":                true,
			"city":               true,
			"cognito_id":         true,
			"indoor":             true,
			"zone_id":            true,
			"SOIL_TEMPERATURE":   true,
		}

		for ki, vi := range v.Measurements {
			_, ok := toIgnore[ki]
			if !ok {
				measurementsData[ki] = append(measurementsData[ki], vi.(float64))
			}
		}
	}

	for k, v := range measurementsData {
		if len(v) > 0 {
			var av float64
			for _, sv := range v {
				av += sv
			}

			if av == 0 {
				result.Latest[k] = float64(0)
			} else {
				result.Latest[k] = av / float64(len(v))
			}
		}
	}

	_, ok := result.Latest["AIR_AQI_RANGE"]
	if ok {
		pm10Val := result.Latest["AIR_PM10"]
		pm25Val := result.Latest["AIR_PM2P5"]

		schemaDefault := s.ExtractVersion("")
		pm25, _ := schemaDefault.ExtractData("AIR_PM2P5")
		pm10, _ := schemaDefault.ExtractData("AIR_PM10")

		pm25lvl := pm25.Result(pm25Val.(float64))
		pm10lvl := pm10.Result(pm10Val.(float64))

		if s.LevelOrder(pm10lvl) <= s.LevelOrder(pm25lvl) {
			result.Latest["AIR_AQI_RANGE"] = s.LevelOrder(pm25lvl)
		} else {
			result.Latest["AIR_AQI_RANGE"] = s.LevelOrder(pm10lvl)
		}

		fmt.Println(result.Latest["AIR_AQI_RANGE"])
	}

	for k, v := range result.Latest {
		v = math.Round(v.(float64))
		result.Latest[k] = v
	}

	return
}

// Qsort impl. for sorting by timestamp
func Qsort(a []map[string]interface{}) []map[string]interface{} {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	// Pick a pivot
	pivotIndex := rand.Int() % len(a)

	// Move the pivot to the right
	a[pivotIndex], a[right] = a[right], a[pivotIndex]

	// Pile elements smaller than the pivot on the left
	for i := range a {
		if a[i]["timestamp_sort"].(float64) > a[right]["timestamp_sort"].(float64) {
			a[i], a[left] = a[left], a[i]
			left++
		}
	}

	// Place the pivot after the last smaller element
	a[left], a[right] = a[right], a[left]

	// Go down the rabbit hole
	Qsort(a[:left])
	Qsort(a[left+1:])

	return a
}

// Get default device data
func Get(city string) (result vm.DeviceGetData) {
	result = GetFrom(time.Now().Add(-time.Hour*2).Unix(), city)

	if len(result.Latest) == 0 {
		result = GetFrom(time.Now().Add(-time.Hour*4).Unix(), city)
	}

	fmt.Println("CITY RESULTS ::: ", result, city)

	if len(result.Latest) == 0 {
		result = GetFrom(time.Now().Add(-time.Hour*12).Unix(), city)
	}

	// Since we did not get any data, get the last successfull state
	if len(result.Latest) == 0 {
		fmt.Println("We didn't get any data, calling get state.")
		result = getState(city)
	} else {
		// We have data, update the state
		saveState(&result)
	}

	return
}

func validateDevice(token string, city string) bool {
	res, err := dal.Get("devices", map[string]*dal.AttributeValue{
		"token": {
			S: aws.String(token),
		},
	})
	if err != nil {
		fmt.Println("Error fetching device")
	}

	model := m.Device{}
	err = res.Unmarshal(&model)
	if err != nil {
		fmt.Println("Error unmarshaling device")
	}

	if model.Meta.Indoor || model.City != city {
		return false
	}

	return true
}
