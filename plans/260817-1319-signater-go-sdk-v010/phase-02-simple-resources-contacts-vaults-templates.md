---
title: "Phase 2: Simple Resources - Contacts Vaults Templates"
status: todo
priority: P1
effort: "1.5d"
dependencies: [1]
---

# Phase 2: Simple Resources - Contacts Vaults Templates

## Overview

Implement the three mostly-CRUD services (27 routes: Contact 7, Vault 13, Template 7 — see [routes-checklist.md](./routes-checklist.md)) on top of the Phase 1 core, locking in the service-pattern conventions that Phase 3 will reuse.

## Requirements

- Functional: every route below implemented with typed params/response; list endpoints get both `List` (raw `Page[T]`) and `ListAutoPaging` (`iter.Seq2`).
- Non-functional: types/fields named from the OpenAPI fragments in `../2026-08-17-sdk-architecture/reference/endpoints/` (each page embeds schema + real examples — use them as test fixtures).

## Architecture

Convention (applies to all services, Stripe-style):

```go
type ContactService struct{ client *Client }
func (s *ContactService) List(ctx context.Context, params *ContactListParams) (*Page[Contact], error)
func (s *ContactService) ListAutoPaging(ctx context.Context, params *ContactListParams) iter.Seq2[Contact, error]
func (s *ContactService) Create(ctx context.Context, params *ContactCreateParams) (*Contact, error)
func (s *ContactService) Get(ctx context.Context, id string) (*Contact, error)
// verbs without response body return error only; optional params are pointers + String()/Bool()/Int() helpers
```

### Routes

Contact (`contact.go`): GET/POST `/v1/ecm/contacts`, GET/PUT/DELETE `/v1/ecm/contacts/{id}`, POST `.../favorite`, POST `.../unfavorite`. List filters: `IsFavorite`, `Search` + ListParams. Contact fields incl. `documentType` enum (e.g. `GenericIdentification`), `phoneIdd`/`phoneNumber`, audit fields (`createdAtUtc`, `createdByName`, ...).

Vault (`vault.go`): POST `/v1/ecm/vaults`, GET/PATCH/PUT/DELETE `/v1/ecm/vaults/{id}`, GET `.../owners`, GET `.../members`, GET `/v1/ecm/vaults/accounts`, GET `/v1/ecm/vaults/user-accounts`, GET `/v1/ecm/vaults/user-accounts/{userAccountId}`, GET `/v1/ecm/vaults/{id}/list` (vault contents: envelopes+templates), POST `.../favorite`, `.../unfavorite`. Access types: UserAccount (private) / Account / UserAccountGroup (member-based).

Template (`template.go`): POST `/v1/ecm/templates`, GET/PUT/DELETE `/v1/ecm/templates/{id}`, POST `.../rename`, `.../move`, `.../restore`.

## Related Code Files

- Create: `contact.go`, `vault.go`, `template.go`, `contact_test.go`, `vault_test.go`, `template_test.go`
- Modify: `signater.go` (wire service fields in `NewClient`)

## Implementation Steps

1. Read the relevant `reference/endpoints/*.md` fragments for each resource before typing its structs; copy real example payloads into `testdata/` fixtures.
2. Implement `contact.go` first (smallest) — it sets the file template: service struct, types, params, methods, doc comments referencing docs URL.
3. Implement `template.go`, then `vault.go` (largest).
4. Vault has TWO update routes: PATCH `{id}` (legacy partial, access type immutable) and PUT `{id}` (newer full-replace, access type mutable, `MemberIds` required for UserAccountGroup). Expose PUT as `Update` and PATCH as `UpdatePartial` — document the difference in both doc comments.
5. Tests per service via `internal/testutil`: assert method, path, query encoding, request body shape, response decoding from fixtures; one auto-paging test on Contacts.List.
6. `go vet` + `go test -race ./...`.

## Success Criteria

- [ ] All 27 routes callable and tested against fixture payloads from the docs; boxes checked in [routes-checklist.md](./routes-checklist.md)
- [ ] `ListAutoPaging` works on every list endpoint
- [ ] Vault dual-update exposed as `Update` (PUT) + `UpdatePartial` (PATCH) with clear doc comments
- [ ] Naming/conventions consistent across the three files (this is the template for Phase 3)

## Risk Assessment

- Doc fragments are Apidog-generated; a schema may disagree with the live API. Where example and schema conflict, trust the example and leave a `// NOTE:` doc comment. Empirical sandbox validation happens in Phase 5 quickstart.
- Vault "list" family has overlapping semantics (account vs user vs contents) — method names must disambiguate: `ListAccount`, `ListMine`, `ListByUser`, `Contents`.
