package signater

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// Option configures a Client. Options are applied by NewClient in order.
type Option func(*Client)

// WithAPIToken sets the API token sent in the x-api-token header. It takes
// precedence over the SIGNATER_API_TOKEN environment variable.
func WithAPIToken(token string) Option {
	return func(c *Client) { c.apiToken = token }
}

// WithBaseURL overrides the API base URL (default https://api.signater.com).
// The URL must be absolute with an http or https scheme; anything else is
// reported as an error by the first API call.
func WithBaseURL(rawURL string) Option {
	return func(c *Client) {
		u, err := url.Parse(rawURL)
		if err != nil {
			c.initErr = fmt.Errorf("signater: invalid base URL %q: %w", rawURL, err)
			return
		}
		if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			c.initErr = fmt.Errorf("signater: base URL %q must be an absolute http(s) URL", rawURL)
			return
		}
		c.baseURL = u
		c.initErr = nil
	}
}

// WithHTTPClient replaces the underlying *http.Client, e.g. to install a
// custom Transport for proxying, instrumentation, or testing. A nil client is
// reported as an error by the first API call.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient == nil {
			c.initErr = fmt.Errorf("signater: WithHTTPClient requires a non-nil *http.Client")
			return
		}
		c.httpClient = httpClient
	}
}

// WithMaxRetries sets how many times a failed request is retried after the
// first attempt (default 2). See the package documentation for the retry
// policy. Negative values are treated as zero.
func WithMaxRetries(n int) Option {
	return func(c *Client) { c.maxRetries = max(n, 0) }
}

// WithRequestTimeout sets the timeout applied to each individual attempt
// (default 30s), independent of the context deadline that governs the whole
// call including retries. Zero disables the per-attempt timeout.
func WithRequestTimeout(d time.Duration) Option {
	return func(c *Client) { c.requestTimeout = max(d, 0) }
}

// WithLogger enables debug logging of requests and retries through the given
// slog logger. By default nothing is logged. Debug entries include full
// request URLs, so query parameters (e.g. Search filters) may surface
// personal data in logs.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) { c.logger = logger }
}

// WithUserAgent overrides the User-Agent header (default "signater-go/<version>").
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}
