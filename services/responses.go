package services

import "github.com/bartekpacia/frycz-eats/messenger"

func MakeQuickReply(recipientPSID string, title string, payload interface{}) *RequestBody {
	recipient := messenger.Person{ID: recipientPSID}

	payload := messenger.Payload{}

	quickReply := messenger.QuickReply{
		ContentType: "text",
		Title:       title,
		Payload:     &payload,
		ImageURL:    "",
	}

	message := messenger.Message{QuickReplies: []*messenger.QuickReply{&quickReply}}
	body := RequestBody{Recipient: &recipient, Message: &message}
	return &body
}
