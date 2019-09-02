package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ZoneListRequest is the request from the client
// FIXME: We list zones by cities ? What if we have zones with cities which are not in db ?
type ZoneListRequest struct {
	CityID string `json:"city_id"`
}

// Validate the request sent from client
func (r *ZoneListRequest) Validate(body string) *ZoneListResponse {
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

	if len(r.CityID) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingID)
		response.Code = es.StatusDeleteZoneError
	}

	return response
}

// ZoneListResponse to the client
type ZoneListResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *ZoneListResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
