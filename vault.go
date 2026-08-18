package signater

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// VaultService handles the Vault endpoints. A vault is the secure container
// that stores envelopes and templates.
type VaultService struct {
	client *Client
}

// VaultType defines who can access a vault's contents.
type VaultType string

// Vault access types.
const (
	// VaultTypeAccount is accessible to all users within the account.
	VaultTypeAccount VaultType = "Account"
	// VaultTypeUserAccount is a personal vault, accessible only to its creator.
	VaultTypeUserAccount VaultType = "UserAccount"
	// VaultTypeUserAccountGroup is accessible to an explicit list of members.
	VaultTypeUserAccountGroup VaultType = "UserAccountGroup"
)

// UserRole is the role of a user within the account.
type UserRole string

// User roles.
const (
	UserRoleAdministrator UserRole = "Administrator"
	UserRoleUser          UserRole = "User"
)

// VaultMember is a user with access to a vault. List responses additionally
// populate Email and IsActive.
type VaultMember struct {
	// ID is the member's user account id.
	ID string `json:"id"`
	// Name is the member's name.
	Name string `json:"name"`
	// Email is the member's email. Only populated by list endpoints.
	Email string `json:"email"`
	// Avatar is the member's avatar URL.
	Avatar string `json:"avatar"`
	// IsActive reports whether the member is active. Only populated by list
	// endpoints.
	IsActive bool `json:"isActive"`
}

// Vault is a secure storage container. List responses additionally populate
// the audit fields, IsFavorite and CanBeRemoved, which Get does not return.
type Vault struct {
	// ID is the vault's unique identifier.
	ID string `json:"id"`
	// Name is the vault's name.
	Name string `json:"name"`
	// Type defines who can access the vault.
	Type VaultType `json:"type"`
	// Description is the vault's description. Only populated by Get.
	Description string `json:"description"`
	// UserAccountMembers are the users with access to the vault.
	UserAccountMembers []VaultMember `json:"userAccountMembers"`
	// IsFavorite reports whether the current user favorited the vault. Only
	// populated by list endpoints.
	IsFavorite bool `json:"isFavorite"`
	// CanBeRemoved reports whether the vault can be removed (only empty
	// vaults can). Only populated by list endpoints.
	CanBeRemoved bool `json:"canBeRemoved"`
	// CreatedAt is when the vault was created (UTC). Only populated by list
	// endpoints.
	CreatedAt time.Time `json:"createdAtUtc"`
	// CreatedByName is the name of the user who created the vault. Only
	// populated by list endpoints.
	CreatedByName string `json:"createdByName"`
	// CreatedByAvatar is the avatar URL of the creating user. Only populated
	// by list endpoints.
	CreatedByAvatar string `json:"createdByAvatar"`
	// UpdatedAt is when the vault was last updated (UTC). Only populated by
	// list endpoints.
	UpdatedAt time.Time `json:"updatedAtUtc"`
	// UpdatedByName is the name of the user who last updated the vault. Only
	// populated by list endpoints.
	UpdatedByName string `json:"updatedByName"`
	// UpdatedByAvatar is the avatar URL of the last updating user. Only
	// populated by list endpoints.
	UpdatedByAvatar string `json:"updatedByAvatar"`
}

// AccountUser is a user returned by the Owners and Members endpoints.
type AccountUser struct {
	// ID is the user's account id.
	ID string `json:"id"`
	// Name is the user's name.
	Name string `json:"name"`
	// IsActive reports whether the user is active.
	IsActive bool `json:"isActive"`
	// IsCurrent reports whether this is the authenticated user.
	IsCurrent bool `json:"isCurrent"`
	// Role is the user's role within the account.
	Role UserRole `json:"role"`
	// Avatar is the user's avatar URL.
	Avatar string `json:"avatar"`
}

// VaultParams are the fields accepted when creating or updating a vault.
//
// Create: Name is required; Type defaults to the API's default when empty;
// MemberIDs is required for VaultTypeUserAccountGroup (add yourself
// explicitly to stay a member). Update (PUT): Name and Type are required and
// the type may be changed. UpdatePartial (PATCH): the type cannot be changed.
type VaultParams struct {
	// Name of the vault (2-200 characters). Required.
	Name string `json:"name"`
	// Description of the vault (up to 400 characters).
	Description *string `json:"description,omitempty"`
	// Type defines who can access the vault.
	Type VaultType `json:"type,omitempty"`
	// MemberIDs are the user account ids granted access (UserAccountGroup
	// vaults only). An empty non-nil slice clears the member list on Update;
	// nil serializes as null.
	MemberIDs []string `json:"memberIds"`
}

// VaultListParams filter and paginate the vault list endpoints.
type VaultListParams struct {
	ListParams
	// IsFavorite filters vaults by favorite status.
	IsFavorite *bool
	// Search matches vaults by id or name.
	Search string
}

