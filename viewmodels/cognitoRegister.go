package viewmodels

import (
	"encoding/json"
	"strconv"
	"time"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	c "github.com/appsmonkey/core.server.functions/integration/cognito"
	m "github.com/appsmonkey/core.server.functions/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	bg "github.com/kjk/betterguid"
)

// CognitoRegisterRequest sent from cognito
type CognitoRegisterRequest struct {
	m.User
}

// Validate the request sent from client
func (r *CognitoRegisterRequest) Validate(body *events.CognitoEventUserPoolsPostConfirmation) *CognitoRegisterResponse {
	response := new(CognitoRegisterResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	r.Attributes = make(map[string]string, 0)
	r.Profile = m.UserProfile{}
	r.Token = bg.New()
	r.SocialID = "none"
	r.SocialType = "none"

	for uak, uav := range body.Request.UserAttributes {
		switch uak {
		case "email":
			r.Email = uav
			r.Attributes[uak] = uav
		case "sub": // Unique Cognito User ID
			r.CognitoID = uav
			r.Attributes[uak] = uav
		default:
			r.Attributes[uak] = uav
		}
	}

	if len(r.CognitoID) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingCognitoID)
		response.Code = es.StatusRegistrationError
	}

	if len(r.Email) == 0 {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingEmail)
		response.Code = es.StatusRegistrationError
	}

	return response
}

// ValidateProfile the request sent from client
func (r *CognitoRegisterRequest) ValidateProfile(body *cognitoidentityprovider.AdminGetUserOutput) *CognitoRegisterResponse {
	response := new(CognitoRegisterResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	r.Attributes = make(map[string]string, 0)
	r.Profile = m.UserProfile{}
	r.Token = bg.New()
	r.SocialID = "none"
	r.SocialType = "none"

	for _, uav := range body.UserAttributes {
		switch *uav.Name {
		case "email":
			r.Email = *uav.Value
			r.Attributes[*uav.Name] = *uav.Value
		case "sub": // Unique Cognito User ID
			r.CognitoID = *uav.Value
			r.Attributes[*uav.Name] = *uav.Value
		default:
			r.Attributes[*uav.Name] = *uav.Value
		}
	}

	if len(r.CognitoID) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingCognitoID)
		response.Code = es.StatusRegistrationError
	}

	if len(r.Email) == 0 {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingEmail)
		response.Code = es.StatusRegistrationError
	}

	return response
}

// ValidateCognito the request sent from client
func (r *CognitoRegisterRequest) ValidateCognito(body *c.CognitoData) *CognitoRegisterResponse {
	response := new(CognitoRegisterResponse)
	response.Code = 0
	response.RequestID = strconv.FormatInt(time.Now().Unix(), 10)

	r.Attributes = make(map[string]string, 0)
	r.Profile = m.UserProfile{}
	r.Token = bg.New()

	for _, uav := range body.UserData.User.Attributes {
		if *uav.Name == "email" {
			r.Email = *uav.Value
			r.Attributes["email"] = r.Email
		} else if *uav.Name == "sub" {
			r.CognitoID = *uav.Value
			r.Attributes["sub"] = r.CognitoID
		} else {
			r.Attributes[*uav.Name] = *uav.Value
		}
	}

	if len(r.CognitoID) == 0 {
		response.Errors = append(response.Errors, es.ErrMissingCognitoID)
		response.Code = es.StatusRegistrationError
	}

	if len(r.Email) == 0 {
		response.Errors = append(response.Errors, es.ErrRegistrationMissingEmail)
		response.Code = es.StatusRegistrationError
	}

	return response
}

// CognitoRegisterResponse to the client
// `Returns a list of all devices assigned to the requestee. Data defained in the *DeviceAddData* struct`
type CognitoRegisterResponse struct {
	BaseResponse
}

// Marshal the response object
func (r *CognitoRegisterResponse) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}

// Marshal the response object
func (r *CognitoRegisterRequest) Marshal() string {
	res, _ := json.Marshal(r)

	return string(res)
}
