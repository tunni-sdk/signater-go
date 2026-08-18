# Code Review — Webhook Empirical Confirmation Pass

Date: 2026-08-18
Scope: uncommitted changes (`git diff HEAD`), advisory only, no files modified.

## Scope

- Files: `webhook/webhook.go`, `webhook/webhook_test.go`, `vault.go`, `envelope.go`,
  `examples/webhook-server/main.go`, `README.md`, `CHANGELOG.md` (7 files, +100/-39)
- No scope drift: changed set matches the stated task exactly. `vault.go` / `envelope.go`
  are comment-only (verified in diff); no code lines touched there.
- Gates run locally: `gofmt -l .` clean, `go vet ./...` clean, `golangci-lint run` → **0 issues**,
  `go test ./...` pass, `webhook` package coverage **95.3%**.
  `go test -race` could not run in this sandbox (`gcc` absent, CGO required); CI runners have it.
  No race risk introduced: `secret`/`apiKey` are read-only strings captured by the handler closure.

## Overall Assessment

The code change is small and correct. `RequestID` is strictly additive and backward compatible,
the example's branch logic has no nil-event or swallowed-error path, and no exported signature or
type changed. The real risk in this diff is **documentation confidence exceeding the evidence**:
the capture confirmed header *presence and shape*, not the HMAC *construction* (which is
unverifiable without a secret), yet the diff removed the honest uncertainty note and promoted a
weaker, bearer-style verification to the documented default without a matching security caveat.

No blocking defects. Two documentation findings should be fixed before landing.

## Critical Issues

None.

## High Priority

### H1 — HMAC construction is asserted as confirmed, but it cannot have been verified (major)

`webhook/webhook.go:7-8`:
```go
// Deliveries carry two authentication headers: X-Hookdeck-Signature (base64
// HMAC-SHA256 of the raw body) and X-Signater-Apikey ...
```
`README.md:128` places the same claim under the section *"encodes behavior verified against the
live sandbox"*.

Evidence it is not verified:
- Signater exposes no signing secret, so `VerifyHookdeckSignature` / `ConstructEvent` were never
  exercised end-to-end against a real delivery. Only the header name and a base64 value were
  observed. Whether the MAC is over the raw body (vs. `timestamp.body`, a webhook-id prefix, or a
  different key) is still inference from Hookdeck's public docs.
- The diff itself is inconsistent: `ConstructEvent`'s godoc (`webhook/webhook.go:144-148`) hedges
  correctly ("matching Hookdeck's published HMAC-SHA256 scheme"), while the package doc and README
  state it flatly.
- `webhook/webhook.go:203-206` **deleted** the accurate qualifier from `sign`'s godoc
  ("the scheme is still VERIFY(sandbox)-unconfirmed"). The uncertainty it described was not
  resolved by this capture; only the co-located signer-id and header-name unknowns were.
- Also verify before shipping: does the doc's "a single base64-encoded **32-byte** value"
  (`webhook/webhook.go:145`) come from actually base64-decoding a captured header, or from
  assuming SHA-256 output length?

Impact: this SDK's credibility rests on its VERIFY discipline — an unverified claim laundered into
a "verified against the live sandbox" list is exactly the failure mode that discipline exists to
prevent. Functionally low risk (signature verification fails closed if the scheme differs), but the
README tells users the opposite of the truth about what was tested.

Suggested wording:
```go
// Captured sandbox deliveries confirm the header name and value shape (a single
// base64 value in X-Hookdeck-Signature). The HMAC construction itself follows
// Hookdeck's published scheme and remains unverified end-to-end: Signater does
// not expose the signing secret, so there is nothing to check it against.
// Verify deliveries with VerifyAPIKey instead.
```
and in README, move the signature line out of the "verified" list or mark the HMAC construction as
inferred.

### H2 — Promoting VerifyAPIKey to the default lacks the matching security caveat (major, docs)

`README.md:102` now recommends `VerifyAPIKey` as the primary path with no weakness note, and
`examples/webhook-server/main.go` defaults to it while serving plain HTTP (`main.go:74`, no TLS).

