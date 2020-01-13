package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// CityAddRequest is the request from the client
type CityAddRequest struct {
	CityID  string `json:"city_id"`
	Country string `json:"country"` // FIXME: maybe we want to have country as an entity too ?
}

// Validate the request sent from client
func (r *CityAddRequest) Validate(body string) *CityAddResponse {
	response := new(CityAddResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusAddCityError
		return response
	}

	if len(r.CityID) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingCityID)
		response.Code = es.StatusCityMissingCountry
	}

	if len(r.Country) == 0 {
		response.Errors = append(response.Errors, es.ErrCityMissingCountry)
		response.Code = es.StatusCityMissingCountry
	}

	return response
}

// CityAddResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type CityAddResponse struct {
	BaseResponse
}

// CityAddData holds the data to be sent to the client for *Device Add*
type CityAddData struct {
	CityID string `json:"city_id"`
}

// Marshal the response object
func (r *CityAddResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
