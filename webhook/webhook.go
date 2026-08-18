// Package webhook parses and verifies Signater webhook deliveries.
//
// Signater delivers webhooks through Hookdeck (retrying exponentially every
// 30s up to 5 times). Configure webhooks at
// https://app.signater.com/account/webhooks.
//
// Deliveries carry two authentication headers: X-Hookdeck-Signature (a
// base64-encoded 32-byte value; Hookdeck documents the scheme as HMAC-SHA256
// of the raw body) and X-Signater-Apikey (a static per-account key, distinct
// from the API token). The webhook configuration page does not currently
// expose the Hookdeck signing secret — which also means the HMAC scheme
// cannot be confirmed end-to-end — so most consumers must verify with the
// api key, pinning its value from a captured delivery:
//
//	payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
//	if err != nil {
//		http.Error(w, "bad request", http.StatusBadRequest)
//		return
//	}
//	if err := webhook.VerifyAPIKey(r.Header, apiKey); err != nil {
//		http.Error(w, "bad request", http.StatusBadRequest) // don't echo err
//		return
//	}
//	event, err := webhook.ParseEvent(payload)
//	if err != nil {
//		http.Error(w, "bad request", http.StatusBadRequest)
//		return
//	}
//	switch event.Type { case webhook.EventEnvelopeSigned: ... }
//
// When a signing secret is available, prefer ConstructEvent, which verifies
// the HMAC signature and parses in one call.
//
// Signater does not guarantee event delivery order, and Hookdeck retries mean
// deliveries are at-least-once: handlers must be idempotent and should dedupe
// (RequestID is the natural key). Neither check covers a timestamp, so a
// captured delivery can be replayed; treat events as hints and query the API
// for the current resource state.
package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// EventType identifies a webhook event. Unknown values parse fine (the type
// is open) so new API events cannot break consumers.
type EventType string

// Webhook event types published by Signater.
const (
	EventEnvelopeCreated             EventType = "envelope.created"
	EventEnvelopePublishScheduled    EventType = "envelope.publish_scheduled"
	EventEnvelopePublished           EventType = "envelope.published"
	EventEnvelopePublishedBySchedule EventType = "envelope.published_by_schedule"
	EventEnvelopeUnpublished         EventType = "envelope.unpublished"
	EventEnvelopeUpdated             EventType = "envelope.updated"
	EventEnvelopeCancelled           EventType = "envelope.cancelled"
	EventEnvelopeExpired             EventType = "envelope.expired"
	EventEnvelopeSigned              EventType = "envelope.signed"
	EventEnvelopeViewedBySigner      EventType = "envelope.viewed_by_signer"
	EventEnvelopeCancelledMfaError   EventType = "envelope.cancelled_due_to_mfa_error_by_signer"
	EventEnvelopeRejectedBySigner    EventType = "envelope.rejected_by_signer"
	EventEnvelopeApprovedBySigner    EventType = "envelope.approved_by_signer"
)

// Event is a webhook delivery payload.
type Event struct {
	// EnvelopeID is the envelope that triggered the event.
	EnvelopeID string `json:"envelope_id"`
	// Type is the event type.
	Type EventType `json:"event_type"`
	// AccountID is the account that triggered the event.
	AccountID string `json:"account_id"`
	// Env is the originating environment: "production" or "sandbox".
	Env string `json:"env"`
	// SignerID is the signer involved, present only on *_by_signer events
	// (key confirmed against captured sandbox deliveries).
	SignerID string `json:"signer_id,omitempty"`
	// Raw is the verbatim delivery payload, for fields not modeled above.
	Raw json.RawMessage `json:"-"`
}

// RequestID returns the delivery's unique id, the natural key for
// deduplicating at-least-once deliveries. It reads the request-id header the
// Signater docs describe, falling back to X-Hookdeck-Eventid — captured
// sandbox deliveries carry only the latter, which is stable across Hookdeck
// retry attempts. Empty when both headers are absent.
func RequestID(h http.Header) string {
	if id := h.Get("Request-Id"); id != "" {
		return id
	}
	return h.Get("X-Hookdeck-Eventid")
}

// Signature and authentication headers.
const (
	// hookdeckSignatureHeader carries the primary HMAC signature.
	hookdeckSignatureHeader = "X-Hookdeck-Signature"
	// hookdeckSignature2Header carries the signature for a rotated secret.
	hookdeckSignature2Header = "X-Hookdeck-Signature-2"
	// apiKeyHeader carries the Signater account api key.
	apiKeyHeader = "X-Signater-Apikey" //nolint:gosec // header name, not a credential
)

