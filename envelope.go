package signater

import (
	"context"
	"net/http"
	"time"
)

// EnvelopeService handles the Envelope endpoints. An envelope groups one or
// more documents with one or more signers and drives the signing workflow
// from Draft through Published to a final state.
type EnvelopeService struct {
	client *Client
}

// EnvelopeStatus is the lifecycle state of an envelope. Draft is the initial
// state; Signed, Cancelled and Rejected are final states.
type EnvelopeStatus string

// Envelope lifecycle states.
const (
	EnvelopeStatusDraft                     EnvelopeStatus = "Draft"
	EnvelopeStatusPublishScheduled          EnvelopeStatus = "PublishScheduled"
	EnvelopeStatusPublished                 EnvelopeStatus = "Published"
	EnvelopeStatusHold                      EnvelopeStatus = "Hold"
	EnvelopeStatusExpired                   EnvelopeStatus = "Expired"
	EnvelopeStatusCancelled                 EnvelopeStatus = "Cancelled"
	EnvelopeStatusCancelledBySignerMfaError EnvelopeStatus = "CancelledBySignerMfaError"
	EnvelopeStatusRejected                  EnvelopeStatus = "Rejected"
	EnvelopeStatusSigned                    EnvelopeStatus = "Signed"
)

// SignerStatus is the state of one signer within an envelope.
type SignerStatus string

// Signer states.
const (
	SignerStatusNone                SignerStatus = "None"
	SignerStatusReadyToReview       SignerStatus = "ReadyToReview"
	SignerStatusQueued              SignerStatus = "Queued"
	SignerStatusApproved            SignerStatus = "Approved"
	SignerStatusRejected            SignerStatus = "Rejected"
	SignerStatusCancelledByMfaError SignerStatus = "CancelledByMfaError"
)

// SignMarkType is the kind of mark a signer places on a document.
type SignMarkType string

// Sign mark types.
const (
	SignMarkInvisible             SignMarkType = "Invisible"
	SignMarkSignature             SignMarkType = "Signature"
	SignMarkLegacyRectangleRubric SignMarkType = "LegacyRectangleRubric"
	SignMarkText                  SignMarkType = "Text"
	SignMarkLongText              SignMarkType = "LongText"
	SignMarkDate                  SignMarkType = "Date"
	SignMarkNumber                SignMarkType = "Number"
	SignMarkCurrency              SignMarkType = "Currency"
	SignMarkCheckbox              SignMarkType = "Checkbox"
	SignMarkRadioButton           SignMarkType = "RadioButton"
	SignMarkDropdown              SignMarkType = "Dropdown"
	SignMarkAttachment            SignMarkType = "Attachment"
	SignMarkRubric                SignMarkType = "Rubric"
)

// CommunicationMode controls how a signer receives messages on one channel.
type CommunicationMode string

// Communication modes.
const (
	CommunicationModeNone           CommunicationMode = "None"
	CommunicationModeInvitationOnly CommunicationMode = "InvitationOnly"
	CommunicationModeFull           CommunicationMode = "Full"
)

// DocumentOrigin tells how an envelope document was created.
type DocumentOrigin string

// Document origins.
const (
	DocumentOriginUpload   DocumentOrigin = "Upload"
	DocumentOriginTemplate DocumentOrigin = "Template"
)

// SignerActionType is the kind of audit action recorded for a signer.
type SignerActionType string

// Signer action types.
const (
	SignerActionView                   SignerActionType = "View"
	SignerActionApproval               SignerActionType = "Approval"
	SignerActionRejection              SignerActionType = "Rejection"
	SignerActionCancellationByMfaError SignerActionType = "CancellationByMfaError"
	SignerActionConfirmSignMark        SignerActionType = "ConfirmSignMark"
	SignerActionUnconfirmSignMark      SignerActionType = "UnconfirmSignMark"
	SignerActionEmailMfaRequest        SignerActionType = "EmailMfaRequest"
	SignerActionSmsMfaRequest          SignerActionType = "SmsMfaRequest"
	SignerActionWhatsAppMfaRequest     SignerActionType = "WhatsAppMfaRequest"
	SignerActionPixMfaRequest          SignerActionType = "PixMfaRequest"
	SignerActionUnsubscribe            SignerActionType = "Unsubscribe"
)