func (p *VaultListParams) values() url.Values {
	if p == nil {
		p = &VaultListParams{}
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

// Create creates a vault and returns its id.
func (s *VaultService) Create(ctx context.Context, params *VaultParams) (string, error) {
	var out struct {
		ID string `json:"id"`
	}
	if err := s.client.do(ctx, http.MethodPost, "/v1/ecm/vaults", nil, params, &out); err != nil {
		return "", err
	}
	return out.ID, nil
}

// Get retrieves a vault by id, including its members.
func (s *VaultService) Get(ctx context.Context, id string) (*Vault, error) {
	var v Vault
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/vaults/"+pathEscape(id), nil, nil, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// Update replaces the full state of a vault (PUT). The access type may be
// changed; when the resulting type is VaultTypeUserAccountGroup, MemberIDs
// must contain at least one member.
func (s *VaultService) Update(ctx context.Context, id string, params *VaultParams) error {
	return s.client.do(ctx, http.MethodPut, "/v1/ecm/vaults/"+pathEscape(id), nil, params, nil)
}

// UpdatePartial updates a vault's properties (legacy PATCH). The access type
// cannot be changed through this endpoint; prefer Update.
func (s *VaultService) UpdatePartial(ctx context.Context, id string, params *VaultParams) error {
	return s.client.do(ctx, http.MethodPatch, "/v1/ecm/vaults/"+pathEscape(id), nil, params, nil)
}

// Remove permanently deletes a vault. Only vaults without envelopes or
// templates can be removed; this operation is irreversible.
func (s *VaultService) Remove(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodDelete, "/v1/ecm/vaults/"+pathEscape(id), nil, nil, nil)
}

// Owners retrieves the users eligible to own vaults.
//
// The path takes no vault id even though the endpoint description mentions
// "a specific vault" — it matches the published OpenAPI fragment and the
// sandbox API serves it account-wide.
func (s *VaultService) Owners(ctx context.Context) ([]AccountUser, error) {
	var out struct {
		Items []AccountUser `json:"items"`
	}
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/vaults/owners", nil, nil, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

// Members retrieves the users with vault access, including administrators.
//
// The path takes no vault id even though the endpoint description mentions
// "a specific vault" — it matches the published OpenAPI fragment and the
// sandbox API serves it account-wide.
func (s *VaultService) Members(ctx context.Context) ([]AccountUser, error) {
	var out struct {
		Items []AccountUser `json:"items"`
	}
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/vaults/members", nil, nil, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

// ListAccount retrieves a page of all vaults in the account.
func (s *VaultService) ListAccount(ctx context.Context, params *VaultListParams) (*Page[Vault], error) {
	return s.listVaults(ctx, "/v1/ecm/vaults/accounts", params)
}

// ListAccountAutoPaging iterates over all account vaults.
func (s *VaultService) ListAccountAutoPaging(ctx context.Context, params *VaultListParams) iter.Seq2[Vault, error] {
	return s.autoPagingVaults(ctx, "/v1/ecm/vaults/accounts", params)
}

// ListMine retrieves a page of the vaults accessible to the current user.
func (s *VaultService) ListMine(ctx context.Context, params *VaultListParams) (*Page[Vault], error) {
	return s.listVaults(ctx, "/v1/ecm/vaults/user-accounts", params)
}

// ListMineAutoPaging iterates over all vaults accessible to the current user.
func (s *VaultService) ListMineAutoPaging(ctx context.Context, params *VaultListParams) iter.Seq2[Vault, error] {
	return s.autoPagingVaults(ctx, "/v1/ecm/vaults/user-accounts", params)
}

// ListByUser retrieves a page of the vaults accessible to a specific user.
func (s *VaultService) ListByUser(ctx context.Context, userAccountID string, params *VaultListParams) (*Page[Vault], error) {
	return s.listVaults(ctx, "/v1/ecm/vaults/user-accounts/"+pathEscape(userAccountID), params)
}

// ListByUserAutoPaging iterates over all vaults accessible to a specific user.
func (s *VaultService) ListByUserAutoPaging(ctx context.Context, userAccountID string, params *VaultListParams) iter.Seq2[Vault, error] {
	return s.autoPagingVaults(ctx, "/v1/ecm/vaults/user-accounts/"+pathEscape(userAccountID), params)
}

func (s *VaultService) listVaults(ctx context.Context, path string, params *VaultListParams) (*Page[Vault], error) {
	var page Page[Vault]
	if err := s.client.do(ctx, http.MethodGet, path, params.values(), nil, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

func (s *VaultService) autoPagingVaults(ctx context.Context, path string, params *VaultListParams) iter.Seq2[Vault, error] {
	var snapshot VaultListParams
	if params != nil {
		snapshot = *params
	}
	return autoPaging(snapshot.PageNumber, func(pageNumber int) (*Page[Vault], error) {
		p := snapshot
		p.PageNumber = pageNumber
		return s.listVaults(ctx, path, &p)
	})
}

// Favorite marks a vault as a favorite for the current user only.
func (s *VaultService) Favorite(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/vaults/"+pathEscape(id)+"/favorite", nil, nil, nil)
}

// Unfavorite removes the favorite status of a vault for the current user only.
func (s *VaultService) Unfavorite(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/vaults/"+pathEscape(id)+"/unfavorite", nil, nil, nil)
}
