package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ZoneDelRequest is the request from the client
type ZoneDelRequest struct {
	ZoneID string `json:"zone_id"`
}

// Validate the request sent from client
func (r *ZoneDelRequest) Validate(body string) *ZoneDelResponse {
	response := new(ZoneDelResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusDeleteZoneError
		return response
	}

	if len(r.ZoneID) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingID)
		response.Code = es.StatusDeleteZoneError
	}

	return response
}

// ZoneDelResponse to the client
type ZoneDelResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *ZoneDelResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
