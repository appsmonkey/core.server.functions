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

// StatusValidateEmailError [could not validate provided email]
var StatusValidateEmailError int64 = 8

// StatusDeleteDeviceError [could not delete the specified device]
var StatusDeleteDeviceError int64 = 9

// StatusChartLiveDeviceError [could not get chart live data for device]
var StatusChartLiveDeviceError int64 = 10

// StatusChartHourDeviceError [could not get chart hour data for device]
var StatusChartHourDeviceError int64 = 11

// StatusChartAllDeviceError [could not get chart all data]
var StatusChartAllDeviceError int64 = 12

// StatusChartHasDataError [could not determined if we have data for a specific chart]
var StatusChartHasDataError int64 = 13

// StatusChartLiveAllError [could not get avg chart live data]
var StatusChartLiveAllError int64 = 14

// StatusForgotPasswordError [could not reset passowrd]
var StatusForgotPasswordError int64 = 15

// City errors

// StatusCityMissingName [Cannot add city, missing name]
var StatusCityMissingName int64 = 16

// StatusCityMissingCountry [Cannot add city, missing country]
var StatusCityMissingCountry int64 = 17

// StatusCityMissingZones [Cannot add city, missing zones or invalid data]
var StatusCityMissingZones int64 = 18

// StatusAddCityError [Could not add new city]
var StatusAddCityError int64 = 19

// StatusDeleteCityError [could not delete the specified city]
var StatusDeleteCityError int64 = 20

// Zone errors

// StatusZoneImportError [could not import zones, invalid data]
var StatusZoneImportError int64 = 21

// StatusDeleteZoneError [could not import zones, invalid data]
var StatusDeleteZoneError int64 = 22
