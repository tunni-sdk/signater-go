---
title: "Signater Go SDK v0.1.0"
description: "Complete handwritten Go SDK for the Signater e-signature API, Stripe-style architecture, module github.com/tunni-sdk/signater-go"
status: pending
priority: P1
effort: "5-8d"
tags: [go, sdk, signater, api-client]
created: 2026-08-17
---

# Signater Go SDK v0.1.0

## Overview

Build `github.com/tunni-sdk/signater-go`: a professional, zero-dependency Go SDK covering all ~50 endpoints of the Signater e-signature API. Architecture approved in the brainstorm report (Stripe-style service pattern, functional options, typed errors, Go 1.23 `iter.Seq2` auto-pagination, retry with backoff, webhook verification subpackage).

- Design doc: `../2026-08-17-sdk-architecture/brainstorm-report.md`
- **Routes coverage checklist (50 routes, the coverage gate): [routes-checklist.md](./routes-checklist.md)**
- Endpoint inventory (50 routes): `../2026-08-17-sdk-architecture/reference/endpoint-inventory.txt`
- Per-endpoint OpenAPI fragments + examples: `../2026-08-17-sdk-architecture/reference/endpoints/*.md`
- Guides (auth, errors, webhooks, lifecycle, concepts): `../2026-08-17-sdk-architecture/reference/*.md`

Constraints: all code/docs/comments in English; Go 1.23+ minimum; stdlib only; project root `/home/geraldo/signater-go`.

## Goals

| # | Goal | Priority |
|---|------|----------|
| 1 | 100% endpoint coverage with typed params/responses and idiomatic DX | P1 |
| 2 | Production-grade core: typed errors, retry/backoff, auto-pagination, telemetry id | P1 |
| 3 | Webhook package with typed events and signature verification | P1 |
| 4 | Docs, examples, CI, release v0.1.0 under github.com/tunni-sdk | P2 |

## Phases

| # | Phase | Status |
|---|-------|--------|
| 1 | [Phase 1: Core Foundation](./phase-01-start.md) | Pending |
| 2 | [Phase 2: Simple Resources - Contacts Vaults Templates](./phase-02-simple-resources-contacts-vaults-templates.md) | Pending |
| 3 | [Phase 3: Documents and Envelopes](./phase-03-documents-and-envelopes.md) | Pending |
| 4 | [Phase 4: Webhook Package](./phase-04-webhook-package.md) | Pending |
| 5 | [Phase 5: Docs Examples CI Release](./phase-05-docs-examples-ci-release.md) | Pending |

Dependencies: strictly sequential 1 → 2 → 3 → 4 → 5 (2-4 all build on core; 4 only needs core; 5 needs everything). Phase 4 may run in parallel with 2-3 if desired — it only depends on Phase 1.

## Success Criteria

- [ ] All 50 routes implemented per [routes-checklist.md](./routes-checklist.md); `go build ./...` clean
- [ ] `go test -race ./...` green; core (errors/retry/pagination/request) with high coverage via httptest
- [ ] `go vet` + `golangci-lint run` clean; zero non-stdlib deps in `go.mod`
- [ ] Webhook `ConstructEvent` verifies signature and returns typed events
- [ ] README quickstart runs against sandbox; pkg.go.dev docs complete
- [ ] Repo pushed to github.com/tunni-sdk/signater-go, tagged v0.1.0

<!-- slug: signater-go-sdk-v010 -->
