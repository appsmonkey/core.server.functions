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

// ErrMissingChart [cannot get data count. Missing chart type]
var ErrMissingChart = Error{Code: 1029, Message: "cannot get data count. Missing chart type"}

// ErrMissingSocialData [Missing social data]
var ErrMissingSocialData = Error{Code: 1030, Message: "Missing social data"}

// ErrMissingCode [Missing code]
var ErrMissingCode = Error{Code: 1031, Message: "Missing code"}

// ErrCouldNotInitiateForgottenPasswordFlow [Forgot Psssword Flow Error]
var ErrCouldNotInitiateForgottenPasswordFlow = Error{Code: 1032, Message: "Forgot password flow error"}

// City errors

// ErrCityMissingName [Cannot add city, missing name]
var ErrCityMissingName = Error{Code: 1033, Message: "Cannot add city missing name"}

// ErrCityMissingCountry [Cannot add city, missing country]
var ErrCityMissingCountry = Error{Code: 1034, Message: "Cannot add city, missing country"}

// ErrCityMissingZones [Cannot add city, missing zones or invalid data]
var ErrCityMissingZones = Error{Code: 1035, Message: "cannot add city, missing zones or invalid data"}

// ErrMissingCityID [cannot get device details. Missing id]
var ErrMissingCityID = Error{Code: 1036, Message: "cannot complete action. Missing city id"}

// ErrCityNotFound [could not find the desired city]
var ErrCityNotFound = Error{Code: 1037, Message: "could not find the desired city"}

// Zone errors

// ErrMissingID [cannot get device details. Missing id]
var ErrMissingID = Error{Code: 1038, Message: "cannot get device details. Missing identifier"}

// ErrZoneNotFound [could not find the desired zone]
var ErrZoneNotFound = Error{Code: 1039, Message: "could not find the desired zone"}

// ErrCityAlreadyExists [Cannot add city, missing country]
var ErrCityAlreadyExists = Error{Code: 1040, Message: "Duplicate, city already exists"}

// VerificationFailed [User verification failed]
var VerificationFailed = Error{Code: 1041, Message: "User not confirmed, verification failed"}

// UserCreationFailedNoToken [User verification failed]
var UserCreationFailedNoToken = Error{Code: 1042, Message: "Registration not completed, no token supplied"}

// UserCreationFailedNoID [User verification failed]
var UserCreationFailedNoID = Error{Code: 1043, Message: "Registration not completed, cognito id not supplied"}

// UserCreationFailedNoUserName [User verification failed]
var UserCreationFailedNoUserName = Error{Code: 1044, Message: "Registration not completed, no username supplied"}

// UserCreationFailedNoPassword [User verification failed]
var UserCreationFailedNoPassword = Error{Code: 1045, Message: "Registration not completed, no password supplied"}

// ForgotPasswordBadRequestNoToken [Forgot password failed failed]
var ForgotPasswordBadRequestNoToken = Error{Code: 1046, Message: "Forgot password deinied, no token supplied"}

// ForgotPasswordBadRequestNoID [User verification failed]
var ForgotPasswordBadRequestNoID = Error{Code: 1047, Message: "Registration not completed, no id supplied"}

// ErrMissingThingCity [cannot update device. Missing name]
var ErrMissingThingCity = Error{Code: 1048, Message: "manage device failed. Missing city"}

// ErrCouldNotResetPasswordForUser [Forgot Psssword Flow Error]
var ErrCouldNotResetPasswordForUser = Error{Code: 1049, Message: "Can not reset password for social login"}

// ErrSocialUserAlreadyExists [SignIn flow error]
var ErrSocialUserAlreadyExists = Error{Code: 1050, Message: "This email is already linked with a social account"}

// ErrProfileMissingEmail [cannot get user's profile. Missing email]
var ErrProfileMissingEmail = Error{Code: 1051, Message: "profile error [missing email]"}

// ErrCouldNotFindUser [User not found Error]
var ErrCouldNotFindUser = Error{Code: 1052, Message: "Could not find user with this email"}

// ErrLoginSocialExistsError signin error [cannot singin user. cognito signin error]
var ErrLoginSocialExistsError = Error{Code: 1053, Message: "This email is already linked with a different social media type."}

// ErrThingTokenInvalid [cannot add device. Token validaiton failed]
var ErrThingTokenInvalid = Error{Code: 1054, Message: "cannot add device. Device token invalid"}