Two properties change materially versus the HMAC path, and only one is documented (in
`VerifyAPIKey`'s godoc, `webhook/webhook.go:185-191`, which README readers will not see):
1. **The shared secret now travels on the wire on every delivery.** With HMAC only the digest
   crossed the network. A static key in a header leaks to any TLS-terminating proxy, request/access
   log, APM header capture, or a plaintext listener.
2. **It does not bind the payload.** Anyone who observes a single delivery (or the value in a log)
   can forge arbitrary events for that account indefinitely — there is no timestamp or body binding
   and no rotation story documented.

The guidance "pin its value from a captured delivery" actively directs users to lift a live
credential out of a capture, which increases the odds it lands in a log file or a commit.

Recommended additions (docs only, no code change):
- README + package doc: "Treat the x-signater-apikey value as a credential equal to the API token:
  serve the endpoint over TLS, never log request headers, keep it in a secret store."
- One line noting the check authenticates the *sender*, not the *payload* — hence the existing
  "treat events as hints and query the API" guidance is mandatory, not optional, in this mode.

## Medium Priority

### M1 — Example log label is now wrong (`examples/webhook-server/main.go:61-62`, `:54`)
```go
log.Printf("event %q envelope=%q env=%q request-id=%q",
    event.Type, event.EnvelopeID, event.Env, webhook.RequestID(r.Header))
```
In practice this prints an `evt_...` Hookdeck event id under a `request-id=` label, contradicting
the README's own statement that deliveries carry no `request-id` header. Rename to
`delivery-id=` (line 61) — the example is the artifact users copy.

### M2 — Empirical captures are not locked in as fixtures (test gap)
Six real deliveries were replayed manually, but nothing in `webhook/` commits that evidence:
there is no `testdata/`, and the new assertions use synthetic values. The `signer_id` key is
covered by the pre-existing `webhook_test.go:33-38`, but the confirmed *header set* is not.
Suggest a redacted golden fixture (payload + header map from a real capture) asserting
`ParseEvent` → `SignerID`/`Type` and `RequestID` → the `evt_` id. That is what makes the
"confirmed against captured sandbox deliveries" comments falsifiable by CI later.

### M3 — Dual-mode precedence is silent (`examples/webhook-server/main.go:37-52`)
Logic is correct, but if a user sets both env vars and the secret is stale, every delivery is
rejected while a valid api key sits unused. Log the active mode once at startup
(`log.Printf("verifying deliveries with %s", mode)`) so the failure is diagnosable.

## Low Priority

- `examples/webhook-server/main.go:50` — `else if err = webhook.VerifyAPIKey(...); err == nil` is
  correct Go and correctly propagates the verification error to the shared `if err != nil` handler,
  but assigning the outer `err` inside an `else if` condition is fragile under future edits. An
  explicit block or a small `verify(r) (*webhook.Event, error)` helper reads better. Not a defect.
- `README.md:105-108` — the snippet still drops the `io.ReadAll` error (pre-existing pattern,
  preserved through the rewrite). The package doc snippet handles it; align the two.
- `webhook/webhook.go:88-92` — `RequestID` now also falls back when `Request-Id` is present but
  empty. Harmless (arguably better), undocumented. Optional one-word doc fix ("absent or empty").
- Deliveries also carry `Idempotency-Key`, which the docs/CHANGELOG never mention. A half-sentence
  on why `X-Hookdeck-Eventid` was chosen over it would prevent a future reviewer re-litigating it.
- `plans/260817-1319-signater-go-sdk-v010/handoff.md:18` still says the `VERIFY(sandbox)` markers
  are the remaining work. Plan hygiene, outside the diff — for the lead, not this review.

## Checklist Verification

| Item | Result |
|---|---|
| 1. `RequestID` backward compatible | **Pass.** Signature unchanged; `Request-Id` wins when non-empty (`webhook.go:94-96`); fallback only fires when it is absent/empty. Doc comment accurate. |
| 2. Example branch logic | **Pass.** `secret != ""` → `ConstructEvent`; otherwise the `else if` always evaluates, so exactly one branch runs. `ConstructEvent`/`ParseEvent` return non-nil `*Event` iff `err == nil` (`webhook.go:128-154`), so no nil-event path reaches line 62. Verification errors are never swallowed: both branches feed the shared `if err != nil` reject at line 53. Error text is logged server-side only, never echoed to the caller. |
| 3. Public contract | **Pass.** No exported func/type/const/field signature changed; `Event`, `RequestID`, `ConstructEvent`, `ParseEvent`, `VerifyAPIKey`, `VerifyHookdeckSignature`, `SignPayload` all identical. `vault.go`/`envelope.go` comment-only. Change is additive → next minor version, correctly filed under `Unreleased / Changed`. |
| 4. Doc claims vs. facts | **Partial** — see H1. `signer_id` (tag at `webhook.go:83`), the absent `request-id` header, the vault paths (`vault.go:201`, `:216`), and `ProcessCertificate`'s `{}` body all match the stated evidence and the code. No leftover claim anywhere in shipping files that the secret is obtainable from the webhook page (grepped `.go`/`.md` excluding `plans/`). No `VERIFY(` markers remain outside `plans/`. |
| 5. Lint/vet/build gates | **Pass.** gofmt clean, `go vet ./...` clean, golangci-lint v2 reports 0 issues, `go test ./...` green. Nothing in the diff uses post-1.23 APIs. `-race` unrunnable locally (no gcc), not a code concern. |
| 6. Tests | **Pass with a gap.** The two additions are meaningful: `webhook_test.go:130-133` proves the fallback, `:136-139` proves precedence, and `h` is correctly reset to a fresh `http.Header` at line 129 so the cases are independent. Header canonicalization is exercised via `Set`, matching what `net/http` produces on real requests. Gap: no fixture from the real captures (M2). |

## Recommended Actions

1. Fix H1: restore honest uncertainty about the HMAC construction in the package doc, `sign`'s
   godoc, and README's "verified" list. Confirm whether the "32-byte" claim was measured.
2. Fix H2: add the TLS / no-header-logging / credential-handling caveat to README and package doc.
3. Fix M1: rename the example's `request-id=` log label to `delivery-id=`.
4. Consider M2: commit one redacted capture as a golden fixture.
5. Optional: M3 startup mode log, and the low-priority doc alignments.

## Unresolved Questions

1. Was `X-Signater-Apikey` empirically compared against the account's API token to justify the
   "distinct from the API token" claim in README/godoc, or is that inferred from format?
2. Was the `X-Hookdeck-Signature` value base64-decoded and measured at 32 bytes, or assumed?
3. Is the api key rotatable from the Signater UI? "Pin its value from a captured delivery" implies
   a credential with no documented rotation path — worth one line if rotation exists.
