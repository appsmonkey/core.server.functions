package cognito

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/lestrrat/go-jwx/jwk"
)

type Cognito struct {
	identityProvider *cognitoidentityprovider.CognitoIdentityProvider
}

type CognitoData struct {
	IDToken      string                                         `json:"id_token"`
	AccessToken  string                                         `json:"access_token"`
	ExpiresIn    int64                                          `json:"expires_in"`
	RefreshToken string                                         `json:"refresh_token,omitempty"`
	UserData     *cognitoidentityprovider.AdminCreateUserOutput `json:"-"`
}

const (
	authFlow = "ADMIN_NO_SRP_AUTH"
)

var (
	region     string
	userPoolID string
	clientID   string
	jwkURL     string
	keySet     *jwk.Set
)

func initialize() {
	region = os.Getenv("COGNITO_REGION")
	userPoolID = os.Getenv("COGNITO_USER_POOL_ID")
	clientID = os.Getenv("COGNITO_CLIENT_ID")
	jwkURL = fmt.Sprintf("https://cognito-idp.%v.amazonaws.com/%v/.well-known/jwks.json", region, userPoolID)

	if err := loadKeySet(); err != nil {
		writeLog("LoadKeySet Error: ", err)
		return
	}
}

// NewCognito creates new instance of cognito and initiates cognito session
func NewCognito() *Cognito {
	initialize()

	c := &Cognito{}

	cred := credentials.NewEnvCredentials()
	s := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: cred,
	}))

	c.identityProvider = cognitoidentityprovider.New(s)
	return c
}

// SignUp register new user
func (c *Cognito) SignUp(username, password, gender, firstname, lastname string) (*CognitoData, error) {
	// Step 1
	adminUserData, err := c.identityProvider.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		Username:          aws.String(username),
		TemporaryPassword: aws.String(password),
		UserPoolId:        aws.String(userPoolID),
		MessageAction:     aws.String(cognitoidentityprovider.MessageActionTypeSuppress),

		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(username),
			},
			{
				Name:  aws.String("gender"),
				Value: aws.String(gender),
			},
			// {
			// 	Name:  aws.String("first_name"),
			// 	Value: aws.String(firstname),
			// },
			// {
			// 	Name:  aws.String("last_name"),
			// 	Value: aws.String(lastname),
			// },
		},
	})

	if err != nil {
		writeLog("AdminCreateUser Error: ", err)
		return nil, err
	}

	// Step 2
	// Attemp login to get session value, which is used to confirm user
	aia := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(authFlow),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		ClientId:   aws.String(clientID),
		UserPoolId: aws.String(userPoolID),
	}

	authresp, autherr := c.identityProvider.AdminInitiateAuth(aia)
	if autherr != nil {
		writeLog("AdminInitiateAuth Error:", autherr)
		return nil, autherr
	}

	// Step 3
	// Set user to confirmed
	artaci := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
		ClientId:      aws.String(clientID),
		UserPoolId:    aws.String(userPoolID),
		ChallengeResponses: map[string]*string{
			"USERNAME":     aws.String(username),
			"NEW_PASSWORD": aws.String(password),
		},
		Session: authresp.Session,
	}

	challangeResponse, err := c.identityProvider.AdminRespondToAuthChallenge(artaci)
	if err != nil {
		writeLog("AdminRespondToAuthChallenge Error:", err)
		return nil, nil
	}

	data := new(CognitoData)
	data.IDToken = aws.StringValue(challangeResponse.AuthenticationResult.IdToken)
	data.AccessToken = aws.StringValue(challangeResponse.AuthenticationResult.AccessToken)
	data.ExpiresIn = aws.Int64Value(challangeResponse.AuthenticationResult.ExpiresIn)
	data.RefreshToken = aws.StringValue(challangeResponse.AuthenticationResult.RefreshToken)
	data.UserData = adminUserData

	return data, nil
}

// Refresh user's tokens based on the provided refresh token
func (c *Cognito) Refresh(token string) (*CognitoData, error) {
	aia := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(token),
		},
		ClientId:   aws.String(clientID),
		UserPoolId: aws.String(userPoolID),
	}
	authresp, autherr := c.identityProvider.AdminInitiateAuth(aia)
	if autherr != nil {
		writeLog("AdminInitiateAuth Error:", autherr)
		return nil, autherr
	}

	data := new(CognitoData)
	data.IDToken = aws.StringValue(authresp.AuthenticationResult.IdToken)
	data.AccessToken = aws.StringValue(authresp.AuthenticationResult.AccessToken)
	data.ExpiresIn = aws.Int64Value(authresp.AuthenticationResult.ExpiresIn)
	data.RefreshToken = aws.StringValue(authresp.AuthenticationResult.RefreshToken)

	return data, nil
}

// Profile returns user's profile based on username
func (c *Cognito) Profile(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	input := new(cognitoidentityprovider.AdminGetUserInput)
	input.UserPoolId = aws.String(userPoolID)
	input.Username = aws.String(username)

	output, err := c.identityProvider.AdminGetUser(input)
	return output, err
}

// SignIn login user based on his username and password
func (c *Cognito) SignIn(username, password string) (*CognitoData, error) {
	authInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(authFlow),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		ClientId:   aws.String(clientID),
		UserPoolId: aws.String(userPoolID),
	}

	authresp, err := c.identityProvider.AdminInitiateAuth(authInput)

	if err != nil {
		writeLog("AdminInitiateAuth Error:", err)
		return nil, err
	}

	data := new(CognitoData)
	data.IDToken = aws.StringValue(authresp.AuthenticationResult.IdToken)
	data.AccessToken = aws.StringValue(authresp.AuthenticationResult.AccessToken)
	data.ExpiresIn = aws.Int64Value(authresp.AuthenticationResult.ExpiresIn)
	data.RefreshToken = aws.StringValue(authresp.AuthenticationResult.RefreshToken)

	return data, nil
}

// ValidateToken checks authorization token
func (c *Cognito) ValidateToken(jwtToken string) (string, string, error) {
	token, err := jwt.Parse(jwtToken, c.getKey)
	if err != nil {
		return "", "", fmt.Errorf("could not parse jwt: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["token_use"] != "access" {
			return "", "", fmt.Errorf("token_use mismatch: %s", claims["token_use"])
		}

		return claims["sub"].(string), claims["username"].(string), nil // valid token
	}
	return "", "", nil // invalid token
}

// getKey returns the key for validating in ValidateToken
func (c *Cognito) getKey(token *jwt.Token) (interface{}, error) {
	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("expecting JWT to have string kid")
	}

	if key := keySet.LookupKeyID(keyID); len(key) == 1 {
		return key[0].Materialize()
	}

	return nil, fmt.Errorf("unable to find key")
}

func writeLog(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err.Error())
	} else {
		fmt.Println(msg)
	}
}

func loadKeySet() error {
	var err error
	keySet, err = jwk.FetchHTTP(jwkURL)
	return err
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading %v\n", err)
	}
}
