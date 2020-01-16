package viewmodels

import "encoding/json"

// VerifyRedirectRequest sent from the client
type VerifyRedirectRequest struct {
	RedirectURL      string `json:"redirect_url"`
	ClientID         string `json:"client_id"`
	UserName         string `json:"user_name"`
	ConfirmationCode string `json:"confirmation_code"`
	Type             string `json:"type"`
	CognitoID        string `json:"cog_id"`
}

// Validate the request sent from client
func (r *VerifyRedirectRequest) Validate(body map[string]string) *VerifyRedirectResponse {
	response := new(VerifyRedirectResponse)
	response.Code = 0

	url, ok := body["redirect_url"]
	if ok {
		r.RedirectURL = url
	}

	reqType, ok := body["type"]
	if ok {
		r.Type = reqType
	}

	cogID, ok := body["cog_id"]
	if ok {
		r.CognitoID = cogID
	}

	cid, ok := body["client_id"]
	if ok {
		r.ClientID = cid
	}

	userName, ok := body["user_name"]
	if ok {
		r.UserName = userName
	}

	cc, ok := body["confirmation_code"]
	if ok {
		r.ConfirmationCode = cc
	}

	return response
}

// VerifyRedirectResponse to the client
type VerifyRedirectResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *VerifyRedirectResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
