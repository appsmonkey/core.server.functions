package defaultDevice

import (
	"encoding/json"

	"github.com/appsmonkey/core.server.functions/dal/access"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
)

func saveState(data *vm.DeviceGetData) {
	state, _ := json.Marshal(data)
	access.SaveState("avg", "last_state", string(state))
}

func getState() vm.DeviceGetData {
	state, ok := access.State("avg", "last_state").(string)

	if !ok || len(state) == 0 {
		return vm.DeviceGetData{}
	}

	var res vm.DeviceGetData
	json.Unmarshal([]byte(state), &res)

	return res
}
