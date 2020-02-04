package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"
)

// LangBundleRequest sent from the client
type LangBundleRequest struct {
	Language string `json:"lang"`
}

// LangBundleResponse to the client
type LangBundleResponse struct {
	BaseResponse
}

// Validate the request sent from client
func (r *LangBundleRequest) Validate(body map[string]string) *LangBundleResponse {
	response := new(LangBundleResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	lang, ok := body["lang"]
	if ok {
		r.Language = lang
	} else {
		r.Language = "BA"
	}

	return response
}

// Marshal the response object
func (r *LangBundleResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
