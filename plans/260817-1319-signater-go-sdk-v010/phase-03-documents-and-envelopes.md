---
title: "Phase 3: Documents and Envelopes"
status: todo
priority: P1
effort: "2d"
dependencies: [2]
---

# Phase 3: Documents and Envelopes

## Overview

The heart of the SDK (~23 routes): multipart document upload, template-based document creation, and the full envelope lifecycle with its deep creation params (documents → signers → sign marks → MFA factors).

## Requirements

- Functional: all Document (4) and Envelope (19) routes; `EnvelopeStatus` constants for the 9 lifecycle states; 302-redirect download endpoints return the pre-signed URL (do not stream the file).
- Non-functional: `envelope_params.go` mirrors `CreateEnvelopeApiRequest` exactly (see `reference/endpoints/create-envelope-12252679e0.md`, 643 lines — the authoritative schema).

## Architecture

Document (`document.go`):
- `Upload(ctx, r io.Reader, params *DocumentUploadParams) (*Document, error)` — multipart POST `/v1/ecm/documents` (extend `request.go` with `doMultipart`). Doc comment: unused documents expire in 24h.
- `CreateFromTemplate(ctx, params)` — POST `/v1/ecm/documents/templates` (template id + field values).
- `OriginalFileURL(ctx, id)` / `SignedFileURL(ctx, id)` — GET `.../original`, `.../signed`.

Envelope (`envelope.go` + `envelope_params.go`):
- CRUD: `Create`, `Get`, `Update` (PUT), `Remove` (DELETE, → trash), `Restore(ctx, id, vaultID)` (POST `.../restore/to/{vaultId}`)
- Lifecycle verbs: `Publish`, `Hold`, `Cancel`, `Unschedule` (all POST `{id}/<verb>`)
- Management: `Rename`, `Move(ctx, id, vaultID)` (POST `{id}/to/{vaultId}`), `ChangeOwner` (POST `{id}/user-account-owner`), `ListAvailableOwners` (GET `{id}/owners`), `ReinviteToReview`
- Signing: `CreateSignatureLink(ctx, envelopeID, signerID)` (GET `.../signature-link`)
- Certificate: `CertificateFileURL` (GET `{id}/certificate`), `ProcessCertificate` (POST `{id}/certificate`)
- Downloads (302 → pre-signed URL): `SignerAttachmentURL(ctx, envelopeID, signMarkID, attachmentExternalID)`, `SignerSelfieURL(ctx, envelopeID, signerID)` — use a non-following `http.Client` (set `CheckRedirect` to stop) and return the `Location` header.

```go
type EnvelopeStatus string
const (
    EnvelopeStatusDraft EnvelopeStatus = "Draft"
    // PublishScheduled, Published, Hold, Expired, Cancelled,
    // CancelledBySignerMfaError, Rejected, Signed — exact strings from Get envelope schema
)

// envelope_params.go (from CreateEnvelopeApiRequest):
// EnvelopeCreateParams{Name, VaultID, PublicDescription, PrivateDescription, EmailMessage,
//   SignatureOrder, ExpiresAtUtc, PublishAtUtc, ReminderSettings,
//   Documents []*EnvelopeDocumentParams,   // name, descriptions, language, markupOrientation
//   Signers   []*EnvelopeSignerParams}     // name, email, phone, document, role, title,
//                                          // auth factors (OTP email/SMS/WhatsApp, passcode, geolocation...),
//                                          // SignMarks []*SignerSignMarkParams (type, page, position)
```

## Related Code Files

- Create: `document.go`, `envelope.go`, `envelope_params.go`, `document_test.go`, `envelope_test.go`
- Modify: `request.go` (add `doMultipart`, add redirect-capture helper), `signater.go` (wire services)

## Implementation Steps

1. Read `reference/endpoints/create-envelope-12252679e0.md` end-to-end and transcribe the full params tree into `envelope_params.go`; likewise `get-envelope` response for the `Envelope` type and status enum strings.
2. Extend `request.go` with `doMultipart` (streaming `multipart.Writer`, no full buffering) and `doRedirectURL` (returns Location of a 302 without following).
3. Implement `document.go` + tests (multipart boundary/field assertions, fixture responses).
4. Implement `envelope.go` verbs (thin, one-liners over `do()`), then CRUD with deep types.
5. Tests: create-envelope round-trip against doc example JSON; each lifecycle verb hits correct method+path; 402 on Publish maps to `PaymentRequiredError` with `ShouldBuy`; redirect endpoints return URL.
6. `go vet` + `go test -race ./...`.

## Success Criteria

- [x] Envelope create with nested signers/sign-marks/MFA marshals byte-identical to the doc example (fixture test)
- [x] All 9 `EnvelopeStatus` constants match API strings
- [x] Multipart upload streams (no `io.ReadAll` of the file) and decodes the returned document id
- [x] 302 endpoints return pre-signed URL without downloading the file
- [x] Publish 402 → `PaymentRequiredError{ShouldBuy}` (tested)

## Risk Assessment

- This phase is ~50% of total effort; the deep params tree is the main defect surface — fixture-based marshal tests are mandatory, not optional.
- MFA factor list includes "upcoming" factors per docs — implement only factors present in the current schema; leave enum open (string type, not closed enum) so new factors don't break decoding.
- `Update envelope` restrictions (no doc changes after approvals) are server-side — surface via `ValidationError`, do not pre-validate client-side (KISS).

## Implementation Notes (as built)

- `Documents.Upload(ctx, fileName, io.Reader)` — single-part multipart per the spec; single attempt (stream not replayable); synchronizes with its writer goroutine so the caller's reader is never touched after return.
- `doRedirectURL` applies the normal GET retry policy and resolves relative `Location` values; expected-redirect violation and missing Location are explicit errors.
- `ProcessCertificate` sends `{}` — the only lifecycle verb whose fragment declares a JSON request body. `VERIFY(sandbox)`.
- `EnvelopeParams.Signers/Documents` and `SignMarks` have no `omitempty`: PUT is a full replace, so an empty non-nil slice must reach the wire to clear a collection (review finding, mirrors `MemberIDs`).
- Public fields use Go initialisms: `SchemaJSON`, `ValueJSON`, `SignMarkSchemaJSON` (tags unchanged).
- Params time fields serialize in RFC 3339 with offset; docs recommend UTC values.
- Review score 8/10 after fixes; field-for-field schema fidelity confirmed against the fragments (all enums verified: 9 statuses, 13 sign-mark types, 11 action types, 6 languages, etc.).
