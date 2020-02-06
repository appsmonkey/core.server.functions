package viewmodels

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ChartHourDeviceRequest is the request from the client
type ChartHourDeviceRequest struct {
	Token     string   `json:"token"`
	Sensor    string   `json:"sensor"`
	SensorAll []string `json:"-"`
	From      string   `json:"from"`
}

// Validate the request sent from client
func (r *ChartHourDeviceRequest) Validate(body map[string]string) *ChartHourDeviceResponse {
	response := new(ChartHourDeviceResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	if len(body) == 0 {
		errData := es.ErrIncorrectRequest
		errData.Data = "no query parameters sent"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartHourDeviceError
		return response
	}

	sensor, ok := body["sensor"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "sensor parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartHourDeviceError
		return response
	}

	from, ok := body["from"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "from parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartHourDeviceError
		return response
	}

	token, ok := body["token"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "token parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartHourDeviceError
		return response
	}

	if len(sensor) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusChartHourDeviceError
	}

	sensorAll := strings.Split(sensor, ",")
	if len(sensorAll) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusChartAllDeviceError
	}

	if len(token) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingToken)
		response.Code = es.StatusChartHourDeviceError
	}

	_, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		response.Errors = append(response.Errors, es.ErrIncorrectTime)
		response.Code = es.StatusChartHourDeviceError
	}

	r.Sensor = sensor
	r.Token = token
	r.SensorAll = sensorAll
	r.From = from

	return response
}

// ChartHourDeviceResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type ChartHourDeviceResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *ChartHourDeviceResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// AddError to the response object
func (r *ChartHourDeviceResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}
