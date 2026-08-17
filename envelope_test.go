package signater

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestEnvelopeCreateBody(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusCreated, `{"id":"cd24af61-630e-4b00-aa59-3d2d416f9b13"}`, nil)

	expires := time.Date(2026, 9, 16, 0, 52, 53, 0, time.UTC)
	id, err := c.Envelopes.Create(context.Background(), &EnvelopeParams{
		VaultID:                 "9efd3e4a-ca6d-4b3c-b6ca-abc74942a4db",
		Name:                    "Example Envelope",
		JurisdictionCountryCode: "BR",
		PrivateDescription:      String("Private description of the envelope"),
		Message:                 String("Message of the envelope"),
		ReviewReminder:          true,
		ExpiresAt:               &expires,
		Language:                LanguageEnUS,
		MarkupOrientation:       MarkupOrientationBottom,
		SignInOrder:             true,
		RedirectURL:             String("https://example.com/callback"),
		Signers: []EnvelopeSignerParams{{
			Name:                         "John Doe",
			Email:                        String("john.doe@example.com"),
			EmailCommunicationMode:       CommunicationModeFull,
			Title:                        String("Lawyer"),
			ShouldEnforceEmailValidation: true,
			DocumentType:                 IdentityDocumentGenericIdentification,
			DocumentValue:                String("1234567890"),
			SignMarks: []SignMarkParams{
				{Type: SignMarkRubric, DocumentID: "7665f556-7cf5-4f89-8a6a-b45c7be33396", Page: 1, X: 100, Y: 50, IsRequired: true},
				{Type: SignMarkSignature, DocumentID: "7665f556-7cf5-4f89-8a6a-b45c7be33396", Page: 2, X: 200, Y: 50, IsRequired: true},
			},
		}},
		Documents: []EnvelopeDocumentParams{{
			ID:   "7665f556-7cf5-4f89-8a6a-b45c7be33396",
			Name: "Example Document",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != "cd24af61-630e-4b00-aa59-3d2d416f9b13" {
		t.Errorf("id = %q", id)
	}

	req := tr.Snapshot()[0]
	if req.Method != http.MethodPost || req.Path != "/v1/ecm/envelopes" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(req.Body, &body); err != nil {
		t.Fatal(err)
	}
	if body["vaultId"] == "" || body["jurisdictionCountryCode"] != "BR" || body["reviewReminder"] != true {
		t.Errorf("body root = %v", body)
	}
	if body["expiresAtUtc"] != "2026-09-16T00:52:53Z" {
		t.Errorf("expiresAtUtc = %v", body["expiresAtUtc"])
	}
	signer := body["signers"].([]any)[0].(map[string]any)
	if signer["emailCommunicationMode"] != "Full" || signer["shouldEnforceEmailValidation"] != true {
		t.Errorf("signer = %v", signer)
	}
	// Non-enforced factors are plain bools and must still be serialized.
	if signer["shouldEnforceSmsValidation"] != false {
		t.Errorf("shouldEnforceSmsValidation = %v", signer["shouldEnforceSmsValidation"])
	}
	marks := signer["signMarks"].([]any)
	mark := marks[1].(map[string]any)
	if mark["type"] != "Signature" || mark["page"] != float64(2) || mark["y"] != float64(50) {
		t.Errorf("sign mark = %v", mark)
	}
	if _, present := mark["id"]; present {
		t.Error("nil sign mark id must be omitted on create")
	}
	doc := body["documents"].([]any)[0].(map[string]any)
	if doc["id"] == "" || doc["name"] != "Example Document" {
		t.Errorf("document = %v", doc)
	}
}

func TestEnvelopeUpdateSerializesIDsAndCollections(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusNoContent, "", nil)

	err := c.Envelopes.Update(context.Background(), "e1", &EnvelopeParams{
		VaultID: "v1", Name: "n", JurisdictionCountryCode: "BR",
		Signers: []EnvelopeSignerParams{{
			ID:   String("signer-1"),
			Name: "John",
			SignMarks: []SignMarkParams{
				{ID: String("mark-1"), Type: SignMarkSignature, DocumentID: "d1", Page: 1},
			},
		}},
		Documents: []EnvelopeDocumentParams{},
	})
	if err != nil {
		t.Fatal(err)
	}

	var body map[string]any
	if err := json.Unmarshal(tr.Snapshot()[0].Body, &body); err != nil {
		t.Fatal(err)
	}
	signer := body["signers"].([]any)[0].(map[string]any)
	if signer["id"] != "signer-1" {
		t.Errorf("signer id = %v (existing ids must be serialized on update)", signer["id"])
	}
	mark := signer["signMarks"].([]any)[0].(map[string]any)
	if mark["id"] != "mark-1" {
		t.Errorf("mark id = %v", mark["id"])
	}
	// PUT is a full replace: an empty non-nil slice must reach the wire as []
	// so the server clears the collection.
	docs, present := body["documents"]
	if !present {
		t.Fatal("empty documents slice was omitted; PUT cannot clear the collection")
	}
	if arr, ok := docs.([]any); !ok || len(arr) != 0 {
		t.Errorf("documents = %v, want []", docs)
	}
}

func TestEnvelopeGet(t *testing.T) {
	const fixture = `{
	  "id": "env-1",
	  "status": "Published",
	  "name": "Example Envelope",
	  "vaultId": "v-1",
	  "language": "PtBr",
	  "signInOrder": true,
	  "isSandbox": true,
	  "expiresAtUtc": "2026-09-16T00:52:53.894Z",
	  "certifiedAtUtc": null,
	  "documents": [{"id": "d-1", "index": 0, "origin": "Upload", "name": "Contract", "originalFileSize": 99, "pageSizes": [{"page": 1, "width": 595, "height": 842}]}],
	  "signers": [{
	    "id": "s-1",
	    "index": 0,
	    "status": "ReadyToReview",
	    "name": "John Doe",
	    "documentType": "BrazilianCpf",
	    "emailCommunicationMode": "Full",
	    "shouldEnforceEmailValidation": true,
	    "signMarks": [{"id": "m-1", "documentId": "d-1", "type": "Signature", "page": 1, "x": 10, "y": 20, "isRequired": true, "schemaJson": null, "valueJson": null}],
	    "actions": [{"id": "a-1", "createdAtUtc": "2026-08-17T10:00:00Z", "type": "View", "ip": "1.2.3.4"}]
	  }]
	}`
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, fixture, nil)

	env, err := c.Envelopes.Get(context.Background(), "env-1")
	if err != nil {
		t.Fatal(err)
	}
	if env.Status != EnvelopeStatusPublished || env.Language != LanguagePtBR || !env.IsSandbox {
		t.Errorf("envelope = %+v", env)
	}
	if env.ExpiresAt == nil || env.CertifiedAt != nil {
		t.Errorf("time pointers: ExpiresAt=%v CertifiedAt=%v", env.ExpiresAt, env.CertifiedAt)
	}
	d := env.Documents[0]
	if d.Origin != DocumentOriginUpload || d.PageSizes[0].Height != 842 {
		t.Errorf("document = %+v", d)
	}
	s := env.Signers[0]
	if s.Status != SignerStatusReadyToReview || s.DocumentType != IdentityDocumentBrazilianCpf ||
		s.Index == nil || *s.Index != 0 {
		t.Errorf("signer = %+v", s)
	}
	if s.SignMarks[0].Type != SignMarkSignature || s.Actions[0].Type != SignerActionView {
		t.Errorf("marks/actions = %+v / %+v", s.SignMarks, s.Actions)
	}
	if tr.Snapshot()[0].Path != "/v1/ecm/envelopes/env-1" {
		t.Errorf("path = %s", tr.Snapshot()[0].Path)
	}
}

func TestEnvelopeLifecycleVerbs(t *testing.T) {
	c, tr := newTestClient(t)
	for range 10 {
		tr.Enqueue(http.StatusNoContent, "", nil)
	}

	ctx := context.Background()
	steps := []struct {
		call         func() error
		method, path string
	}{
		{func() error {
			return c.Envelopes.Update(ctx, "e1", &EnvelopeParams{VaultID: "v1", Name: "n", JurisdictionCountryCode: "BR"})
		}, http.MethodPut, "/v1/ecm/envelopes/e1"},
		{func() error { return c.Envelopes.Remove(ctx, "e1") }, http.MethodDelete, "/v1/ecm/envelopes/e1"},
		{func() error { return c.Envelopes.Restore(ctx, "e1", "v2") }, http.MethodPost, "/v1/ecm/envelopes/e1/restore/to/v2"},
		{func() error { return c.Envelopes.Publish(ctx, "e1") }, http.MethodPost, "/v1/ecm/envelopes/e1/publish"},
		{func() error { return c.Envelopes.Unschedule(ctx, "e1") }, http.MethodPost, "/v1/ecm/envelopes/e1/unschedule"},
		{func() error { return c.Envelopes.Hold(ctx, "e1") }, http.MethodPost, "/v1/ecm/envelopes/e1/hold"},
		{func() error { return c.Envelopes.Cancel(ctx, "e1") }, http.MethodPost, "/v1/ecm/envelopes/e1/cancel"},
		{func() error { return c.Envelopes.Move(ctx, "e1", "v3") }, http.MethodPost, "/v1/ecm/envelopes/e1/to/v3"},
		{func() error { return c.Envelopes.ReinviteToReview(ctx, "e1") }, http.MethodPost, "/v1/ecm/envelopes/e1/manual-reinvite-to-review"},
		{func() error { return c.Envelopes.ProcessCertificate(ctx, "e1") }, http.MethodPost, "/v1/ecm/envelopes/e1/certificate"},
	}
	for i, step := range steps {
		if err := step.call(); err != nil {
			t.Fatalf("step %d: %v", i, err)
		}
		req := tr.Snapshot()[i]
		if req.Method != step.method || req.Path != step.path {
			t.Errorf("step %d = %s %s, want %s %s", i, req.Method, req.Path, step.method, step.path)
		}
	}
}

func TestEnvelopeRenameAndChangeOwner(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusNoContent, "", nil)
	tr.Enqueue(http.StatusNoContent, "", nil)

	ctx := context.Background()
	if err := c.Envelopes.Rename(ctx, "e1", "New envelope name"); err != nil {
		t.Fatal(err)
	}
	if err := c.Envelopes.ChangeOwner(ctx, "e1", "owner-9"); err != nil {
		t.Fatal(err)
	}

	reqs := tr.Snapshot()
	var rename, owner map[string]string
	if err := json.Unmarshal(reqs[0].Body, &rename); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(reqs[1].Body, &owner); err != nil {
		t.Fatal(err)
	}
	if rename["name"] != "New envelope name" || owner["ownerId"] != "owner-9" {
		t.Errorf("bodies = %v / %v", rename, owner)
	}
	if reqs[1].Path != "/v1/ecm/envelopes/e1/user-account-owner" {
		t.Errorf("path = %s", reqs[1].Path)
	}
}

