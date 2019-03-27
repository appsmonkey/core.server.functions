package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// RefreshRequest sent from the client
type RefreshRequest struct {
	Token string `json:"refresh_token"`
}

// RefreshResponse to the client
type RefreshResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *RefreshRequest) Validate(body string) *RefreshResponse {
	response := new(RefreshResponse)
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

	if len(r.Token) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingRefreshToken)
		response.Code = es.StatusSignInError
	}

	return response
}

// Marshal the response object
func (r *RefreshResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
