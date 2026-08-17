package signater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// telemetryHeader carries the unique operation id Signater attaches to every
// response. Provide it to Signater support when investigating a request.
const telemetryHeader = "X-Signater-Telemetry-Operation-Id"

// APIError is the base error returned for any non-2xx API response. More
// specific error types embed it; use errors.As to match either the concrete
// type or *APIError itself:
//
//	var apiErr *signater.APIError
//	if errors.As(err, &apiErr) {
//		log.Println(apiErr.StatusCode, apiErr.TelemetryOperationID)
//	}
type APIError struct {
	// StatusCode is the HTTP status code of the response.
	StatusCode int
	// Message is a human-readable description extracted from the response body.
	Message string
	// TelemetryOperationID is the value of the X-Signater-Telemetry-Operation-Id
	// response header, useful when contacting Signater support.
	TelemetryOperationID string
	// RawBody is the raw response body, kept for debugging and forward
	// compatibility with fields this SDK does not model yet.
	RawBody []byte
}

// Error implements the error interface.
func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = http.StatusText(e.StatusCode)
	}
	if e.TelemetryOperationID != "" {
		return fmt.Sprintf("signater: HTTP %d: %s (telemetry operation id: %s)", e.StatusCode, msg, e.TelemetryOperationID)
	}
	return fmt.Sprintf("signater: HTTP %d: %s", e.StatusCode, msg)
}

// FieldError is a single validation failure inside a 400 Bad Request response.
type FieldError struct {
	// Message describes the validation failure.
	Message string `json:"message"`
	// Metadata carries optional structured details about the failure, or null.
	Metadata json.RawMessage `json:"metadata"`
}

// ValidationError is returned for 400 Bad Request responses. Errors lists the
// individual validation failures reported by the API.
type ValidationError struct {
	APIError
	Errors []FieldError
}

// Error implements the error interface, appending the count of additional
// validation failures beyond the first.
func (e *ValidationError) Error() string {
	if n := len(e.Errors); n > 1 {
		return fmt.Sprintf("%s (and %d more validation errors)", e.APIError.Error(), n-1)
	}
	return e.APIError.Error()
}

// Unwrap allows errors.As to match *APIError.
func (e *ValidationError) Unwrap() error { return &e.APIError }

// AuthenticationError is returned for 401 Unauthorized responses: the token is
// missing, invalid, inactive, or belongs to an inactive user.
type AuthenticationError struct {
	APIError
}

// Unwrap allows errors.As to match *APIError.
func (e *AuthenticationError) Unwrap() error { return &e.APIError }

// ShouldBuy identifies the product required to proceed when the API returns
// 402 Payment Required.
type ShouldBuy string

// Possible ShouldBuy values documented by the Signater API.
const (
	ShouldBuyEnvelopes        ShouldBuy = "Envelopes"
	ShouldBuyAPIEnvelopes     ShouldBuy = "ApiEnvelopes"
	ShouldBuySignerMfaCredits ShouldBuy = "SignerMfaCredits"
)

// PaymentRequiredError is returned for 402 Payment Required responses,
// currently only from envelope publication. ShouldBuy indicates which product
// must be purchased at https://app.signater.com/billing/plans.
type PaymentRequiredError struct {
	APIError
	ShouldBuy ShouldBuy
}

// Unwrap allows errors.As to match *APIError.
func (e *PaymentRequiredError) Unwrap() error { return &e.APIError }

// PermissionError is returned for 403 Forbidden responses: the authenticated
// user lacks the permission required by the endpoint.
type PermissionError struct {
	APIError
}

// Unwrap allows errors.As to match *APIError.
func (e *PermissionError) Unwrap() error { return &e.APIError }

// NotFoundError is returned for 404 Not Found responses on endpoints that
// document them.
type NotFoundError struct {
	APIError
}

// Unwrap allows errors.As to match *APIError.
func (e *NotFoundError) Unwrap() error { return &e.APIError }

// RateLimitError is returned for 429 Too Many Requests responses. The API
// enforces a limit of 1,000 requests per minute. RetryAfter is the wait time
// suggested by the Retry-After header, or zero when the header is absent.
type RateLimitError struct {
	APIError
	RetryAfter time.Duration
}

// Unwrap allows errors.As to match *APIError.
func (e *RateLimitError) Unwrap() error { return &e.APIError }

// newAPIError maps a non-2xx response to the matching typed error.
func newAPIError(resp *http.Response, body []byte) error {
	base := APIError{
		StatusCode:           resp.StatusCode,
		TelemetryOperationID: resp.Header.Get(telemetryHeader),
		RawBody:              body,
	}
	switch resp.StatusCode {
	case http.StatusBadRequest:
		e := &ValidationError{APIError: base}
		var payload struct {
			Errors []FieldError `json:"errors"`
		}
		_ = json.Unmarshal(body, &payload)
		e.Errors = payload.Errors
		if len(e.Errors) > 0 {
			e.Message = e.Errors[0].Message
		} else {
			e.Message = "invalid request"
		}
		return e
	case http.StatusUnauthorized:
		base.Message = messageFrom(body, "authentication failed")
		return &AuthenticationError{APIError: base}
	case http.StatusPaymentRequired:
		e := &PaymentRequiredError{APIError: base}
		var payload struct {
			ShouldBuy ShouldBuy `json:"ShouldBuy"`
		}
		_ = json.Unmarshal(body, &payload)
		e.ShouldBuy = payload.ShouldBuy
		e.Message = "insufficient credits"
		if e.ShouldBuy != "" {
			e.Message = "insufficient credits: should buy " + string(e.ShouldBuy)
		}
		return e
	case http.StatusForbidden:
		base.Message = messageFrom(body, "forbidden")
		return &PermissionError{APIError: base}
	case http.StatusNotFound:
		base.Message = messageFrom(body, "not found")
		return &NotFoundError{APIError: base}
	case http.StatusTooManyRequests:
		e := &RateLimitError{APIError: base}
		e.Message = "rate limit exceeded"
		e.RetryAfter = retryAfterFrom(resp.Header)
		return e
	default:
		base.Message = messageFrom(body, "")
		return &base
	}
}

// messageFrom extracts a {"message": ...} body, falling back to def.
func messageFrom(body []byte, def string) string {
	var payload struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.Message != "" {
		return payload.Message
	}
	return def
}

// retryAfterFrom parses a Retry-After header in either of its RFC 9110 forms:
// delay seconds or an HTTP date.
func retryAfterFrom(h http.Header) time.Duration {
	raw := h.Get("Retry-After")
	if raw == "" {
		return 0
	}
	if secs, err := strconv.Atoi(raw); err == nil {
		if secs < 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(raw); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}
