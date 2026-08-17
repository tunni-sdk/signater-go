package signater

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// VaultItemType discriminates the entries returned by VaultService.Contents.
type VaultItemType string

// Vault content item types.
const (
	VaultItemTypeEnvelope VaultItemType = "Envelope"
	VaultItemTypeTemplate VaultItemType = "Template"
)

// VaultSummary describes the vault that Contents was called on.
type VaultSummary struct {
	// ID is the vault's unique identifier.
	ID string `json:"id"`
	// IsSandbox reports whether the vault lives in the sandbox environment.
	IsSandbox bool `json:"isSandbox"`
	// Name is the vault's name.
	Name string `json:"name"`
	// Type defines who can access the vault.
	Type VaultType `json:"type"`
	// CanBeRemoved reports whether the vault can be removed (only empty
	// vaults can).
	CanBeRemoved bool `json:"canBeRemoved"`
	// UserAccountMembers are the users with access to the vault.
	UserAccountMembers []VaultMember `json:"userAccountMembers"`
}

// VaultEnvelopeSummary is the envelope shape returned inside vault contents.
type VaultEnvelopeSummary struct {
	// ID is the envelope's unique identifier.
	ID string `json:"id"`
	// Status is the envelope's lifecycle state.
	Status EnvelopeStatus `json:"status"`
	// Name is the envelope's name.
	Name string `json:"name"`
	// PrivateDescription is the envelope's private description.
	PrivateDescription string `json:"privateDescription"`
	// CanBeRemoved reports whether the envelope can be removed (published
	// envelopes must be held first).
	CanBeRemoved bool `json:"canBeRemoved"`
	// Signers lists the envelope's signers (name and email only).
	Signers []VaultEnvelopeSigner `json:"signers"`
	// CreatedAt is when the envelope was created (UTC).
	CreatedAt time.Time `json:"createdAtUtc"`
	// CreatedByID is the id of the user who created the envelope.
	CreatedByID string `json:"createdById"`
	// CreatedByName is the name of the user who created the envelope.
	CreatedByName string `json:"createdByName"`
	// CreatedByAvatar is the avatar URL of the creating user.
	CreatedByAvatar string `json:"createdByAvatar"`
	// UpdatedAt is when the envelope was last updated (UTC).
	UpdatedAt time.Time `json:"updatedAtUtc"`
	// LastUpdateByID is the id of the user who last updated the envelope.
	LastUpdateByID string `json:"lastUpdateById"`
	// UpdatedByName is the name of the user who last updated the envelope.
	UpdatedByName string `json:"updatedByName"`
	// UpdatedByAvatar is the avatar URL of the last updating user.
	UpdatedByAvatar string `json:"updatedByAvatar"`
}

// VaultEnvelopeSigner is the signer shape returned inside vault contents.
type VaultEnvelopeSigner struct {
	// Name is the signer's name.
	Name string `json:"name"`
	// Email is the signer's email.
	Email string `json:"email"`
}

// VaultTemplateSummary is the template shape returned inside vault contents.
type VaultTemplateSummary struct {
	// ID is the template's unique identifier.
	ID string `json:"id"`
	// Name is the template's name.
	Name string `json:"name"`
	// CreatedAt is when the template was created (UTC).
	CreatedAt time.Time `json:"createdAtUtc"`
	// CreatedByID is the id of the user who created the template.
	CreatedByID string `json:"createdById"`
	// CreatedByName is the name of the user who created the template.
	CreatedByName string `json:"createdByName"`
	// CreatedByAvatar is the avatar URL of the creating user.
	CreatedByAvatar string `json:"createdByAvatar"`
	// UpdatedAt is when the template was last updated (UTC).
	UpdatedAt time.Time `json:"updatedAtUtc"`
	// LastUpdateByID is the id of the user who last updated the template.
	LastUpdateByID string `json:"lastUpdateById"`
	// UpdatedByName is the name of the user who last updated the template.
	UpdatedByName string `json:"updatedByName"`
	// UpdatedByAvatar is the avatar URL of the last updating user.
	UpdatedByAvatar string `json:"updatedByAvatar"`
}

// VaultItem is one entry of a vault's contents: either an envelope or a
// template, indicated by Type.
type VaultItem struct {
	// Type discriminates which of Envelope or Template is set.
	Type VaultItemType `json:"type"`
	// Envelope is set when Type is VaultItemTypeEnvelope.
	Envelope *VaultEnvelopeSummary `json:"envelope"`
	// Template is set when Type is VaultItemTypeTemplate.
	Template *VaultTemplateSummary `json:"template"`
}

// VaultContents is one page of a vault's contents.
type VaultContents struct {
	// Vault describes the vault itself.
	Vault VaultSummary `json:"vault"`
	// Items are the envelopes and templates in this page.
	Items []VaultItem `json:"items"`
	// Pagination locates this page within the full result set.
	Pagination Pagination `json:"pagination"`
}

// VaultContentsParams filter and paginate Contents.
type VaultContentsParams struct {
	ListParams
	// ExcludeEnvelopes omits envelopes from the results.
	ExcludeEnvelopes *bool
	// ExcludeTemplates omits templates from the results.
	ExcludeTemplates *bool
	// Search matches items by id or name.
	Search string
}

func (p *VaultContentsParams) values() url.Values {
	if p == nil {
		p = &VaultContentsParams{}
	}
	v := p.ListParams.values()
	if p.ExcludeEnvelopes != nil {
		v.Set("ExcludeEnvelopes", strconv.FormatBool(*p.ExcludeEnvelopes))
	}
	if p.ExcludeTemplates != nil {
		v.Set("ExcludeTemplates", strconv.FormatBool(*p.ExcludeTemplates))
	}
	if p.Search != "" {
		v.Set("Search", p.Search)
	}
	return v
}

// Contents retrieves a single page of a vault's envelopes and templates. Use
// ContentsAutoPaging to iterate across all pages.
func (s *VaultService) Contents(ctx context.Context, vaultID string, params *VaultContentsParams) (*VaultContents, error) {
	var out VaultContents
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/vaults/"+pathEscape(vaultID)+"/list", params.values(), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ContentsAutoPaging iterates over all items of a vault's contents, fetching
// pages as needed.
func (s *VaultService) ContentsAutoPaging(ctx context.Context, vaultID string, params *VaultContentsParams) iter.Seq2[VaultItem, error] {
	var snapshot VaultContentsParams
	if params != nil {
		snapshot = *params
	}
	return autoPaging(snapshot.PageNumber, func(pageNumber int) (*Page[VaultItem], error) {
		p := snapshot
		p.PageNumber = pageNumber
		contents, err := s.Contents(ctx, vaultID, &p)
		if err != nil {
			return nil, err
		}
		return &Page[VaultItem]{Items: contents.Items, Pagination: contents.Pagination}, nil
	})
}
