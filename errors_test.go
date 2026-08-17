package signater

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

// respond builds an *http.Response suitable for newAPIError.
func respond(status int, header http.Header) *http.Response {
	if header == nil {
		header = http.Header{}
	}
	header.Set(telemetryHeader, "123e4567-e89b-12d3-a456-426614174000")
	return &http.Response{StatusCode: status, Header: header}
}

func TestValidationError(t *testing.T) {
	// Payload from the Signater error-handling documentation.
	body := []byte(`{"errors":[{"message":"The length of 'Name' must be at least 2 characters. You entered 1 characters.","metadata":null}]}`)
	err := newAPIError(respond(http.StatusBadRequest, nil), body)

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("want *ValidationError, got %T", err)
	}
	if len(ve.Errors) != 1 || ve.Errors[0].Message == "" {
		t.Errorf("Errors = %+v", ve.Errors)
	}
	if ve.Message != ve.Errors[0].Message {
		t.Errorf("Message = %q", ve.Message)
	}
	var base *APIError
	if !errors.As(err, &base) {
		t.Fatal("errors.As should match *APIError via Unwrap")
	}
	if base.StatusCode != 400 || base.TelemetryOperationID == "" {
		t.Errorf("base = %+v", base)
	}
}

func TestValidationErrorMalformedBody(t *testing.T) {
	err := newAPIError(respond(http.StatusBadRequest, nil), []byte("not json"))
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("want *ValidationError, got %T", err)
	}
	if ve.Message != "invalid request" {
		t.Errorf("Message = %q", ve.Message)
	}
}

func TestAuthenticationError(t *testing.T) {
	err := newAPIError(respond(http.StatusUnauthorized, nil), []byte(`{"message":"Authentication failed."}`))
	var ae *AuthenticationError
	if !errors.As(err, &ae) {
		t.Fatalf("want *AuthenticationError, got %T", err)
	}
	if ae.Message != "Authentication failed." {
		t.Errorf("Message = %q", ae.Message)
	}
}

func TestPaymentRequiredError(t *testing.T) {
	// The 402 payload uses PascalCase, unlike other endpoints.
	err := newAPIError(respond(http.StatusPaymentRequired, nil), []byte(`{"ShouldBuy":"ApiEnvelopes"}`))
	var pe *PaymentRequiredError
	if !errors.As(err, &pe) {
		t.Fatalf("want *PaymentRequiredError, got %T", err)
	}
	if pe.ShouldBuy != ShouldBuyAPIEnvelopes {
		t.Errorf("ShouldBuy = %q", pe.ShouldBuy)
	}
}

func TestPermissionError(t *testing.T) {
	err := newAPIError(respond(http.StatusForbidden, nil), []byte(`{"message":"User should be an administrator."}`))
	var pe *PermissionError
	if !errors.As(err, &pe) {
		t.Fatalf("want *PermissionError, got %T", err)
	}
	if pe.Message != "User should be an administrator." {
		t.Errorf("Message = %q", pe.Message)
	}
}

func TestNotFoundError(t *testing.T) {
	err := newAPIError(respond(http.StatusNotFound, nil), nil)
	var ne *NotFoundError
	if !errors.As(err, &ne) {
		t.Fatalf("want *NotFoundError, got %T", err)
	}
	if ne.Message != "not found" {
		t.Errorf("Message = %q", ne.Message)
	}
}

func TestRateLimitError(t *testing.T) {
	h := http.Header{}
	h.Set("Retry-After", "60")
	err := newAPIError(respond(http.StatusTooManyRequests, h), nil)
	var re *RateLimitError
	if !errors.As(err, &re) {
		t.Fatalf("want *RateLimitError, got %T", err)
	}
	if re.RetryAfter != 60*time.Second {
		t.Errorf("RetryAfter = %v", re.RetryAfter)
	}
}

func TestRateLimitErrorNoHeader(t *testing.T) {
	err := newAPIError(respond(http.StatusTooManyRequests, nil), nil)
	var re *RateLimitError
	if !errors.As(err, &re) {
		t.Fatalf("want *RateLimitError, got %T", err)
	}
	if re.RetryAfter != 0 {
		t.Errorf("RetryAfter = %v, want 0", re.RetryAfter)
	}
}

func TestUnmappedStatusReturnsBaseError(t *testing.T) {
	err := newAPIError(respond(http.StatusBadGateway, nil), []byte(`{"message":"upstream broke"}`))
	var base *APIError
	if !errors.As(err, &base) {
		t.Fatalf("want *APIError, got %T", err)
	}
	if base.Message != "upstream broke" || base.StatusCode != 502 {
		t.Errorf("base = %+v", base)
	}
}

func TestErrorString(t *testing.T) {
	e := &APIError{StatusCode: 400, Message: "bad", TelemetryOperationID: "tid"}
	want := "signater: HTTP 400: bad (telemetry operation id: tid)"
	if e.Error() != want {
		t.Errorf("Error() = %q, want %q", e.Error(), want)
	}
	e2 := &APIError{StatusCode: 500}
	if e2.Error() != "signater: HTTP 500: Internal Server Error" {
		t.Errorf("Error() = %q", e2.Error())
	}
}
