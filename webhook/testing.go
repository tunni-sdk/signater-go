package webhook

import "net/http"

// SignPayload returns headers carrying a valid Hookdeck signature for
// payload, so consumers can unit-test their webhook handlers without real
// secrets or deliveries:
//
//	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payload))
//	req.Header = webhook.SignPayload(payload, "test-secret")
func SignPayload(payload []byte, secret string) http.Header {
	h := http.Header{}
	h.Set(hookdeckSignatureHeader, Sign(payload, secret))
	h.Set("Content-Type", "application/json; charset=utf-8")
	return h
}
