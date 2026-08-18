# Post-Release Cleanup: CI Lint Breakage, Empirical Webhook Capture, and a Caught Truthfulness Bug

**Date**: 2026-08-18 10:04–10:42
**Severity**: Medium
**Component**: CI pipeline (`golangci-lint`), `webhook` package, `vault.go`/`envelope.go` docs, release/repo visibility
**Status**: Resolved (one item still open: sandbox token rotation)

## What Happened

Day two after tagging v0.1.0. Four commits closed out everything left in the handoff except one: `4c9b18b` (CI lint fix), `2ec023c` (drop a stray backup file), `36f13a3` (webhook empirical findings encoded into the SDK), `24ae83f` (plan/handoff bookkeeping).

The first GitHub Actions run on the new public repo failed outright. `golangci/golangci-lint-action@v6` pinned a `golangci-lint` binary built with go1.24, run against a go1.26 toolchain on the Actions runner. The version skew produced garbage: false "missing type in composite literal" errors and stdlib typecheck failures that had nothing to do with the actual code. Fixed by moving to `golangci-lint-action@v8`, which pulls golangci-lint v2, plus a full `.golangci.yml` migration (`version: "2"`, restructured `exclusions.presets`/`rules` instead of the old `issues.exclude-rules`).

With a working linter, it found 14 real findings, all fixed in `4c9b18b`:
- `contact_test.go` / `signater_test.go`: `err != errEmptyResourceID` / `err != errMissingToken` → `errors.Is(...)` (errorlint).
- `pagination_test.go`: two unused `page` params → `_` (revive).
- `bodyclose` false positives on `request.go`'s `attempt`/`redirectAttempt` (both close the body internally) → `//nolint:bodyclose` with justification comments.
- `gosec` G101/G404 false positives on `envAPIToken`, `apiKeyHeader` (header/env *names*, not secrets) and `backoffDelay`'s `rand.N` (retry jitter, no security requirement) → scoped `//nolint`.
- `examples/webhook-server/main.go`: switched `%s` to `%q` in log lines to escape control characters in untrusted webhook fields, added `ReadHeaderTimeout: 10 * time.Second` to the `http.Server`.
- gosec's G706 doesn't model `%q` as sanitization, so it kept flagging the examples logging — handled with a scoped config exclusion (`path: ^examples/`) rather than inline nolints, since it applied to every log line in that file.

Separately, the last two `VERIFY(sandbox)` markers in the codebase (webhook signature scheme, `signer_id` payload key) got closed with real captured data. The dashboard's documented "CLI Settings" webhook flow no longer exists in the current Signater UI, so the capture used a regular webhook pointed at an SSH tunnel (localhost.run) instead. First attempt used a personal `hookdeck listen` tunnel, which silently overwrote Signater's original `X-Hookdeck-*` headers with the tunnel's own — that capture round was misleading and had to be discarded once noticed. The direct SSH tunnel preserved the raw headers Signater actually sent.

Findings encoded into the SDK (`36f13a3`):
- `signer_id` is confirmed as the payload key on `*_by_signer` events.
- Real deliveries carry `X-Hookdeck-Signature` (base64, 32 bytes) but the HMAC scheme itself can't be verified end-to-end — Signater's webhook config page exposes no signing secret. `VerifyAPIKey` (static `X-Signater-Apikey` header check) is now the documented default, with an explicit trade-off note: the key doesn't cover the payload and travels in cleartext on every delivery.
- No `request-id` header exists on real deliveries. `webhook.RequestID` now falls back to `X-Hookdeck-Eventid`, which is stable across Hookdeck retries.
- `Vaults.Owners`/`Members` re-confirmed as account-wide (no vault id in path); `Envelopes.ProcessCertificate` re-confirmed to accept `{}`.
- All 6 captured real deliveries replay cleanly through the SDK.

A code-reviewer pass on this diff (`plans/reports/code-reviewer-260818-1026-webhook-empirical.md`) caught something worth remembering: the package doc and README initially stated the HMAC scheme flatly as fact under a section titled "verified against the live sandbox," while `ConstructEvent`'s own godoc still correctly hedged it as "matching Hookdeck's published scheme." Only the header *name and shape* (base64, 32 bytes) were actually observed — the MAC construction (raw body vs. `timestamp.body`, key derivation, etc.) was never exercised because no signing secret is available to verify against. The docs were corrected to separate observed fact from Hookdeck-documented assumption before landing.

Post-release, the repo was flipped from private to public (user decision), pkg.go.dev indexed it, `go get github.com/tunni-sdk/signater-go@v0.1.0` resolves cleanly from a fresh module, and a GitHub Release was cut from the CHANGELOG. Go Report Card is gone — the badge just says "retired" — so that check is dropped from the release process entirely.

