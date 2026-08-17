package signater

import (
	"context"
	"net/http"
)

// TemplateService handles the Template endpoints. Templates hold reusable
// document structure and predefined fields so new documents can be generated
// quickly and consistently.
type TemplateService struct {
	client *Client
}

// TemplateFieldType is the data type of a template field when creating or
// updating a template.
type TemplateFieldType string

// Template field types accepted when creating or updating templates.
const (
	TemplateFieldBrCep    TemplateFieldType = "BrCep"
	TemplateFieldBrCpf    TemplateFieldType = "BrCpf"
	TemplateFieldBrCnpj   TemplateFieldType = "BrCnpj"
	TemplateFieldDateOnly TemplateFieldType = "DateOnly"
	TemplateFieldTimeOnly TemplateFieldType = "TimeOnly"
	TemplateFieldDateTime TemplateFieldType = "DateTime"
	TemplateFieldEmail    TemplateFieldType = "Email"
	TemplateFieldNumber   TemplateFieldType = "Number"
	TemplateFieldText     TemplateFieldType = "Text"
	TemplateFieldCustom   TemplateFieldType = "Custom"
)

// TemplateFieldDisplayType is the presentation type of a template field as
// returned by Get. It uses a different vocabulary from TemplateFieldType
// (the create/update request enum) — both are published as-is by the API.
type TemplateFieldDisplayType string

// Template field display types returned by Get.
const (
	TemplateFieldDisplayText        TemplateFieldDisplayType = "Text"
	TemplateFieldDisplayLongText    TemplateFieldDisplayType = "LongText"
	TemplateFieldDisplayNumber      TemplateFieldDisplayType = "Number"
	TemplateFieldDisplayCurrency    TemplateFieldDisplayType = "Currency"
	TemplateFieldDisplayDate        TemplateFieldDisplayType = "Date"
	TemplateFieldDisplayCheckbox    TemplateFieldDisplayType = "Checkbox"
	TemplateFieldDisplayRadioButton TemplateFieldDisplayType = "RadioButton"
	TemplateFieldDisplayDropdown    TemplateFieldDisplayType = "Dropdown"
)

// Template is a document template with its fields configuration.
type Template struct {
	// ID is the template's unique identifier.
	ID string `json:"id"`
	// VaultID is the vault the template is stored in.
	VaultID string `json:"vaultId"`
	// VaultName is the name of that vault.
	VaultName string `json:"vaultName"`
	// Name is the template's name.
	Name string `json:"name"`
	// Description is the template's description.
	Description string `json:"description"`
	// Fields are the template's predefined fields.
	Fields []TemplateField `json:"fields"`
}

// TemplateField is one predefined field of a template, as returned by Get.
type TemplateField struct {
	// ID is the field's unique identifier.
	ID string `json:"id"`
	// Alias uniquely identifies the field when filling it.
	Alias string `json:"alias"`
	// IsRequired reports whether the field must be filled.
	IsRequired bool `json:"isRequired"`
	// Type is the field's presentation type.
	Type TemplateFieldDisplayType `json:"type"`
	// Schema defines the field's label, placeholder, options and constraints.
	Schema *TemplateFieldSchema `json:"schema"`
	// Index is the field's 0-based position.
	Index int `json:"index"`
}

// TemplateFieldSchema defines a field's label, placeholder, options and
// constraints.
type TemplateFieldSchema struct {
	// Placeholder shown while the field is empty.
	Placeholder string `json:"placeholder"`
	// Label shown next to the field.
	Label string `json:"label"`
	// Description of the field.
	Description string `json:"description"`
	// Options for choice fields (dropdown, radio).
	Options []string `json:"options"`
	// Currency code for currency fields.
	Currency string `json:"currency"`
	// DefaultValue prefilled into the field.
	DefaultValue string `json:"defaultValue"`
	// Constraints validate the field's value.
	Constraints *TemplateFieldConstraints `json:"constraints"`
	// DateFormat for date fields.
	DateFormat string `json:"dateFormat"`
}

// TemplateFieldProfile is a predefined validation profile for a template
// field.
type TemplateFieldProfile string

// Template field validation profiles.
const (
	TemplateFieldProfileNone  TemplateFieldProfile = "None"
	TemplateFieldProfileEmail TemplateFieldProfile = "Email"
	TemplateFieldProfilePhone TemplateFieldProfile = "Phone"
	TemplateFieldProfileCpf   TemplateFieldProfile = "Cpf"
	TemplateFieldProfileCnpj  TemplateFieldProfile = "Cnpj"
	TemplateFieldProfileCep   TemplateFieldProfile = "Cep"
)

