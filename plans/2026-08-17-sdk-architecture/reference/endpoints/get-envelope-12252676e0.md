# Get envelope

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/envelopes/{envelopeId}:
    get:
      summary: Get envelope
      deprecated: false
      description: >-
        Retrieves detailed information about a specific envelope using its
        unique identifier. The response includes the envelope's status,
        associated documents, signers, and other relevant metadata. Ensure the
        provided envelope ID exists and is valid before making the request.
      tags:
        - Envelope
        - Envelope
      parameters:
        - name: envelopeId
          in: path
          description: Envelope ID
          required: true
          example: ''
          schema:
            type: string
            description: Envelope ID
            format: uuid
      responses:
        '200':
          description: Envelope retrieved successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GetEnvelopeApiResponse'
              example:
                id: e0b0ef0a-b2d8-4921-b996-a21777779881
                createdAtUtc: '2026-08-17T00:52:53.8625128Z'
                createdById: 64d86dc2-c160-480b-9c1d-8388d159a20f
                createdByName: John Doe
                createdByAvatar: https://example.com/avatar.jpg
                updatedAtUtc: '2026-08-17T00:52:53.8625143Z'
                lastUpdateById: 83d5177d-ba62-435d-8c86-119176ae55b2
                updatedByName: John Doe
                updatedByAvatar: https://example.com/avatar.jpg
                ownerId: 5e8c7615-06fa-40a7-93a3-2e8efe111a36
                ownerName: John Doe
                ownerAvatar: https://example.com/avatar.jpg
                language: EnUs
                markupOrientation: Bottom
                status: Published
                hasScheduledPublish: false
                toBePublishedAtUtc: null
                vaultId: b1905a70-f046-4435-b315-f2e3bc877c41
                vaultName: Example Vault
                name: Example Envelope
                privateDescription: Private description
                publicDescription: Public description
                message: Message
                reviewReminder: false
                expiresAtUtc: null
                expirationReminder: false
                isRemoved: false
                removedById: null
                removedBy: null
                removedAtUtc: null
                signInOrder: false
                certifiedAtUtc: null
                isSandbox: false
                hasActionsBeingProcessed: false
                isAiChatForSignersEnabled: false
                isLtvEnabled: true
                redirectUrl: null
                jurisdictionCountryCode: null
                documents:
                  - id: 2894e56d-523a-45d9-8bf9-d6c1a374ff93
                    index: 0
                    origin: Upload
                    name: Example Document
                    privateDescription: Private description
                    publicDescription: Public description
                    originalFileSize: 0
                    pageSizes:
                      - page: 1
                        width: 100
                        height: 100
                signers:
                  - id: ee7e9105-222f-4623-9f67-27966ad1aa5b
                    index: null
                    status: ReadyToReview
                    hadApproved: false
                    approvedAtUtc: null
                    hadRejected: false
                    rejectedAtUtc: null
                    name: John Doe
                    email: john.doe@example.com
                    title: Title
                    passcode: '123456'
                    passcodeHint: 1 to 6
                    phoneIdd: 1
                    phoneNumber: '123456789'
                    documentType: GenericIdentification
                    documentValue: '123456789'
                    shouldEnforceSmsValidation: true
                    shouldEnforceWhatsAppValidation: true
                    shouldEnforcePixValidation: true
                    shouldEnforceCustomDigitalCertificateValidation: true
                    shouldEnforcePasscodeValidation: true
                    shouldEnforceEmailValidation: true
                    shouldEnforceSimpleSelfieValidation: true
                    emailCommunicationMode: Full
                    smsCommunicationMode: Full
                    whatsAppCommunicationMode: Full
                    signMarks:
                      - id: 88fee637-10d4-4ce2-9944-d02c19155b09
                        documentId: 2894e56d-523a-45d9-8bf9-d6c1a374ff93
                        type: Signature
                        page: 1
                        x: 10
                        'y': 10
                        rotation: 0
                        width: 0
                        height: 0
                        isRequired: false
                        schemaJson: null
                        valueJson: null
                    actions:
                      - id: 8c5be8f5-2eec-4fe5-b1ef-80694b946118
                        createdAtUtc: '2026-08-17T00:49:53.8625258Z'
                        type: View
                        ip: 127.0.0.1
                        userAgent: Mozilla/5.0
                        location: New York, USA
                        signMarkId: null
                        signMarkType: null
                        signMarkSchemaJson: null
                        valueJson: null
                      - id: 2daba658-f986-4e5f-bd98-a9c272006c6f
                        createdAtUtc: '2026-08-17T00:49:53.8625279Z'
                        type: EmailMfaRequest
                        ip: 127.0.0.1
                        userAgent: Mozilla/5.0
                        location: New York, USA
                        signMarkId: null
                        signMarkType: null
                        signMarkSchemaJson: null
                        valueJson: null
          headers: {}
          x-apidog-name: OK
        '400':
          description: Invalid request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BadRequestApiResponse'
              example:
                errors:
                  - message: >-
                      The length of 'Name' must be at least 2 characters. You
                      entered 1 characters.
                    metadata: null
          headers: {}
          x-apidog-name: Bad Request
        '401':
          description: Unauthorized request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UnauthorizedApiResponse'
              example:
                message: User should be authenticated.
          headers: {}
          x-apidog-name: Unauthorized
        '403':
          description: Forbidden request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ForbiddenApiResponse'
              example:
                message: User should be an administrator.
          headers: {}
          x-apidog-name: Forbidden
        '404':
          description: Resource not found
          headers: {}
          x-apidog-name: Record Not Found
      security:
        - apikey-header-x-api-token: []
      x-apidog-folder: Envelope
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252676-run
components:
  schemas:
    GetEnvelopeApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the envelope
          format: uuid
        createdAtUtc:
          type: string
          description: Date and time when the envelope was created in UTC
          format: date-time
        createdById:
          type: string
          description: Id of the user who created the envelope
          format: uuid
        createdByName:
          type: string
          description: Name of the user who created the envelope
          nullable: true
        createdByAvatar:
          type: string
          description: Url of the avatar of the user who created the envelope
          nullable: true
        updatedAtUtc:
          type: string
          description: Date and time when the envelope was last updated in UTC
          format: date-time
        lastUpdateById:
          type: string
          description: Id of the user who last updated the envelope
          format: uuid
        updatedByName:
          type: string
          description: Name of the user who last updated the envelope
          nullable: true
        updatedByAvatar:
          type: string
          description: Url of the avatar of the user who last updated the envelope
          nullable: true
        ownerId:
          type: string
          description: Id of the owner of the envelope
          format: uuid
        ownerName:
          type: string
          description: Name of the owner of the envelope
          nullable: true
        ownerAvatar:
          type: string
          description: Url of the avatar of the owner of the envelope
          nullable: true
        language:
          enum:
            - EnUs
            - PtBr
            - EsEs
            - FrFr
            - DeDe
            - ItIt
          type: string
          description: Language of the envelope
        markupOrientation:
          enum:
            - None
            - Bottom
            - Top
            - Left
            - Right
          type: string
          description: Orientation of the markup
        status:
          enum:
            - Draft
            - PublishScheduled
            - Published
            - Hold
            - Cancelled
            - CancelledBySignerMfaError
            - Rejected
            - Signed
            - Expired
          type: string
          description: Status of the envelope
        hasScheduledPublish:
          type: boolean
          description: Indicates if the envelope has scheduled publish
        toBePublishedAtUtc:
          type: string
          description: Date and time when the envelope will be published in UTC
          format: date-time
          nullable: true
        vaultId:
          type: string
          description: Id of the vault
          format: uuid
        vaultName:
          type: string
          description: Name of the vault
          nullable: true
        name:
          type: string
          description: Name of the envelope
          nullable: true
        privateDescription:
          type: string
          description: Private description of the envelope
          nullable: true
        publicDescription:
          type: string
          description: Public description of the envelope
          nullable: true
        message:
          type: string
          description: Message of the envelope, which will be send to the signers via email
          nullable: true
        reviewReminder:
          type: boolean
          description: Indicates if review reminders should be sent to the signers
        expiresAtUtc:
          type: string
          description: Date and time when the envelope will expire in UTC
          format: date-time
          nullable: true
        expirationReminder:
          type: boolean
          description: Indicates if expiration reminders should be sent to the signers
        isRemoved:
          type: boolean
          description: Indicates if the envelope is removed
        removedById:
          type: string
          description: Id of the user who removed the envelope
          format: uuid
          nullable: true
        removedBy:
          type: string
          description: Name of the user who removed the envelope
          nullable: true
        removedAtUtc:
          type: string
          description: Date and time when the envelope was removed in UTC
          format: date-time
          nullable: true
        signInOrder:
          type: boolean
          description: Indicates if the signers should sign in order
        certifiedAtUtc:
          type: string
          description: >-
            Date and time when the envelope certification was last updated in
            UTC
          format: date-time
          nullable: true
        isSandbox:
          type: boolean
          description: Indicates if the envelope was created in sandbox mode
        hasActionsBeingProcessed:
          type: boolean
          description: Indicates if the envelope has some action being processed
        isAiChatForSignersEnabled:
          type: boolean
          description: Indicates if the AI chat for signers is enabled
        isLtvEnabled:
          type: boolean
          description: Indicates if long-term validation (LTV) is enabled for the envelope
        redirectUrl:
          type: string
          description: >-
            URL to redirect the signer after completing their action (approve,
            reject, cancel). The system will inject 'envelopeId' and 'signerId'
            as query params. If the URL already contains these params (in any
            casing, e.g. ENVELOPEID), they will be overwritten. Additional
            custom query params (e.g. envelope_id) are accepted and preserved.
          nullable: true
        jurisdictionCountryCode:
          type: string
          description: >-
            Jurisdiction country of the envelope as ISO 3166-1 alpha-2 code
            (e.g. 'BR', 'US').
          nullable: true
        documents:
          type: array
          items:
            $ref: '#/components/schemas/GetEnvelopeApiResponseDocument'
          description: List of documents of the envelope
          nullable: true
        signers:
          type: array
          items:
            $ref: '#/components/schemas/GetEnvelopeApiResponseSigner'
          description: List of signers of the envelope
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - createdAtUtc
        - createdById
        - createdByName
        - createdByAvatar
        - updatedAtUtc
        - lastUpdateById
        - updatedByName
        - updatedByAvatar
        - ownerId
        - ownerName
        - ownerAvatar
        - language
        - markupOrientation
        - status
        - hasScheduledPublish
        - toBePublishedAtUtc
        - vaultId
        - vaultName
        - name
        - privateDescription
        - publicDescription
        - message
        - reviewReminder
        - expiresAtUtc
        - expirationReminder
        - isRemoved
        - removedById
        - removedBy
        - removedAtUtc
        - signInOrder
        - certifiedAtUtc
        - isSandbox
        - hasActionsBeingProcessed
        - isAiChatForSignersEnabled
        - isLtvEnabled
        - redirectUrl
        - jurisdictionCountryCode
        - documents
        - signers
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    GetEnvelopeApiResponseSigner:
      type: object
      properties:
        id:
          type: string
          description: Id of the signer
          format: uuid
        index:
          type: integer
          description: >-
            Index of the signer, if the envelope is signed in order, starting
            from 0
          format: int32
          nullable: true
        status:
          enum:
            - None
            - ReadyToReview
            - Queued
            - Approved
            - Rejected
            - CancelledByMfaError
          type: string
          description: Status of the signer
        hadApproved:
          type: boolean
          description: Indicates if the signer had approved the envelope
        approvedAtUtc:
          type: string
          description: Date and time when the signer approved the envelope in UTC
          format: date-time
          nullable: true
        hadRejected:
          type: boolean
          description: Indicates if the signer had rejected the envelope
        rejectedAtUtc:
          type: string
          description: Date and time when the signer rejected the envelope in UTC
          format: date-time
          nullable: true
        name:
          type: string
          description: Name of the signer
          nullable: true
        email:
          type: string
          description: Email of the signer
          nullable: true
        title:
          type: string
          description: Title of the signer
          nullable: true
        passcode:
          type: string
          description: Passcode of the signer
          nullable: true
        passcodeHint:
          type: string
          description: Passcode hint of the signer
          nullable: true
        phoneIdd:
          type: integer
          description: Phone IDD of the signer
          format: int32
          nullable: true
        phoneNumber:
          type: string
          description: Phone number of the signer
          nullable: true
        documentType:
          enum:
            - GenericIdentification
            - BrazilianCpf
            - BrazilianIdentity
            - Passport
          type: string
          description: Document type of the signer
        documentValue:
          type: string
          description: Document value of the signer
          nullable: true
        shouldEnforceSmsValidation:
          type: boolean
          description: >-
            Indicates if the SMS OTP validation should be enforced for the
            signer
        shouldEnforceWhatsAppValidation:
          type: boolean
          description: >-
            Indicates if the WhatsApp OTP validation should be enforced for the
            signer
        shouldEnforcePixValidation:
          type: boolean
          description: >-
            Indicates if the PIX OTP validation should be enforced for the
            signer
        shouldEnforceCustomDigitalCertificateValidation:
          type: boolean
          description: >-
            Indicates if the custom digital certificate validation should be
            enforced for the signer
        shouldEnforcePasscodeValidation:
          type: boolean
          description: >-
            Indicates if the passcode validation should be enforced for the
            signer
        shouldEnforceEmailValidation:
          type: boolean
          description: >-
            Indicates if the email OTP validation should be enforced for the
            signer
        shouldEnforceSimpleSelfieValidation:
          type: boolean
          description: >-
            Indicates if the simple selfie validation should be enforced for the
            signer
        emailCommunicationMode:
          enum:
            - None
            - InvitationOnly
            - Full
          type: string
          description: Indicates how the signer prefers to receive email communications
        smsCommunicationMode:
          enum:
            - None
            - InvitationOnly
            - Full
          type: string
          description: Indicates how the signer prefers to receive SMS communications
        whatsAppCommunicationMode:
          enum:
            - None
            - InvitationOnly
            - Full
          type: string
          description: Indicates how the signer prefers to receive WhatsApp communications
        signMarks:
          type: array
          items:
            $ref: '#/components/schemas/GetEnvelopeApiResponseSignMark'
          description: List of sign marks of the signer
          nullable: true
        actions:
          type: array
          items:
            $ref: '#/components/schemas/GetEnvelopeApiResponseSignerAction'
          description: List of actions of the signer
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - index
        - status
        - hadApproved
        - approvedAtUtc
        - hadRejected
        - rejectedAtUtc
        - name
        - email
        - title
        - passcode
        - passcodeHint
        - phoneIdd
        - phoneNumber
        - documentType
        - documentValue
        - shouldEnforceSmsValidation
        - shouldEnforceWhatsAppValidation
        - shouldEnforcePixValidation
        - shouldEnforceCustomDigitalCertificateValidation
        - shouldEnforcePasscodeValidation
        - shouldEnforceEmailValidation
        - shouldEnforceSimpleSelfieValidation
        - emailCommunicationMode
        - smsCommunicationMode
        - whatsAppCommunicationMode
        - signMarks
        - actions
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    GetEnvelopeApiResponseSignerAction:
      type: object
      properties:
        id:
          type: string
          description: Id of the signer action
          format: uuid
        createdAtUtc:
          type: string
          description: Date and time when the signer action was created in UTC
          format: date-time
        type:
          enum:
            - View
            - Approval
            - Rejection
            - CancellationByMfaError
            - ConfirmSignMark
            - EmailMfaRequest
            - SmsMfaRequest
            - Unsubscribe
            - WhatsAppMfaRequest
            - PixMfaRequest
            - UnconfirmSignMark
          type: string
          description: Type of the signer action
        ip:
          type: string
          description: IP address of the signer action
          nullable: true
        userAgent:
          type: string
          description: User agent of the signer action
          nullable: true
        location:
          type: string
          description: Location of the signer action
          nullable: true
        signMarkId:
          type: string
          description: >-
            Id of the SignMark that this action targeted
            (Confirm/UnconfirmSignMark only)
          format: uuid
          nullable: true
        signMarkType:
          enum:
            - Invisible
            - Signature
            - LegacyRectangleRubric
            - Text
            - LongText
            - Date
            - Number
            - Currency
            - Checkbox
            - RadioButton
            - Dropdown
            - Attachment
            - Rubric
          type: string
          description: Type of the SignMark that this action targeted
          nullable: true
        signMarkSchemaJson:
          type: string
          description: Schema JSON of the SignMark that this action targeted
          nullable: true
        valueJson:
          type: string
          description: Snapshot of the value at action time, for audit display
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - createdAtUtc
        - type
        - ip
        - userAgent
        - location
        - signMarkId
        - signMarkType
        - signMarkSchemaJson
        - valueJson
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    GetEnvelopeApiResponseSignMark:
      type: object
      properties:
        id:
          type: string
          description: Id of the sign mark
          format: uuid
        documentId:
          type: string
          description: Id of the document
          format: uuid
        type:
          enum:
            - Invisible
            - Signature
            - LegacyRectangleRubric
            - Text
            - LongText
            - Date
            - Number
            - Currency
            - Checkbox
            - RadioButton
            - Dropdown
            - Attachment
            - Rubric
          type: string
          description: Type of the sign mark
        page:
          type: integer
          description: Page number
          format: int32
        x:
          type: integer
          description: X coordinate in pixels of the sign mark top left corner
          format: int32
        'y':
          type: integer
          description: Y coordinate in pixels of the sign mark top left corner
          format: int32
        rotation:
          type: integer
          description: Rotation of the sign mark
          format: int32
        width:
          type: integer
          description: Width of the sign mark in pixels
          format: int32
        height:
          type: integer
          description: Height of the sign mark in pixels
          format: int32
        isRequired:
          type: boolean
          description: Indicates if the sign mark is required
        schemaJson:
          type: string
          description: Schema of the sign mark as JSON
          nullable: true
        valueJson:
          type: string
          description: Value of the sign mark as JSON
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - documentId
        - type
        - page
        - x
        - 'y'
        - rotation
        - width
        - height
        - isRequired
        - schemaJson
        - valueJson
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    GetEnvelopeApiResponseDocument:
      type: object
      properties:
        id:
          type: string
          description: Id of the document
          format: uuid
        index:
          type: integer
          description: Index of the document, starting from 0
          format: int32
        origin:
          enum:
            - Upload
            - Template
          type: string
          description: Origin of the document
        name:
          type: string
          description: Name of the document
          nullable: true
        privateDescription:
          type: string
          description: Private description of the document
          nullable: true
        publicDescription:
          type: string
          description: Public description of the document
          nullable: true
        originalFileSize:
          type: integer
          description: Size of the original file in bytes
          format: int64
        pageSizes:
          type: array
          items:
            $ref: '#/components/schemas/GetEnvelopeApiResponseDocumentPageSize'
          description: List of page sizes of the document
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - index
        - origin
        - name
        - privateDescription
        - publicDescription
        - originalFileSize
        - pageSizes
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    GetEnvelopeApiResponseDocumentPageSize:
      type: object
      properties:
        page:
          type: integer
          description: Page number starting from 1
          format: int32
        width:
          type: number
          description: Width of the page in pixels
          format: double
        height:
          type: number
          description: Height of the page in pixels
          format: double
      additionalProperties: false
      x-apidog-orders:
        - page
        - width
        - height
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    BadRequestApiResponse:
      type: object
      properties:
        errors:
          type: array
          items:
            $ref: '#/components/schemas/BadRequestApiResponseError'
          description: List of errors that caused the request to fail
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - errors
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    BadRequestApiResponseError:
      type: object
      properties:
        message:
          type: string
          description: Error message
          nullable: true
        metadata:
          type: object
          additionalProperties:
            type: string
            nullable: true
          description: Optional error metadata dictionary
          x-apidog-orders: []
          properties: {}
          x-apidog-ignore-properties: []
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - message
        - metadata
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    UnauthorizedApiResponse:
      type: object
      properties:
        message:
          type: string
          description: Message containing the reason why the request was unauthorized
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - message
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ForbiddenApiResponse:
      type: object
      properties:
        message:
          type: string
          description: Message containing the reason why the request was forbidden
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - message
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
  securitySchemes:
    Api Key:
      type: apikey
      description: >-
        Please enter a header called x-api-token with the api token created in
        the application
      name: Authorization
      in: header
    apikey-header-x-api-token:
      type: apiKey
      in: header
      name: x-api-token
servers:
  - url: https://api.signater.com
    description: Amb. de prod.
security:
  - apikey-header-x-api-token: []

```