// Envelope is the full envelope representation returned by Get.
type Envelope struct {
	// ID is the envelope's unique identifier.
	ID string `json:"id"`
	// Status is the envelope's lifecycle state.
	Status EnvelopeStatus `json:"status"`
	// Name is the envelope's name.
	Name string `json:"name"`
	// VaultID is the vault the envelope is stored in.
	VaultID string `json:"vaultId"`
	// VaultName is the name of that vault.
	VaultName string `json:"vaultName"`
	// PrivateDescription is visible only inside the account.
	PrivateDescription string `json:"privateDescription"`
	// PublicDescription is visible to signers.
	PublicDescription string `json:"publicDescription"`
	// Message is emailed to signers.
	Message string `json:"message"`
	// Language of the envelope.
	Language Language `json:"language"`
	// MarkupOrientation controls where document information is stamped.
	MarkupOrientation MarkupOrientation `json:"markupOrientation"`
	// SignInOrder makes signers sign sequentially.
	SignInOrder bool `json:"signInOrder"`
	// ReviewReminder enables review reminder emails to signers.
	ReviewReminder bool `json:"reviewReminder"`
	// ExpirationReminder enables expiration reminder emails to signers.
	ExpirationReminder bool `json:"expirationReminder"`
	// ExpiresAt is when the envelope expires (UTC), if set.
	ExpiresAt *time.Time `json:"expiresAtUtc"`
	// HasScheduledPublish reports whether a publication is scheduled.
	HasScheduledPublish bool `json:"hasScheduledPublish"`
	// ToBePublishedAt is when the envelope is scheduled to publish (UTC).
	ToBePublishedAt *time.Time `json:"toBePublishedAtUtc"`
	// JurisdictionCountryCode is the ISO 3166-1 alpha-2 jurisdiction country.
	JurisdictionCountryCode string `json:"jurisdictionCountryCode"`
	// RedirectURL receives the signer after completing their action.
	RedirectURL string `json:"redirectUrl"`
	// IsAiChatForSignersEnabled enables the AI chat for signers.
	IsAiChatForSignersEnabled bool `json:"isAiChatForSignersEnabled"`
	// IsLtvEnabled enables long-term validation (LTV).
	IsLtvEnabled bool `json:"isLtvEnabled"`
	// IsSandbox reports whether the envelope lives in the sandbox environment.
	IsSandbox bool `json:"isSandbox"`
	// HasActionsBeingProcessed reports whether actions are being processed.
	HasActionsBeingProcessed bool `json:"hasActionsBeingProcessed"`
	// CertifiedAt is when the certificate was last processed (UTC), if ever.
	CertifiedAt *time.Time `json:"certifiedAtUtc"`
	// IsRemoved reports whether the envelope is in the trash.
	IsRemoved bool `json:"isRemoved"`
	// RemovedByID is the id of the user who removed the envelope.
	RemovedByID string `json:"removedById"`
	// RemovedBy is the name of the user who removed the envelope.
	RemovedBy string `json:"removedBy"`
	// RemovedAt is when the envelope was removed (UTC), if ever.
	RemovedAt *time.Time `json:"removedAtUtc"`
	// OwnerID is the id of the envelope's owner.
	OwnerID string `json:"ownerId"`
	// OwnerName is the name of the envelope's owner.
	OwnerName string `json:"ownerName"`
	// OwnerAvatar is the avatar URL of the envelope's owner.
	OwnerAvatar string `json:"ownerAvatar"`
	// CreatedAt is when the envelope was created (UTC).
	CreatedAt time.Time `json:"createdAtUtc"`
	// CreatedByID is the id of the creating user.
	CreatedByID string `json:"createdById"`
	// CreatedByName is the name of the creating user.
	CreatedByName string `json:"createdByName"`
	// CreatedByAvatar is the avatar URL of the creating user.
	CreatedByAvatar string `json:"createdByAvatar"`
	// UpdatedAt is when the envelope was last updated (UTC).
	UpdatedAt time.Time `json:"updatedAtUtc"`
	// LastUpdateByID is the id of the last updating user.
	LastUpdateByID string `json:"lastUpdateById"`
	// UpdatedByName is the name of the last updating user.
	UpdatedByName string `json:"updatedByName"`
	// UpdatedByAvatar is the avatar URL of the last updating user.
	UpdatedByAvatar string `json:"updatedByAvatar"`
	// Documents are the envelope's documents.
	Documents []EnvelopeDocument `json:"documents"`
	// Signers are the envelope's signers.
	Signers []EnvelopeSigner `json:"signers"`
}

