package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Port is the server's address
const Port = 5000
// VerifyToken is the token used by Facebook to authorize the webhook
var VerifyToken string
// AccessToken is ...idk lol
var AccessToken string


func main() {
	godotenv.Load()
	var env string
	flag.StringVar(&env, "env", "dev", "env: dev | prod")
	flag.Parse()
	fmt.Println("env:", env)

	addr := ":" + strconv.Itoa(Port)
	if env == "dev" {
		addr = "localhost:" + strconv.Itoa(Port)
	}

	var ok bool
	VerifyToken, ok = os.LookupEnv("verify_token")
	if !ok {
		log.Fatalln("Error reading verify_token")
	}
	AccessToken, ok = os.LookupEnv("access_token")
	if !ok {
		log.Fatalln("Error reading access_token")
	}

	fmt.Println("access_token:", AccessToken)

	http.HandleFunc("/", HandleMain)
	http.HandleFunc("/webhook", HandleWebhook)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// HandleMain handles "/" route.
func HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Frycz Eats API homepage. Hi!")
}

// HandleWebhook handles "/webhook" route.
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		mode := r.URL.Query().Get("hub.mode")
		token := r.URL.Query().Get("hub.verify_token")
		challenge := r.URL.Query().Get("hub.challenge")

		fmt.Printf("mode: %s, token: %s, challenge: %s\n", mode, token, challenge)
		if mode != "" && token != "" {
			if mode == "subscribe" && token == VerifyToken {
				log.Println("Webhook successfully verified. 200")
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

	if r.Method == "POST" {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
