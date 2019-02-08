package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// SigninRequest sent from the client
type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SigninResponse to the client
type SigninResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *SigninRequest) Validate(body string) *SigninResponse {
	response := new(SigninResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusSignInError
		return response
	}

	if !validateEmail(r.Email) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingEmail)
		response.Code = es.StatusSignInError
	}

	if !validatePassword(r.Password) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingPass)
		response.Code = es.StatusSignInError
	}

	return response
}

// Marshal the response object
func (r *SigninResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
