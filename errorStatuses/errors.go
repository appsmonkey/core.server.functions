package errorstatuses

// Error object that will contain details of the error
type Error struct {
	Code    int64  `json:"error-code"`
	Message string `json:"error-message"`
	Data    string `json:"error-data"`
}

// ErrUnknown [could not determine the error]
var ErrUnknown = Error{Code: 1001, Message: "unknown error"}

// ErrRegistrationMissingFirstName [cannot register user. Missing first name]
var ErrRegistrationMissingFirstName = Error{Code: 1002, Message: "registration error [missing user's first name]"}

// ErrRegistrationMissingLastName [cannot register user. Missing last name]
var ErrRegistrationMissingLastName = Error{Code: 1003, Message: "registration error [missing user's last name]"}

// ErrRegistrationInvalidGender [cannot register user. Invalid gender value]
var ErrRegistrationInvalidGender = Error{Code: 1004, Message: "registration error [invalid gender value]"}

// ErrRegistrationMissingEmail [cannot register user. Missing email]
var ErrRegistrationMissingEmail = Error{Code: 1005, Message: "registration error [missing email]"}

// ErrRegistrationMissingPass [cannot register user. Missing password]
var ErrRegistrationMissingPass = Error{Code: 1006, Message: "registration error [missing password]"}

// ErrRegistrationIncorrectRequest [cannot register user. incorrect request]
var ErrRegistrationIncorrectRequest = Error{Code: 1007, Message: "registration error [could not read request]"}

// ErrRegistrationCognitoSignupError [cannot register user. cognito signup error]
var ErrRegistrationCognitoSignupError = Error{Code: 1008, Message: "registration error [cannot register user. cognito signup error]"}

// ErrRegistrationSignInError signin error [cannot singin user. cognito signin error]
var ErrRegistrationSignInError = Error{Code: 1009, Message: "signin error [cannot singin user. cognito signin error]"}

// ErrProfileMissingEmail [cannot get user's profile. Missing email]
var ErrProfileMissingEmail = Error{Code: 1005, Message: "profile error [missing email]"}
