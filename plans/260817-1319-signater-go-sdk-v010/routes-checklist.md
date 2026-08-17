# Routes Checklist — 100% API Coverage Gate

Source of truth: `../2026-08-17-sdk-architecture/reference/llms.txt` (50 endpoints).
Verified 2026-08-17: llms.txt ≡ site navigation ≡ downloaded fragments (`reference/endpoints/`). Totals: Contact 7, Document 4, Envelope 19, Template 7, Vault 13 = **50**.

Rules for checking a box: method implemented + typed params/response from the endpoint's OpenAPI fragment + test asserting HTTP method, path, and body/response shape against the doc example.

## Contact — `contact.go` (Phase 2) — 7 routes

- [x] `GET /v1/ecm/contacts` → `Contacts.List` / `Contacts.ListAutoPaging` — [list-contacts-12252667e0](../2026-08-17-sdk-architecture/reference/endpoints/list-contacts-12252667e0.md)
- [x] `POST /v1/ecm/contacts` → `Contacts.Create` — [create-contact-12252668e0](../2026-08-17-sdk-architecture/reference/endpoints/create-contact-12252668e0.md)
- [x] `GET /v1/ecm/contacts/{contactId}` → `Contacts.Get` — [get-contact-12252669e0](../2026-08-17-sdk-architecture/reference/endpoints/get-contact-12252669e0.md)
- [x] `PUT /v1/ecm/contacts/{contactId}` → `Contacts.Update` — [update-contact-12252670e0](../2026-08-17-sdk-architecture/reference/endpoints/update-contact-12252670e0.md)
- [x] `DELETE /v1/ecm/contacts/{contactId}` → `Contacts.Remove` — [remove-contact-12252671e0](../2026-08-17-sdk-architecture/reference/endpoints/remove-contact-12252671e0.md)
- [x] `POST /v1/ecm/contacts/{contactId}/favorite` → `Contacts.Favorite` — [favorite-contact-12252672e0](../2026-08-17-sdk-architecture/reference/endpoints/favorite-contact-12252672e0.md)
- [x] `POST /v1/ecm/contacts/{contactId}/unfavorite` → `Contacts.Unfavorite` — [unfavorite-contact-12252673e0](../2026-08-17-sdk-architecture/reference/endpoints/unfavorite-contact-12252673e0.md)

## Document — `document.go` (Phase 3) — 4 routes

- [x] `POST /v1/ecm/documents` (multipart) → `Documents.Upload` — [upload-document-12252674e0](../2026-08-17-sdk-architecture/reference/endpoints/upload-document-12252674e0.md)
- [x] `POST /v1/ecm/documents/templates` → `Documents.CreateFromTemplate` — [create-document-from-template-12252675e0](../2026-08-17-sdk-architecture/reference/endpoints/create-document-from-template-12252675e0.md)
- [x] `GET /v1/ecm/documents/{documentId}/original` → `Documents.OriginalFileURL` — [get-document-original-file-url-13038397e0](../2026-08-17-sdk-architecture/reference/endpoints/get-document-original-file-url-13038397e0.md)
- [x] `GET /v1/ecm/documents/{documentId}/signed` → `Documents.SignedFileURL` — [get-document-signed-file-url-13038398e0](../2026-08-17-sdk-architecture/reference/endpoints/get-document-signed-file-url-13038398e0.md)

## Envelope — `envelope.go` (Phase 3) — 19 routes

