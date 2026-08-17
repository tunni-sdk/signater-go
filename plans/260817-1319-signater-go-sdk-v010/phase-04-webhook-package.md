---
title: "Phase 4: Webhook Package"
status: todo
priority: P1
effort: "0.5-1d"
dependencies: [1]
---

# Phase 4: Webhook Package

## Overview

`webhook/` subpackage (Stripe pattern): parse Signater webhook payloads into typed events and verify request authenticity (Hookdeck signature and/or `x-signater-apikey` match). Depends only on Phase 1 (can run parallel to 2-3).

## Requirements

- Functional: parse `{envelope_id, event_type, account_id, env}` (+ signer id on `_by_signer` events); 13 `EventType` constants; verification helpers; test-payload generator.
- Non-functional: stdlib only (`crypto/hmac`, `crypto/sha256`, `encoding/base64`); usable from any HTTP framework (operates on `[]byte` body + `http.Header`).

## Architecture

```go
package webhook // github.com/tunni-sdk/signater-go/webhook

type EventType string
const (
    EventEnvelopeCreated              EventType = "envelope.created"
    EventEnvelopePublishScheduled     EventType = "envelope.publish_scheduled"
    EventEnvelopePublished            EventType = "envelope.published"
    EventEnvelopePublishedBySchedule  EventType = "envelope.published_by_schedule"
    EventEnvelopeUnpublished          EventType = "envelope.unpublished"
    EventEnvelopeUpdated              EventType = "envelope.updated"
    EventEnvelopeCancelled            EventType = "envelope.cancelled"
    EventEnvelopeExpired              EventType = "envelope.expired"
    EventEnvelopeSigned               EventType = "envelope.signed"
    EventEnvelopeViewedBySigner       EventType = "envelope.viewed_by_signer"
    EventEnvelopeCancelledMfaError    EventType = "envelope.cancelled_due_to_mfa_error_by_signer"
    EventEnvelopeRejectedBySigner     EventType = "envelope.rejected_by_signer"
    EventEnvelopeApprovedBySigner     EventType = "envelope.approved_by_signer"
)

type Event struct {
    EnvelopeID string    `json:"envelope_id"`
    Type       EventType `json:"event_type"`
    AccountID  string    `json:"account_id"`
    Env        string    `json:"env"`              // "production" | "sandbox"
    SignerID   string    `json:"signer_id,omitempty"` // only on *_by_signer events — verify exact JSON key empirically
}

func ParseEvent(payload []byte) (*Event, error)                              // no verification
func ConstructEvent(payload []byte, h http.Header, secret string) (*Event, error) // verify + parse
func VerifyHookdeckSignature(payload []byte, h http.Header, secret string) error
    // Hookdeck: x-hookdeck-signature (and rotated x-hookdeck-signature-2) =
    // base64(HMAC-SHA256(secret, raw body)); compare with hmac.Equal against both headers
func VerifyAPIKey(h http.Header, apiKey string) error // constant-time match of x-signater-apikey

// webhook/testing.go
func SignPayload(payload []byte, secret string) http.Header // for consumer tests
```

## Related Code Files

- Create: `webhook/webhook.go`, `webhook/testing.go`, `webhook/webhook_test.go`
- Reference: `../2026-08-17-sdk-architecture/reference/webhooks-799774m0.md`

## Implementation Steps

1. Implement types + `ParseEvent` (reject unknown/empty `event_type`? No — keep open: unknown types parse fine, `Type` is a string; document this forward-compat choice).
2. Implement `VerifyHookdeckSignature` per Hookdeck's documented scheme (base64 HMAC-SHA256 of raw body; check both signature headers); `ConstructEvent` = verify + parse.
3. Implement `VerifyAPIKey` with `subtle.ConstantTimeCompare`.
4. `testing.go` `SignPayload` helper; tests: valid sig passes, tampered body fails, rotated second header passes, api-key match/mismatch, payload parsing incl. `_by_signer`.
5. Doc comment on the package: how to get the secret (Hookdeck/Signater webhook config page), recommend async processing + re-query API (docs: event order not guaranteed).

## Success Criteria

- [x] All 13 event constants exactly match doc strings
- [x] `ConstructEvent` rejects tampered payloads and accepts either signature header
- [x] `SignPayload` helper lets consumers unit-test their handlers without real secrets
- [x] Package compiles standalone (`go build ./webhook`), stdlib only

## Risk Assessment

- **Empirical unknowns (must validate in sandbox during Phase 5):** exact Hookdeck signature header name/format for Signater's setup, and the JSON key for signer id (docs say "signer ID will be included" without showing the key). Mark both with `// VERIFY(sandbox):` comments; quickstart validation in Phase 5 closes them.
- If real payloads diverge, only `webhook.go` changes — API surface (`ConstructEvent`) is stable.
