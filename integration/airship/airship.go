package airship

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var username = os.Getenv("USER_NAME")
var password = os.Getenv("PASSWORD")

// ChanelType represents a channel in UA
type ChanelType int8

const (
	// Good level (everything below Sensitive)
	Good ChanelType = 0

	// Sensitive chanel
	Sensitive ChanelType = 1

	// Unhealthy chanel
	Unhealthy ChanelType = 2

	// VeryUnhealthy chanel
	VeryUnhealthy ChanelType = 3

	// Hazardous chanel
	Hazardous ChanelType = 4
)

// Airship main object
type Airship struct {
	auth        string
	accept      string
	contentType string
	url         string
	push        string
}

// New Airhip objecct
func New() (as *Airship) {
	as = new(Airship)
	as.accept = "application/vnd.urbanairship+json; version=3"
	as.contentType = "application/json"
	as.auth = fmt.Sprintf("Basic %s", basicAuth())
	as.url = "https://go.urbanairship.com"
	as.push = "/api/push"

	return
}

// Send the message
func (as *Airship) Send(value string, lvl ChanelType) {
	msgs := newMessage(value, lvl)

	for _, msg := range msgs {
		sendMessage(msg, as)
	}
}

func sendMessage(msg *Message, as *Airship) error {
	toBeSent, err := msg.Marshal()
	if err != nil {
		fmt.Println("Error sending to Airship [could not unmarshal message]", err)
		return err
	}

	req, err := http.NewRequest("POST", as.url+as.push, bytes.NewBuffer(toBeSent))
	req.Header.Set("Authorization", as.auth)
	req.Header.Set("Accept", as.accept)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not send data to Aitship [client.Do]", err)
		return err
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}

func basicAuth() string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
