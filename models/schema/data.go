package schema

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	"github.com/aws/aws-sdk-go/aws"
)

func init() {
	defaultData = make(map[string]*Data, 0)
	type schemaData struct {
		Data      map[string]*Data
		Version   string
		Heartbeat int
	}

	schemaRes, err := dal.Get("schema", map[string]*dal.AttributeValue{
		"version": {
			S: aws.String("1"),
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	schema := new(schemaData)
	err = schemaRes.Unmarshal(schema)
	if err != nil {
		fmt.Println("Error unmarshaling schema ::. ", err)
	}

	defaultData = schema.Data

	// defaultData["1"] = &Data{Name: "Temperature", Unit: "℃", DefaultValue: "OK"}
	// defaultData["2"] = &Data{Name: "Humidity", Unit: "g/m3", DefaultValue: "OK"}
	// defaultData["3"] = &Data{Name: "Temperature Feel", Unit: "℃", DefaultValue: "OK"}
	// defaultData["4"] = &Data{Name: "Pressure", Unit: "Pa", DefaultValue: "OK"}
	// defaultData["5"] = &Data{Name: "Altitude", Unit: "m", DefaultValue: "OK"}
	// defaultData["6"] = &Data{Name: "PM 1", Unit: "μg/m³", DefaultValue: "OK"}
	// defaultData["9"] = &Data{Name: "API Range", Unit: "?", DefaultValue: "OK"}
	// defaultData["10"] = &Data{Name: "PM 2.5 Range", Unit: "μg/m³", DefaultValue: "OK"}
	// defaultData["11"] = &Data{Name: "PM 10 Range", Unit: "μg/m³", DefaultValue: "OK"}
	// defaultData["12"] = &Data{Name: "Light Lux", Unit: "℃", DefaultValue: "OK"}
	// defaultData["13"] = &Data{Name: "Eco 2", Unit: "℃", DefaultValue: "OK"}
	// defaultData["14"] = &Data{Name: "TVOC", Unit: "℃", DefaultValue: "OK"}
	// defaultData["15"] = &Data{Name: "Soil Temperature", Unit: "℃", DefaultValue: "OK"}
	// defaultData["16"] = &Data{Name: "Soil Moisture", Unit: "g/m3", DefaultValue: "OK"}
	// defaultData["17"] = &Data{Name: "Unix Time", Unit: "ms", DefaultValue: "OK"}
	// defaultData["18"] = &Data{Name: "Water Level", Unit: "m", DefaultValue: "OK"}
	// defaultData["19"] = &Data{Name: "Motion", Unit: "?", DefaultValue: "OK"}
	// defaultData["7"] = &Data{
	// 	Name:         "PM 2.5",
	// 	Unit:         "μg/m³",
	// 	DefaultValue: "OK",
	// 	CalcSteps: []*CalcStep{
	// 		&CalcStep{
	// 			From:   0,
	// 			To:     12,
	// 			Result: "Great",
	// 		},
	// 		&CalcStep{
	// 			From:   12,
	// 			To:     35,
	// 			Result: "OK",
	// 		},
	// 		&CalcStep{
	// 			From:   35,
	// 			To:     55,
	// 			Result: "Sensitive beware",
	// 		},
	// 		&CalcStep{
	// 			From:   55,
	// 			To:     150,
	// 			Result: "Unhealthy",
	// 		},
	// 		&CalcStep{
	// 			From:   150,
	// 			To:     250,
	// 			Result: "Very Unhealthy",
	// 		},
	// 		&CalcStep{
	// 			From:   250,
	// 			To:     math.MaxFloat64,
	// 			Result: "Hazardous",
	// 		},
	// 	},
	// }
	// defaultData["8"] = &Data{
	// 	Name:         "PM 10",
	// 	Unit:         "μg/m³",
	// 	DefaultValue: "OK",
	// 	CalcSteps: []*CalcStep{
	// 		&CalcStep{
	// 			From:   0,
	// 			To:     54,
	// 			Result: "Great",
	// 		},
	// 		&CalcStep{
	// 			From:   54,
	// 			To:     154,
	// 			Result: "OK",
	// 		},
	// 		&CalcStep{
	// 			From:   154,
	// 			To:     254,
	// 			Result: "Sensitive beware",
	// 		},
	// 		&CalcStep{
	// 			From:   254,
	// 			To:     354,
	// 			Result: "Unhealthy",
	// 		},
	// 		&CalcStep{
	// 			From:   354,
	// 			To:     424,
	// 			Result: "Very Unhealthy",
	// 		},
	// 		&CalcStep{
	// 			From:   424,
	// 			To:     math.MaxFloat64,
	// 			Result: "Hazardous",
	// 		},
	// 	},
	// }
}
