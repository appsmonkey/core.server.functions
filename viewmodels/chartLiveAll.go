package viewmodels

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ChartLiveAllRequest is the request from the client
type ChartLiveAllRequest struct {
	Sensor    string   `json:"sensor"`
	From      int64    `json:"from"`
	SensorAll []string `json:"-"`
}

// Validate the request sent from client
func (r *ChartLiveAllRequest) Validate(body map[string]string) *ChartLiveAllResponse {
	response := new(ChartLiveAllResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	if len(body) == 0 {
		errData := es.ErrIncorrectRequest
		errData.Data = "no query parameters sent"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartLiveAllError
		return response
	}

	sensor, ok := body["sensor"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "sensor parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartLiveAllError
		return response
	}

	from, ok := body["from"]
	if !ok {
		errData := es.ErrIncorrectRequest
		errData.Data = "from parameter is missing"
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusChartLiveAllError
		return response
	}

	if len(sensor) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusChartLiveAllError
	}

	sensorAll := strings.Split(sensor, ",")
	if len(sensorAll) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingSensorType)
		response.Code = es.StatusChartAllDeviceError
	}

	f, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		response.Errors = append(response.Errors, es.ErrIncorrectTime)
		response.Code = es.StatusChartLiveAllError
	}

	r.Sensor = sensor
	r.From = f
	r.SensorAll = sensorAll

	return response
}

// ChartLiveAllResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type ChartLiveAllResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *ChartLiveAllResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// AddError to the response object
func (r *ChartLiveAllResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}
