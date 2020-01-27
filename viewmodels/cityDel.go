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
func (r *CityDelRequest) Validate(body map[string]string) *CityDelResponse {
	response := new(CityDelResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	cid, ok := body["city_id"]
	if ok {
		r.CityID = cid
	}

	if len(r.CityID) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingID)
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
