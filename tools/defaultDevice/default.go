package defaultDevice

import (
	"fmt"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
)

// GetMinimal default device data with minimal data
func GetMinimal() (result vm.DeviceGetDataMinimal) {
	result.DeviceID = ""
	result.Model = "BOXY"
	result.Name = "Sarajevo Air"
	result.Indoor = false
	result.Active = true
	result.DefaultDevice = true

	return
}

// GetFrom will return the default device from the specific time
func GetFrom(from int64) (result vm.DeviceGetData) {
	result.DefaultDevice = true
	result.Active = true
	result.DeviceID = ""
	result.Indoor = false
	result.Location = m.Location{}
	result.MapMeta = make(map[string]m.MapMeta, 0)
	result.Latest = make(map[string]interface{}, 0)
	result.Mine = true
	result.Model = "BOXY"
	result.Name = "Sarajevo Air"
	result.Timestamp = float64(time.Now().Unix())

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

	data := make(map[string][]float64, 0)
	for _, v := range dbData {
		for ki, vi := range v {
			if ki != "timestamp" && ki != "token" && ki != "timestamp_sort" && ki != "ttl" {
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
func Get() (result vm.DeviceGetData) {
	result = GetFrom(time.Now().Add(-time.Hour * 3).Unix())

	if len(result.Latest) == 0 {
		result = GetFrom(time.Now().Add(-time.Hour * 6).Unix())
	}

	if len(result.Latest) == 0 {
		result = GetFrom(time.Now().Add(-time.Hour * 24).Unix())
	}

	// Since we did not get any data, get the last successfull state
	if len(result.Latest) == 0 {
		result = getState()
	} else {
		// We have data, update the state
		saveState(&result)
	}

	return
}
