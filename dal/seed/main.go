package data

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
)

// RowCollection represents one seed collection
type RowCollection struct {
	Table string
	Rows  []interface{}
}

// Save the seed into the DB
func (rc *RowCollection) Save() int {
	if len(rc.Rows) > 0 && len(rc.Table) > 0 {
		for _, r := range rc.Rows {
			dal.Insert(rc.Table, r)
			// time.Sleep(500 * time.Millisecond) // wait for half a second after each insert to not excert the write capacity on DynamoDB :D
		}

		return len(rc.Rows)
	}

	return 0
}

// Storage holds all possible seeds stgored by a unique identifier
var Storage map[string]*RowCollection

// Run a particular seed.
// Use `all` for all seeds
func Run(seed string) {
	if seed == "all" {
		for _, rc := range Storage {
			count := rc.Save()
			fmt.Printf("Saved %v rows for seed %v\n", count, seed)
		}

		return
	}

	rc, ok := Storage[seed]
	if ok {
		count := rc.Save()
		fmt.Printf("Saved %v rows for seed %v\n", count, seed)
	}
}
