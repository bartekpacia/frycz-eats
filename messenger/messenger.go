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

type RequestBody struct {
	Sender    *Person   `json:"sender,omitempty"`
	Recipient *Person   `json:"recipient,omitempty"`
	Timestamp *int64    `json:"timestamp,omitempty"`
	Message   *Message  `json:"message,omitempty"`
	Postback  *Postback `json:"postback,omitempty"`
}

// Message represents a textual message
// https://developers.facebook.com/docs/messenger-platform/reference/webhook-events/messages
type Message struct {
	Text        string        `json:"text"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	Attachment  *Attachment   `json:"attachment,omitempty"`
	Mid         string        `json:"mid,omitempty"`
}

type Attachment struct {
	Type    string   `json:"type"`
	Payload *Payload `json:"payload,omitempty"`
}

type Payload struct {
	URL          string     `json:"url,omitempty"`
	TemplateType string     `json:"template_type,omitempty"`
	Elements     []*Element `json:"elements,omitempty"`
}

type Element struct {
	Title    string    `json:"title,omitempty"`
	Subtitle string    `json:"subtitle,omitempty"`
	ImageURL string    `json:"image_url,omitempty"`
	Buttons  []*Button `json:"buttons,omitempty"`
}

type Button struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Payload string `json:"payload,omitempty"`
}

// Postback represents a postback, which is the action that occurs
// when some button in Messenger is tapped.
// https://developers.facebook.com/docs/messenger-platform/reference/webhook-events/messaging_postbacks
type Postback struct {
	Title    string    `json:"title"`
	Payload  string    `json:"payload"`
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
	ProfilePicURL string `json:"profile_pic"`
}

// SendMessage sends a text message on Messenger to recipientPSID
func SendMessage(recipientPSID string, text string, accessToken string) error {
	recipient := Person{ID: recipientPSID}
	message := Message{Text: text}
	body := RequestBody{Recipient: &recipient, Message: &message}

	err := send(body, accessToken)
	return err
}

// SendPostback sends a postback message on Messenger to recipientPSID
func SendPostback(recipientPSID string, accessToken string) error {
	buttonYes := Button{Type: "postback", Title: "Yes!", Payload: "yes"}
	buttonNo := Button{Type: "postback", Title: "No!", Payload: "no"}
	element := Element{Title: "Title", Subtitle: "Subtitle", Buttons: []*Button{&buttonYes, &buttonNo}}
	payload := Payload{TemplateType: "generic", Elements: []*Element{&element}}
	attachment := Attachment{Type: "template", Payload: &payload}
	message := Message{ Attachment: &attachment}

	body := RequestBody{Recipient: &Person{ID: recipientPSID}, Message: &message}

	err := send(body, accessToken)
	return err
}

func send(body RequestBody, accessToken string) error {
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
		errMsg := fmt.Sprintln("failed to send. Response:", response)
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

	return &user, nil
}
