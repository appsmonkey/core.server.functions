package errorstatuses

// StatusUnknown [could not determine the cause of the issue]
var StatusUnknown int64 = 1

// StatusRegistrationError [could not register new user]
var StatusRegistrationError int64 = 2

// StatusSignInError [could not signin user]
var StatusSignInError int64 = 3

// StatusAddDeviceError [could not add a new device]
var StatusAddDeviceError int64 = 4

// StatusGetDeviceError [could not get device details]
var StatusGetDeviceError int64 = 5

// StatusMapError [could not get map details]
var StatusMapError int64 = 6

// StatusProfileUpdateError [could not update user's profile details]
var StatusProfileUpdateError int64 = 7
