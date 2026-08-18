# Handoff — signater-go (continuation guide)

Date: 2026-08-17. Status: **published** — `main` + tag `v0.1.0` live at
https://github.com/tunni-sdk/signater-go (remote: SSH; the fine-grained PAT
cannot create/push to the org — use SSH for git operations).

## Where things stand

- Plan: 4/5 phases complete, Phase 5 at ~90% (`ak plan status plans/260817-1319-signater-go-sdk-v010`).
- All 50 API routes implemented and contract-tested ([routes-checklist.md](./routes-checklist.md), 50/50 checked).
- Gates: `gofmt` / `go vet` clean, `go test -race` green, coverage 87.2% (root) / 95.1% (webhook), zero non-stdlib deps.
- Every phase passed adversarial code review + test validation by subagents; all majors fixed (details in each `phase-*.md` "Implementation Notes").
- Sandbox-validated end-to-end with a real token: auth, vaults, contacts, error mapping, multipart upload, envelope draft lifecycle. Live-API quirks discovered are encoded in the SDK and listed in README "API knowledge baked in".

## Open items (in priority order)

All 2026-08-17 items closed on 2026-08-18 (commits 4c9b18b..36f13a3) except token rotation:

1. **Rotate the sandbox token** `15a3...8223` — pasted in chat twice now; treat as exposed. Generate a new one at app.signater.com → API Tokens. **Still open, user action.**
2. ~~Webhook empirical capture~~ — done; all `VERIFY(sandbox)` markers closed. Findings (also in README "API knowledge baked in" and phase-05 notes): `signer_id` confirmed; no signing secret exposed by the UI so `VerifyAPIKey` is the practical check; real deliveries have no `request-id` header → `RequestID` falls back to `X-Hookdeck-Eventid`. The dashboard no longer offers "CLI Settings"; capture used a regular webhook pointed at a tunnel.
3. ~~Post-release checks~~ — CI green (lint fixed by migrating to golangci-lint v2/action@v8, which surfaced and closed 14 findings), repo made **public**, pkg.go.dev indexed, `go get @v0.1.0` resolves from a clean module. Go Report Card no longer exists (service retired).
4. ~~GitHub Release~~ — created from the CHANGELOG: https://github.com/tunni-sdk/signater-go/releases/tag/v0.1.0
5. New (left for a future session): the sandbox webhook named "SDK capture" in the dashboard points at a dead tunnel URL — deactivate or repoint it. Unreleased CHANGELOG entries exist; cut v0.1.1 when convenient.

## Key decisions to not re-litigate (verified; see review-audit rules)

- `Create` methods return `(id string, error)` — deliberate; API returns only `{id}` (phase-02 notes).
- POST is never retried on 5xx/network (no idempotency keys server-side); 429 retried for all methods.
- `Signers`/`Documents`/`SignMarks`/`MemberIDs` have no `omitempty` — PUT is full-replace, empty slice must clear.
- `Vaults.Owners`/`Members` take no vault id — matches the published spec, confirmed working in sandbox.
- `IdentityDocumentType` tolerates numeric `0` decoding — live API returns numbers for unset enums.

## Map

- SDK root: `signater.go` (client), `request.go` (transport/retry), `errors.go`, `pagination.go`, one file per resource, `webhook/`.
- Plan dir: `plans/260817-1319-signater-go-sdk-v010/` (this file, plan.md, phase files, routes-checklist, reports/).
- API reference: `plans/2026-08-17-sdk-architecture/reference/` (llms.txt index + 50 OpenAPI fragments + guides) — source of truth for types; re-fetch `https://docs.api.signater.com/llms.txt` and diff before future releases.
- Design rationale: `plans/2026-08-17-sdk-architecture/brainstorm-report.md`.
