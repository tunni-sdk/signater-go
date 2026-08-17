// Package testutil provides shared helpers for signater-go tests.
package testutil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// RecordedRequest captures the relevant parts of a request seen by Transport.
type RecordedRequest struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   []byte
}

// Step is one scripted round trip: either a response or a transport error.
type Step struct {
	Status int
	Header http.Header
	Body   string
	// Err, when non-nil, makes RoundTrip fail with this transport-level error
	// instead of returning a response.
	Err error
}

// Transport is an http.RoundTripper that records incoming requests and replays
// a scripted queue of steps. The zero value is ready to use; an exhausted
// queue fails the request loudly so over-fetching cannot pass silently.
type Transport struct {
	mu       sync.Mutex
	queue    []Step
	requests []RecordedRequest
}

// Enqueue appends a scripted response.
func (t *Transport) Enqueue(status int, body string, header http.Header) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.queue = append(t.queue, Step{Status: status, Body: body, Header: header})
}

// EnqueueError appends a scripted transport-level failure.
func (t *Transport) EnqueueError(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.queue = append(t.queue, Step{Err: err})
}

// Count returns how many requests were seen.
func (t *Transport) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.requests)
}

// Snapshot returns a copy of the requests seen so far.
func (t *Transport) Snapshot() []RecordedRequest {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]RecordedRequest(nil), t.requests...)
}

// RoundTrip implements http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Mirror net/http: a request whose context is already done never reaches
	// the wire.
	if err := req.Context().Err(); err != nil {
		return nil, err
	}

	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		req.Body.Close()
	}

	t.mu.Lock()
	t.requests = append(t.requests, RecordedRequest{
		Method: req.Method,
		Path:   req.URL.Path,
		Query:  req.URL.RawQuery,
		Header: req.Header.Clone(),
		Body:   body,
	})
	n := len(t.requests)
	if len(t.queue) == 0 {
		t.mu.Unlock()
		return nil, fmt.Errorf("testutil: unexpected request #%d: %s %s (response queue exhausted)", n, req.Method, req.URL)
	}
	step := t.queue[0]
	t.queue = t.queue[1:]
	t.mu.Unlock()

	if step.Err != nil {
		return nil, step.Err
	}
	header := step.Header
	if header == nil {
		header = http.Header{}
	}
	return &http.Response{
		StatusCode: step.Status,
		Header:     header,
		Body:       io.NopCloser(bytes.NewReader([]byte(step.Body))),
		Request:    req,
	}, nil
}
