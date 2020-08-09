package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/bartekpacia/frycz-eats/webhooks"
)

const facebookURL string = "https://graph.facebook.com/v2.6/me/messages"

type RequestBody struct {
	Recipient *webhooks.Person
	Message   *Message          
}

type Message struct {
	Text string 
}

// User represents a Messenger user obtained through the Facebook's Graph API
type User struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
}

// SendMessage sends a text message on Messenger to recipientPsid
func SendMessage(recipientPsid string, text string, accessToken string) error {
	recipient := webhooks.Person{ID: recipientPsid}
	message := Message{Text: text}
	body := RequestBody{&recipient, &message}

	bodyJSON, err := json.MarshalIndent(body, "", "    ")
	if err != nil {
		return err
	}
	fmt.Printf("body: %v\n", body)
	fmt.Printf("bodyJSON: %s\n", string(bodyJSON))

	fullURL, err := url.Parse(facebookURL)
	if err != nil {
		return err
	}

	qs := fullURL.Query()
	qs.Set("access_token", accessToken)
	fullURL.RawQuery = qs.Encode()
	// fmt.Printf("fullURL: %s\n", fullURL.String())

	response, err := http.Post(fullURL.String(), "application/json", bytes.NewBuffer(bodyJSON))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		errorText := fmt.Sprintf("failed to send a message. status=%d\n", response.StatusCode)
		err = errors.New(errorText)
	}

	return err
}

// GetUser obtains the user of PSID from the Facebook Graph API
func GetUser(PSID string, accessToken string) (*User, error) {
	baseURL := fmt.Sprint("https://graph.facebook.com/", PSID)
	fullURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	qs := fullURL.Query()
	qs.Set("access_token", accessToken)
	qs.Set("fields", "first_name,last_name,profile_pic")
	fullURL.RawQuery = qs.Encode()
	fmt.Printf("fullURL: %s\n", fullURL.String())

	res, err := http.Get(fullURL.String())
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		return nil, err
	}

	fmt.Print("User.FirstName:", user.FirstName)

	return &user, nil
}
