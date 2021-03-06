package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bartekpacia/frycz-eats/messenger"
	"github.com/bartekpacia/frycz-eats/webhooks"
	"github.com/joho/godotenv"
)

var (
	host        string
	port        string
	verifyToken string
	accessToken string
)

func init() {
	flag.StringVar(&host, "host", "", "hostname to listen on")
	flag.StringVar(&port, "port", "5000", "port to listen on")

	godotenv.Load()
	var ok bool
	verifyToken, ok = os.LookupEnv("verify_token")
	if !ok {
		log.Fatalln("verifyToken is empty")
	}
	accessToken, ok = os.LookupEnv("access_token")
	if !ok {
		log.Fatalln("accessToken is empty")
	}

}

func main() {
	flag.Parse()

	http.HandleFunc("/", HandleMain)
	http.HandleFunc("/webhook", HandleWebhook)

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// HandleMain handles "/" route.
func HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Frycz Eats API homepage. Hi!")
}

// HandleWebhook handles "/webhook" route.
func HandleWebhook(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		log.Printf("GET /webhook")
		webhooks.VerifyWebhook(w, req, verifyToken)
	}

	if req.Method == "POST" {
		log.Println("POST /webhook")

		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("Error reading request body", err)
		}

		entries, err := webhooks.UnmarshallEntries(bodyBytes)
		if err != nil {
			log.Fatalf("Error unmarshalling webhookEvents: %e\n", err)
		}

		for _, event := range entries {
			webhookData := event.WebhookData[0]

			if webhookData.Message != nil {
				// responseText, err := webhookData.HandleMessage(accessToken)
				// if err != nil {
				// 	fmt.Println("Error handling message:", err)
				// }

				// user, err := messenger.GetUser(webhookData.Sender.ID, accessToken)
				// if err != nil {
				// 	fmt.Println("Error getting user:", err)
				// }
				
				// msgText := fmt.Sprintf("Cześć, %s %s! %s \n", user.FirstName, user.LastName, responseText)

				// err = messenger.SendMessage(webhookData.Sender.ID, msgText, accessToken)
				err = messenger.SendPostback(webhookData.Sender.ID, accessToken)
				if err != nil {
					fmt.Println("Error sending message:", err)
				}

				// err = database.SaveOrder(webhookData, webhookData.Message.Text)
				// if err != nil {
				// 	fmt.Println("Error saving order to firestore:", err)
				// }

			} else if webhookData.Postback != nil {
				msgText, err := webhookData.HandlePostback()
				if err != nil {
					fmt.Println("Error handling postback")
				}

				err = messenger.SendMessage(webhookData.Sender.ID, msgText, accessToken)
				if err != nil {
					fmt.Println("Error sending message:", err)
				}
			}

		}

		s, _ := json.MarshalIndent(entries, "", "    ")
		fmt.Printf("webhookEvents: %s\n", string(s))

		w.WriteHeader(http.StatusOK)
	}
}
