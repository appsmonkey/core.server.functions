package viewmodels

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
)

// RegisterFillUserDataRequest sent from cognito
type RegisterFillUserDataRequest struct {
	UserProfile m.UserProfile `json:"user_profile"`
	Token       string        `json:"token"`
	UserName    string        `json:"user_name,omitempty"`
	CognitoID   string        `json:"cognito_id,omitempty"`
	Password    string        `json:"password,omitempty"`
}

// Validate the request sent from client
func (r *RegisterFillUserDataRequest) Validate(body string) *RegisterFillUserDataResponse {
	response := new(RegisterFillUserDataResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		fmt.Println("Incorrect request error!")
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusProfileUpdateError
		return response
	}

	if len(r.CognitoID) == 0 {
		response.Errors = append(response.Errors, es.UserCreationFailedNoID)
		response.Code = es.StatusProfileUpdateError
	}

	if len(r.Token) == 0 {
		response.Errors = append(response.Errors, es.UserCreationFailedNoToken)
		response.Code = es.StatusProfileUpdateError
	}

	if len(r.UserName) == 0 {
		response.Errors = append(response.Errors, es.UserCreationFailedNoUserName)
		response.Code = es.StatusProfileUpdateError
	}

	if len(r.Password) == 0 {
		response.Errors = append(response.Errors, es.UserCreationFailedNoPassword)
		response.Code = es.StatusProfileUpdateError
	}

	return response
}

// RegisterFillUserDataResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type RegisterFillUserDataResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *RegisterFillUserDataResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// AddError to the response
func (r *RegisterFillUserDataResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}

// Marshal the response object
func (r *RegisterFillUserDataRequest) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
