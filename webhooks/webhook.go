package webhooks

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bartekpacia/frycz-eats/messenger"
)

/*
Package webhooks provides types and functions for webhooks.
Read more here: https://developers.facebook.com/docs/messenger-platform/reference/webhook-events/
*/

// WholeBody represents the body of a new messenger message (?)
type wholeBody struct {
	Object  string
	Entries []*Entries `json:"entry"`
}

type Entries struct {
	ID          string
	Time        int64
	WebhookData []WebhookData `json:"messaging"`
}

// WebhookData is data that is sent to us when a Messenger /POST webhoook is activated
type WebhookData struct {
	Sender    *messenger.Person
	Recipient *messenger.Person
	Timestamp int64
	Message   *messenger.Message
	Postback  *messenger.Postback
}

func (wd WebhookData) HandleMessage(accessToken string) (responseText string, err error) {
	fmt.Println("New message:", wd.Message.Text)

	if wd.Message.Attachments != nil {
		responseText = "Po co wysyłasz nam zdjęcia? Przestań plz."
		return responseText, nil
	}

	if wd.Message.Text == "" {
		err = errors.New("message text is empty")
		return "", err
	}

	return "Twoje zamówienie zostało zapisane!", nil
}

func (wd WebhookData) HandlePostback() {
	fmt.Println("New postback:", wd.Postback.Title)
}

// UnmarshallEntries parses response body and returns the list of entries
func UnmarshallEntries(data []byte) ([]*Entries, error) {
	var wholeBody wholeBody
	err := json.Unmarshal(data, &wholeBody)
	if err != nil {
		return nil, err
	}

	return wholeBody.Entries, err
}
