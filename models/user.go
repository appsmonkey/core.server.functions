package models

// User data model
type User struct {
	Email      string            `json:"email"`
	CognitoID  string            `json:"cognito_id"`
	Token      string            `json:"token"`
	Attributes map[string]string `json:"attributes"`
	Profile    UserProfile       `json:"profile"`
}

// UserProfile meta data
type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Mantra    string `json:"mantra"`
	Bio       string `json:"bio"`
	City      string `json:"City"`
	Website   string `json:"website"`
	Birthday  int64  `json:"birthday"`
	Gender    string `json:"gender"`
}
