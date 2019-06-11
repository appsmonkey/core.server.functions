package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// Social data
type Social struct {
	Type  string `json:"type"`
	Token string `json:"token"`
	ID    string `json:"id"`
}

// SigninRequest sent from the client
type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Social   Social `json:"social"`
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

	if !r.Social.Validate() {
		response.Errors = append(response.Errors, es.ErrMissingSocialData)
		response.Code = es.StatusSignInError
	}

	if !validateEmail(r.Email) {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingEmail)
		response.Code = es.StatusSignInError
	}

	if !r.Social.HasData() && !validatePassword(r.Password) {
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

// Validate Social data
func (s Social) Validate() bool {
	if s.HasData() && (len(s.Type) == 0 || len(s.Token) == 0 || len(s.ID) == 0) {
		return false
	}

	return true
}

// HasData indicates if social data has been passed
func (s Social) HasData() bool {
	// No data provided, normal login
	if len(s.Token) > 0 && len(s.Type) > 0 && len(s.ID) > 0 {
		return true
	}

	return false
}
