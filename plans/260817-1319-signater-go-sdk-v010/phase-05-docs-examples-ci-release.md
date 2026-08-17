---
title: "Phase 5: Docs Examples CI Release"
status: todo
priority: P2
effort: "1d"
dependencies: [3, 4]
---

# Phase 5: Docs Examples CI Release

## Overview

Ship it: runnable examples, README, CI pipeline, sandbox validation of the empirical unknowns, GitHub org/repo creation, and the v0.1.0 release.

## Requirements

- Functional: quickstart example runs end-to-end against sandbox (given `SIGNATER_API_TOKEN`); webhook-server example verifies a signed test payload.
- Non-functional: pkg.go.dev renders full docs; CI green on push; no secrets committed (sandbox token via env only).

## Related Code Files

- Create: `README.md`, `CHANGELOG.md`, `doc.go` (package overview), `examples/quickstart/main.go`, `examples/webhook-server/main.go`, `.github/workflows/ci.yml`, `.golangci.yml`

## Implementation Steps

1. `doc.go`: package-level godoc — auth, sandbox note (token decides environment), error handling with `errors.As`, pagination, retry policy summary.
2. `examples/quickstart/main.go`: upload PDF → create envelope (1 signer, email OTP) → publish → poll `Envelopes.Get` until status change; reads `SIGNATER_API_TOKEN` from env; clear log output.
3. `examples/webhook-server/main.go`: `net/http` server using `webhook.ConstructEvent`; instructions for Hookdeck CLI local debugging (per docs).
4. **Sandbox validation pass** (needs a sandbox token from the user): run quickstart; capture one real webhook via Hookdeck CLI; close the `// VERIFY(sandbox):` items from Phase 4 (signature header format, signer id JSON key); confirm the gateway accepts Go's canonicalized `X-Api-Token` header casing (Phase 1 deferred item); fix any schema mismatches found in Phases 2-3 fixtures.
5. README: badges (CI, Go Reference, Go Report Card), install, quickstart snippet, auth/sandbox, pagination (both styles), error handling table, retry policy, webhook verification, link to Signater docs. English.
6. CI (`.github/workflows/ci.yml`): `go vet`, `golangci-lint run`, `go test -race -cover ./...` on go 1.23 + stable matrix. `.golangci.yml` with a sane default linter set.
7. CHANGELOG.md (v0.1.0, keep-a-changelog format). Conventional commits, no AI references.
8. Create GitHub org `tunni-sdk` (user action if org creation needs interactive auth) + repo `signater-go`; push; tag `v0.1.0`; create GitHub release; verify `go get github.com/tunni-sdk/signater-go@v0.1.0` resolves and pkg.go.dev indexes.

## Success Criteria

- [ ] Quickstart runs green against sandbox end-to-end
- [ ] Real sandbox webhook verified by `ConstructEvent` (empirical unknowns closed)
- [ ] CI green; README complete; pkg.go.dev renders
- [ ] `go get github.com/tunni-sdk/signater-go@v0.1.0` works

## Risk Assessment

- Sandbox token availability is a user dependency — request it before starting step 4; steps 1-3 and 5-7 proceed without it.
- Org name `tunni-sdk` could be taken on GitHub — fallback: confirm alternative with user BEFORE first push (import path is baked into go.mod).
- Envelope publish consumes credits even in sandbox flows? Docs say sandbox needs no paid plan — verify; if publish fails with 402 in sandbox, document and stub that step in quickstart.
