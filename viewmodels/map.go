package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
)

// MapRequest is the request from the client
type MapRequest struct {
	Sensor string `json:"sensor"`
}

// Validate the request sent from client
func (r *MapRequest) Validate(body map[string]string) *MapResponse {
	response := new(MapResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	if len(body) == 0 {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = "no query parameters sent"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusMapError
		return response
	}

	sensor, ok := body["sensor"]
	if !ok {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = "sensor parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusMapError
		return response
	}

	r.Sensor = sensor

	if len(r.Sensor) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusMapError
	}

	return response
}

// MapResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type MapResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *MapResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// AddError to the response object
func (r *MapResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}

// MapResponseData to be returned
type MapResponseData struct {
	Zones   []m.Zone    `json:"zones"`
	Devices []m.MapMeta `json:"devices"`
}