## The Brutal Truth

The CI failure on the very first run of a freshly public repo is not a great look, and it wasn't even our bug — it was `golangci-lint-action@v6` resolving a lint binary built for a Go version older than the runner's toolchain. Wasted time chasing "missing type in composite literal" errors that were pure noise before realizing the tool itself was broken, not the code.

The webhook capture redo is the more expensive lesson. Trusting `hookdeck listen` to hand back "the real headers" without checking whether the tunnel itself mutates them wasted a full capture cycle. If the review hadn't caught the doc overclaim on the HMAC scheme, this SDK would have shipped a README claiming "verified against the live sandbox" for something that was never actually verified — exactly the kind of confidently wrong documentation that erodes trust in an SDK once someone tries to build a signature verifier against it and it doesn't match. Good thing the VERIFY(sandbox) discipline extended to the review step and not just the initial claim.

## Technical Details

- `.golangci.yml` migrated to schema v2: `version: "2"`, `linters.exclusions.presets: [comments, common-false-positives, legacy, std-error-handling]`, plus two scoped `exclusions.rules` (gosec G706 on `^examples/`, bodyclose/errcheck/gosec on `_test\.go`).
- `.github/workflows/ci.yml`: `golangci-lint-action@v6` → `@v8`.
- `webhook.RequestID`: `h.Get("Request-Id")` with fallback to `h.Get("X-Hookdeck-Eventid")`.
- Coverage on `webhook` package after the change: 95.3% (per code-reviewer report); `go test -race` couldn't run locally (no `gcc`/CGO in this sandbox) — deferred to CI.
- Sandbox token `15a3...8223` has now been pasted into chat twice and is still not rotated.

## What We Tried

- Debugging the "missing type in composite literal" CI error as a real typecheck problem before recognizing it was a lint-binary/Go-toolchain version mismatch (v6 action pinning an old golangci-lint build).
- First webhook capture via `hookdeck listen` tunnel — discarded after noticing `X-Hookdeck-*` headers didn't match what Signater's docs and the second (SSH-tunnel) capture showed.

## Root Cause Analysis

CI failure: `golangci-lint-action@v6`'s `version: latest` resolved a golangci-lint release built against an older Go than the `stable` setup-go toolchain on the runner — a toolchain/linter-binary compatibility gap that only `@v8` (golangci-lint v2) closes correctly.

Webhook capture confusion: using a third-party tunnel tool (`hookdeck listen`) that terminates and re-signs/re-headers requests introduced an intermediary that silently diverged from what the origin service actually sends — a classic "the test harness is part of the system under test" trap.

Doc overclaim: writing the README/package-doc summary before double-checking it matched the narrowest, most-hedged statement already present in the code (`ConstructEvent`'s godoc). Summarizing tends to strip hedges; that's exactly where overclaiming creeps in.

## Lessons Learned

- Pin CI linter actions to major versions that track the same toolchain as `setup-go: stable`, not `version: latest` inside an old action major — the "latest linter, arbitrary Go" combination is a real compatibility surface, not a hypothetical.
- Don't trust an intermediary tunnel/proxy to pass through headers verbatim without confirming it — verify against the most direct path available (raw SSH tunnel here) before trusting captured data.
- When promoting a hedged, VERIFY(sandbox)-style claim into README/package-doc prose, re-derive it from the narrowest confirmed statement in the code, not from memory of what was "basically confirmed." A reviewer catching this before merge is the discipline working as designed — don't skip that step under release-day time pressure.
- gosec taint analysis doesn't understand `%q` formatting as sanitization; a scoped path-based config exclusion is more maintainable than per-line nolints when the same false positive repeats across every log line in a file.

## Next Steps

1. Rotate the sandbox API token (`15a3...8223`) — exposed in chat twice now. Owner: user. Timeline: immediate, before any further sandbox work.
2. Deactivate or repoint the sandbox webhook named "SDK capture" in the dashboard — it still points at the dead tunnel URL from this session.
3. Cut v0.1.1 when convenient to ship the `Unreleased` CHANGELOG entries (RequestID fallback, VerifyAPIKey default) currently sitting uncut.

## AgentWiki

`agentwiki` CLI is not configured in this environment and no AgentWiki MCP tools are exposed in this session. AgentWiki publish skipped — local file is the source of truth.

---
Status: DONE
Summary: CI lint pipeline fixed after a golangci-lint-action version/toolchain mismatch, the last two webhook VERIFY(sandbox) markers were closed with real sandbox captures (with a tunnel-header gotcha and a caught doc-overclaim along the way), and post-release checks (repo public, pkg.go.dev, GitHub Release) are done except sandbox token rotation.
Concerns/Blockers: sandbox token still exposed and unrotated (user action required); AgentWiki publishing unavailable in this environment.
