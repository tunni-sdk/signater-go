// Command webhook-server is a minimal HTTP server that receives and verifies
// Signater webhook deliveries.
//
// Usage:
//
//	export SIGNATER_WEBHOOK_SECRET=<Hookdeck signing secret>
//	go run ./examples/webhook-server -addr :8080
//
// For local debugging, forward deliveries with the Hookdeck CLI (see the
// "Webhooks CLI Debugging" section of the Signater docs):
//
//	hookdeck listen 8080
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tunni-sdk/signater-go/webhook"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	secret := os.Getenv("SIGNATER_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("SIGNATER_WEBHOOK_SECRET not set")
	}

	http.HandleFunc("POST /webhook", func(w http.ResponseWriter, r *http.Request) {
		payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		event, err := webhook.ConstructEvent(payload, r.Header, secret)
		if err != nil {
			log.Printf("rejected delivery %q: %v", webhook.RequestID(r.Header), err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Deliveries are at-least-once and unordered: dedupe by request id
		// and query the API for the envelope's current state before acting.
		log.Printf("event %q envelope=%q env=%q request-id=%q",
			event.Type, event.EnvelopeID, event.Env, webhook.RequestID(r.Header))

		switch event.Type {
		case webhook.EventEnvelopeSigned:
			log.Printf("envelope %q fully signed!", event.EnvelopeID)
		case webhook.EventEnvelopeRejectedBySigner:
			log.Printf("envelope %q rejected by signer %q", event.EnvelopeID, event.SignerID)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Printf("listening on %s", *addr)
	server := &http.Server{Addr: *addr, ReadHeaderTimeout: 10 * time.Second}
	log.Fatal(server.ListenAndServe())
}