// EnvelopeDocument is one document inside an envelope, as returned by Get.
type EnvelopeDocument struct {
	// ID is the document's unique identifier.
	ID string `json:"id"`
	// Index is the document's 0-based position.
	Index int `json:"index"`
	// Origin tells whether the document was uploaded or template-generated.
	Origin DocumentOrigin `json:"origin"`
	// Name is the document's name.
	Name string `json:"name"`
	// PrivateDescription is visible only inside the account.
	PrivateDescription string `json:"privateDescription"`
	// PublicDescription is visible to signers.
	PublicDescription string `json:"publicDescription"`
	// OriginalFileSize is the size of the original file in bytes.
	OriginalFileSize int64 `json:"originalFileSize"`
	// PageSizes lists the pixel dimensions of each page.
	PageSizes []PageSize `json:"pageSizes"`
}

// EnvelopeSigner is one signer inside an envelope, as returned by Get.
type EnvelopeSigner struct {
	// ID is the signer's unique identifier.
	ID string `json:"id"`
	// Index is the signer's 0-based signing position when SignInOrder is set.
	Index *int `json:"index"`
	// Status is the signer's state.
	Status SignerStatus `json:"status"`
	// Name is the signer's name.
	Name string `json:"name"`
	// Email is the signer's email, if any.
	Email string `json:"email"`
	// Title is the signer's title, e.g. "Lawyer".
	Title string `json:"title"`
	// HadApproved reports whether the signer approved the envelope.
	HadApproved bool `json:"hadApproved"`
	// ApprovedAt is when the signer approved (UTC), if they did.
	ApprovedAt *time.Time `json:"approvedAtUtc"`
	// HadRejected reports whether the signer rejected the envelope.
	HadRejected bool `json:"hadRejected"`
	// RejectedAt is when the signer rejected (UTC), if they did.
	RejectedAt *time.Time `json:"rejectedAtUtc"`
	// PhoneIDD is the signer's international dialing code.
	PhoneIDD *int `json:"phoneIdd"`
	// PhoneNumber is the signer's phone number excluding the IDD.
	PhoneNumber string `json:"phoneNumber"`
	// DocumentType is the signer's identity document kind.
	DocumentType IdentityDocumentType `json:"documentType"`
	// DocumentValue is the signer's identity document number.
	DocumentValue string `json:"documentValue"`
	// Passcode challenges the signer, when passcode validation is enforced.
	Passcode string `json:"passcode"`
	// PasscodeHint is shown to the signer during signing.
	PasscodeHint string `json:"passcodeHint"`
	// EmailCommunicationMode controls email messages to the signer.
	EmailCommunicationMode CommunicationMode `json:"emailCommunicationMode"`
	// SmsCommunicationMode controls SMS messages to the signer.
	SmsCommunicationMode CommunicationMode `json:"smsCommunicationMode"`
	// WhatsAppCommunicationMode controls WhatsApp messages to the signer.
	WhatsAppCommunicationMode CommunicationMode `json:"whatsAppCommunicationMode"`
	// ShouldEnforceEmailValidation requires email OTP validation.
	ShouldEnforceEmailValidation bool `json:"shouldEnforceEmailValidation"`
	// ShouldEnforceSmsValidation requires SMS OTP validation.
	ShouldEnforceSmsValidation bool `json:"shouldEnforceSmsValidation"`
	// ShouldEnforceWhatsAppValidation requires WhatsApp OTP validation.
	ShouldEnforceWhatsAppValidation bool `json:"shouldEnforceWhatsAppValidation"`
	// ShouldEnforcePixValidation requires PIX validation.
	ShouldEnforcePixValidation bool `json:"shouldEnforcePixValidation"`
	// ShouldEnforceCustomDigitalCertificateValidation requires a custom
	// digital certificate.
	ShouldEnforceCustomDigitalCertificateValidation bool `json:"shouldEnforceCustomDigitalCertificateValidation"`
	// ShouldEnforcePasscodeValidation requires passcode validation.
	ShouldEnforcePasscodeValidation bool `json:"shouldEnforcePasscodeValidation"`
	// ShouldEnforceSimpleSelfieValidation requires a simple selfie capture.
	ShouldEnforceSimpleSelfieValidation bool `json:"shouldEnforceSimpleSelfieValidation"`
	// SignMarks are the signer's marks on the documents.
	SignMarks []SignMark `json:"signMarks"`
	// Actions is the signer's audit trail.
	Actions []SignerAction `json:"actions"`
}

