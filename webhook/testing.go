package webhook

import "net/http"

// SignPayload returns headers carrying valid Hookdeck signatures for payload,
// so consumers can unit-test their webhook handlers without real secrets or
// deliveries. The first secret fills the primary signature header; a second
// secret, when given, fills the rotated-secret header. The returned header
// replaces (not augments) any existing headers:
//
//	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payload))
//	req.Header = webhook.SignPayload(payload, "test-secret")
func SignPayload(payload []byte, secrets ...string) http.Header {
	h := http.Header{}
	if len(secrets) > 0 {
		h.Set(hookdeckSignatureHeader, sign(payload, secrets[0]))
	}
	if len(secrets) > 1 {
		h.Set(hookdeckSignature2Header, sign(payload, secrets[1]))
	}
	h.Set("Content-Type", "application/json; charset=utf-8")
	return h
}
