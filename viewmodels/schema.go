package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"
)

// SchemaInitResponse the response
func SchemaInitResponse() *SchemaResponse {
	response := new(SchemaResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	return response
}

// SchemaResponse to the client
// `Returns detailed data for a specific device. Data defained in the *DeviceGetData* struct`
type SchemaResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *SchemaResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
