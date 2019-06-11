package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ForgotPasswordEndRequest sent from the client
type ForgotPasswordEndRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
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

	if !validateEmail(r.Email) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingEmail)
		response.Code = es.StatusForgotPasswordError
	}

	if len(r.Code) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingCode)
		response.Code = es.StatusForgotPasswordError
	}

	if !validatePassword(r.Password) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingPass)
		response.Code = es.StatusForgotPasswordError
	}

	return response
}

// Marshal the response object
func (r *ForgotPasswordEndResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
