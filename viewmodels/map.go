package viewmodels

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// MapRequest is the request from the client
type MapRequest struct {
	Zone   []string `json:"zone_data"`
	Sensor []string `json:"device_data"`
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

	sensor, ok := body["device_data"]
	if !ok {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = "device_data parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusMapError
		return response
	}

	zone, ok := body["zone_data"]
	if !ok {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = "zone_data parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusMapError
		return response
	}

	r.Sensor = strings.Split(sensor, ",")
	r.Zone = strings.Split(zone, ",")

	if len(r.Sensor) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusMapError
	}

	if len(r.Zone) == 0 {
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