func TestEnvelopePublish402(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusPaymentRequired, `{"ShouldBuy":"ApiEnvelopes"}`, nil)

	err := c.Envelopes.Publish(context.Background(), "e1")
	var pe *PaymentRequiredError
	if !errors.As(err, &pe) || pe.ShouldBuy != ShouldBuyAPIEnvelopes {
		t.Fatalf("want PaymentRequiredError{ApiEnvelopes}, got %v", err)
	}
}

func TestEnvelopeSignatureLinkAndOwners(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, `{"url":"https://app.signater.com/sign/abc"}`, nil)
	tr.Enqueue(http.StatusOK, `{"items":[{"id":"u1","name":"John","isActive":true,"isCurrent":false,"role":"User"}]}`, nil)

	ctx := context.Background()
	link, err := c.Envelopes.CreateSignatureLink(ctx, "e1", "s1")
	if err != nil || link != "https://app.signater.com/sign/abc" {
		t.Fatalf("link = %q, %v", link, err)
	}
	owners, err := c.Envelopes.ListAvailableOwners(ctx, "e1")
	if err != nil || len(owners) != 1 || owners[0].Role != UserRoleUser {
		t.Fatalf("owners = %+v, %v", owners, err)
	}

	reqs := tr.Snapshot()
	if reqs[0].Path != "/v1/ecm/envelopes/e1/signers/s1/signature-link" || reqs[1].Path != "/v1/ecm/envelopes/e1/owners" {
		t.Errorf("paths = %s, %s", reqs[0].Path, reqs[1].Path)
	}
}

