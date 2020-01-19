package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ForgotPasswordEndRequest sent from the client
type ForgotPasswordEndRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	CognitoID string `json:"cognito_id"`
}

// ForgotPasswordEndResponse to the client
type ForgotPasswordEndResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *ForgotPasswordEndRequest) Validate(body string) *ForgotPasswordEndResponse {
	response := new(ForgotPasswordEndResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusForgotPasswordError
		return response
	}

	if len(r.Token) == 0 {
		response.Errors = append(response.Errors, es.ForgotPasswordBadRequestNoToken)
		response.Code = es.StatusForgotPasswordError
	}

	if len(r.CognitoID) == 0 {
		response.Errors = append(response.Errors, es.ForgotPasswordBadRequestNoID)
		response.Code = es.StatusForgotPasswordError
	}

	if !validateEmail(r.Email) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingEmail)
		response.Code = es.StatusForgotPasswordError
	}

	if !validatePassword(r.Password) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingPass)
		response.Code = es.StatusForgotPasswordError
	}

	return response
}

// AddError to the response object
func (r *ForgotPasswordEndResponse) AddError(err *es.Error) {
	r.Errors = append(r.Errors, *err)
}

// Marshal the response object
func (r *ForgotPasswordEndResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
