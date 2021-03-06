package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// ProfileRequest sent from the client
type ProfileRequest struct {
	Email string `json:"email"`
}

// ProfileResponse to the client
type ProfileResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *ProfileRequest) Validate(body string) *ProfileResponse {
	response := new(ProfileResponse)
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

	return response
}

// Marshal the response object
func (r *ProfileResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
