package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ValidateEmailRequest sent from the client
type ValidateEmailRequest struct {
	Email string `json:"email"`
}

// ValidateEmailResponse to the client
type ValidateEmailResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *ValidateEmailRequest) Validate(body string) *ValidateEmailResponse {
	response := new(ValidateEmailResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	err := json.Unmarshal([]byte(body), r)
	if err != nil {
		errData := es.ErrRegistrationIncorrectRequest
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)

		response.Code = es.StatusValidateEmailError
		return response
	}

	if !validateEmail(r.Email) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingEmail)
		response.Code = es.StatusValidateEmailError
	}

	return response
}

// Marshal the response object
func (r *ValidateEmailResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
