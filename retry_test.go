package signater

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestRetryGetOn5xx(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusInternalServerError, "", nil)
	tr.Enqueue(http.StatusServiceUnavailable, "", nil)
	tr.Enqueue(http.StatusOK, `{}`, nil)

	err := c.do(context.Background(), http.MethodGet, "/v1/ecm/contacts", nil, nil, nil)
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if tr.Count() != 3 {
		t.Errorf("requests = %d, want 3", tr.Count())
	}
}

func TestNoRetryPostOn5xx(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusInternalServerError, "", nil)

	err := c.do(context.Background(), http.MethodPost, "/v1/ecm/envelopes", nil, nil, nil)
	var base *APIError
	if !errors.As(err, &base) || base.StatusCode != 500 {
		t.Fatalf("want 500 APIError, got %v", err)
	}
	if tr.Count() != 1 {
		t.Errorf("requests = %d, want 1 (POST must not retry on 5xx)", tr.Count())
	}
}

func TestRetryPostOn429(t *testing.T) {
	c, tr := newTestClient(t)
	h := http.Header{}
	h.Set("Retry-After", "0") // header present but zero: falls back to jittered backoff
	tr.Enqueue(http.StatusTooManyRequests, "", h)
	tr.Enqueue(http.StatusCreated, `{}`, nil)

	err := c.do(context.Background(), http.MethodPost, "/v1/ecm/envelopes", nil, nil, nil)
	if err != nil {
		t.Fatalf("expected success after 429 retry, got %v", err)
	}
	if tr.Count() != 2 {
		t.Errorf("requests = %d, want 2", tr.Count())
	}
}

func TestRetryExhaustedReturnsTypedError(t *testing.T) {
	c, tr := newTestClient(t) // maxRetries = 2 → 3 attempts
	for range 3 {
		tr.Enqueue(http.StatusTooManyRequests, "", nil)
	}

	err := c.do(context.Background(), http.MethodGet, "/v1/ecm/contacts", nil, nil, nil)
	var re *RateLimitError
	if !errors.As(err, &re) {
		t.Fatalf("want *RateLimitError, got %v", err)
	}
	if tr.Count() != 3 {
		t.Errorf("requests = %d, want 3", tr.Count())
	}
}

func TestRetryGetOnNetworkError(t *testing.T) {
	c, tr := newTestClient(t)
	tr.EnqueueError(errors.New("connection reset"))
	tr.Enqueue(http.StatusOK, `{}`, nil)

	if err := c.do(context.Background(), http.MethodGet, "/v1/ecm/contacts", nil, nil, nil); err != nil {
		t.Fatalf("expected success after network retry, got %v", err)
	}
	if tr.Count() != 2 {
		t.Errorf("requests = %d, want 2", tr.Count())
	}
}

func TestNoRetryPostOnNetworkError(t *testing.T) {
	c, tr := newTestClient(t)
	tr.EnqueueError(errors.New("connection reset"))

	err := c.do(context.Background(), http.MethodPost, "/v1/ecm/envelopes", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if tr.Count() != 1 {
		t.Errorf("requests = %d, want 1 (POST must not retry on network error)", tr.Count())
	}
}

func TestRetryHonorsContextCancellation(t *testing.T) {
	c, tr := newTestClient(t)
	h := http.Header{}
	h.Set("Retry-After", "5") // long enough that the context wins
	tr.Enqueue(http.StatusTooManyRequests, "", h)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := c.do(ctx, http.MethodGet, "/v1/ecm/contacts", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("call took %v; context cancellation did not interrupt the retry sleep", elapsed)
	}
	if tr.Count() != 1 {
		t.Errorf("requests = %d, want 1", tr.Count())
	}
}

func TestBackoffDelay(t *testing.T) {
	if d := backoffDelay(0, time.Second, 42*time.Second); d != 42*time.Second {
		t.Errorf("retryAfter should win, got %v", d)
	}
	if d := backoffDelay(0, time.Second, 24*time.Hour); d != maxRetryAfter {
		t.Errorf("retryAfter must be capped at %v, got %v", maxRetryAfter, d)
	}
	for attempt := range 5 {
		d := backoffDelay(attempt, 500*time.Millisecond, 0)
		ceiling := min(500*time.Millisecond<<attempt, maxRetryDelay)
		if d < 0 || d > ceiling {
			t.Errorf("attempt %d: delay %v outside [0, %v]", attempt, d, ceiling)
		}
	}
	// Large attempt values must not overflow the shift and panic.
	for _, attempt := range []int{35, 63, 1000} {
		d := backoffDelay(attempt, 500*time.Millisecond, 0)
		if d < 0 || d > maxRetryDelay {
			t.Errorf("attempt %d: delay %v outside [0, %v]", attempt, d, maxRetryDelay)
		}
	}
}

func TestRetryAfterFrom(t *testing.T) {
	h := http.Header{}
	h.Set("Retry-After", "30")
	if d := retryAfterFrom(h); d != 30*time.Second {
		t.Errorf("seconds form = %v", d)
	}
	h.Set("Retry-After", time.Now().Add(90*time.Second).UTC().Format(http.TimeFormat))
	if d := retryAfterFrom(h); d < 80*time.Second || d > 90*time.Second {
		t.Errorf("HTTP-date form = %v, want ~90s", d)
	}
	for _, raw := range []string{"-5", "garbage", ""} {
		h.Set("Retry-After", raw)
		if d := retryAfterFrom(h); d != 0 {
			t.Errorf("Retry-After %q = %v, want 0", raw, d)
		}
	}
}

func TestShouldRetryStatus(t *testing.T) {
	cases := []struct {
		method string
		status int
		want   bool
	}{
		{http.MethodGet, 500, true},
		{http.MethodPut, 503, true},
		{http.MethodDelete, 502, true},
		{http.MethodPost, 500, false},
		{http.MethodPost, 429, true},
		{http.MethodGet, 429, true},
		{http.MethodGet, 400, false},
		{http.MethodGet, 401, false},
		{http.MethodPatch, 500, false},
		// Non-transient 5xx never retry: a 501 will not succeed on retry.
		{http.MethodGet, 501, false},
		{http.MethodGet, 505, false},
	}
	for _, tc := range cases {
		if got := shouldRetryStatus(tc.method, tc.status); got != tc.want {
			t.Errorf("shouldRetryStatus(%s, %d) = %v, want %v", tc.method, tc.status, got, tc.want)
		}
	}
}