// SignMark is one mark a signer places on a document, as returned by Get.
type SignMark struct {
	// ID is the sign mark's unique identifier.
	ID string `json:"id"`
	// DocumentID is the document the mark belongs to.
	DocumentID string `json:"documentId"`
	// Type is the kind of mark.
	Type SignMarkType `json:"type"`
	// Page is the document page number.
	Page int `json:"page"`
	// X is the pixel X coordinate of the mark's top-left corner.
	X int `json:"x"`
	// Y is the pixel Y coordinate of the mark's top-left corner.
	Y int `json:"y"`
	// Rotation of the mark.
	Rotation int `json:"rotation"`
	// Width of the mark in pixels.
	Width int `json:"width"`
	// Height of the mark in pixels.
	Height int `json:"height"`
	// IsRequired reports whether the mark must be completed.
	IsRequired bool `json:"isRequired"`
	// SchemaJSON is the mark's schema as JSON (Text, Dropdown, etc.).
	SchemaJSON string `json:"schemaJson"`
	// ValueJSON is the mark's value as JSON.
	ValueJSON string `json:"valueJson"`
}

// SignerAction is one audit trail entry of a signer.
type SignerAction struct {
	// ID is the action's unique identifier.
	ID string `json:"id"`
	// CreatedAt is when the action happened (UTC).
	CreatedAt time.Time `json:"createdAtUtc"`
	// Type is the kind of action.
	Type SignerActionType `json:"type"`
	// IP is the signer's IP address at action time.
	IP string `json:"ip"`
	// UserAgent is the signer's user agent at action time.
	UserAgent string `json:"userAgent"`
	// Location is the signer's resolved location at action time.
	Location string `json:"location"`
	// SignMarkID is the targeted sign mark (Confirm/UnconfirmSignMark only).
	SignMarkID string `json:"signMarkId"`
	// SignMarkType is the type of the targeted sign mark.
	SignMarkType SignMarkType `json:"signMarkType"`
	// SignMarkSchemaJSON is the schema of the targeted sign mark.
	SignMarkSchemaJSON string `json:"signMarkSchemaJSON"`
	// ValueJSON is a snapshot of the value at action time, for audit display.
	ValueJSON string `json:"valueJson"`
}

// Create creates an envelope with its documents and signers and returns its
// id. The envelope starts in Draft.
func (s *EnvelopeService) Create(ctx context.Context, params *EnvelopeParams) (string, error) {
	var out struct {
		ID string `json:"id"`
	}
	if err := s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes", nil, params, &out); err != nil {
		return "", err
	}
	return out.ID, nil
}

// Get retrieves an envelope by id, including documents, signers, sign marks
// and the signers' audit trails.
func (s *EnvelopeService) Get(ctx context.Context, id string) (*Envelope, error) {
	var e Envelope
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/envelopes/"+pathEscape(id), nil, nil, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// Update replaces an envelope's details, documents and signers. Documents
// cannot change once any signer approved, and approved signers cannot be
// changed.
func (s *EnvelopeService) Update(ctx context.Context, id string, params *EnvelopeParams) error {
	return s.client.do(ctx, http.MethodPut, "/v1/ecm/envelopes/"+pathEscape(id), nil, params, nil)
}

// Remove moves an envelope to the trash without changing its status. It can
// be recovered with Restore.
func (s *EnvelopeService) Remove(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodDelete, "/v1/ecm/envelopes/"+pathEscape(id), nil, nil, nil)
}

// Restore recovers a removed envelope from the trash into the given vault,
// keeping its status, documents and signers unchanged.
func (s *EnvelopeService) Restore(ctx context.Context, id, vaultID string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/restore/to/"+pathEscape(vaultID), nil, nil, nil)
}

// Publish makes an envelope available to its signers, who receive access
// links. Requires at least one signer and one document, and sufficient
// credits — insufficient credits surface as a *PaymentRequiredError.
func (s *EnvelopeService) Publish(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/publish", nil, nil, nil)
}

// Unschedule cancels a scheduled publication, moving the envelope to Hold so
// it can be rescheduled or published directly later.
func (s *EnvelopeService) Unschedule(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/unschedule", nil, nil, nil)
}

