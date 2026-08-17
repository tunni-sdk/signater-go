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

1. **Rotate the sandbox token** `15a3...8223` — it was pasted in a chat session; treat as exposed. Generate a new one at app.signater.com → API Tokens.
2. **Webhook empirical capture** — the only remaining `VERIFY(sandbox)` markers (grep `VERIFY(sandbox)` in `webhook/webhook.go` and `vault.go`):
   - Confirm Hookdeck signature header name/format (assumed `X-Hookdeck-Signature[-2]`, base64 HMAC-SHA256 of raw body).
   - Confirm the signer id JSON key on `*_by_signer` events (assumed `signer_id`; `Event.Raw` preserves the payload as a hedge).
   - How: configure a sandbox webhook at app.signater.com/account/webhooks, run `hookdeck listen 8080` + `examples/webhook-server`, publish a sandbox envelope, inspect the delivery.
3. **Post-release checks**: CI run on GitHub Actions (first push triggered it), pkg.go.dev indexing (`go get github.com/tunni-sdk/signater-go@v0.1.0` from a clean module), Go Report Card.
4. Optional: create a GitHub Release for the `v0.1.0` tag using the CHANGELOG entry.

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
