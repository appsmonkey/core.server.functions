package models

// User data model
type User struct {
	Email      string            `json:"email"`
	CognitoID  string            `json:"cognito_id"`
	GroupID    string            `json:"group_id"`
	IsGroup    bool              `json:"is_group"`
	Token      string            `json:"token"`
	Profile    UserProfile       `json:"profile"`
	Attributes map[string]string `json:"attributes"`
	Devices    []string          `json:"devices"`
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
