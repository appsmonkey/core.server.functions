package data

import (
	"fmt"

	m "github.com/appsmonkey/core.server.functions/models"
	s "github.com/appsmonkey/core.server.functions/models/schema"

	z "github.com/appsmonkey/core.server.functions/tools/zones"
	// Loading the sarajevo map
	_ "github.com/appsmonkey/core.server.functions/tools/zones/sarajevo"
)

func loadZones() {
	data := new(RowCollection)
	data.Table = "zones"
	data.Rows = make([]interface{}, 0)
	zones := z.All()
	sensors := s.ExtractVersion("1")

	fmt.Println("SENSORS:", sensors)

	for zk := range zones {
		for sk := range sensors {
			ld, ln := s.SensorReading("1", sk, 0)
			zd := m.Zone{
				ZoneID:   "Sarajevo@" + zk,
				CityID:   "Sarajevo",
				SensorID: sk,
				Data: m.ZoneMeta{
					SensorID:    sk,
					Name:        zk,
					Level:       ln,
					Value:       0,
					Measurement: ld.Name,
					Unit:        ld.Unit,
				},
			}

			data.Rows = append(data.Rows, zd)
		}
	}

	Storage["zones"] = data
}
