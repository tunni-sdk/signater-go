package signater

import (
	"context"
	"math/rand/v2"
	"net/http"
	"time"
)

// maxRetryDelay caps the computed exponential backoff delay.
const maxRetryDelay = 8 * time.Second

// maxRetryAfter caps how long a server-provided Retry-After header is honored,
// so a hostile or misconfigured response cannot park a goroutine indefinitely.
const maxRetryAfter = time.Minute

// Retry policy. The Signater API does not support idempotency keys, so a POST
// whose outcome is unknown (network error, 5xx, timeout) must not be blindly
// retried — it could, for example, create a duplicate envelope. A 429 is safe
// to retry for any method because the request was rejected before processing.
//
//	network error / 5xx  → retried for idempotent methods only
//	429 Too Many Requests → retried for all methods, honoring Retry-After
func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions:
		return true
	}
	return false
}

// shouldRetryStatus reports whether a response status warrants a retry for the
// given method. Only transient 5xx statuses qualify: 501/505-style failures
// will not succeed on retry.
func shouldRetryStatus(method string, status int) bool {
	switch status {
	case http.StatusTooManyRequests:
		return true
	case http.StatusInternalServerError, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return isIdempotent(method)
	}
	return false
}

// backoffDelay computes the wait before retry attempt (0-based). A positive
// retryAfter (from a Retry-After header) is honored up to maxRetryAfter;
// otherwise the delay is exponential with full jitter: random in
// [0, base*2^attempt], capped at maxRetryDelay.
func backoffDelay(attempt int, base, retryAfter time.Duration) time.Duration {
	if retryAfter > 0 {
		return min(retryAfter, maxRetryAfter)
	}
	if base <= 0 {
		base = defaultRetryBaseDelay
	}
	// Guard the shift: large attempt values would overflow into a negative
	// ceiling and panic rand.N.
	ceiling := maxRetryDelay
	if attempt < 30 {
		if d := base << attempt; d > 0 && d < ceiling {
			ceiling = d
		}
	}
	return rand.N(ceiling + 1) //nolint:gosec // retry jitter has no security requirement
}

// sleepCtx sleeps for d or until ctx is done, returning ctx.Err() in that case.
func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
