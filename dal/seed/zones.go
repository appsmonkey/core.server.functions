package data

import (
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

	for zk := range zones {
		for sk := range sensors {
			ld, ln := s.SensorReading("1", sk, 0)
			zd := m.Zone{
				ZoneID:   zk,
				SensorID: sk,
				Data: m.ZoneMeta{
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
