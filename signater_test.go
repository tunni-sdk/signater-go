package signater

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tunni-sdk/signater-go/internal/testutil"
)

// newTestClient returns a client wired to a scripted transport with fast
// retries and a fixed token.
func newTestClient(t *testing.T, opts ...Option) (*Client, *testutil.Transport) {
	t.Helper()
	tr := &testutil.Transport{}
	base := []Option{
		WithAPIToken("test-token"),
		WithHTTPClient(&http.Client{Transport: tr}),
	}
	c := NewClient(append(base, opts...)...)
	c.retryBaseDelay = time.Millisecond
	return c, tr
}

func TestNewClientDefaults(t *testing.T) {
	t.Setenv(envAPIToken, "")
	c := NewClient()
	if got := c.baseURL.String(); got != defaultBaseURL {
		t.Errorf("baseURL = %q, want %q", got, defaultBaseURL)
	}
	if c.maxRetries != defaultMaxRetries {
		t.Errorf("maxRetries = %d, want %d", c.maxRetries, defaultMaxRetries)
	}
	if c.requestTimeout != defaultRequestTimeout {
		t.Errorf("requestTimeout = %v, want %v", c.requestTimeout, defaultRequestTimeout)
	}
	if want := "signater-go/" + Version; c.userAgent != want {
		t.Errorf("userAgent = %q, want %q", c.userAgent, want)
	}
	if c.apiToken != "" {
		t.Errorf("apiToken = %q, want empty", c.apiToken)
	}
}

func TestNewClientEnvTokenFallback(t *testing.T) {
	t.Setenv(envAPIToken, "env-token")
	if c := NewClient(); c.apiToken != "env-token" {
		t.Errorf("apiToken = %q, want env-token", c.apiToken)
	}
	// Explicit option wins over the environment.
	if c := NewClient(WithAPIToken("explicit")); c.apiToken != "explicit" {
		t.Errorf("apiToken = %q, want explicit", c.apiToken)
	}
}

func TestOptionsApply(t *testing.T) {
	logger := slog.Default()
	hc := &http.Client{}
	c := NewClient(
		WithAPIToken("tok"),
		WithBaseURL("https://sandbox.example.com"),
		WithHTTPClient(hc),
		WithMaxRetries(5),
		WithRequestTimeout(10*time.Second),
		WithLogger(logger),
		WithUserAgent("custom-agent/1.0"),
	)
	if c.baseURL.Host != "sandbox.example.com" {
		t.Errorf("baseURL host = %q", c.baseURL.Host)
	}
	if c.httpClient != hc || c.maxRetries != 5 || c.requestTimeout != 10*time.Second ||
		c.logger != logger || c.userAgent != "custom-agent/1.0" || c.apiToken != "tok" {
		t.Error("options were not all applied")
	}
}

func TestWithBaseURLInvalidSurfacesOnCall(t *testing.T) {
	c := NewClient(WithAPIToken("tok"), WithBaseURL("://bad url"))
	err := c.do(context.Background(), http.MethodGet, "/v1/ecm/contacts", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error from invalid base URL, got nil")
	}
}

func TestDoMissingToken(t *testing.T) {
	t.Setenv(envAPIToken, "")
	c := NewClient()
	err := c.do(context.Background(), http.MethodGet, "/v1/ecm/contacts", nil, nil, nil)
	if !errors.Is(err, errMissingToken) {
		t.Fatalf("err = %v, want errMissingToken", err)
	}
}

func TestDoSendsHeadersAndDecodes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: Header.Get canonicalizes, so this cannot assert the on-wire
		// casing of "x-api-token" (Go sends "X-Api-Token"; HTTP headers are
		// case-insensitive per RFC 9110).
		if got := r.Header.Get("x-api-token"); got != "test-token" {
			t.Errorf("x-api-token = %q", got)
		}
		if got := r.Header.Get("User-Agent"); got != "signater-go/"+Version {
			t.Errorf("User-Agent = %q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"abc-123"}`))
	}))
	defer srv.Close()

	c := NewClient(WithAPIToken("test-token"), WithBaseURL(srv.URL))
	var out struct {
		ID string `json:"id"`
	}
	in := map[string]string{"name": "test"}
	if err := c.do(context.Background(), http.MethodPost, "/v1/ecm/contacts", nil, in, &out); err != nil {
		t.Fatal(err)
	}
	if out.ID != "abc-123" {
		t.Errorf("decoded id = %q, want abc-123", out.ID)
	}
}

func TestDoEmptyResponseBodyWithOut(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusNoContent, "", nil)
	var out struct{ ID string }
	if err := c.do(context.Background(), http.MethodGet, "/v1/ecm/contacts/x", nil, nil, &out); err != nil {
		t.Fatal(err)
	}
	if out.ID != "" {
		t.Errorf("out must stay zero-valued on an empty body, got %+v", out)
	}
}

func TestWithBaseURLRejectsNonAbsolute(t *testing.T) {
	for _, raw := range []string{"api.signater.com", "localhost:8080", "not-a-url", "ftp://x.com"} {
		c := NewClient(WithAPIToken("tok"), WithBaseURL(raw))
		err := c.do(context.Background(), http.MethodGet, "/v1/ecm/contacts", nil, nil, nil)
		if err == nil {
			t.Errorf("base URL %q: expected init error, got nil", raw)
		}
	}
	// A valid URL after an invalid one clears the recorded error.
	c, tr := newTestClient(t, WithBaseURL("not-a-url"), WithBaseURL("https://api.example.com"))
	tr.Enqueue(http.StatusOK, `{}`, nil)
	if err := c.do(context.Background(), http.MethodGet, "/v1/ecm/contacts", nil, nil, nil); err != nil {
		t.Errorf("valid base URL after invalid one must recover, got %v", err)
	}
}

func TestWithHTTPClientNil(t *testing.T) {
	c := NewClient(WithAPIToken("tok"), WithHTTPClient(nil))
	if err := c.do(context.Background(), http.MethodGet, "/v1/ecm/contacts", nil, nil, nil); err == nil {
		t.Error("expected init error for nil http client, got nil")
	}
}

func TestPathEscape(t *testing.T) {
	if got := pathEscape("../../../v1/admin"); got == "../../../v1/admin" {
		t.Error("pathEscape must neutralize traversal sequences")
	}
	if got := pathEscape("a/b"); got != "a%2Fb" {
		t.Errorf("pathEscape(a/b) = %q", got)
	}
}
