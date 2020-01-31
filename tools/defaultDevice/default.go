package defaultDevice

import (
	"fmt"
	"math/rand"
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

	// sort data by timestamp
	dbData = Qsort(dbData)

	fmt.Println(dbData, "DB DATA")
	// make data distinc
	var distinctData []map[string]interface{}
	var keyList = make(map[string]bool)
	for _, v := range dbData {
		if _, ok := keyList[v["token"].(string)]; ok {
			continue
		} else {
			keyList[v["token"].(string)] = true
			distinctData = append(distinctData, v)
		}
	}
	fmt.Println("Filtered key list ::: ", keyList)
	fmt.Println("DISTINCT DATA :::", distinctData)

	data := make(map[string][]float64, 0)
	for _, v := range distinctData {
		if v["city"] == city {
			result.ActiveCount++
		}

		if v["indoor"] == true || v["indoor"] == "true" || v["city"] != city {
			continue
		}

		for ki, vi := range v {
			if ki != "timestamp" && ki != "token" && ki != "timestamp_sort" && ki != "ttl" && ki != "city" && ki != "cognito_id" && ki != "indoor" && ki != "zone_id" {

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
		if a[i]["timestamp"].(float64) > a[right]["timestamp"].(float64) {
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

	if len(result.Latest) == 0 {
		result = GetFrom(time.Now().Add(-time.Hour*12).Unix(), city)
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
