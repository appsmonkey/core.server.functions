package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
)

// CognitoProfileUpdateRequest sent from cognito
type CognitoProfileUpdateRequest struct {
	m.UserProfile
}

// Validate the request sent from client
func (r *CognitoProfileUpdateRequest) Validate(body string) *CognitoProfileUpdateResponse {
	response := new(CognitoProfileUpdateResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusProfileUpdateError
		return response
	}

	return response
}

// CognitoProfileUpdateResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type CognitoProfileUpdateResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *CognitoProfileUpdateResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// AddError to the response object
func (r *CognitoProfileUpdateResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}

// Marshal the response object
func (r *CognitoProfileUpdateRequest) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
