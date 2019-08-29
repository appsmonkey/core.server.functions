package viewmodels

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// CityAddRequest is the request from the client
type CityAddRequest struct {
	CityID  string   `json:"id"`
	Name    string   `json:"name"`
	Country string   `json:"country"` // ID
	Zones   []string `json:"zones"`   // IDs
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

	if len(r.Name) == 0 {
		response.Errors = append(response.Errors, es.ErrCityMissingName)
		response.Code = es.StatusCityMissingName
	}

	if len(r.Country) == 0 {
		response.Errors = append(response.Errors, es.ErrCityMissingCountry)
		response.Code = es.StatusCityMissingCountry
	}

	fmt.Sprintln(r.Zones, "Zones")
	rt := reflect.TypeOf(r.Zones)
	if len(r.Zones) == 0 && rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		response.Errors = append(response.Errors, es.ErrCityMissingZones)
		response.Code = es.StatusCityMissingZones
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
	Token string `json:"token"`
}

// Marshal the response object
func (r *CityAddResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
