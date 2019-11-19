package cognito

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	dal "github.com/appsmonkey/core.server.functions/dal/access"
	m "github.com/appsmonkey/core.server.functions/models"
	bg "github.com/kjk/betterguid"
	"google.golang.org/api/oauth2/v2"
)

// Google login
func (c *Cognito) Google(id, token, inEmail string, client *http.Client) (*CognitoData, error) {
	uti, err := verifyIDTokenGoogle(token, client)
	if err != nil {
		return nil, err
	}

	if uti.UserId != id || uti.Email != inEmail {
		fmt.Println("invalid token received", "Received ID", id, "Received email", inEmail, "Got ID", uti.UserId, "Got Email", uti.Email)
		return nil, errors.New("invalid token received")
	}

	_, _, _, _, _, err = dal.CheckSocial(id)
	if err != nil {
		// Insert temp data into DB
		// err := dal.AddTempUser(inEmail, id, "G")
		// if err != nil {
		// 	fmt.Println("Could not save a temp user", err)
		// 	return nil, err
		// }

		// Register
		cd, err := c.SignUp(inEmail, "@aA"+id, "male", "fn", "ln")
		if err != nil {
			fmt.Println("Could not register at cognito", err)
			return nil, err
		}

		// Get userdata
		p, err := c.Profile(inEmail)
		if err != nil {
			fmt.Println("Could not get user's profile from cognito", err)
			return nil, err
		}

		r := &m.User{}
		r.Attributes = make(map[string]string, 0)
		r.Profile = m.UserProfile{}
		r.Token = bg.New()
		r.SocialID = id
		r.SocialType = "G"

		for _, uav := range p.UserAttributes {
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
			fmt.Println("Could not get user's profile from cognito [CognitoID]", err)
			return nil, err
		}

		if len(r.Email) == 0 {
			fmt.Println("Could not get user's profile from cognito [Email]", err)
			return nil, err
		}

		err = dal.AddUser(r)
		if err != nil {
			fmt.Println("Could not save user into cognito", err)
			return nil, err
		}

		return cd, nil
	}

	return c.SignIn(inEmail, "@aA"+id)
}

// Facebook login
func (c *Cognito) Facebook(id, token, inEmail string) (*CognitoData, error) {
	clientID := os.Getenv("FB_CLIENT_ID")
	clientSecret := os.Getenv("FB_CLIENT_SECRET")

	fmt.Println("FB_CLIENT:" + clientID)

	appLink := fmt.Sprintf(`https://graph.facebook.com/oauth/access_token?client_id=%v&client_secret=%v&grant_type=client_credentials`, clientID, clientSecret)

	appTokenData := request(appLink)
	fmt.Println("appTokenData")
	fmt.Println(appTokenData)
	appToken := appTokenData["access_token"].(string)

	link := fmt.Sprintf(`https://graph.facebook.com/debug_token?input_token=%v&access_token=%v`, token, appToken)
	userData := request(link)["data"].(map[string]interface{})
	user_id := userData["user_id"].(string)

	if user_id != id {
		fmt.Println("invalid token received", "Received ID", id, "Received email", inEmail, "Got ID", userData, "Got Email", "none")
		return nil, errors.New("invalid token received")
	}

	_, _, _, _, _, err := dal.CheckSocial(id)
	if err != nil {
		// Register
		cd, err := c.SignUp(inEmail, "@aA"+id, "male", "fn", "ln")
		if err != nil {
			fmt.Println("Could not register at cognito", err)
			return nil, err
		}

		// Get userdata
		p, err := c.Profile(inEmail)
		if err != nil {
			fmt.Println("Could not get user's profile from cognito", err)
			return nil, err
		}

		r := &m.User{}
		r.Attributes = make(map[string]string, 0)
		r.Profile = m.UserProfile{}
		r.Token = bg.New()
		r.SocialID = id
		r.SocialType = "G"

		for _, uav := range p.UserAttributes {
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
			fmt.Println("Could not get user's profile from cognito [CognitoID]", err)
			return nil, err
		}

		if len(r.Email) == 0 {
			fmt.Println("Could not get user's profile from cognito [Email]", err)
			return nil, err
		}

		err = dal.AddUser(r)
		if err != nil {
			fmt.Println("Could not save user into cognito", err)
			return nil, err
		}

		return cd, nil
	}

	return c.SignIn(inEmail, "@aA"+id)
}

func verifyIDTokenGoogle(idToken string, client *http.Client) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(client)
	if err != nil {
		fmt.Println("verifyIDTokenGoogle New oauth2 client", err.Error())
		return nil, err
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		fmt.Println("verifyIDTokenGoogle Token Error", err.Error())
		return nil, err
	}

	m, err := tokenInfo.MarshalJSON()
	if err != nil {
		fmt.Println("verifyIDTokenGoogle Token Unmarshal Error", err.Error())
		return nil, err
	}

	fmt.Println("verifyIDTokenGoogle Result", string(m))

	return tokenInfo, nil
}

func request(url string) map[string]interface{} {
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not send data to FB [client.Do]", err)
		return make(map[string]interface{}, 0)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Could not send data to FB [json.Ubnmarshal]", err)
		return make(map[string]interface{}, 0)
	}

	return res
}
