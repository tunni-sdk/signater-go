package signater

import "time"

// EnvelopeParams are the fields accepted when creating or updating an
// envelope. VaultID, Name and JurisdictionCountryCode are required. The ID
// fields on signers and sign marks are only meaningful on Update, where they
// address existing records; leave them nil on Create.
type EnvelopeParams struct {
	// VaultID is the vault to store the envelope in. Required.
	VaultID string `json:"vaultId"`
	// Name of the envelope (2-200 characters). Required.
	Name string `json:"name"`
	// JurisdictionCountryCode is the ISO 3166-1 alpha-2 jurisdiction country
	// (e.g. "BR", "US"). Required; on Update it can only change while the
	// envelope is in Draft.
	JurisdictionCountryCode string `json:"jurisdictionCountryCode"`
	// PrivateDescription is visible only inside the account (up to 4000
	// characters).
	PrivateDescription *string `json:"privateDescription,omitempty"`
	// PublicDescription is visible to signers (up to 4000 characters).
	PublicDescription *string `json:"publicDescription,omitempty"`
	// Message is emailed to signers (up to 4000 characters).
	Message *string `json:"message,omitempty"`
	// ReviewReminder enables review reminder emails to signers.
	ReviewReminder bool `json:"reviewReminder"`
	// ExpirationReminder enables expiration reminder emails to signers.
	ExpirationReminder bool `json:"expirationReminder"`
	// ExpiresAt is when the envelope expires. Times are serialized in
	// RFC 3339 with their offset; prefer UTC values (time.Time.UTC).
	ExpiresAt *time.Time `json:"expiresAtUtc,omitempty"`
	// ToBePublishedAt schedules the publication. Times are serialized in
	// RFC 3339 with their offset; prefer UTC values (time.Time.UTC).
	ToBePublishedAt *time.Time `json:"toBePublishedAtUtc,omitempty"`
	// Language of the envelope. REQUIRED by the live API despite the spec
	// marking it optional (verified in sandbox: omitting it returns 400
	// "'Language' has a range of values which does not include '0'").
	Language Language `json:"language,omitempty"`
	// MarkupOrientation controls where document information is stamped.
	MarkupOrientation MarkupOrientation `json:"markupOrientation,omitempty"`
	// SignInOrder makes signers sign sequentially.
	SignInOrder bool `json:"signInOrder"`
	// IsAiChatForSignersEnabled enables the AI chat for signers.
	IsAiChatForSignersEnabled bool `json:"isAiChatForSignersEnabled"`
	// IsLtvEnabled enables long-term validation (LTV).
	IsLtvEnabled bool `json:"isLtvEnabled"`
	// RedirectURL receives the signer after completing their action; the
	// system injects envelopeId and signerId query params (up to 2000
	// characters).
	RedirectURL *string `json:"redirectUrl,omitempty"`
	// Signers of the envelope. At least one is required to publish. On
	// Update (a full replace) an empty non-nil slice clears the list; nil
	// serializes as null.
	Signers []EnvelopeSignerParams `json:"signers"`
	// Documents of the envelope, referencing previously created document ids.
	// At least one is required to publish.
	Documents []EnvelopeDocumentParams `json:"documents"`
}

// EnvelopeDocumentParams attach one previously created document to an
// envelope.
type EnvelopeDocumentParams struct {
	// ID of the document (from DocumentService.Upload or
	// CreateFromTemplate). Required. A document can belong to only one
	// envelope, even one in the trash — upload a new document per envelope.
	ID string `json:"id"`
	// Name of the document (2-200 characters). Required.
	Name string `json:"name"`
	// PrivateDescription of the document (up to 4000 characters).
	PrivateDescription *string `json:"privateDescription,omitempty"`
	// PublicDescription of the document (up to 4000 characters).
	PublicDescription *string `json:"publicDescription,omitempty"`
	// Language overrides the envelope's language for this document.
	Language Language `json:"language,omitempty"`
	// MarkupOrientation overrides the envelope's markup orientation for this
	// document.
	MarkupOrientation MarkupOrientation `json:"markupOrientation,omitempty"`
}

