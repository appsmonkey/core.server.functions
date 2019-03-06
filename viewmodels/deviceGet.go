package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// DeviceGetRequest is the request from the client
type DeviceGetRequest struct {
	Token string `json:"token"`
}

// Validate the request sent from client
func (r *DeviceGetRequest) Validate(body string) *DeviceAddResponse {
	response := new(DeviceAddResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusGetDeviceError
		return response
	}

	if len(r.Token) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingToken)
		response.Code = es.StatusGetDeviceError
	}

	return response
}

// DeviceGetResponse to the client
// `Returns detailed data for a specific device. Data defained in the *DeviceGetData* struct`
type DeviceGetResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *DeviceGetResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
