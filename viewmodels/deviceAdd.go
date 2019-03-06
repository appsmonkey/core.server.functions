package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
)

// DeviceAddRequest is the request from the client
type DeviceAddRequest struct {
	m.Metadata
}

// Validate the request sent from client
func (r *DeviceAddRequest) Validate(body string) *DeviceAddResponse {
	response := new(DeviceAddResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusAddDeviceError
		return response
	}

	if len(r.Model) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingModel)
		response.Code = es.StatusAddDeviceError
	}

	if len(r.Coordinates.Long) == 0 || len(r.Coordinates.Lat) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingLocation)
		response.Code = es.StatusAddDeviceError
	}

	if len(r.Name) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingName)
		response.Code = es.StatusAddDeviceError
	}

	return response
}

// DeviceAddResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type DeviceAddResponse struct {
	BaseResponse
}

// DeviceAddData holds the data to be sent to the client for *Device Add*
type DeviceAddData struct {
	Token string `json:"token"`
}

// Marshal the response object
func (r *DeviceAddResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
