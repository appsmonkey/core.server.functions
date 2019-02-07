package viewmodels

import (
	"encoding/json"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// RegisterRequest sent from the client
type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse to the client
type RegisterResponse struct {
	BaseResponse
}

// RegisterData holding request/response specific data to be returned to the client
type RegisterData struct {
}

// Validate the request sent from client
func (r *RegisterRequest) Validate(body string) *RegisterResponse {
	response := new(RegisterResponse)
	response.Code = 0

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

	if len(r.FullName) == 0 {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingName)
		response.Code = es.StatusRegistrationError
	}

	return response
}

// Marshal the response object
func (r *RegisterResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
