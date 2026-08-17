package signater

import (
	"context"
	"io"
	"net/http"
)

// DocumentService handles the Document endpoints. Documents are uploaded (or
// generated from templates) first, then attached to envelopes by id. A
// document that is not attached to an envelope within 24 hours is
// automatically removed.
type DocumentService struct {
	client *Client
}

// Language of a document or envelope, used in markup stamps and signer
// communication.
type Language string

// Supported languages.
const (
	LanguageEnUS Language = "EnUs"
	LanguagePtBR Language = "PtBr"
	LanguageEsES Language = "EsEs"
	LanguageFrFR Language = "FrFr"
	LanguageDeDE Language = "DeDe"
	LanguageItIT Language = "ItIt"
)

// MarkupOrientation controls where the document markup (Signater logo,
// document id, envelope id, public access link, QR code) is stamped.
type MarkupOrientation string

// Markup orientations.
const (
	MarkupOrientationNone   MarkupOrientation = "None"
	MarkupOrientationBottom MarkupOrientation = "Bottom"
	MarkupOrientationTop    MarkupOrientation = "Top"
	MarkupOrientationLeft   MarkupOrientation = "Left"
	MarkupOrientationRight  MarkupOrientation = "Right"
)

// PageSize is the pixel dimensions of one document page.
type PageSize struct {
	// Page is the 1-based page number.
	Page int `json:"page"`
	// Width in pixels.
	Width float64 `json:"width"`
	// Height in pixels.
	Height float64 `json:"height"`
}

// Document is an uploaded or template-generated document. Name and the
// description fields are only populated by CreateFromTemplate.
type Document struct {
	// ID is the document's unique identifier, used when composing envelopes.
	ID string `json:"id"`
	// FileName is the original file name.
	FileName string `json:"fileName"`
	// Name is the document's name. Only populated by CreateFromTemplate.
	Name string `json:"name"`
	// PrivateDescription is the document's private description. Only
	// populated by CreateFromTemplate.
	PrivateDescription string `json:"privateDescription"`
	// PublicDescription is the document's public description. Only populated
	// by CreateFromTemplate.
	PublicDescription string `json:"publicDescription"`
	// OriginalFileSize is the size of the original file in bytes.
	OriginalFileSize int64 `json:"originalFileSize"`
	// PageSizes lists the pixel dimensions of each page.
	PageSizes []PageSize `json:"pageSizes"`
}

// TemplateFieldValue is the typed value for one template field. Set exactly
// one member — the one matching the field's schema type.
type TemplateFieldValue struct {
	// Text value, for text-like fields.
	Text *string `json:"text,omitempty"`
	// Date value as a string, for date fields (use the field schema's
	// dateFormat).
	Date *string `json:"date,omitempty"`
	// Checked value, for checkbox fields.
	Checked *bool `json:"checked,omitempty"`
	// Selected options, for dropdown/radio fields.
	Selected []string `json:"selected,omitempty"`
	// Number value, for numeric fields.
	Number *float64 `json:"number,omitempty"`
	// Amount value, for currency fields.
	Amount *float64 `json:"amount,omitempty"`
}

// DocumentFieldParams fills one template field when generating a document.
type DocumentFieldParams struct {
	// FieldID is the template field's id (see TemplateService.Get). Required.
	FieldID string `json:"fieldId"`
	// Value is the field's typed value.
	Value *TemplateFieldValue `json:"value,omitempty"`
}

// DocumentFromTemplateParams are the fields accepted by CreateFromTemplate.
type DocumentFromTemplateParams struct {
	// TemplateID identifies the template to apply. Required.
	TemplateID string `json:"templateId"`
	// Name of the document; defaults to the original file name.
	Name *string `json:"name,omitempty"`
	// PrivateDescription of the document (up to 4000 characters).
	PrivateDescription *string `json:"privateDescription,omitempty"`
	// PublicDescription of the document (up to 4000 characters).
	PublicDescription *string `json:"publicDescription,omitempty"`
	// Language used in the markup, if configured.
	Language Language `json:"language,omitempty"`
	// MarkupOrientation controls where document information is stamped.
	MarkupOrientation MarkupOrientation `json:"markupOrientation,omitempty"`
	// Fields fill the template's predefined fields.
	Fields []DocumentFieldParams `json:"fields,omitempty"`
}

// Upload streams a file to the API and returns the created document. Accepted
// formats include PDF, DOCX, DOC, RTF, HTML, TXT, ODT, EPUB, XLSX and common
// image formats; non-PDF files are converted to PDF. The upload is performed
// in a single attempt (no automatic retry) because the stream cannot be
// replayed. The document expires in 24 hours if not attached to an envelope.
func (s *DocumentService) Upload(ctx context.Context, fileName string, file io.Reader) (*Document, error) {
	var d Document
	if err := s.client.doMultipart(ctx, http.MethodPost, "/v1/ecm/documents", "File", fileName, file, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// CreateFromTemplate generates a document from a template, filling its fields,
// and returns the created document. The document expires in 24 hours if not
// attached to an envelope.
func (s *DocumentService) CreateFromTemplate(ctx context.Context, params *DocumentFromTemplateParams) (*Document, error) {
	var d Document
	if err := s.client.do(ctx, http.MethodPost, "/v1/ecm/documents/templates", nil, params, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

// OriginalFileURL returns the URL to download the document's original file.
func (s *DocumentService) OriginalFileURL(ctx context.Context, id string) (string, error) {
	var out struct {
		URL string `json:"url"`
	}
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/documents/"+pathEscape(id)+"/original", nil, nil, &out); err != nil {
		return "", err
	}
	return out.URL, nil
}

// SignedFileURL returns the URL to download the document's signed file.
func (s *DocumentService) SignedFileURL(ctx context.Context, id string) (string, error) {
	var out struct {
		URL string `json:"url"`
	}
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/documents/"+pathEscape(id)+"/signed", nil, nil, &out); err != nil {
		return "", err
	}
	return out.URL, nil
}
