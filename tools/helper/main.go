package helper

import (
	"fmt"
	"math/rand"
	"strings"
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
func GetLang(lang string) string {

	switch lang {
	case "BA":
		return "https://cityos-universal-links.s3.amazonaws.com/lang/BA-lng.txt"
	default:
		return "https://cityos-universal-links.s3.amazonaws.com/lang/BA-lng.txt"
	}
}

// TransformCityString replaces some unwanted chars with altrenatives
func TransformCityString(s string) string {
	fmt.Printf("Original: %s\n", s)
	cs := strings.Map(normalize, s)
	fmt.Printf("Cleaned: %s\n", cs)

	return cs
}

func normalize(in rune) rune {
	switch in {
	case 'ć':
		return 'c'
	case 'č':
		return 'c'
	case 'ž':
		return 'z'
	case 'š':
		return 's'
	case 'đ':
		return 'd'
	}
	return in
}

// MapCity maps city param from clients to match out needs
func MapCity(in string) string {
	fmt.Printf("Original: %s\n", in)

	switch in {
	case "Novo Sarajevo":
		return "Sarajevo"
	case "Vogosca":
		return "Sarajevo"
	case "Ilidza":
		return "Sarajevo"
	case "Novi Grad":
		return "Sarajevo"
	case "Hadzici":
		return "Sarajevo"
	case "Ilijas":
		return "Sarajevo"
	case "Trnovo":
		return "Sarajevo"
	case "Luzani":
		return "Sarajevo"
	case "Kanton Sarajevo":
		return "Sarajevo"
	case "Canton Sarajevo":
		return "Sarajevo"
	case "Sarajevo Canton":
		return "Sarajevo"
	}

	return in
}
