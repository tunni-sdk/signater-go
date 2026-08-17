package signater

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Version is the version of this SDK, sent in the User-Agent header.
const Version = "0.1.0"

const (
	defaultBaseURL        = "https://api.signater.com"
	defaultMaxRetries     = 2
	defaultRequestTimeout = 30 * time.Second
	defaultRetryBaseDelay = 500 * time.Millisecond

	apiTokenHeader = "x-api-token"
	// envAPIToken is the environment variable read when WithAPIToken is not used.
	envAPIToken = "SIGNATER_API_TOKEN"
)

// Client is the entry point to the Signater API. Create one with NewClient.
// A Client is safe for concurrent use by multiple goroutines.
type Client struct {
	// Contacts accesses the Contact endpoints.
	Contacts *ContactService
	// Documents accesses the Document endpoints.
	Documents *DocumentService
	// Envelopes accesses the Envelope endpoints.
	Envelopes *EnvelopeService
	// Templates accesses the Template endpoints.
	Templates *TemplateService
	// Vaults accesses the Vault endpoints.
	Vaults *VaultService

	httpClient     *http.Client
	baseURL        *url.URL
	apiToken       string
	maxRetries     int
	requestTimeout time.Duration
	retryBaseDelay time.Duration
	logger         *slog.Logger
	userAgent      string

	// initErr records a configuration error (e.g. invalid base URL) and is
	// returned by every API call so misconfiguration cannot pass silently.
	initErr error
}

// NewClient creates a Client.
//
// The API token is taken from WithAPIToken or, when absent, from the
// SIGNATER_API_TOKEN environment variable. Defaults: base URL
// https://api.signater.com, 2 retries after the first attempt, 30s per-attempt
// timeout.
func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient:     &http.Client{},
		maxRetries:     defaultMaxRetries,
		requestTimeout: defaultRequestTimeout,
		retryBaseDelay: defaultRetryBaseDelay,
		userAgent:      "signater-go/" + Version,
	}
	c.baseURL, _ = url.Parse(defaultBaseURL)
	for _, opt := range opts {
		opt(c)
	}
	if c.apiToken == "" {
		c.apiToken = os.Getenv(envAPIToken)
	}
	c.Contacts = &ContactService{client: c}
	c.Documents = &DocumentService{client: c}
	c.Envelopes = &EnvelopeService{client: c}
	c.Templates = &TemplateService{client: c}
	c.Vaults = &VaultService{client: c}
	return c
}
