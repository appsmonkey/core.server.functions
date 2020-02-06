package viewmodels

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ChartLiveDeviceRequest is the request from the client
type ChartLiveDeviceRequest struct {
	Token     string   `json:"token"`
	Sensor    string   `json:"sensor"`
	SensorAll []string `json:"-"`
	From      string   `json:"from"`
}

// Validate the request sent from client
func (r *ChartLiveDeviceRequest) Validate(body map[string]string) *ChartLiveDeviceResponse {
	response := new(ChartLiveDeviceResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	if len(body) == 0 {
		errData := es.ErrIncorrectRequest
		errData.Data = "no query parameters sent"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartLiveDeviceError
		return response
	}

	sensor, ok := body["sensor"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "sensor parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartLiveDeviceError
		return response
	}

	from, ok := body["from"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "from parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartLiveDeviceError
		return response
	}

	token, ok := body["token"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "token parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartLiveDeviceError
		return response
	}

	if len(sensor) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusChartLiveDeviceError
	}

	sensorAll := strings.Split(sensor, ",")
	if len(sensorAll) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusChartAllDeviceError
	}

	if len(token) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingToken)
		response.Code = es.StatusChartLiveDeviceError
	}

	_, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		response.Errors = append(response.Errors, es.ErrIncorrectTime)
		response.Code = es.StatusChartLiveDeviceError
	}

	r.Sensor = sensor
	r.Token = token
	r.From = from
	r.SensorAll = sensorAll
	return response
}

// ChartLiveDeviceResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type ChartLiveDeviceResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *ChartLiveDeviceResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// AddError to the response object
func (r *ChartLiveDeviceResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}
