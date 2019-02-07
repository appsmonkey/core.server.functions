package viewmodels

import (
	"regexp"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
)

// BaseResponse structure that will be on all responses from all functions
type BaseResponse struct {
	Code      int64       `json:"Code"`
	Errors    []es.Error  `json:"Errors"`
	Data      interface{} `json:"Data"`
	RequestID string      `json:"RequestID"`
}

var regx = struct {
	Password *regexp.Regexp
	Email    *regexp.Regexp
}{
	Password: regexp.MustCompile(`^\S{8,20}$`),
	Email:    regexp.MustCompile(`^[a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`),
}

func validateEmail(email string) bool {
	return regx.Email.MatchString(email)
}

// ValidatePassword will check if password match to regexp
func validatePassword(password string) bool {
	return regx.Password.MatchString(password)
}
