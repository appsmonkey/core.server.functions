package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// DeviceListResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceListData* struct`
type DeviceListResponse struct {
	BaseResponse
}

// DeviceListData holds the data to be sent to the client for the *Device List*
// `List of Device meta data`
type DeviceListData struct {
}

// Marshal the response object
func (r *DeviceListResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// Init the response object
func (r *DeviceListResponse) Init() {
	r.Code = 0
	r.RequestID = strconv.FormatInt(time.Now().Unix(), 10)
}

// AddError to the response object
func (r *DeviceListResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}
