---
title: "Phase 1: Core Foundation"
status: todo
priority: P1
effort: "1.5d"
dependencies: []
---

# Phase 1: Core Foundation

## Overview

Bootstrap the repo and build the transport core every service depends on: client + functional options, HTTP request pipeline, typed error hierarchy, retry with backoff, and generic pagination with Go 1.23 iterators.

## Requirements

- Functional: `NewClient()` works with `SIGNATER_API_TOKEN` env fallback; `do()` sends `x-api-token`, decodes JSON, maps non-2xx to typed errors carrying `X-Signater-Telemetry-Operation-Id`.
- Non-functional: stdlib only; Go 1.23; every exported identifier has a doc comment; `go test -race` green.

## Architecture

```go
// signater.go
type Client struct {
    Contacts  *ContactService
    Documents *DocumentService
    Envelopes *EnvelopeService
    Templates *TemplateService
    Vaults    *VaultService
    // unexported: httpClient *http.Client, baseURL *url.URL, apiToken string,
    // maxRetries int, requestTimeout time.Duration, logger *slog.Logger, userAgent string
}
func NewClient(opts ...Option) *Client   // token: WithAPIToken > SIGNATER_API_TOKEN env
type Option func(*Client)
// options.go: WithAPIToken, WithBaseURL, WithHTTPClient, WithMaxRetries,
//             WithRequestTimeout, WithLogger, WithUserAgent

// errors.go
type Error struct { StatusCode int; Message string; TelemetryOperationID string; RawBody []byte }
func (e *Error) Error() string
type ValidationError struct { Error; Errors []FieldError }       // 400: {errors:[{message, metadata}]}
type FieldError struct { Message string; Metadata json.RawMessage }
type AuthenticationError struct{ Error }                          // 401: {message}
type PaymentRequiredError struct { Error; ShouldBuy string }      // 402: {ShouldBuy}; constants ShouldBuyEnvelopes|ApiEnvelopes|SignerMfaCredits
type PermissionError struct{ Error }                              // 403
type NotFoundError struct{ Error }                                // 404 (documented per endpoint)
type RateLimitError struct { Error; RetryAfter time.Duration }    // 429 + Retry-After header

// pagination.go
type ListParams struct { PageSize int; PageNumber int; OrderByDirection string } // defaults 10 / 1 / DESC; consts OrderAsc/OrderDesc
type Pagination struct { TotalItems, PageSize, PageNumber, PageItems int }
type Page[T any] struct { Items []T; Pagination Pagination }
// per-service: List returns Page[T]; ListAutoPaging returns iter.Seq2[T, error]
// generic helper: func autoPaging[T any](ctx, fetch func(page int) (*Page[T], error)) iter.Seq2[T, error]
// stop condition: pageItems < pageSize OR pageNumber*pageSize >= totalItems

// retry.go — policy (from brainstorm, API has NO Idempotency-Key support):
//   retryable: network errors, 5xx, 429 → GET/PUT/DELETE/HEAD only
//   POST retried ONLY on 429 (rejected before processing)
//   default 3 attempts; exponential backoff (500ms base, x2, full jitter, cap 8s); honor Retry-After
```

`request.go` implements `(c *Client) do(ctx, method, path string, query url.Values, body, out any) error` plus a multipart variant hook for Phase 3. Per-attempt timeout via `requestTimeout` (separate context wrap) while ctx governs the whole call.

## Related Code Files

- Create: `go.mod`, `LICENSE` (MIT), `.gitignore`, `signater.go`, `options.go`, `request.go`, `errors.go`, `retry.go`, `pagination.go`
- Create tests: `signater_test.go`, `errors_test.go`, `retry_test.go`, `pagination_test.go`, `internal/testutil/testutil.go`

## Implementation Steps

1. `git init`; write `go.mod` (`module github.com/tunni-sdk/signater-go`, `go 1.23`), MIT LICENSE, `.gitignore` (coverage files, `.idea`, etc.).
2. Implement `errors.go` with the hierarchy above; a single `newAPIError(resp *http.Response, body []byte) error` factory switching on status code. All wrap/embed base `Error`; support `errors.As` for both base and concrete types.
3. Implement `signater.go` + `options.go`; base URL default `https://api.signater.com`; `User-Agent: signater-go/<version>`; version const in `signater.go`.
4. Implement `request.go` (`do()`): build URL, JSON encode body when non-nil, set headers, execute with retry loop from `retry.go`, capture telemetry header, decode into `out` when non-nil; treat 204/empty body.
5. Implement `retry.go` per policy; unit-test with a counting `RoundTripper`: 429+Retry-After honored, 500 GET retried, 500 POST NOT retried, 429 POST retried, context cancellation aborts.
6. Implement `pagination.go` with generic `Page[T]` + `autoPaging` iterator; test multi-page, early-break, error propagation mid-iteration.
7. `internal/testutil`: `NewServer(t, handler)` returning a `*Client` pointed at `httptest.Server`; fixture loader.
8. Run `go vet ./...`, `go test -race ./...`.

## Success Criteria

- [x] `NewClient(WithAPIToken("x"))` and env-var fallback both work (tested)
- [x] Each status code 400/401/402/403/404/429 maps to its typed error with telemetry id populated (tested against doc example payloads)
- [x] Retry policy behaves exactly per table above (tested)
- [x] `autoPaging` iterates across pages and stops correctly (tested)
- [x] `go vet` + `go test -race` clean; no external deps

## Risk Assessment

- 402 payload uses PascalCase `ShouldBuy` while other payloads are camelCase — decode with exact struct tags per doc example, do not assume a global naming convention.
- Retry-After may be absent on 429 — fall back to computed backoff.

## Implementation Notes (as built)

Accepted deviations from the sketch above, plus hardening applied after code review:

- Base error type is `APIError` (not `Error`) to avoid colliding with the `error` interface; all concrete types embed it and implement `Unwrap()` so `errors.As` matches both.
- `autoPaging(firstPage int, fetch)` — no ctx param; cancellation flows through the fetch closure. Stop conditions: nil page, empty page, missing pagination metadata (`PageSize <= 0`), short page, or `page*PageSize >= TotalItems` (saves the extra empty fetch on exact boundaries).
- `internal/testutil` is a scripted `Transport` (RoundTripper), not `NewServer(t, handler)` — the planned helper would import-cycle with in-package tests. It honors pre-cancelled request contexts and fails loudly on queue exhaustion.
- Retry hardening: `Retry-After` honored for any retryable status (not only 429), parsed in both seconds and HTTP-date forms, capped at 60s; 5xx retries narrowed to 500/502/503/504; backoff shift guarded against overflow.
- `WithBaseURL` requires absolute http(s) URLs (error surfaces on first call via `initErr`; a later valid URL clears it); `WithHTTPClient(nil)` is an error.
- `pathEscape` helper added now so Phases 2-3 escape every dynamic path segment (traversal guard).
- Deferred (revisit before v1.0): expose the telemetry operation id on successful responses; empirically confirm the gateway accepts canonicalized `X-Api-Token` casing (Phase 5 sandbox validation).
