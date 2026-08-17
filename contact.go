package signater

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ContactService handles the Contact endpoints. Contacts store reusable
// signer information (name, email, phone, identity document) so envelopes can
// be assembled without re-entering it.
type ContactService struct {
	client *Client
}

// IdentityDocumentType is the kind of identity document attached to a contact
// or signer.
type IdentityDocumentType string

// Identity document types accepted by the API.
const (
	IdentityDocumentGenericIdentification IdentityDocumentType = "GenericIdentification"
	IdentityDocumentBrazilianCpf          IdentityDocumentType = "BrazilianCpf"
	IdentityDocumentBrazilianIdentity     IdentityDocumentType = "BrazilianIdentity"
	IdentityDocumentPassport              IdentityDocumentType = "Passport"
)

// Contact is a stored contact. List responses additionally populate the audit
// fields (CreatedAt, UpdatedAt, ...) and IsFavorite, which Get does not return.
type Contact struct {
	// ID is the contact's unique identifier.
	ID string `json:"id"`
	// Name is the contact's name.
	Name string `json:"name"`
	// Title is the contact's title, e.g. "Company Director".
	Title string `json:"title"`
	// Email is the contact's email address.
	Email string `json:"email"`
	// PhoneIDD is the international dialing code, e.g. 55 for Brazil.
	PhoneIDD int `json:"phoneIdd"`
	// PhoneNumber is the phone number excluding the IDD.
	PhoneNumber string `json:"phoneNumber"`
	// DocumentType is the kind of identity document.
	DocumentType IdentityDocumentType `json:"documentType"`
	// DocumentValue is the identity document number.
	DocumentValue string `json:"documentValue"`
	// Description is a free-form description of the contact. Only populated
	// by Get.
	Description string `json:"description"`
	// IsFavorite reports whether the current user favorited the contact.
	// Only populated by List.
	IsFavorite bool `json:"isFavorite"`
	// CreatedAt is when the contact was created (UTC). Only populated by List.
	CreatedAt time.Time `json:"createdAtUtc"`
	// CreatedByName is the name of the user who created the contact. Only
	// populated by List.
	CreatedByName string `json:"createdByName"`
	// CreatedByAvatar is the avatar URL of the creating user. Only populated
	// by List.
	CreatedByAvatar string `json:"createdByAvatar"`
	// UpdatedAt is when the contact was last updated (UTC). Only populated by
	// List.
	UpdatedAt time.Time `json:"updatedAtUtc"`
	// UpdatedByName is the name of the user who last updated the contact.
	// Only populated by List.
	UpdatedByName string `json:"updatedByName"`
	// UpdatedByAvatar is the avatar URL of the last updating user. Only
	// populated by List.
	UpdatedByAvatar string `json:"updatedByAvatar"`
}

// ContactParams are the fields accepted when creating or updating a contact.
// Name is required; optional fields use pointers (see String, Int, Bool
// helpers) and are omitted when nil.
type ContactParams struct {
	// Name of the contact (2-200 characters). Required.
	Name string `json:"name"`
	// Title of the contact.
	Title *string `json:"title,omitempty"`
	// Email address of the contact.
	Email *string `json:"email,omitempty"`
	// PhoneIDD is the international dialing code, e.g. 55 for Brazil.
	PhoneIDD *int `json:"phoneIdd,omitempty"`
	// PhoneNumber excluding the IDD.
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	// DocumentType is the kind of identity document.
	DocumentType IdentityDocumentType `json:"documentType,omitempty"`
	// DocumentValue is the identity document number.
	DocumentValue *string `json:"documentValue,omitempty"`
	// Description of the contact (up to 2000 characters).
	Description *string `json:"description,omitempty"`
}

// ContactListParams filter and paginate List.
type ContactListParams struct {
	ListParams
	// IsFavorite filters contacts by favorite status.
	IsFavorite *bool
	// Search matches contacts by id, name, email, or description.
	Search string
}

func (p *ContactListParams) values() url.Values {
	if p == nil {
		p = &ContactListParams{}
	}
	v := p.ListParams.values()
	if p.IsFavorite != nil {
		v.Set("IsFavorite", strconv.FormatBool(*p.IsFavorite))
	}
	if p.Search != "" {
		v.Set("Search", p.Search)
	}
	return v
}

// List retrieves a single page of contacts. Use ListAutoPaging to iterate
// across all pages.
func (s *ContactService) List(ctx context.Context, params *ContactListParams) (*Page[Contact], error) {
	var page Page[Contact]
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/contacts", params.values(), nil, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// ListAutoPaging iterates over all contacts, fetching pages as needed,
// starting at params.PageNumber (or the first page).
func (s *ContactService) ListAutoPaging(ctx context.Context, params *ContactListParams) iter.Seq2[Contact, error] {
	var snapshot ContactListParams
	if params != nil {
		snapshot = *params
	}
	return autoPaging(snapshot.PageNumber, func(pageNumber int) (*Page[Contact], error) {
		p := snapshot
		p.PageNumber = pageNumber
		return s.List(ctx, &p)
	})
}

// Create creates a contact and returns its id.
func (s *ContactService) Create(ctx context.Context, params *ContactParams) (string, error) {
	var out struct {
		ID string `json:"id"`
	}
	if err := s.client.do(ctx, http.MethodPost, "/v1/ecm/contacts", nil, params, &out); err != nil {
		return "", err
	}
	return out.ID, nil
}

// Get retrieves a contact by id.
func (s *ContactService) Get(ctx context.Context, id string) (*Contact, error) {
	var c Contact
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/contacts/"+pathEscape(id), nil, nil, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Update replaces the details of an existing contact.
func (s *ContactService) Update(ctx context.Context, id string, params *ContactParams) error {
	return s.client.do(ctx, http.MethodPut, "/v1/ecm/contacts/"+pathEscape(id), nil, params, nil)
}

// Remove permanently deletes a contact. This operation is irreversible.
func (s *ContactService) Remove(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodDelete, "/v1/ecm/contacts/"+pathEscape(id), nil, nil, nil)
}

// Favorite marks a contact as a favorite for the current user only.
func (s *ContactService) Favorite(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/contacts/"+pathEscape(id)+"/favorite", nil, nil, nil)
}

// Unfavorite removes the favorite status of a contact for the current user only.
func (s *ContactService) Unfavorite(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/contacts/"+pathEscape(id)+"/unfavorite", nil, nil, nil)
}
