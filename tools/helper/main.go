package helper

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
