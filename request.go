package signater

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// errMissingToken is returned when no API token was configured.
var errMissingToken = errors.New("signater: missing API token: use WithAPIToken or set " + envAPIToken)

// pathEscape escapes a dynamic value (e.g. a resource id) for safe use as a
// single URL path segment when building API paths.
func pathEscape(s string) string { return url.PathEscape(s) }

// do executes an API request with the client's retry policy.
//
// path is the API path (e.g. "/v1/ecm/contacts"); any dynamic segment (such
// as a resource id) MUST be escaped with pathEscape before interpolation, or
// values containing "/" or ".." would change the request target. query is the
// optional query string, body the optional JSON payload (an untyped nil or a
// non-nil value — a typed nil pointer would serialize as JSON null), and out
// the optional destination for the decoded JSON response.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	if c.initErr != nil {
		return c.initErr
	}
	if c.apiToken == "" {
		return errMissingToken
	}

	var payload []byte
	if body != nil {
		var err error
		if payload, err = json.Marshal(body); err != nil {
			return fmt.Errorf("signater: encoding request body: %w", err)
		}
	}

	u := c.baseURL.JoinPath(path)
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	for attempt := 0; ; attempt++ {
		resp, respBody, err := c.attempt(ctx, method, u.String(), payload)
		if err != nil {
			// Transport-level failure: the request may or may not have reached
			// the server, so only idempotent methods are retried.
			if ctx.Err() == nil && attempt < c.maxRetries && isIdempotent(method) {
				if serr := c.retrySleep(ctx, attempt, 0, method, err.Error()); serr == nil {
					continue
				}
			}
			return err
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if out != nil && len(respBody) > 0 {
				if err := json.Unmarshal(respBody, out); err != nil {
					return fmt.Errorf("signater: decoding response: %w", err)
				}
			}
			return nil
		}

		apiErr := newAPIError(resp, respBody)
		if attempt < c.maxRetries && shouldRetryStatus(method, resp.StatusCode) {
			// Retry-After may accompany any retryable status (e.g. 503), not
			// only 429, so it is read from the response directly.
			retryAfter := retryAfterFrom(resp.Header)
			if serr := c.retrySleep(ctx, attempt, retryAfter, method, apiErr.Error()); serr == nil {
				continue
			}
		}
		return apiErr
	}
}

// attempt performs a single HTTP round trip, applying the per-attempt timeout,
// and returns the response with its fully read body.
func (c *Client) attempt(ctx context.Context, method, fullURL string, payload []byte) (*http.Response, []byte, error) {
	attemptCtx := ctx
	if c.requestTimeout > 0 {
		var cancel context.CancelFunc
		attemptCtx, cancel = context.WithTimeout(ctx, c.requestTimeout)
		defer cancel()
	}

	var bodyReader io.Reader
	if payload != nil {
		bodyReader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(attemptCtx, method, fullURL, bodyReader)
	if err != nil {
		return nil, nil, fmt.Errorf("signater: building request: %w", err)
	}
	req.Header.Set(apiTokenHeader, c.apiToken)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.logger != nil {
		c.logger.DebugContext(ctx, "signater request", "method", method, "url", fullURL)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("signater: request failed: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("signater: reading response body: %w", err)
	}
	return resp, respBody, nil
}

// retrySleep logs and waits the backoff delay before the next attempt.
func (c *Client) retrySleep(ctx context.Context, attempt int, retryAfter time.Duration, method, cause string) error {
	delay := backoffDelay(attempt, c.retryBaseDelay, retryAfter)
	if c.logger != nil {
		c.logger.WarnContext(ctx, "signater retrying request",
			"method", method, "attempt", attempt+1, "delay", delay, "cause", cause)
	}
	return sleepCtx(ctx, delay)
}
