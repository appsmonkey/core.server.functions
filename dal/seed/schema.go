package data

import (
	"math"

	s "github.com/appsmonkey/core.server.functions/models/schema"
)

func loadSchema() {
	defaultData := make(map[string]*s.Data, 0)
	defaultData["1"] = &s.Data{Name: "Temperature", Unit: "℃", DefaultValue: "OK"}
	defaultData["2"] = &s.Data{Name: "Humidity", Unit: "g/m3", DefaultValue: "OK"}
	defaultData["3"] = &s.Data{Name: "Temperature Feel", Unit: "℃", DefaultValue: "OK"}
	defaultData["4"] = &s.Data{Name: "Pressure", Unit: "Pa", DefaultValue: "OK"}
	defaultData["5"] = &s.Data{Name: "Altitude", Unit: "m", DefaultValue: "OK"}
	defaultData["6"] = &s.Data{Name: "PM 1", Unit: "μg/m³", DefaultValue: "OK"}
	defaultData["9"] = &s.Data{Name: "API Range", Unit: "?", DefaultValue: "OK"}
	defaultData["10"] = &s.Data{Name: "PM 2.5 Range", Unit: "μg/m³", DefaultValue: "OK"}
	defaultData["11"] = &s.Data{Name: "PM 10 Range", Unit: "μg/m³", DefaultValue: "OK"}
	defaultData["12"] = &s.Data{Name: "Light Lux", Unit: "℃", DefaultValue: "OK"}
	defaultData["13"] = &s.Data{Name: "Eco 2", Unit: "℃", DefaultValue: "OK"}
	defaultData["14"] = &s.Data{Name: "TVOC", Unit: "℃", DefaultValue: "OK"}
	defaultData["15"] = &s.Data{Name: "Soil Temperature", Unit: "℃", DefaultValue: "OK"}
	defaultData["16"] = &s.Data{Name: "Soil Moisture", Unit: "g/m3", DefaultValue: "OK"}
	defaultData["17"] = &s.Data{Name: "Unix Time", Unit: "ms", DefaultValue: "OK"}
	defaultData["18"] = &s.Data{Name: "Water Level", Unit: "m", DefaultValue: "OK"}
	defaultData["19"] = &s.Data{Name: "Motion", Unit: "?", DefaultValue: "OK"}
	defaultData["7"] = &s.Data{
		Name:         "PM 2.5",
		Unit:         "μg/m³",
		DefaultValue: "OK",
		CalcSteps: []*s.CalcStep{
			&s.CalcStep{
				From:   0,
				To:     12,
				Result: "Great",
			},
			&s.CalcStep{
				From:   12,
				To:     35,
				Result: "OK",
			},
			&s.CalcStep{
				From:   35,
				To:     55,
				Result: "Sensitive beware",
			},
			&s.CalcStep{
				From:   55,
				To:     150,
				Result: "Unhealthy",
			},
			&s.CalcStep{
				From:   150,
				To:     250,
				Result: "Very Unhealthy",
			},
			&s.CalcStep{
				From:   250,
				To:     math.MaxFloat32,
				Result: "Hazardous",
			},
		},
	}
	defaultData["8"] = &s.Data{
		Name:         "PM 10",
		Unit:         "μg/m³",
		DefaultValue: "OK",
		CalcSteps: []*s.CalcStep{
			&s.CalcStep{
				From:   0,
				To:     54,
				Result: "Great",
			},
			&s.CalcStep{
				From:   54,
				To:     154,
				Result: "OK",
			},
			&s.CalcStep{
				From:   154,
				To:     254,
				Result: "Sensitive beware",
			},
			&s.CalcStep{
				From:   254,
				To:     354,
				Result: "Unhealthy",
			},
			&s.CalcStep{
				From:   354,
				To:     424,
				Result: "Very Unhealthy",
			},
			&s.CalcStep{
				From:   424,
				To:     math.MaxFloat32,
				Result: "Hazardous",
			},
		},
	}

	type version struct {
		Version string      `json:"version"`
		Data    interface{} `json:"data"`
	}

	rc := new(RowCollection)
	rc.Table = "schema"
	rc.Rows = []interface{}{version{Version: "1", Data: defaultData}}

	// for k, v := range defaultData {
	// 	ver := &version{
	// 		Version: "1",
	// 		Sensor:  k,
	// 		Data:    v,
	// 	}

	// 	rc.Rows = append(rc.Rows, ver)
	// }

	Storage["schema"] = rc
}
