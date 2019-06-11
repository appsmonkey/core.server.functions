package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ChartHasDataRequest is the request from the client
type ChartHasDataRequest struct {
	Chart  string `json:"chart"`
	Sensor string `json:"sensor"`
	Device bool   `json:"-"`
	From   int64  `json:"from"`
	Token  string `json:"token"`
}

// Validate the request sent from client
func (r *ChartHasDataRequest) Validate(body map[string]string) *ChartHasDataResponse {
	response := new(ChartHasDataResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	if len(body) == 0 {
		errData := es.ErrIncorrectRequest
		errData.Data = "no query parameters sent"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartHasDataError
		return response
	}

	chart, ok := body["chart"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "chart parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartHasDataError
		return response
	}

	sensor, ok := body["sensor"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "sensor parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartHasDataError
		return response
	}

	from, ok := body["from"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "from parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartHasDataError
		return response
	}

	device := body["token"]

	if len(chart) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingChart)
		response.Code = es.StatusChartHasDataError
	}

	if len(sensor) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusChartHasDataError
	}

	f, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		response.Errors = append(response.Errors, es.ErrIncorrectTime)
		response.Code = es.StatusChartHasDataError
	}

	if len(device) > 0 {
		r.Device = true
	}

	r.Chart = chart
	r.Sensor = sensor
	r.From = f
	r.Token = device

	return response
}

// ChartHasDataResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type ChartHasDataResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *ChartHasDataResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// AddError to the response object
func (r *ChartHasDataResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}
