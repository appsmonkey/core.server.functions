package errorstatuses

// Error object that will contain details of the error
type Error struct {
	Code    int64  `json:"error-code"`
	Message string `json:"error-message"`
	Data    string `json:"error-data"`
}

// ErrNo [There is no error]
var ErrNo = Error{Code: 1000, Message: "no error"}

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

// ErrMissingCognitoID [cannot save user. Missing cognito user ID]
var ErrMissingCognitoID = Error{Code: 1010, Message: "profile error [missing cognito ID]"}

// ErrMissingThingModel [cannot add device. Missing model]
var ErrMissingThingModel = Error{Code: 1011, Message: "cannot add device. Missing model"}

// ErrMissingThingName [cannot add device. Missing name]
var ErrMissingThingName = Error{Code: 1012, Message: "cannot add device. Missing name"}

// ErrMissingThingLocation [cannot add device. Missing location]
var ErrMissingThingLocation = Error{Code: 1013, Message: "cannot add device. Missing location"}

// ErrMissingThingToken [cannot get device details. Missing token]
var ErrMissingThingToken = Error{Code: 1014, Message: "cannot get device details. Missing token"}

// ErrMissingSensorType [cannot get map details. Missing sensor]
var ErrMissingSensorType = Error{Code: 1015, Message: "cannot get map details. Missing sensor"}

// ErrMissingBio [cannot updateProfile details. Missing bio]
var ErrMissingBio = Error{Code: 1016, Message: "cannot update profile details. Missing bio"}

// ErrMissingFirstname [cannot updateProfile details. Missing first name]
var ErrMissingFirstname = Error{Code: 1017, Message: "cannot update profile details. Missing first name"}

// ErrMissingLastname [cannot updateProfile details. Missing last name]
var ErrMissingLastname = Error{Code: 1018, Message: "cannot update profile details. Missing last name"}

// ErrMissingMantra [cannot updateProfile details. Missing mantra]
var ErrMissingMantra = Error{Code: 1019, Message: "cannot update profile details. Missing mantra"}

// ErrMissingCity [cannot updateProfile details. Missing city]
var ErrMissingCity = Error{Code: 1020, Message: "cannot update profile details. Missing City"}

// ErrMissingWebsite [cannot updateProfile details. Missing website]
var ErrMissingWebsite = Error{Code: 1021, Message: "cannot update profile details. Missing website"}

// ErrMissingBirthday [cannot updateProfile details. Missing birthday]
var ErrMissingBirthday = Error{Code: 1022, Message: "cannot update profile details. Missing birthday"}

// ErrNotYours [cannot retrieve device. It does not belong to you]
var ErrNotYours = Error{Code: 1023, Message: "cannot retrieve device. It does not belong to you"}

// ErrMissingRefreshToken [could not refresh identity tokens. Missing refresh token]
var ErrMissingRefreshToken = Error{Code: 1024, Message: "could not refresh identity tokens. Missing refresh token"}

// ErrDeviceNotFound [could not find the desired device]
var ErrDeviceNotFound = Error{Code: 1025, Message: "could not find the desired device"}

// ErrSchemaNotFound [could not find the desired schena]
var ErrSchemaNotFound = Error{Code: 1026, Message: "could not find the desired schema"}

// ErrIncorrectRequest [could not understand the request]
var ErrIncorrectRequest = Error{Code: 1027, Message: "could not understand the request"}

// ErrIncorrectTime [could not understand the time]
var ErrIncorrectTime = Error{Code: 1028, Message: "could not understand the time"}
