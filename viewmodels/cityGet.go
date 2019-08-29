package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"
)

// CityGetRequest is the request from the client
type CityGetRequest struct {
	CityID string `json:"city_id"`
}

// Validate the request sent from client
func (r *CityGetRequest) Validate(body map[string]string) *DeviceAddResponse {
	response := new(DeviceAddResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	id, ok := body["city_id"]
	if ok {
		r.CityID = id
	}

	return response
}

// CityGetResponse to the client
// `Returns detailed data for a specific device. Data defained in the *DeviceGetData* struct`
type CityGetResponse struct {
	BaseResponse
}

// CityGetData returned to user
type CityGetData struct {
	CityID    string   `json:"device_id"`
	Name      string   `json:"name"`
	Country   string   `json:"country"`
	Zones     []string `json:"zones"`
	Timestamp float64  `json:"timestamp"`
}

// CityGetDataMinimal returned to user
type CityGetDataMinimal struct {
	CityID  string `json:"device_id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

// Marshal the response object
func (r *CityGetResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
