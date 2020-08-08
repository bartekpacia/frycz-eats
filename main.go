package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

		mode := req.URL.Query().Get("hub.mode")
		token := req.URL.Query().Get("hub.verify_token")
		challenge := req.URL.Query().Get("hub.challenge")

		fmt.Printf("mode: %s, token: %s, challenge: %s\n", mode, token, challenge)
		if mode != "" && token != "" {
			if mode == "subscribe" && token == verifyToken {
				log.Println("Webhook verified successfully. 200")
				fmt.Fprintf(w, challenge)
			} else {
				log.Println("Webhook verification failed. 403")
				w.WriteHeader(http.StatusForbidden)
			}
		} else {
			log.Println("Webhook verification failed. 400")
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	if req.Method == "POST" {
		log.Println("POST /webhook")

		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("Error reading request body", err)
		}

		var body WholeBody
		if err = json.Unmarshal(bodyBytes, &body); err != nil {
			fmt.Println("Error parsing json", err)
		}
		fmt.Printf("text: %s\n", body.Entries[0].Messaging[0].Message.Text)

		dst := bytes.Buffer{}
		if err := json.Indent(&dst, bodyBytes, "", "  "); err != nil {
			panic(err)
		}

		s, _ := json.MarshalIndent(body, "", "    ")
		fmt.Printf("body: %s", string(s))

		w.WriteHeader(http.StatusOK)
	}
}

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
	Referral Referral
}

type Referral struct {
	Ref    string
	Source string
	Type   string
}