// Verification errors.
var (
	// ErrMissingSignature is returned when no signature header is present.
	ErrMissingSignature = errors.New("webhook: missing Hookdeck signature header")
	// ErrInvalidSignature is returned when no signature header matches the
	// payload.
	ErrInvalidSignature = errors.New("webhook: signature verification failed")
	// ErrAPIKeyMismatch is returned when the x-signater-apikey header does
	// not match the expected key.
	ErrAPIKeyMismatch = errors.New("webhook: x-signater-apikey mismatch")
	// ErrNoSecret is returned when the signing secret (or expected api key)
	// is empty — verification must fail closed, never accept forgeable
	// signatures computed with an empty key.
	ErrNoSecret = errors.New("webhook: signing secret is empty")
)

// ParseEvent decodes a webhook payload WITHOUT verifying its authenticity.
// Prefer ConstructEvent whenever a signing secret is available.
func ParseEvent(payload []byte) (*Event, error) {
	var e Event
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, fmt.Errorf("webhook: decoding payload: %w", err)
	}
	if e.Type == "" {
		return nil, errors.New("webhook: payload has no event_type")
	}
	e.Raw = append(json.RawMessage(nil), payload...)
	return &e, nil
}

// ConstructEvent verifies the Hookdeck signature of a delivery and returns
// the parsed event. payload must be the raw request body and h the request
// headers.
//
// Captured sandbox deliveries confirm the header name and format only: a
// single base64-encoded 32-byte value in X-Hookdeck-Signature, consistent
// with Hookdeck's published HMAC-SHA256 scheme (X-Hookdeck-Signature-2
// covers a rotated secret). The HMAC construction itself is unconfirmed —
// the Signater webhook page does not expose the signing secret, so until it
// does, verify deliveries with VerifyAPIKey instead.
func ConstructEvent(payload []byte, h http.Header, secret string) (*Event, error) {
	if err := VerifyHookdeckSignature(payload, h, secret); err != nil {
		return nil, err
	}
	return ParseEvent(payload)
}

// VerifyHookdeckSignature checks the delivery's HMAC signature against
// secret, accepting either the primary or the rotated-secret header. An empty
// secret is rejected (ErrNoSecret): HMAC with an empty key is forgeable.
//
// Header lookup is canonicalized — build headers with http.Header.Set (as
// net/http does), not a map literal with lowercase keys.
func VerifyHookdeckSignature(payload []byte, h http.Header, secret string) error {
	if secret == "" {
		return ErrNoSecret
	}
	sigs := make([]string, 0, 2)
	if s := h.Get(hookdeckSignatureHeader); s != "" {
		sigs = append(sigs, s)
	}
	if s := h.Get(hookdeckSignature2Header); s != "" {
		sigs = append(sigs, s)
	}
	if len(sigs) == 0 {
		return ErrMissingSignature
	}
	expected := sign(payload, secret)
	for _, sig := range sigs {
		if hmac.Equal([]byte(sig), []byte(expected)) {
			return nil
		}
	}
	return ErrInvalidSignature
}

// VerifyAPIKey compares the delivery's x-signater-apikey header with the
// expected account webhook key, in constant time. The key is static per
// account, distinct from the API token, and sent on every delivery — pin its
// value from a captured delivery. It is a weaker check than
// VerifyHookdeckSignature (it does not cover the payload), but it is the
// strongest one available while Signater does not expose the signing secret.
// An empty expected key is rejected (ErrNoSecret).
//
// Because the key travels in cleartext inside every request, anyone who can
// observe a delivery (TLS-terminating proxies, header-logging middleware)
// can replay it on forged events: terminate TLS as close to the handler as
// possible, keep the key out of logs, and always re-query the API for the
// resource's current state before acting on an event.
func VerifyAPIKey(h http.Header, apiKey string) error {
	if apiKey == "" {
		return ErrNoSecret
	}
	got := h.Get(apiKeyHeader)
	if got == "" || subtle.ConstantTimeCompare([]byte(got), []byte(apiKey)) != 1 {
		return ErrAPIKeyMismatch
	}
	return nil
}

// sign computes the base64 HMAC-SHA256 signature of payload with secret, per
// Hookdeck's published scheme (unconfirmed against real deliveries — the
// signing secret is not exposed). Unexported: consumers verifying by string
// comparison would reintroduce the timing side channel this package avoids.
// Use SignPayload to generate test headers.
func sign(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
