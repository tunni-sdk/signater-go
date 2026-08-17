package webhook

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

// docPayload is the documented example webhook payload.
const docPayload = `{
  "envelope_id": "696d0b7f-8e2f-4fb1-9605-72b197ee154d",
  "event_type": "envelope.created",
  "account_id": "00cae3c2-8b7c-41b6-b934-6724808fe7f7",
  "env": "production"
}`

func TestParseEvent(t *testing.T) {
	e, err := ParseEvent([]byte(docPayload))
	if err != nil {
		t.Fatal(err)
	}
	if e.EnvelopeID != "696d0b7f-8e2f-4fb1-9605-72b197ee154d" || e.Type != EventEnvelopeCreated ||
		e.AccountID != "00cae3c2-8b7c-41b6-b934-6724808fe7f7" || e.Env != "production" {
		t.Errorf("event = %+v", e)
	}
	if e.SignerID != "" {
		t.Errorf("SignerID = %q, want empty", e.SignerID)
	}
}

func TestParseEventBySigner(t *testing.T) {
	payload := `{"envelope_id":"e1","event_type":"envelope.approved_by_signer","account_id":"a1","env":"sandbox","signer_id":"s1"}`
	e, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatal(err)
	}
	if e.Type != EventEnvelopeApprovedBySigner || e.SignerID != "s1" || e.Env != "sandbox" {
		t.Errorf("event = %+v", e)
	}
}

func TestParseEventUnknownTypeIsAccepted(t *testing.T) {
	e, err := ParseEvent([]byte(`{"envelope_id":"e1","event_type":"envelope.some_future_event"}`))
	if err != nil {
		t.Fatal(err)
	}
	if e.Type != "envelope.some_future_event" {
		t.Errorf("Type = %q", e.Type)
	}
}

func TestParseEventInvalid(t *testing.T) {
	if _, err := ParseEvent([]byte("not json")); err == nil {
		t.Error("malformed payload must error")
	}
	if _, err := ParseEvent([]byte(`{"envelope_id":"e1"}`)); err == nil {
		t.Error("payload without event_type must error")
	}
}

func TestConstructEventValidSignature(t *testing.T) {
	payload := []byte(docPayload)
	h := SignPayload(payload, "secret-1")

	e, err := ConstructEvent(payload, h, "secret-1")
	if err != nil {
		t.Fatal(err)
	}
	if e.Type != EventEnvelopeCreated {
		t.Errorf("Type = %q", e.Type)
	}
}

func TestConstructEventTamperedPayload(t *testing.T) {
	payload := []byte(docPayload)
	h := SignPayload(payload, "secret-1")

	// Flip one character inside envelope_id: still valid JSON, semantically
	// different — the MAC must cover meaningful content.
	tampered := []byte(strings.Replace(string(payload), "696d0b7f", "696d0b7e", 1))
	if _, err := ConstructEvent(tampered, h, "secret-1"); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("err = %v, want ErrInvalidSignature", err)
	}
}

func TestConstructEventEmptySecretFailsClosed(t *testing.T) {
	// HMAC with an empty key is computable by anyone; verification must
	// reject it instead of accepting a forgeable signature.
	payload := []byte(docPayload)
	h := SignPayload(payload, "")
	if _, err := ConstructEvent(payload, h, ""); !errors.Is(err, ErrNoSecret) {
		t.Errorf("err = %v, want ErrNoSecret", err)
	}
}

func TestVerifyHookdeckSignatureLowercaseMapLiteralFailsClosed(t *testing.T) {
	// Headers built as map literals with lowercase keys are not canonical
	// and are deliberately invisible to lookup — must fail closed.
	payload := []byte(docPayload)
	h := http.Header{"x-hookdeck-signature": {sign(payload, "s")}}
	if err := VerifyHookdeckSignature(payload, h, "s"); !errors.Is(err, ErrMissingSignature) {
		t.Errorf("err = %v, want ErrMissingSignature", err)
	}
}

func TestParseEventPreservesRaw(t *testing.T) {
	payload := `{"envelope_id":"e1","event_type":"envelope.created","signerId":"unmodeled-key"}`
	e, err := ParseEvent([]byte(payload))
	if err != nil {
		t.Fatal(err)
	}
	if string(e.Raw) != payload {
		t.Errorf("Raw = %s", e.Raw)
	}
}

func TestRequestID(t *testing.T) {
	h := http.Header{}
	h.Set("request-id", "|6b1bc96c.09a8251c.")
	if got := RequestID(h); got != "|6b1bc96c.09a8251c." {
		t.Errorf("RequestID = %q", got)
	}
	if got := RequestID(http.Header{}); got != "" {
		t.Errorf("RequestID(empty) = %q", got)
	}
}

func TestConstructEventWrongSecret(t *testing.T) {
	payload := []byte(docPayload)
	h := SignPayload(payload, "secret-1")
	if _, err := ConstructEvent(payload, h, "other-secret"); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("err = %v, want ErrInvalidSignature", err)
	}
}

func TestConstructEventMissingSignature(t *testing.T) {
	if _, err := ConstructEvent([]byte(docPayload), http.Header{}, "secret-1"); !errors.Is(err, ErrMissingSignature) {
		t.Errorf("err = %v, want ErrMissingSignature", err)
	}
}

func TestVerifyHookdeckSignatureRotatedSecondHeader(t *testing.T) {
	payload := []byte(docPayload)
	h := http.Header{}
	h.Set("X-Hookdeck-Signature", "invalid-old-signature")
	h.Set("X-Hookdeck-Signature-2", sign(payload, "rotated-secret"))

	if err := VerifyHookdeckSignature(payload, h, "rotated-secret"); err != nil {
		t.Errorf("rotated secret in second header must verify, got %v", err)
	}
}

func TestVerifyAPIKey(t *testing.T) {
	h := http.Header{}
	h.Set("x-signater-apikey", "35198a2a282a41e5b1c64d55bc55c81c")

	if err := VerifyAPIKey(h, "35198a2a282a41e5b1c64d55bc55c81c"); err != nil {
		t.Errorf("matching key must verify, got %v", err)
	}
	if err := VerifyAPIKey(h, "wrong"); !errors.Is(err, ErrAPIKeyMismatch) {
		t.Errorf("err = %v, want ErrAPIKeyMismatch", err)
	}
	if err := VerifyAPIKey(http.Header{}, "x"); !errors.Is(err, ErrAPIKeyMismatch) {
		t.Errorf("missing header: err = %v, want ErrAPIKeyMismatch", err)
	}
}