// TemplateFieldConstraints validate a template field's value. Nil pointers
// mean the constraint is not set.
type TemplateFieldConstraints struct {
	// MinLength is the minimum value length.
	MinLength *int `json:"minLength"`
	// MaxLength is the maximum value length.
	MaxLength *int `json:"maxLength"`
	// Pattern is a regular expression the value must match.
	Pattern string `json:"pattern"`
	// Min is the minimum numeric value.
	Min *float64 `json:"min"`
	// Max is the maximum numeric value.
	Max *float64 `json:"max"`
	// DecimalPlaces is the number of decimal places allowed.
	DecimalPlaces *int `json:"decimalPlaces"`
	// AllowNegative permits negative numbers.
	AllowNegative *bool `json:"allowNegative"`
	// AllowMultiple permits selecting multiple options.
	AllowMultiple *bool `json:"allowMultiple"`
	// MinDate is the earliest accepted date.
	MinDate string `json:"minDate"`
	// MaxDate is the latest accepted date.
	MaxDate string `json:"maxDate"`
	// Profile applies a predefined validation profile.
	Profile TemplateFieldProfile `json:"profile"`
}

// TemplateParams are the fields accepted when creating or updating a template.
type TemplateParams struct {
	// VaultID is the vault to store the template in. Required.
	VaultID string `json:"vaultId"`
	// Name of the template (2-200 characters). Required.
	Name string `json:"name"`
	// Description of the template (up to 4000 characters).
	Description *string `json:"description,omitempty"`
	// Payload is the template's document payload.
	Payload *string `json:"payload,omitempty"`
	// Fields are the template's predefined fields.
	Fields []TemplateFieldParams `json:"fields,omitempty"`
}

// TemplateFieldParams define one template field when creating or updating a
// template.
type TemplateFieldParams struct {
	// ID of an existing field. Only meaningful on update; leave nil to create
	// a new field.
	ID *string `json:"id,omitempty"`
	// Name of the field (2-200 characters). Required.
	Name string `json:"name"`
	// Alias uniquely identifies the field; alphanumerics, dashes and
	// underscores only. Required.
	Alias string `json:"alias"`
	// Description of the field (up to 4000 characters).
	Description *string `json:"description,omitempty"`
	// IsRequired makes the field mandatory.
	IsRequired bool `json:"isRequired"`
	// Type is the field's data type.
	Type TemplateFieldType `json:"type,omitempty"`
	// CustomRegex validates the value when Type is TemplateFieldCustom.
	CustomRegex *string `json:"customRegex,omitempty"`
}

// TemplateRenameParams are the fields accepted by Rename.
type TemplateRenameParams struct {
	// Name is the template's new name (2-200 characters). Required.
	Name string `json:"name"`
	// Description is the template's new description.
	Description *string `json:"description,omitempty"`
}

// Create creates a template and returns its id.
func (s *TemplateService) Create(ctx context.Context, params *TemplateParams) (string, error) {
	var out struct {
		ID string `json:"id"`
	}
	if err := s.client.do(ctx, http.MethodPost, "/v1/ecm/templates", nil, params, &out); err != nil {
		return "", err
	}
	return out.ID, nil
}

// Get retrieves a template by id, including its fields configuration.
func (s *TemplateService) Get(ctx context.Context, id string) (*Template, error) {
	var t Template
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/templates/"+pathEscape(id), nil, nil, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// Update replaces a template's content, structure and fields.
func (s *TemplateService) Update(ctx context.Context, id string, params *TemplateParams) error {
	return s.client.do(ctx, http.MethodPut, "/v1/ecm/templates/"+pathEscape(id), nil, params, nil)
}

// Remove moves a template to the trash. It can be recovered with Restore.
func (s *TemplateService) Remove(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodDelete, "/v1/ecm/templates/"+pathEscape(id), nil, nil, nil)
}

// Rename changes a template's name and/or description without altering its
// content, structure or fields.
func (s *TemplateService) Rename(ctx context.Context, id string, params *TemplateRenameParams) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/templates/"+pathEscape(id)+"/rename", nil, params, nil)
}

// Move relocates a template to another vault.
func (s *TemplateService) Move(ctx context.Context, id, vaultID string) error {
	body := struct {
		VaultID string `json:"vaultId"`
	}{VaultID: vaultID}
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/templates/"+pathEscape(id)+"/move", nil, body, nil)
}

// Restore recovers a previously removed template from the trash into the
// given vault.
func (s *TemplateService) Restore(ctx context.Context, id, vaultID string) error {
	body := struct {
		VaultID string `json:"vaultId"`
	}{VaultID: vaultID}
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/templates/"+pathEscape(id)+"/restore", nil, body, nil)
}
