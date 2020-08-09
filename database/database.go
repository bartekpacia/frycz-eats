package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/bartekpacia/frycz-eats/webhooks"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var client *firestore.Client

// OrderData represents order that can be placed by a person.
type OrderData struct {
	Psid      string      `json:"orderer_psid"`
	Order     string      `json:"order"`
	OrderedAt interface{} `json:"ordered_at"`
	Completed bool        `json:"completed"`
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

// SaveOrder saves the new order to database.
func SaveOrder(wd webhooks.WebhookData, order string) (err error) {
	var orderData = OrderData{
		Psid:      wd.Sender.ID,
		Order:     order,
		OrderedAt: firestore.ServerTimestamp,
		Completed: false,
	}

	_, _, err = client.Collection("orders").Add(context.Background(), map[string]interface{}{
		"orderer_psid": orderData.Psid,
		"order":        orderData.Order,
		"ordered_at":   firestore.ServerTimestamp,
		"completed":    orderData.Completed})

	return err
}

// GetRecentOrders returns "orderCount" recent orders.
func GetRecentOrders(orderCount int) (orders []OrderData, err error) {
	iter := client.Collection("orders").OrderBy(
		"ordered_at", firestore.Desc).Limit(orderCount).Documents(context.Background())

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		bajty, err := json.Marshal(doc.Data())
		if err != nil {
			return nil, err
		}

		var order OrderData
		err = json.Unmarshal(bajty, &order)
		if err != nil {
			return nil, err
		}

		fmt.Printf("order: %v\n", order)
	}

	return nil, nil
}

// AcceptOrder assign the order to the person who wish to complete it.
func AcceptOrder() {

}
