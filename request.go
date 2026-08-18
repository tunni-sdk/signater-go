package signater

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// errMissingToken is returned when no API token was configured.
var errMissingToken = errors.New("signater: missing API token: use WithAPIToken or set " + envAPIToken)

// errEmptyResourceID is returned when a required id would produce an empty
// path segment, which would silently target a different endpoint.
var errEmptyResourceID = errors.New("signater: empty resource id in request path")

// pathEscape escapes a dynamic value (e.g. a resource id) for safe use as a
// single URL path segment when building API paths. All-dot segments are
// percent-encoded so URL path cleaning cannot collapse the path onto a
// different endpoint.
func pathEscape(s string) string {
	e := url.PathEscape(s)
	if e != "" && strings.Trim(e, ".") == "" {
		e = strings.ReplaceAll(e, ".", "%2E")
	}
	return e
}

// validatePath rejects paths where a dynamic segment came in empty (trailing
// slash or double slash), which would target a different route.
func validatePath(path string) error {
	if strings.HasSuffix(path, "/") || strings.Contains(path, "//") {
		return errEmptyResourceID
	}
	return nil
}

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
	if err := validatePath(path); err != nil {
		return err
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
		resp, respBody, err := c.attempt(ctx, method, u.String(), payload) //nolint:bodyclose // attempt fully reads and closes the body
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

// doMultipart executes a multipart/form-data upload with a single attempt (the
// body stream cannot be replayed, so the retry policy does not apply).
func (c *Client) doMultipart(ctx context.Context, method, path, fieldName, fileName string, r io.Reader, out any) error {
	if c.initErr != nil {
		return c.initErr
	}
	if c.apiToken == "" {
		return errMissingToken
	}
	if err := validatePath(path); err != nil {
		return err
	}

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)
		part, err := mw.CreateFormFile(fieldName, fileName)
		if err == nil {
			_, err = io.Copy(part, r)
		}
		if err == nil {
			err = mw.Close()
		}
		pw.CloseWithError(err)
	}()
	// The caller's reader must not be touched after this function returns.
	defer func() {
		pr.CloseWithError(errors.New("signater: upload aborted"))
		<-writerDone
	}()

	attemptCtx := ctx
	if c.requestTimeout > 0 {
		var cancel context.CancelFunc
		attemptCtx, cancel = context.WithTimeout(ctx, c.requestTimeout)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(attemptCtx, method, c.baseURL.JoinPath(path).String(), pr)
	if err != nil {
		return fmt.Errorf("signater: building request: %w", err)
	}
	req.Header.Set(apiTokenHeader, c.apiToken)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", mw.FormDataContentType())

	if c.logger != nil {
		c.logger.DebugContext(ctx, "signater request", "method", method, "url", req.URL.String())
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("signater: request failed: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("signater: reading response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp, respBody)
	}
	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("signater: decoding response: %w", err)
		}
	}
	return nil
}

// doRedirectURL performs a GET expected to answer with a redirect and returns
// the resolved Location target (a pre-signed URL) without following it. It
// applies the client's normal GET retry policy.
func (c *Client) doRedirectURL(ctx context.Context, path string) (string, error) {
	if c.initErr != nil {
		return "", c.initErr
	}
	if c.apiToken == "" {
		return "", errMissingToken
	}
	if err := validatePath(path); err != nil {
		return "", err
	}

	// Clone the client so redirects are surfaced instead of followed.
	hc := *c.httpClient
	hc.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	fullURL := c.baseURL.JoinPath(path).String()

	for attempt := 0; ; attempt++ {
		resp, respBody, err := c.redirectAttempt(ctx, &hc, fullURL) //nolint:bodyclose // redirectAttempt fully reads and closes the body
		if err != nil {
			if ctx.Err() == nil && attempt < c.maxRetries {
				if serr := c.retrySleep(ctx, attempt, 0, http.MethodGet, err.Error()); serr == nil {
					continue
				}
			}
			return "", err
		}
		switch {
		case resp.StatusCode >= 300 && resp.StatusCode < 400:
			loc := resp.Header.Get("Location")
			if loc == "" {
				return "", fmt.Errorf("signater: HTTP %d redirect without Location header", resp.StatusCode)
			}
			// Location may be relative per RFC 9110; resolve it.
			if u, perr := resp.Request.URL.Parse(loc); perr == nil {
				return u.String(), nil
			}
			return loc, nil
		case resp.StatusCode >= 200 && resp.StatusCode < 300:
			return "", fmt.Errorf("signater: expected a redirect, got HTTP %d", resp.StatusCode)
		default:
			apiErr := newAPIError(resp, respBody)
			if attempt < c.maxRetries && shouldRetryStatus(http.MethodGet, resp.StatusCode) {
				if serr := c.retrySleep(ctx, attempt, retryAfterFrom(resp.Header), http.MethodGet, apiErr.Error()); serr == nil {
					continue
				}
			}
			return "", apiErr
		}
	}
}

// redirectAttempt performs one round trip of doRedirectURL.
func (c *Client) redirectAttempt(ctx context.Context, hc *http.Client, fullURL string) (*http.Response, []byte, error) {
	attemptCtx := ctx
	if c.requestTimeout > 0 {
		var cancel context.CancelFunc
		attemptCtx, cancel = context.WithTimeout(ctx, c.requestTimeout)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(attemptCtx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("signater: building request: %w", err)
	}
	req.Header.Set(apiTokenHeader, c.apiToken)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := hc.Do(req)
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
