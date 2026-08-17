// Package signater is a Go SDK for the Signater e-signature API
// (https://docs.api.signater.com), covering all published endpoints:
// contacts, documents, envelopes, templates and vaults, plus webhook
// verification in the companion package
// github.com/tunni-sdk/signater-go/webhook.
//
// # Authentication
//
// Create a token in the API Tokens section of the Signater App and pass it
// with WithAPIToken, or set the SIGNATER_API_TOKEN environment variable:
//
//	client := signater.NewClient() // reads SIGNATER_API_TOKEN
//
// Sandbox needs no configuration: sandbox tokens automatically route to
// sandbox resources.
//
// # A complete signing flow
//
//	doc, err := client.Documents.Upload(ctx, "contract.pdf", file)
//	envID, err := client.Envelopes.Create(ctx, &signater.EnvelopeParams{
//		VaultID:                 vaultID,
//		Name:                    "Service Agreement",
//		JurisdictionCountryCode: "BR",
//		Language:                signater.LanguagePtBR,
//		Signers: []signater.EnvelopeSignerParams{{
//			Name:                      "Jane Doe",
//			Email:                     signater.String("jane@example.com"),
//			EmailCommunicationMode:    signater.CommunicationModeFull,
//			SmsCommunicationMode:      signater.CommunicationModeNone,
//			WhatsAppCommunicationMode: signater.CommunicationModeNone,
//		}},
//		Documents: []signater.EnvelopeDocumentParams{{ID: doc.ID, Name: "contract.pdf"}},
//	})
//	err = client.Envelopes.Publish(ctx, envID)
//
// # Errors
//
// Non-2xx responses map to typed errors carrying the response's telemetry
// operation id (quote it to Signater support). Match them with errors.As:
//
//	var rateErr *signater.RateLimitError
//	var payErr *signater.PaymentRequiredError
//	var apiErr *signater.APIError // matches any API error
//
// # Pagination
//
// List endpoints return one Page at a time; the AutoPaging variants return a
// Go 1.23 iterator that fetches pages transparently:
//
//	for contact, err := range client.Contacts.ListAutoPaging(ctx, nil) {
//		if err != nil { return err }
//		...
//	}
//
// # Retries
//
// Failed requests are retried up to WithMaxRetries times (default 2) with
// exponential backoff and jitter: network errors and transient 5xx on
// idempotent methods, and 429 on any method, honoring Retry-After (capped at
// one minute). POST is never retried on 5xx or network errors — the API has
// no idempotency keys, and a blind retry could duplicate an envelope.
// Multipart uploads are attempted exactly once.
package signater
