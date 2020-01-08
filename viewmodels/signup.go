package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// SignupRequest sent from the client
type SignupRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	ClientID  string `json:"client_id"`
}

// SignupResponse to the client
type SignupResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *SignupRequest) Validate(body string) *SignupResponse {
	response := new(SignupResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusRegistrationError
		return response
	}

	if !validateEmail(r.Email) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingEmail)
		response.Code = es.StatusRegistrationError
	}

	if !validatePassword(r.Password) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingPass)
		response.Code = es.StatusRegistrationError
	}

	return response
}

// Marshal the response object
func (r *SignupResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
