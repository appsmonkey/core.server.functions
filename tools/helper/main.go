package helper

import (
	"math/rand"
	"net/http"
)

// CognitoIDZeroValue is the zero value for a Cognito ID
const CognitoIDZeroValue = "none"

// IsCognitoIDEmpty indicates if the cognito ID is empty.
// returns `true` if empty
func IsCognitoIDEmpty(cid string) bool {
	if len(cid) == 0 || cid == CognitoIDZeroValue {
		return true
	}

	return false
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq generated random string with spec. length
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetLang translations
func GetLang(lang string) (interface{}, error) {

	switch lang {
	case "BA":
		return http.Get("https://cityos-universal-links.s3.amazonaws.com/lang/BA-lng.txt")
	default:
		return http.Get("https://cityos-universal-links.s3.amazonaws.com/lang/BA-lng.txt")
	}
}
