package messenger

import "github.com/bartekpacia/frycz-eats/services"

const facebookURL string = "https://graph.facebook.com/v2.6/me/messages"

// Message represents a textual message
// https://developers.facebook.com/docs/messenger-platform/reference/webhook-events/messages
type Message struct {
	Text         string        `json:"text"`
	Attachments  []*Attachment `json:"attachments,omitempty"`
	Attachment   *Attachment   `json:"attachment,omitempty"`
	Mid          string        `json:"mid,omitempty"`
	QuickReplies []*QuickReply `json:"quick_replies,omitempty"`
}

type QuickReply struct {
	ContentType string   `json:"content_type"`
	Title       string   `json:"title"`
	Payload     *Payload `json:"payload"`
	ImageURL    string   `json:"image_url"`
}

type Attachment struct {
	Type    string   `json:"type"`
	Payload *Payload `json:"payload,omitempty"`
}

type Payload struct {
	URL          string     `json:"url,omitempty"`
	TemplateType string     `json:"template_type,omitempty"`
	Elements     []*Element `json:"elements,omitempty"` // Used only in responses
	Text         string     `json:"text"` // Used only in requests
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

// SendMessage sends a text message on Messenger to recipientPSID
func SendMessage(recipientPSID string, text string, accessToken string) error {
	recipient := Person{ID: recipientPSID}
	message := Message{Text: text}
	body := services.RequestBody{Recipient: &recipient, Message: &message}

	err := services.Send(body, accessToken)
	return err
}

// SendPostback sends a postback message on Messenger to recipientPSID
func SendPostback(recipientPSID string, accessToken string) error {
	buttonOrder := Button{Type: "postback", Title: "Złóż zamówienie", Payload: "submit_order"}
	buttonList := Button{Type: "postback", Title: "Pokaż listę zamówień", Payload: "show_list"}
	element := Element{Title: "Title", Subtitle: "Subtitle", Buttons: []*Button{&buttonOrder, &buttonList}}
	payload := Payload{TemplateType: "generic", Elements: []*Element{&element}}
	attachment := Attachment{Type: "template", Payload: &payload}
	message := Message{Attachment: &attachment}

	body := services.RequestBody{Recipient: &Person{ID: recipientPSID}, Message: &message}

	err := services.Send(body, accessToken)
	return err
}
