package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// CognitoProfileListResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type CognitoProfileListResponse struct {
	BaseResponse
}

// Init the request sent from client
func (r *CognitoProfileListResponse) Init() {
	r.Code = 0
	r.RequestID = strconv.FormatInt(time.Now().Unix(), 10)
}

// AddError to the response object
func (r *CognitoProfileListResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}

// Marshal the response object
func (r *CognitoProfileListResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
