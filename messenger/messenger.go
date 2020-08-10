package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const facebookURL string = "https://graph.facebook.com/v2.6/me/messages"

type requestBody struct {
	Recipient *Person `json:"recipient"`
	Message   *Message         `json:"message"`
}

// Message represents a textual message
type Message struct {
	Text        string  `json:"text"`
	Attachments []*Attachment  `json:"attachments,omitempty"`
	Mid         string `json:"mid,omitempty"`
}

type Attachment struct {
	Type    string `json:"type"`
	Payload *Payload `json:"payload,omitempty"`
}

type Payload struct {
	URL string `json:"url,omitempty"`
}

// Postback represents a postback, which is the action that occurs
// when some button in Messenger is tapped.
// https://developers.facebook.com/docs/messenger-platform/reference/webhook-events/messaging_postbacks
type Postback struct {
	Title    string `json:"title"`
	Payload  string `json:"payload"`
	Referral *Referral `json:"referral"`
}

type Referral struct {
	Ref    string `json:"ref"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

// Person represents the actual person using Messenger (usually sender or recipient)
type Person struct {
	ID string `json:"id"`
}

// User represents a Messenger user obtained through the Facebook's Graph API
type User struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
}


// SendMessage sends a text message on Messenger to recipientPsid
func SendMessage(recipientPsid string, text string, accessToken string) error {
	recipient := Person{ID: recipientPsid}
	message := Message{Text: text}
	body := requestBody{&recipient, &message}

	bodyJSON, err := json.MarshalIndent(body, "", "    ")
	if err != nil {
		return err
	}

	fullURL, err := url.Parse(facebookURL)
	if err != nil {
		return err
	}

	qs := fullURL.Query()
	qs.Set("access_token", accessToken)
	fullURL.RawQuery = qs.Encode()

	response, err := http.Post(fullURL.String(), "application/json", bytes.NewBuffer(bodyJSON))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		errMsg := fmt.Sprintln("failed to send a message. Response:", response)
		return errors.New(errMsg)
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