// Hold pauses an envelope's approval process, preventing signers from
// approving or rejecting until it is published again.
func (s *EnvelopeService) Hold(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/hold", nil, nil, nil)
}

// Cancel terminates an envelope's approval process irreversibly.
func (s *EnvelopeService) Cancel(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/cancel", nil, nil, nil)
}

// Rename changes an envelope's name without affecting its contents or
// signers.
func (s *EnvelopeService) Rename(ctx context.Context, id, name string) error {
	body := struct {
		Name string `json:"name"`
	}{Name: name}
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/rename", nil, body, nil)
}

// Move relocates an envelope to another vault.
func (s *EnvelopeService) Move(ctx context.Context, id, vaultID string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/to/"+pathEscape(vaultID), nil, nil, nil)
}

// ListAvailableOwners retrieves the users eligible to own the envelope.
func (s *EnvelopeService) ListAvailableOwners(ctx context.Context, id string) ([]AccountUser, error) {
	var out struct {
		Items []AccountUser `json:"items"`
	}
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/envelopes/"+pathEscape(id)+"/owners", nil, nil, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

// ChangeOwner transfers the envelope's ownership — and all its notifications
// and responsibilities — to another user.
func (s *EnvelopeService) ChangeOwner(ctx context.Context, id, ownerID string) error {
	body := struct {
		OwnerID string `json:"ownerId"`
	}{OwnerID: ownerID}
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/user-account-owner", nil, body, nil)
}

// ReinviteToReview manually re-sends the reminder email to signers still in
// ReadyToReview. Signers without email receive nothing — use
// CreateSignatureLink for those.
func (s *EnvelopeService) ReinviteToReview(ctx context.Context, id string) error {
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/manual-reinvite-to-review", nil, nil, nil)
}

// CreateSignatureLink generates a time-limited link a signer can use to
// review and sign, useful for signers without a registered email.
func (s *EnvelopeService) CreateSignatureLink(ctx context.Context, envelopeID, signerID string) (string, error) {
	var out struct {
		URL string `json:"url"`
	}
	path := "/v1/ecm/envelopes/" + pathEscape(envelopeID) + "/signers/" + pathEscape(signerID) + "/signature-link"
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return "", err
	}
	return out.URL, nil
}

// CertificateFileURL returns the URL to download the envelope's certificate,
// which contains the audit logs (approvals, timestamps and related details).
func (s *EnvelopeService) CertificateFileURL(ctx context.Context, id string) (string, error) {
	var out struct {
		URL string `json:"url"`
	}
	if err := s.client.do(ctx, http.MethodGet, "/v1/ecm/envelopes/"+pathEscape(id)+"/certificate", nil, nil, &out); err != nil {
		return "", err
	}
	return out.URL, nil
}

// ProcessCertificate requests regeneration of the envelope's certificate so
// it reflects the latest state. The system also regenerates it automatically
// whenever necessary.
func (s *EnvelopeService) ProcessCertificate(ctx context.Context, id string) error {
	// The spec declares an (empty) JSON request body for this verb — the only
	// lifecycle verb that does — so send {} to satisfy body-required bindings;
	// the sandbox API accepts the empty object.
	return s.client.do(ctx, http.MethodPost, "/v1/ecm/envelopes/"+pathEscape(id)+"/certificate", nil, struct{}{}, nil)
}

// SignerAttachmentURL returns the pre-signed URL of a file a signer uploaded
// to an Attachment sign mark, without downloading it.
func (s *EnvelopeService) SignerAttachmentURL(ctx context.Context, envelopeID, signMarkID, attachmentExternalID string) (string, error) {
	path := "/v1/ecm/envelopes/" + pathEscape(envelopeID) + "/sign-marks/" + pathEscape(signMarkID) +
		"/attachments/" + pathEscape(attachmentExternalID) + "/download"
	return s.client.doRedirectURL(ctx, path)
}

// SignerSelfieURL returns the short-lived pre-signed URL of the simple-selfie
// image captured during signer validation, without downloading it. Only
// available for signers whose envelope enforced simple-selfie validation.
func (s *EnvelopeService) SignerSelfieURL(ctx context.Context, envelopeID, signerID string) (string, error) {
	path := "/v1/ecm/envelopes/" + pathEscape(envelopeID) + "/signers/" + pathEscape(signerID) + "/simple-selfie"
	return s.client.doRedirectURL(ctx, path)
}
