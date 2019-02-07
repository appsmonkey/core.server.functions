package errorstatuses

// Error object that will contain details of the error
type Error struct {
	Code    int64  `json:"error-code"`
	Message string `json:"error-message"`
	Data    string `json:"error-data"`
}

// ErrUnknown [could not determine the error]
var ErrUnknown = Error{Code: 1001, Message: "unknown error"}

// ErrRegistrationMissingName [cannot register user. Missing full name]
var ErrRegistrationMissingName = Error{Code: 1002, Message: "registration error [missing full user's name]"}

// ErrRegistrationMissingEmail [cannot register user. Missing email]
var ErrRegistrationMissingEmail = Error{Code: 1003, Message: "registration error [missing email]"}

// ErrRegistrationMissingPass [cannot register user. Missing password]
var ErrRegistrationMissingPass = Error{Code: 1003, Message: "registration error [missing password]"}

// ErrRegistrationIncorrectRequest [cannot register user. incorrect request]
var ErrRegistrationIncorrectRequest = Error{Code: 1004, Message: "registration error [could not read request]"}