func TestEnvelopeCertificateFileURL(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, `{"url":"https://cdn.example.com/cert.pdf"}`, nil)

	u, err := c.Envelopes.CertificateFileURL(context.Background(), "e1")
	if err != nil || u != "https://cdn.example.com/cert.pdf" {
		t.Fatalf("url = %q, %v", u, err)
	}
	if req := tr.Snapshot()[0]; req.Method != http.MethodGet {
		t.Errorf("method = %s", req.Method)
	}
}

func TestEnvelopeRedirectDownloads(t *testing.T) {
	c, tr := newTestClient(t)
	h := http.Header{}
	h.Set("Location", "https://s3.example.com/presigned-attachment")
	tr.Enqueue(http.StatusFound, "", h)
	h2 := http.Header{}
	h2.Set("Location", "https://s3.example.com/presigned-selfie")
	tr.Enqueue(http.StatusFound, "", h2)

	ctx := context.Background()
	u, err := c.Envelopes.SignerAttachmentURL(ctx, "e1", "m1", "att-1")
	if err != nil || u != "https://s3.example.com/presigned-attachment" {
		t.Fatalf("attachment url = %q, %v", u, err)
	}
	u2, err := c.Envelopes.SignerSelfieURL(ctx, "e1", "s1")
	if err != nil || u2 != "https://s3.example.com/presigned-selfie" {
		t.Fatalf("selfie url = %q, %v", u2, err)
	}

	reqs := tr.Snapshot()
	if reqs[0].Path != "/v1/ecm/envelopes/e1/sign-marks/m1/attachments/att-1/download" {
		t.Errorf("attachment path = %s", reqs[0].Path)
	}
	if reqs[1].Path != "/v1/ecm/envelopes/e1/signers/s1/simple-selfie" {
		t.Errorf("selfie path = %s", reqs[1].Path)
	}
}

func TestEnvelopeRedirectDownloadErrors(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, "", nil) // not a redirect
	tr.Enqueue(http.StatusNotFound, `{"message":"not found"}`, nil)

	ctx := context.Background()
	if _, err := c.Envelopes.SignerSelfieURL(ctx, "e1", "s1"); err == nil {
		t.Error("non-redirect 200 must be an error")
	}
	_, err := c.Envelopes.SignerSelfieURL(ctx, "e1", "s1")
	var nfe *NotFoundError
	if !errors.As(err, &nfe) {
		t.Errorf("want NotFoundError, got %v", err)
	}
}
