package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bartekpacia/frycz-eats/messenger"
	"io/ioutil"
	"net/http"
	"net/url"
)

const GraphApiURL string = "https://graph.facebook.com/v2.6/me/messages"

type RequestBody struct {
	Sender    *messenger.Person   `json:"sender,omitempty"`
	Recipient *messenger.Person   `json:"recipient,omitempty"`
	Timestamp *int64              `json:"timestamp,omitempty"`
	Message   *messenger.Message  `json:"message,omitempty"`
	Postback  *messenger.Postback `json:"postback,omitempty"`
}

// User represents a Messenger user obtained through the Facebook's Graph API
type User struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	ProfilePicURL string `json:"profile_pic"`
}

// Send calls SendAPI to send a body
func Send(body RequestBody, accessToken string) error {
	bodyJSON, err := json.MarshalIndent(body, "", "    ")
	if err != nil {
		return err
	}

	fullURL, err := url.Parse(GraphApiURL)
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
