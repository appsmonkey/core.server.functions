package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
)

// DeviceUpdateMetaRequest is the request from the client
type DeviceUpdateMetaRequest struct {
	m.Metadata
	Token string `json:"token"`
	City  string `json:"city"`
}

// Validate the request sent from client
func (r *DeviceUpdateMetaRequest) Validate(body string) *DeviceUpdateMetaResponse {
	response := new(DeviceUpdateMetaResponse)
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

	if r.Coordinates.Lng == 0 && r.Coordinates.Lat == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingLocation)
		response.Code = es.StatusAddDeviceError
	}

	if len(r.Name) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingName)
		response.Code = es.StatusAddDeviceError
	}

	if len(r.City) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingCity)
		response.Code = es.StatusAddDeviceError
	}

	return response
}

// DeviceUpdateMetaResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type DeviceUpdateMetaResponse struct {
	BaseResponse
}

// DeviceUpdateMetaData holds the data to be sent to the client for *Device Add*
type DeviceUpdateMetaData struct {
	Success bool `json:"success"`
}

// Marshal the response object
func (r *DeviceUpdateMetaResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// AddError to the response object
func (r *DeviceUpdateMetaResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}
