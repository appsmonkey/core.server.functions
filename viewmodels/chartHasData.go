package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ChartHasDataRequest is the request from the client
type ChartHasDataRequest struct {
	Chart string `json:"chart"`
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

	if len(chart) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingChart)
		response.Code = es.StatusChartHasDataError
	}

	r.Chart = chart

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
