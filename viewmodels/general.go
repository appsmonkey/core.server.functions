package viewmodels

import (
	"encoding/json"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// GeneralRequest sent from the client
type GeneralRequest struct {
	Temp     float64 `json:"temp"`
	Humidity float64 `json:"humidity"`
	PM10     float64 `json:"pm10"`
	PM25     float64 `json:"pm25"`
	Date     int64   `json:"date"`
}

// GeneralResponse to the client
type GeneralResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *GeneralRequest) Validate(body string) *GeneralResponse {
	response := new(GeneralResponse)
	response.Code = 0

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusRegistrationError
		return response
	}

	return response
}

// Marshal the response object
func (r *GeneralResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
