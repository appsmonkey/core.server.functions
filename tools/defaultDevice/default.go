package defaultDevice

import (
	"fmt"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
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

	res, err := dal.ListNoProjection("live", dal.Name("timestamp").GreaterThanEqual(dal.Value(from)))

	if err != nil {
		fmt.Println("could not retirieve data")
		return
	}

	var dbData []map[string]interface{}
	err = res.Unmarshal(&dbData)
	if err != nil {
		fmt.Println("could not unmarshal data from the DB")
		return
	}

	fmt.Println("Fetched db data ::: ", dbData, "Count: ", len(dbData))

	data := make(map[string][]float64, 0)
	for _, v := range dbData {
		if v["indoor"] == true || v["indoor"] == "true" || v["city"] != city {
			continue
		}

		result.ActiveCount++

		for ki, vi := range v {
			if ki != "timestamp" && ki != "token" && ki != "timestamp_sort" && ki != "ttl" && ki != "city" && ki != "cognito_id" && ki != "indoor" && ki != "zone_id" {

				// if ki == "AIR_TEMPERATURE" || ki == "AIR_TEMPERATURE_FEEL" {
				// if v["indoor"] == true || v["indoor"] == "true" {
				// 	fmt.Println("Indoor device skipping sensor: ki")
				// 	continue
				// }
				// }

				data[ki] = append(data[ki], vi.(float64))
			}
		}
	}

	for k, v := range data {
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

	return
}

// Get default device data
func Get(city string) (result vm.DeviceGetData) {
	result = GetFrom(time.Now().Add(-time.Hour*3).Unix(), city)

	if len(result.Latest) == 0 {
		result = GetFrom(time.Now().Add(-time.Hour*6).Unix(), city)
	}

	if len(result.Latest) == 0 {
		result = GetFrom(time.Now().Add(-time.Hour*24).Unix(), city)
	}

	// Since we did not get any data, get the last successfull state
	if len(result.Latest) == 0 {
		fmt.Println("We didn't get any data, calling get state.")
		result = getState()
	} else {
		// We have data, update the state
		saveState(&result)
	}

	return
}
