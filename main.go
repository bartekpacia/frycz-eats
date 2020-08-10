package main

import (
	"context"
	"flag"
	"fmt"
	messenger "github.com/bartekpacia/facebook-messenger"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	PhotoURL string = "https://firebasestorage.googleapis.com/v0/b/discoverrudy.appspot.com/o/static%2Frudy.jpg?alt=media"
)

var (
	host        string
	port        string
	verifyToken string
	accessToken string
	pageID      string
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

	pageID, ok = os.LookupEnv("page_id")
	if !ok {
		log.Fatalln("pageID is empty")
	}
}

func main() {
	flag.Parse()

	bot := &messenger.Messenger{
		AccessToken: accessToken,
		VerifyToken: verifyToken,
		PageID:      pageID,
	}

	bot.MessageReceived = messageReceived
	bot.PostbackReceived = postbackReceived
	// bot.DeliveryReceived = deliveryReceived // comment/delete if not used

	http.HandleFunc("/", handleMain)
	http.Handle("/bot", bot)

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// HandleMain handles "/" route.
func handleMain(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Frycz Eats API homepage. Hi!")
	if err != nil {
		log.Println("Error writing text to homepage ResponseWriter:", err)
	}
}

func messageReceived(ctx context.Context, msng *messenger.Messenger, userID int64, m messenger.FacebookMessage) {
	switch m.Text {
	case "hi":
		msng.SendTextMessage(userID, "Hello there")

	case "postback":
		gm := msng.NewGenericMessage(userID)
		btn2 := msng.NewPostbackButton("Ok", "postback")
		gm.AddNewElement("Title", "Subtitle", "https://google.pl", PhotoURL, []messenger.Button{btn2})

		msng.SendMessage(gm)
	case "button":
		gm := msng.NewButtonMessage(userID, "message")

		btnMakeOrder := msng.NewWebURLButton("title", "https://odkryjrudy.pl")
		btnShowOrders := msng.NewWebURLButton("TITL", "https://odkryjrudy.pl")
		gm.AddNewButton(btnMakeOrder)
		gm.AddNewButton(btnShowOrders)

		msng.SendMessage(gm)

	default:
		// upthere we haven't check for errors and responses for cleaner example code
		// but keep in mind that SendMessage returns FacebookResponse struct and error
		// errors are received from Facebook if sometnihg went wrong with message sending
		resp, err := msng.SendTextMessage(userID, m.Text) // echo, send back to user the same text he sent to us
		if err != nil {
			log.Println(err)
			return // if there is an error, resp is empty struct, useless
		}
		log.Println("Message ID", resp.MessageID, "sent to user", resp.RecipientID)
		// store resp.MessageID if you want to track delivery reports that will be sent later from Facebook
	}
}

func postbackReceived(ctx context.Context, msng *messenger.Messenger, userID int64, p messenger.FacebookPostback) {
	if p.Payload == "postback" {
		// user just clicked Ok button from previouse example, lets just send him a message
		msng.SendTextMessage(userID, "Ok, I'm always online, chat with me anytime :)")
	}
}
