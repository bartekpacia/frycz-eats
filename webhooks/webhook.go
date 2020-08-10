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

// WholeBody represents the body of a request that is sent to the webhook when it is hit
type wholeBody struct {
	Object  string     `json:"object"`
	Entries []*Entries `json:"entry"`
}

// WebhookData represents data that is sent to the webhook when it is hit
type WebhookData messenger.RequestBody

type Entries struct {
	ID          string        `json:"id"`
	Time        int64         `json:"time"`
	WebhookData []WebhookData `json:"messaging"`
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

func (wd WebhookData) HandlePostback() (message string, err error) {
	fmt.Println("New postback:", wd.Postback.Title)

	if wd.Postback.Payload == "yes" {
		return "You selected yes", nil
	} else if wd.Postback.Payload == "no" {
		return "You selected no", nil
	} else {
		return "", errors.New("user selected invalid response")
	}
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
