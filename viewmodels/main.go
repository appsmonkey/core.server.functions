package viewmodels

import (
	"regexp"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// BaseResponse structure that will be on all responses from all functions
type BaseResponse struct {
	Code      int64       `json:"code"`
	Errors    []es.Error  `json:"errors"`
	Data      interface{} `json:"data"`
	Groups    interface{} `json:"user_groups"`
	RequestID string      `json:"request_id"`
}

var regx = struct {
	Password *regexp.Regexp
	Email    *regexp.Regexp
}{
	Password: regexp.MustCompile(`^\S{8,20}$`),
	Email:    regexp.MustCompile(`^[a-zA-Z0-9._%-+]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`),
}

// ValidateEmail will check if email match to regexp
func validateEmail(email string) bool {
	return regx.Email.MatchString(email)
}

// ValidatePassword will check if password match to regexp
func validatePassword(password string) bool {
	return regx.Password.MatchString(password)
}

// ValidateGender will check if gender match to:
// - Cognito User Pool Attributes https://amzn.to/2DIhn1W and
// - OpenID Connect Specification https://bit.ly/2Dh1502
func validateGender(gender string) bool {
	return gender == "male" || gender == "female"
}

// Headers for cors
func (r *BaseResponse) Headers() map[string]string {
	hdr := make(map[string]string, 0)
	hdr["Access-Control-Allow-Origin"] = "*"
	hdr["Access-Control-Allow-Methods"] = "GET,OPTIONS,PUT,POST"
	hdr["Access-Control-Allow-Headers"] = "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,AccessToken,Accept-Version"

	return hdr
}
