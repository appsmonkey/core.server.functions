package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// CityDelRequest is the request from the client
type CityDelRequest struct {
	CityID string `json:"city_id"`
}

// Validate the request sent from client
func (r *CityDelRequest) Validate(body string) *CityDelResponse {
	response := new(CityDelResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusDeleteCityError
		return response
	}

	if len(r.CityID) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingThingToken)
		response.Code = es.StatusDeleteCityError
	}

	return response
}

// CityDelResponse to the client
type CityDelResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *CityDelResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
