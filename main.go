package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	host string
	port string
	verifyToken string
	accessToken string
)

func init()  {
	flag.StringVar(&host, "host", "", "hostname to listen on")
	flag.StringVar(&port, "port", "5000", "port to listen on")

	godotenv.Load()
	var ok bool
	verifyToken, ok = os.LookupEnv("verify_token")
	if !ok {
		log.Fatalln("Error reading verify_token")
	}
	accessToken, ok = os.LookupEnv("access_token")
	if !ok {
		log.Fatalln("Error reading access_token")
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
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Printf("GET /webhook")

		mode := r.URL.Query().Get("hub.mode")
		token := r.URL.Query().Get("hub.verify_token")
		challenge := r.URL.Query().Get("hub.challenge")

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

	if r.Method == "POST" {
		log.Printf("POST /webhook")
		
		w.WriteHeader(http.StatusNotImplemented)
	}
}
