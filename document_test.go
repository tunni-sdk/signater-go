package signater

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestDocumentUpload(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusCreated, `{"id":"doc-1","fileName":"contract.pdf","originalFileSize":1234,"pageSizes":[{"page":1,"width":595.5,"height":842.0}]}`, nil)

	doc, err := c.Documents.Upload(context.Background(), "contract.pdf", strings.NewReader("%PDF-1.7 fake"))
	if err != nil {
		t.Fatal(err)
	}
	if doc.ID != "doc-1" || doc.FileName != "contract.pdf" || doc.OriginalFileSize != 1234 ||
		len(doc.PageSizes) != 1 || doc.PageSizes[0].Width != 595.5 {
		t.Errorf("document = %+v", doc)
	}

	req := tr.Snapshot()[0]
	if req.Method != http.MethodPost || req.Path != "/v1/ecm/documents" {
		t.Errorf("request = %s %s", req.Method, req.Path)
	}
	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data; boundary=") {
		t.Errorf("Content-Type = %q", ct)
	}
	body := string(req.Body)
	if !strings.Contains(body, `name="File"; filename="contract.pdf"`) || !strings.Contains(body, "%PDF-1.7 fake") {
		t.Errorf("multipart body missing file part:\n%s", body)
	}
}

// trackingReader records whether it was read to EOF.
type trackingReader struct {
	r   *strings.Reader
	eof bool
}

func (t *trackingReader) Read(p []byte) (int, error) {
	n, err := t.r.Read(p)
	if err != nil {
		t.eof = true
	}
	return n, err
}

func TestDocumentUploadConsumesReaderBeforeReturning(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusCreated, `{"id":"doc-1"}`, nil)

	reader := &trackingReader{r: strings.NewReader(strings.Repeat("x", 1<<16))}
	if _, err := c.Documents.Upload(context.Background(), "big.pdf", reader); err != nil {
		t.Fatal(err)
	}
	// Upload must synchronize with its writer goroutine so the caller can
	// safely close/reuse the reader immediately after return.
	if !reader.eof {
		t.Error("Upload returned before fully consuming the reader")
	}
}

func TestDocumentUploadNoRetry(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusInternalServerError, "", nil)

	_, err := c.Documents.Upload(context.Background(), "a.pdf", strings.NewReader("x"))
	var base *APIError
	if !errors.As(err, &base) || base.StatusCode != 500 {
		t.Fatalf("want 500 APIError, got %v", err)
	}
	if tr.Count() != 1 {
		t.Errorf("requests = %d, want 1 (multipart must not retry)", tr.Count())
	}
}

func TestDocumentCreateFromTemplate(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusCreated, `{"id":"doc-2","fileName":"generated.pdf","name":"Example document","originalFileSize":10}`, nil)

	doc, err := c.Documents.CreateFromTemplate(context.Background(), &DocumentFromTemplateParams{
		TemplateID:        "3e533f62-031a-4ea6-8e13-5be21f6ded7b",
		Name:              String("Example document"),
		Language:          LanguageEnUS,
		MarkupOrientation: MarkupOrientationBottom,
		Fields: []DocumentFieldParams{
			{FieldID: "178712be-0599-4387-b68c-d8376cf9a5a3", Value: &TemplateFieldValue{Text: String("Example field value")}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if doc.ID != "doc-2" || doc.Name != "Example document" {
		t.Errorf("document = %+v", doc)
	}

	req := tr.Snapshot()[0]
	if req.Path != "/v1/ecm/documents/templates" {
		t.Errorf("path = %s", req.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(req.Body, &body); err != nil {
		t.Fatal(err)
	}
	if body["templateId"] == "" || body["language"] != "EnUs" || body["markupOrientation"] != "Bottom" {
		t.Errorf("body = %v", body)
	}
	fields := body["fields"].([]any)
	value := fields[0].(map[string]any)["value"].(map[string]any)
	if value["text"] != "Example field value" {
		t.Errorf("field value = %v", value)
	}
	if _, present := value["checked"]; present {
		t.Error("unset value members must be omitted")
	}
}

func TestDocumentFileURLs(t *testing.T) {
	c, tr := newTestClient(t)
	tr.Enqueue(http.StatusOK, `{"url":"https://cdn.example.com/original.pdf"}`, nil)
	tr.Enqueue(http.StatusOK, `{"url":"https://cdn.example.com/signed.pdf"}`, nil)

	ctx := context.Background()
	orig, err := c.Documents.OriginalFileURL(ctx, "doc-1")
	if err != nil || orig != "https://cdn.example.com/original.pdf" {
		t.Fatalf("original = %q, %v", orig, err)
	}
	signed, err := c.Documents.SignedFileURL(ctx, "doc-1")
	if err != nil || signed != "https://cdn.example.com/signed.pdf" {
		t.Fatalf("signed = %q, %v", signed, err)
	}

	reqs := tr.Snapshot()
	if reqs[0].Path != "/v1/ecm/documents/doc-1/original" || reqs[1].Path != "/v1/ecm/documents/doc-1/signed" {
		t.Errorf("paths = %s, %s", reqs[0].Path, reqs[1].Path)
	}
}
