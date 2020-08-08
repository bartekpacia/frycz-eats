package webhooks

import (
	"encoding/json"
	"fmt"
)

/*
Package api provides types for messenger-related stuff, such as webhooks.
Read more here: https://developers.facebook.com/docs/messenger-platform/reference/webhook-events/
*/

// WholeBody represents the body of a new messenger message (?)
type wholeBody struct {
	Object  string
	Entries []*Entries `json:"entry"`
}

type Entries struct {
	ID        string
	Time      int64
	WebhookData []WebhookData `json:"messaging"`
}

type WebhookData struct {
	Sender    Person
	Recipient Person
	Timestamp int64
	Message   *Message
	Postback  *Postback
}

func (wd WebhookData) HandleMessage() {
	fmt.Printf("Handling message saying: %s\n", wd.Message.Text)
}

func (wd WebhookData) HandlePostback() {
	fmt.Printf("Handling postback of title: %s\n", wd.Postback.Title)
}

// Person represents the actual person (usually sender or recipient)
type Person struct {
	ID string
}

// Message represents a textual message
type Message struct {
	Mid  string
	Text string
	Attachments []*Attachment
}

type Attachment struct {
	Type string
	Payload *Payload
}

type Payload struct {
	Url string
}

// Postback represents a postback, which is the action that occurs
// when some button in Messenger is tapped.
// https://developers.facebook.com/docs/messenger-platform/reference/webhook-events/messaging_postbacks
type Postback struct {
	Title    string
	Payload  string
	Referral *Referral
}

type Referral struct {
	Ref    string
	Source string
	Type   string
}

func UnmarshallEntries(data []byte) ([]*Entries, error) {
	var wholeBody wholeBody
	err := json.Unmarshal(data, &wholeBody)
	if err != nil {
		return nil, err
	}

	return wholeBody.Entries, err
}