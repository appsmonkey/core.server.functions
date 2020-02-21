package viewmodels

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ChartHourAllRequest is the request from the client
type ChartHourAllRequest struct {
	Sensor    string   `json:"sensor"`
	From      string   `json:"from"`
	SensorAll []string `json:"-"`
	City      string   `json:"city"`
}

// Validate the request sent from client
func (r *ChartHourAllRequest) Validate(body map[string]string) *ChartHourAllResponse {
	response := new(ChartHourAllResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	if len(body) == 0 {
		errData := es.ErrIncorrectRequest
		errData.Data = "no query parameters sent"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartAllDeviceError
		return response
	}

	sensor, ok := body["sensor"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "sensor parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartAllDeviceError
		return response
	}

	from, ok := body["from"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "from parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartAllDeviceError
		return response
	}

	city, ok := body["city"]
	if !ok {
		city = "Sarajevo"
	}

	if len(sensor) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusChartAllDeviceError
	}

	sensorAll := strings.Split(sensor, ",")
	if len(sensorAll) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusChartAllDeviceError
	}

	_, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		response.Errors = append(response.Errors, es.ErrIncorrectTime)
		response.Code = es.StatusChartAllDeviceError
	}

	r.Sensor = sensor
	r.From = from
	r.SensorAll = sensorAll
	r.City = city

	return response
}

// ChartHourAllResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type ChartHourAllResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *ChartHourAllResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// AddError to the response object
func (r *ChartHourAllResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}
