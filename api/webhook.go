package api

// WholeBody represents the body of a new messenger message (?)
type WholeBody struct {
	Object  string
	Entries []WebhookEvent `json:"entry"`
}

type WebhookEvent struct {
	ID        string
	Time      int64
	Messaging []WebhookData
}

type WebhookData struct {
	Sender    Person
	Recipient Person
	Timestamp int64
	Message   *Message
	Postback  *Postback
}

type Person struct {
	ID string
}

type Message struct {
	Mid  string
	Text string
}

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