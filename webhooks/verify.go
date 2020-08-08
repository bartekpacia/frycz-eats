package webhooks

import (
	"fmt"
	"log"
	"net/http"
)

// VerifyWebhook performs initial authentication with Facebook
func VerifyWebhook(w http.ResponseWriter, req *http.Request, verifyToken string) {
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
