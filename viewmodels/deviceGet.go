package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	m "github.com/appsmonkey/core.server.functions/models"
)

// DeviceGetRequest is the request from the client
type DeviceGetRequest struct {
	Token string `json:"token"`
}

// Validate the request sent from client
func (r *DeviceGetRequest) Validate(body map[string]string) *DeviceAddResponse {
	response := new(DeviceAddResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	token, ok := body["token"]
	if ok {
		r.Token = token
	}

	return response
}

// DeviceGetResponse to the client
// `Returns detailed data for a specific device. Data defained in the *DeviceGetData* struct`
type DeviceGetResponse struct {
	BaseResponse
}

// DeviceGetData returned to user
type DeviceGetData struct {
	DeviceID  string                 `json:"device_id"`
	Name      string                 `json:"name"`
	Active    bool                   `json:"active"`
	Model     string                 `json:"model"`
	Indoor    bool                   `json:"indoor"`
	Mine      bool                   `json:"mine"`
	Location  m.Location             `json:"location"`
	MapMeta   map[string]m.MapMeta   `json:"map_meta"`
	Latest    map[string]interface{} `json:"latest"`
	Timestamp float64                `json:"timestamp"`
}

// DeviceGetDataMinimal returned to user
type DeviceGetDataMinimal struct {
	DeviceID string `json:"device_id"`
	Name     string `json:"name"`
	Active   bool   `json:"active"`
	Model    string `json:"model"`
	Indoor   bool   `json:"indoor"`
}

// Marshal the response object
func (r *DeviceGetResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
