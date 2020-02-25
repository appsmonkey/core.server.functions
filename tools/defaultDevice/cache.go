package defaultDevice

import (
	"encoding/json"
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal/access"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
)

func saveState(data *vm.DeviceGetData) {
	fmt.Println("Saving state for ::: ", data)
	state, _ := json.Marshal(data)
	access.SaveState("avg_"+data.City, "last_state", string(state))
}

func getState(city string) vm.DeviceGetData {
	state, ok := access.State("avg_"+city, "last_state").(string)

	if !ok || len(state) == 0 {
		return vm.DeviceGetData{}
	}

	var res vm.DeviceGetData
	json.Unmarshal([]byte(state), &res)

	return res
}
