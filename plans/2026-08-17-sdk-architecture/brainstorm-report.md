# Signater Go SDK — Architecture Brainstorm Report

Date: 2026-08-17
Status: APPROVED by user
Module: `github.com/tunni-sdk/signater-go` (GitHub org `tunni-sdk`, owned by user's personal account)
Language: all code, docs, comments, errors in English. Go 1.23+ minimum. Zero external dependencies (stdlib only).

## Problem Statement

Build a complete, professional Go SDK for the Signater e-signature API (https://docs.api.signater.com), modeled on the best SDKs on the market (stripe-go, openai-go, go-github).

## API Surface (scouted from docs)

- Base URL `https://api.signater.com`, paths `/v1/ecm/...`
- Auth: `x-api-token` header. Sandbox switched by token type (no separate base URL).
- Rate limit: 1,000 req/min; 429 + `Retry-After`.
- Telemetry: every response carries `X-Signater-Telemetry-Operation-Id` header.
- Pagination: page-number style — query `PageSize` (1–100, default 10), `PageNumber` (1-based), `OrderByDirection` (ASC/DESC); response `{items: [], pagination: {totalItems, pageSize, pageNumber, pageItems}}`.
- Errors: 400 `{errors:[{message, metadata}]}`, 401 `{message}`, 402 `{ShouldBuy: Envelopes|ApiEnvelopes|SignerMfaCredits}`, 403 `{message}`, 429.
- Resources (~50 endpoints): Contact (7), Document (4, incl. multipart upload + create-from-template; docs expire in 24h if unused), Envelope (~18, lifecycle: Draft→PublishScheduled→Published→{Hold, Expired, Cancelled, CancelledBySignerMfaError, Rejected, Signed} + trash/restore), Template (8), Vault (13).
- Webhooks: account-level, delivered via Hookdeck (retry: exponential, every 30s up to 5x). 13 `envelope.*` event types. Payload: `{envelope_id, event_type, account_id, env}`; `_by_signer` events add signer id. Headers include `x-signater-apikey`. Auth type default `HookdeckSignature`.
- No published single OpenAPI file, but every endpoint doc page (`<slug>.md`) embeds its OpenAPI fragment → full spec reconstructable. `llms.txt` indexes all pages. Copies in `./reference/`.

## Evaluated Approaches

1. **Handwritten Stripe-style service pattern — CHOSEN.** Full DX control, idiomatic, ~50 endpoints manageable by hand. Matches user's prior art (`~/bunny-sdk-go`).
2. Codegen (oapi-codegen) + manual layer — rejected: DX hostage to Apidog-generated spec quality, two layers to maintain.
3. Stainless-style (openai-go `param.Field`/`option.RequestOption`) — rejected: verbose to maintain by hand without the generator; `param.Opt[T]` needs Go 1.24.
4. AWS/GCP multi-module split — rejected: overkill for 5 resources (YAGNI).

Market research (subagent, 13 dimensions × 5 SDKs): consensus = functional options client, service-struct-per-resource, typed API error with request id + `errors.As`, Go 1.23 `iter.Seq2` auto-pagination (stripe-go new client), auto-retry w/ backoff+jitter, injectable `http.RoundTripper` testing, single-module semver, README + examples/ + full godoc.

## Final Architecture

```
signater-go/
├── go.mod                  # module github.com/tunni-sdk/signater-go, go 1.23
├── README.md  CHANGELOG.md  LICENSE (MIT)
├── signater.go             # Client struct + NewClient(opts ...Option) + service fields
├── options.go              # WithAPIToken, WithBaseURL, WithHTTPClient, WithMaxRetries,
│                           # WithRequestTimeout (per-attempt), WithLogger(*slog.Logger), WithUserAgent
├── request.go              # unexported do(): build/encode/decode, x-api-token, telemetry id capture
├── errors.go               # Error base + typed errors (see below)
├── retry.go                # exponential backoff + jitter, honors Retry-After
├── pagination.go           # ListParams{PageSize, PageNumber, OrderByDirection}, Page[T], iter.Seq2 helpers
├── contact.go              # ContactService + Contact + params
├── document.go             # multipart upload, create-from-template, original/signed file URL
├── envelope.go             # CRUD + lifecycle verbs (Publish, Hold, Cancel, Restore, Move, Rename,
│                           # ChangeOwner, Reinvite, CreateSignatureLink, certificate, downloads)
├── envelope_params.go      # deep create/update params (documents, signers, signmarks, MFA factors)
├── template.go
├── vault.go
├── webhook/
│   └── webhook.go          # ConstructEvent(payload, headers, secret) → typed Event;
│                           # Hookdeck signature verification + x-signater-apikey match; test helper
├── internal/testutil/      # httptest helpers, fixture loading
├── examples/
│   ├── quickstart/main.go  # upload → create envelope → publish → poll status
│   └── webhook-server/main.go
└── .github/workflows/ci.yml  # go test -race, golangci-lint
```

### Key decisions

| Dimension | Decision | Source |
|---|---|---|
| Client | `NewClient(opts ...Option)`; env fallback `SIGNATER_API_TOKEN` | openai-go |
| Services | `client.Envelopes.Create(ctx, params)` style; ctx always first arg | stripe-go / go-github |
| Params typing | pointer fields + `signater.String()/Bool()/Int()` helpers | stripe-go |
| Errors | base `*signater.Error{StatusCode, Message, TelemetryOperationID}`; typed: `ValidationError` (400 errors[]), `AuthenticationError` (401), `PaymentRequiredError` (402, `ShouldBuy`), `PermissionError` (403), `RateLimitError` (429, `RetryAfter`); all `errors.As`-compatible | stripe-go |
| Pagination | `List(ctx, params)` returns raw `Page[T]`; `ListAutoPaging(ctx, params)` returns `iter.Seq2[T, error]` incrementing PageNumber until `pageItems < pageSize` | stripe-go new client |
| Retry | default 3 attempts: network errors, 5xx, 429 (honor Retry-After) on idempotent methods (GET/PUT/DELETE). POST retried **only on 429** (rejected pre-processing) — API has no documented Idempotency-Key, blind POST retry could duplicate envelopes. Configurable. | stripe-go, adapted |
| Timeouts | overall ctx deadline + separate per-attempt `WithRequestTimeout` | openai-go |
| Webhook | `webhook.ConstructEvent` → typed event w/ 13 `envelope.*` constants; Hookdeck HMAC verification (validate exact format empirically in sandbox) | stripe-go webhook pkg |
| Sandbox | zero config (token decides); document it | — |
| Testing | injectable `http.RoundTripper` + `httptest`; fixtures from real response examples embedded in doc pages | openai-go |
| Versioning | single module, semver, v0.x → v1.0.0 at stability | consensus |
| Logging | optional `WithLogger(*slog.Logger)` for request/retry debug | modern idiom |

## Risks / Notes

- Envelope is the heavy resource (~18 endpoints, deep create params: signers → signmarks → MFA factors) — ~50% of effort. Rest is simple CRUD.
- Reconstruct full OpenAPI spec from the 50 endpoint `.md` pages as internal type source-of-truth (not shipped). Index: `reference/llms.txt`.
- Hookdeck signature format must be validated empirically against sandbox before finalizing `webhook` pkg.
- Two "Update vault" endpoints exist (12252702e0 partial-era + 36651167e0 full-replace) — check which is current when implementing.

## Implementation Phases (input for ak:plan)

1. Core: repo bootstrap, client, options, request, errors, retry, pagination + unit tests.
2. Simple resources: Contacts, Vaults, Templates.
3. Documents (multipart) + Envelopes (CRUD + lifecycle).
4. `webhook/` package.
5. examples/, README, CI, doc comments, release v0.1.0.

## Success Metrics

- 100% endpoint coverage (~50) with typed params/responses.
- `go vet` + `golangci-lint` clean, `go test -race` green, high coverage on core (errors/retry/pagination).
- Zero non-stdlib deps. pkg.go.dev renders complete docs. Quickstart example runs against sandbox.
