package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"
)

// ZoneGetRequest is the request from the client
type ZoneGetRequest struct {
	ZoneID string `json:"zone_id"`
}

// ZoneGetResponse to the client
type ZoneGetResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *ZoneGetRequest) Validate(body map[string]string) *ZoneGetResponse {
	response := new(ZoneGetResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	id, ok := body["zone_id"]
	if ok {
		r.ZoneID = id
	}

	return response
}

// CityGetResponse to the client
// `Returns detailed data for a specific device. Data defained in the *DeviceGetData* struct`
// type CityGetResponse struct {
// 	BaseResponse
// }

// CityGetData returned to user
// type CityGetData struct {
// 	// ?? what do we need on the client ??
// }

// Marshal the response object
func (r *ZoneGetResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
