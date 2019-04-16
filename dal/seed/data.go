package data

func init() {
	Storage = make(map[string]*RowCollection, 0)
	loadSchema()
	loadZones()
}