// EnvelopeSignerParams define one signer of an envelope, including the
// authentication factors enforced on them and their sign marks.
type EnvelopeSignerParams struct {
	// ID of an existing signer. Only meaningful on Update; leave nil to add
	// a new signer.
	ID *string `json:"id,omitempty"`
	// Name of the signer (2-200 characters). Required.
	Name string `json:"name"`
	// Email of the signer. Signers without email must sign via
	// EnvelopeService.CreateSignatureLink.
	Email *string `json:"email,omitempty"`
	// Title of the signer, e.g. "Lawyer".
	Title *string `json:"title,omitempty"`
	// EmailCommunicationMode controls email messages to the signer.
	// REQUIRED by the live API (verified in sandbox: omitting any
	// communication mode returns 400 "range of values which does not
	// include '0'"). Use CommunicationModeNone to disable the channel.
	EmailCommunicationMode CommunicationMode `json:"emailCommunicationMode,omitempty"`
	// SmsCommunicationMode controls SMS messages to the signer. REQUIRED by
	// the live API; use CommunicationModeNone to disable the channel.
	SmsCommunicationMode CommunicationMode `json:"smsCommunicationMode,omitempty"`
	// WhatsAppCommunicationMode controls WhatsApp messages to the signer.
	// REQUIRED by the live API; use CommunicationModeNone to disable the
	// channel.
	WhatsAppCommunicationMode CommunicationMode `json:"whatsAppCommunicationMode,omitempty"`
	// PhoneIDD is the signer's international dialing code, e.g. 55 for
	// Brazil.
	PhoneIDD *int `json:"phoneIdd,omitempty"`
	// PhoneNumber excluding the IDD.
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	// DocumentType is the signer's identity document kind.
	DocumentType IdentityDocumentType `json:"documentType,omitempty"`
	// DocumentValue is the signer's identity document number.
	DocumentValue *string `json:"documentValue,omitempty"`
	// ShouldEnforceEmailValidation requires email OTP validation.
	ShouldEnforceEmailValidation bool `json:"shouldEnforceEmailValidation"`
	// ShouldEnforceSmsValidation requires SMS OTP validation.
	ShouldEnforceSmsValidation bool `json:"shouldEnforceSmsValidation"`
	// ShouldEnforceWhatsAppValidation requires WhatsApp OTP validation.
	ShouldEnforceWhatsAppValidation bool `json:"shouldEnforceWhatsAppValidation"`
	// ShouldEnforcePixValidation requires PIX validation.
	ShouldEnforcePixValidation bool `json:"shouldEnforcePixValidation"`
	// ShouldEnforceCustomDigitalCertificateValidation requires a custom
	// digital certificate (eCPF, eCNPJ).
	ShouldEnforceCustomDigitalCertificateValidation bool `json:"shouldEnforceCustomDigitalCertificateValidation"`
	// ShouldEnforceSimpleSelfieValidation requires a simple selfie capture.
	ShouldEnforceSimpleSelfieValidation bool `json:"shouldEnforceSimpleSelfieValidation"`
	// ShouldEnforcePasscodeValidation requires passcode validation.
	ShouldEnforcePasscodeValidation bool `json:"shouldEnforcePasscodeValidation"`
	// Passcode the signer must provide (up to 100 characters).
	Passcode *string `json:"passcode,omitempty"`
	// PasscodeHint shown to the signer during signing (up to 200 characters).
	PasscodeHint *string `json:"passcodeHint,omitempty"`
	// SignMarks the signer must complete on the documents.
	SignMarks []SignMarkParams `json:"signMarks"`
}

// SignMarkParams define one mark a signer must place on a document.
type SignMarkParams struct {
	// ID of an existing sign mark. Only meaningful on Update; leave nil to
	// add a new mark.
	ID *string `json:"id,omitempty"`
	// Type is the kind of mark.
	Type SignMarkType `json:"type,omitempty"`
	// DocumentID is the document the mark belongs to. Required.
	DocumentID string `json:"documentId"`
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
	// IsRequired makes the mark mandatory.
	IsRequired bool `json:"isRequired"`
	// SchemaJSON is the mark's schema as JSON (for Text, Dropdown, etc., up
	// to 10000 characters).
	SchemaJSON *string `json:"schemaJson,omitempty"`
}
