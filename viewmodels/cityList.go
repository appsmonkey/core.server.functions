package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// CityListResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceListData* struct`
type CityListResponse struct {
	BaseResponse
}

// CityListData holds the data to be sent to the client for the *City List*
type CityListData struct {
}

// Marshal the response object
func (r *CityListResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// Init the response object
func (r *CityListResponse) Init() {
	r.Code = 0
	r.RequestID = strconv.FormatInt(time.Now().Unix(), 10)
}

// AddError to the response object
func (r *CityListResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}