- [x] `POST /v1/ecm/envelopes` → `Envelopes.Create` — [create-envelope-12252679e0](../2026-08-17-sdk-architecture/reference/endpoints/create-envelope-12252679e0.md)
- [x] `GET /v1/ecm/envelopes/{envelopeId}` → `Envelopes.Get` — [get-envelope-12252676e0](../2026-08-17-sdk-architecture/reference/endpoints/get-envelope-12252676e0.md)
- [x] `PUT /v1/ecm/envelopes/{envelopeId}` → `Envelopes.Update` — [update-envelope-12252677e0](../2026-08-17-sdk-architecture/reference/endpoints/update-envelope-12252677e0.md)
- [x] `DELETE /v1/ecm/envelopes/{envelopeId}` → `Envelopes.Remove` (to trash) — [remove-envelope-12252678e0](../2026-08-17-sdk-architecture/reference/endpoints/remove-envelope-12252678e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/restore/to/{vaultId}` → `Envelopes.Restore` — [restore-envelope-12883372e0](../2026-08-17-sdk-architecture/reference/endpoints/restore-envelope-12883372e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/publish` → `Envelopes.Publish` (402 → `PaymentRequiredError`) — [publish-envelope-12252685e0](../2026-08-17-sdk-architecture/reference/endpoints/publish-envelope-12252685e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/unschedule` → `Envelopes.Unschedule` — [unschedule-envelope-publication-12252686e0](../2026-08-17-sdk-architecture/reference/endpoints/unschedule-envelope-publication-12252686e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/hold` → `Envelopes.Hold` — [hold-envelope-12252687e0](../2026-08-17-sdk-architecture/reference/endpoints/hold-envelope-12252687e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/cancel` → `Envelopes.Cancel` — [cancel-envelope-12252688e0](../2026-08-17-sdk-architecture/reference/endpoints/cancel-envelope-12252688e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/rename` → `Envelopes.Rename` — [rename-envelope-12252683e0](../2026-08-17-sdk-architecture/reference/endpoints/rename-envelope-12252683e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/to/{vaultId}` → `Envelopes.Move` — [move-envelope-12252684e0](../2026-08-17-sdk-architecture/reference/endpoints/move-envelope-12252684e0.md)
- [x] `GET /v1/ecm/envelopes/{envelopeId}/owners` → `Envelopes.ListAvailableOwners` — [list-available-envelope-owners-12252681e0](../2026-08-17-sdk-architecture/reference/endpoints/list-available-envelope-owners-12252681e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/user-account-owner` → `Envelopes.ChangeOwner` — [change-envelope-user-account-owner-12252682e0](../2026-08-17-sdk-architecture/reference/endpoints/change-envelope-user-account-owner-12252682e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/manual-reinvite-to-review` → `Envelopes.ReinviteToReview` — [manual-reinvite-to-review-envelope-12252690e0](../2026-08-17-sdk-architecture/reference/endpoints/manual-reinvite-to-review-envelope-12252690e0.md)
- [x] `GET /v1/ecm/envelopes/{envelopeId}/signers/{signerId}/signature-link` → `Envelopes.CreateSignatureLink` — [create-signature-link-12252691e0](../2026-08-17-sdk-architecture/reference/endpoints/create-signature-link-12252691e0.md)
- [x] `GET /v1/ecm/envelopes/{envelopeId}/certificate` → `Envelopes.CertificateFileURL` — [get-envelope-certificate-file-13038399e0](../2026-08-17-sdk-architecture/reference/endpoints/get-envelope-certificate-file-13038399e0.md)
- [x] `POST /v1/ecm/envelopes/{envelopeId}/certificate` → `Envelopes.ProcessCertificate` — [process-envelope-certificate-12252680e0](../2026-08-17-sdk-architecture/reference/endpoints/process-envelope-certificate-12252680e0.md)
- [x] `GET /v1/ecm/envelopes/{envelopeId}/sign-marks/{signMarkId}/attachments/{attachmentExternalId}/download` (302) → `Envelopes.SignerAttachmentURL` — [download-an-attachment-uploaded-by-a-signer-35224678e0](../2026-08-17-sdk-architecture/reference/endpoints/download-an-attachment-uploaded-by-a-signer-35224678e0.md)
- [x] `GET /v1/ecm/envelopes/{envelopeId}/signers/{signerId}/simple-selfie` (302) → `Envelopes.SignerSelfieURL` — [download-a-signers-simple-selfie-image-35224679e0](../2026-08-17-sdk-architecture/reference/endpoints/download-a-signers-simple-selfie-image-35224679e0.md)

## Template — `template.go` (Phase 2) — 7 routes

- [x] `POST /v1/ecm/templates` → `Templates.Create` — [create-template-12252697e0](../2026-08-17-sdk-architecture/reference/endpoints/create-template-12252697e0.md)
- [x] `GET /v1/ecm/templates/{templateId}` → `Templates.Get` — [get-template-12252694e0](../2026-08-17-sdk-architecture/reference/endpoints/get-template-12252694e0.md)
- [x] `PUT /v1/ecm/templates/{templateId}` → `Templates.Update` — [update-template-12252695e0](../2026-08-17-sdk-architecture/reference/endpoints/update-template-12252695e0.md)
- [x] `DELETE /v1/ecm/templates/{templateId}` → `Templates.Remove` — [remove-template-12252696e0](../2026-08-17-sdk-architecture/reference/endpoints/remove-template-12252696e0.md)
- [x] `POST /v1/ecm/templates/{templateId}/rename` → `Templates.Rename` — [rename-template-12252698e0](../2026-08-17-sdk-architecture/reference/endpoints/rename-template-12252698e0.md)
- [x] `POST /v1/ecm/templates/{templateId}/move` → `Templates.Move` — [move-template-12252699e0](../2026-08-17-sdk-architecture/reference/endpoints/move-template-12252699e0.md)
- [x] `POST /v1/ecm/templates/{templateId}/restore` → `Templates.Restore` — [restore-template-12252700e0](../2026-08-17-sdk-architecture/reference/endpoints/restore-template-12252700e0.md)

## Vault — `vault.go` (Phase 2) — 13 routes

- [x] `POST /v1/ecm/vaults` → `Vaults.Create` — [create-vault-12252701e0](../2026-08-17-sdk-architecture/reference/endpoints/create-vault-12252701e0.md)
- [x] `GET /v1/ecm/vaults/{vaultId}` → `Vaults.Get` — [get-vault-12252703e0](../2026-08-17-sdk-architecture/reference/endpoints/get-vault-12252703e0.md)
- [x] `PUT /v1/ecm/vaults/{vaultId}` (full replace, newer) → `Vaults.Update` — [update-vault-36651167e0](../2026-08-17-sdk-architecture/reference/endpoints/update-vault-36651167e0.md)
- [x] `PATCH /v1/ecm/vaults/{vaultId}` (partial, legacy) → `Vaults.UpdatePartial` — [update-vault-12252702e0](../2026-08-17-sdk-architecture/reference/endpoints/update-vault-12252702e0.md)
- [x] `DELETE /v1/ecm/vaults/{vaultId}` → `Vaults.Remove` — [remove-vault-12252704e0](../2026-08-17-sdk-architecture/reference/endpoints/remove-vault-12252704e0.md)
- [x] `GET /v1/ecm/vaults/owners` (spec has no {vaultId} despite the description; VERIFY(sandbox)) → `Vaults.Owners` — [get-vault-owners-12252705e0](../2026-08-17-sdk-architecture/reference/endpoints/get-vault-owners-12252705e0.md)
- [x] `GET /v1/ecm/vaults/members` (spec has no {vaultId} despite the description; VERIFY(sandbox)) → `Vaults.Members` — [get-vault-members-12252706e0](../2026-08-17-sdk-architecture/reference/endpoints/get-vault-members-12252706e0.md)
- [x] `GET /v1/ecm/vaults/accounts` → `Vaults.ListAccount` — [list-account-vaults-12252707e0](../2026-08-17-sdk-architecture/reference/endpoints/list-account-vaults-12252707e0.md)
- [x] `GET /v1/ecm/vaults/user-accounts` → `Vaults.ListMine` — [list-current-user-vaults-12252708e0](../2026-08-17-sdk-architecture/reference/endpoints/list-current-user-vaults-12252708e0.md)
- [x] `GET /v1/ecm/vaults/user-accounts/{userAccountId}` → `Vaults.ListByUser` — [list-user-vaults-12252709e0](../2026-08-17-sdk-architecture/reference/endpoints/list-user-vaults-12252709e0.md)
- [x] `GET /v1/ecm/vaults/{vaultId}/list` → `Vaults.Contents` (envelopes + templates) — [list-vault-12252710e0](../2026-08-17-sdk-architecture/reference/endpoints/list-vault-12252710e0.md)
- [x] `POST /v1/ecm/vaults/{vaultId}/favorite` → `Vaults.Favorite` — [favorite-vault-12252711e0](../2026-08-17-sdk-architecture/reference/endpoints/favorite-vault-12252711e0.md)
- [x] `POST /v1/ecm/vaults/{vaultId}/unfavorite` → `Vaults.Unfavorite` — [unfavorite-vault-12252712e0](../2026-08-17-sdk-architecture/reference/endpoints/unfavorite-vault-12252712e0.md)

## Coverage Gate (Phase 5 exit criterion)

- [ ] All 50 boxes above checked
- [ ] Re-fetch `https://docs.api.signater.com/llms.txt` before release and diff against this list — if Signater added endpoints since 2026-08-17, add them or explicitly defer with a note
- [ ] Path params in exact `{camelCase}` names above match implemented URL builders (grep test)
