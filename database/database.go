package database

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/bartekpacia/frycz-eats/webhooks"
	"google.golang.org/api/option"
)

var client *firestore.Client

// OrderData represents order that can be placed by a person.
type OrderData struct {
	Psid      string `json:"orderer_psid"`
	Order     string
	OrderedAt interface{} `json:"orderer_at"`
	Completed bool
}

func init() {
	sa := option.WithCredentialsFile("./key.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalf("Error initializing App: %v\n", err)
	}

	client, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Database: %v\n", err)
	}
}

func SaveToDatabase(wd webhooks.WebhookData, order string) (err error) {
	var orderData = OrderData{
		Psid:      wd.Sender.ID,
		Order:     order,
		OrderedAt: firestore.ServerTimestamp,
		Completed: false,
	}

	_, _, err = client.Collection("orders").Add(context.Background(), map[string]interface{}{
		"ordered_psid": orderData.Psid,
		"content":      orderData.Order,
		"ordered_at":   firestore.ServerTimestamp,
		"completed":    orderData.Completed})

	return err
}
