package general

import (
	"fmt"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
)

// Get default device data
func Get() (result vm.DeviceGetData) {
	result.Active = true
	result.DeviceID = "DEFAULT"
	result.Indoor = false
	result.Location = m.Location{}
	result.MapMeta = make(map[string]m.MapMeta, 0)
	result.Latest = make(map[string]interface{}, 0)
	result.Mine = true
	result.Model = "BOXY"
	result.Name = "Sarajevo Air"
	result.Timestamp = float64(time.Now().Unix())

	from := time.Now().Add(-time.Hour * 3).Unix()
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
