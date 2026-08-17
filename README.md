# signater-go

[![CI](https://github.com/tunni-sdk/signater-go/actions/workflows/ci.yml/badge.svg)](https://github.com/tunni-sdk/signater-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tunni-sdk/signater-go.svg)](https://pkg.go.dev/github.com/tunni-sdk/signater-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/tunni-sdk/signater-go)](https://goreportcard.com/report/github.com/tunni-sdk/signater-go)

A Go SDK for the [Signater](https://docs.api.signater.com) e-signature API. Covers the complete published API surface — contacts, documents, envelopes, templates and vaults (50 endpoints) — plus webhook signature verification. Zero dependencies outside the standard library. Requires Go 1.23+.

## Installation

```bash
go get github.com/tunni-sdk/signater-go
```

## Quickstart

```go
import signater "github.com/tunni-sdk/signater-go"

client := signater.NewClient() // reads SIGNATER_API_TOKEN
// or: signater.NewClient(signater.WithAPIToken(token))

// 1. Upload a document (expires in 24h if not attached to an envelope).
doc, err := client.Documents.Upload(ctx, "contract.pdf", file)

// 2. Create an envelope with a signer.
envID, err := client.Envelopes.Create(ctx, &signater.EnvelopeParams{
    VaultID:                 vaultID,
    Name:                    "Service Agreement",
    JurisdictionCountryCode: "BR",
    Language:                signater.LanguagePtBR,
    Signers: []signater.EnvelopeSignerParams{{
        Name:                      "Jane Doe",
        Email:                     signater.String("jane@example.com"),
        EmailCommunicationMode:    signater.CommunicationModeFull,
        SmsCommunicationMode:      signater.CommunicationModeNone,
        WhatsAppCommunicationMode: signater.CommunicationModeNone,
    }},
    Documents: []signater.EnvelopeDocumentParams{{ID: doc.ID, Name: "contract.pdf"}},
})

// 3. Publish: signers receive access links by email.
err = client.Envelopes.Publish(ctx, envID)
```

Runnable versions: [`examples/quickstart`](examples/quickstart) and [`examples/webhook-server`](examples/webhook-server).

## Authentication and sandbox

Create a token in the **API Tokens** section of the [Signater App](https://app.signater.com) and export it as `SIGNATER_API_TOKEN` (or pass `WithAPIToken`). Sandbox needs no extra configuration — sandbox tokens automatically route to sandbox resources, and sandbox usage requires no paid plan.

## Services

| Service | Endpoints |
|---|---|
| `client.Contacts` | List, Create, Get, Update, Remove, Favorite, Unfavorite |
| `client.Documents` | Upload (multipart), CreateFromTemplate, OriginalFileURL, SignedFileURL |
| `client.Envelopes` | Create, Get, Update, Remove, Restore, Publish, Unschedule, Hold, Cancel, Rename, Move, ChangeOwner, ListAvailableOwners, ReinviteToReview, CreateSignatureLink, CertificateFileURL, ProcessCertificate, SignerAttachmentURL, SignerSelfieURL |
| `client.Templates` | Create, Get, Update, Remove, Rename, Move, Restore |
| `client.Vaults` | Create, Get, Update, UpdatePartial, Remove, Owners, Members, ListAccount, ListMine, ListByUser, Contents, Favorite, Unfavorite |

## Pagination

List endpoints return one page at a time; the `AutoPaging` variants return a Go 1.23 iterator that fetches pages transparently:

```go
for contact, err := range client.Contacts.ListAutoPaging(ctx, nil) {
    if err != nil {
        return err
    }
    fmt.Println(contact.Name)
}
```

## Error handling

Every non-2xx response maps to a typed error carrying the `X-Signater-Telemetry-Operation-Id` (quote it to Signater support). Match with `errors.As`:

| Status | Error type | Extras |
|---|---|---|
| 400 | `*signater.ValidationError` | `Errors []FieldError` |
| 401 | `*signater.AuthenticationError` | |
| 402 | `*signater.PaymentRequiredError` | `ShouldBuy` (product to purchase) |
| 403 | `*signater.PermissionError` | |
| 404 | `*signater.NotFoundError` | |
| 429 | `*signater.RateLimitError` | `RetryAfter` |
| any | `*signater.APIError` | matches every API error |

```go
var payErr *signater.PaymentRequiredError
if errors.As(err, &payErr) {
    fmt.Println("buy:", payErr.ShouldBuy) // e.g. ApiEnvelopes
}
```

## Retries

Transient failures are retried automatically (default 2 retries, exponential backoff with jitter, `Retry-After` honored): network errors and transient 5xx on idempotent methods, 429 on any method. POST is never retried on 5xx/network errors — the API has no idempotency keys, so a blind retry could duplicate an envelope. Tune with `WithMaxRetries`, `WithRequestTimeout`; observe with `WithLogger`.

## Webhooks

The [`webhook`](webhook) package verifies and parses deliveries (sent via Hookdeck):

```go
import "github.com/tunni-sdk/signater-go/webhook"

payload, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
event, err := webhook.ConstructEvent(payload, r.Header, secret)
switch event.Type {
case webhook.EventEnvelopeSigned:
    // ...
}
```

Deliveries are at-least-once and unordered: dedupe by `webhook.RequestID(r.Header)` and query the API for the current resource state. Test your handler with `webhook.SignPayload`.

## API knowledge baked in

The SDK encodes behavior verified against the live sandbox that the API reference does not spell out:

- `EnvelopeParams.Language` and each signer's three `*CommunicationMode` fields are required by the live API.
- A document can belong to only one envelope (even a trashed one) — upload one document per envelope.
- Unset signer `documentType` comes back as the number `0`; the SDK decodes it as the empty string.
- Rate limit: 1,000 requests/minute.

## License

[MIT](LICENSE)
