# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- `webhook.RequestID` falls back to the `X-Hookdeck-Eventid` header: captured
  sandbox deliveries carry no `request-id` header, and the Hookdeck event id
  is stable across retry attempts, making it the practical dedupe key.
- Webhook docs and the `webhook-server` example now verify deliveries with
  `VerifyAPIKey` by default: the Signater webhook page does not currently
  expose the Hookdeck signing secret, so the static `x-signater-apikey`
  header is the strongest check available. Sandbox captures confirmed the
  `X-Hookdeck-Signature` header name/format and the `signer_id` payload key
  on `*_by_signer` events.

## [0.1.0] - 2026-08-17

### Added

- Complete coverage of the Signater API's 50 published endpoints across five
  services: `Contacts`, `Documents`, `Envelopes`, `Templates`, `Vaults`.
- Typed error hierarchy (`ValidationError`, `AuthenticationError`,
  `PaymentRequiredError`, `PermissionError`, `NotFoundError`,
  `RateLimitError`) carrying the telemetry operation id, all matchable via
  `errors.As` against the concrete type or `*APIError`.
- Automatic retries with exponential backoff and jitter; `Retry-After`
  support (seconds and HTTP-date); POST retried only on 429.
- Go 1.23 `iter.Seq2` auto-pagination on every list endpoint.
- Streaming multipart document upload; pre-signed URL capture for
  redirect-based downloads.
- `webhook` package: typed events for the 13 published event types, Hookdeck
  HMAC signature verification (fail-closed), api-key check, and a
  `SignPayload` helper for consumer tests.
- Functional client options: `WithAPIToken`, `WithBaseURL`, `WithHTTPClient`,
  `WithMaxRetries`, `WithRequestTimeout`, `WithLogger`, `WithUserAgent`.

[Unreleased]: https://github.com/tunni-sdk/signater-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/tunni-sdk/signater-go/releases/tag/v0.1.0
