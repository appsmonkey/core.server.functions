package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ForgotPasswordStartRequest sent from the client
type ForgotPasswordStartRequest struct {
	Email string `json:"email"`
}

// ForgotPasswordStartResponse to the client
type ForgotPasswordStartResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *ForgotPasswordStartRequest) Validate(body string) *ForgotPasswordStartResponse {
	response := new(ForgotPasswordStartResponse)
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

	return response
}

// Marshal the response object
func (r *ForgotPasswordStartResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
