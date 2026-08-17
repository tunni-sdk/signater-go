# Update envelope

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/envelopes/{envelopeId}:
    put:
      summary: Update envelope
      deprecated: false
      description: >-
        Updates the details of an existing envelope, such as adding or removing
        documents, modifying signers, or changing other associated metadata. Be
        aware that modifications to documents will not be allowed if any signer
        has already approved the envelope, nor will you be able to change
        signers who have already approved the envelope. The updated envelope
        will retain its unique identifier, and changes can be tracked in
        subsequent operations. Ensure the provided envelope ID exists and the
        request contains the necessary changes.
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
      requestBody:
        content:
          application/json:
            schema:
              allOf:
                - $ref: '#/components/schemas/UpdateEnvelopeApiRequest'
            example:
              vaultId: 2f56960c-c2ff-40c4-98e6-6186b7c24bcf
              name: Example Envelope
              privateDescription: Private description of the envelope
              publicDescription: Public description of the envelope
              message: Message of the envelope
              reviewReminder: true
              expiresAtUtc: '2026-09-16T00:52:53.872149Z'
              expirationReminder: true
              toBePublishedAtUtc: '2026-08-18T00:52:53.8721956Z'
              language: EnUs
              markupOrientation: Bottom
              signInOrder: true
              isAiChatForSignersEnabled: true
              isLtvEnabled: true
              redirectUrl: https://example.com/callback
              jurisdictionCountryCode: BR
              signers:
                - id: d3df5c3b-267d-41fa-8a14-111953e230bd
                  name: John Doe
                  email: john.doe@example.com
                  emailCommunicationMode: Full
                  title: Lawyer
                  shouldEnforceEmailValidation: true
                  shouldEnforceSmsValidation: true
                  shouldEnforceWhatsAppValidation: true
                  shouldEnforcePixValidation: true
                  shouldEnforceCustomDigitalCertificateValidation: true
                  shouldEnforceSimpleSelfieValidation: true
                  phoneIdd: 1
                  phoneNumber: '5555555555'
                  smsCommunicationMode: Full
                  whatsAppCommunicationMode: Full
                  documentType: GenericIdentification
                  documentValue: '1234567890'
                  shouldEnforcePasscodeValidation: true
                  passcode: '123456'
                  passcodeHint: 1 to 6
                  signMarks:
                    - id: 59c6adec-57ea-4bcc-90eb-050bd1d9e518
                      type: Rubric
                      documentId: 7cf41112-335f-474d-888f-e824450bb788
                      page: 1
                      x: 150
                      'y': 50
                      rotation: 0
                      width: 0
                      height: 0
                      isRequired: true
                      schemaJson: null
                    - id: null
                      type: Signature
                      documentId: 7cf41112-335f-474d-888f-e824450bb788
                      page: 2
                      x: 200
                      'y': 50
                      rotation: 0
                      width: 0
                      height: 0
                      isRequired: true
                      schemaJson: null
              documents:
                - id: 7cf41112-335f-474d-888f-e824450bb788
                  name: Example Document
                  privateDescription: Private description of the document
                  publicDescription: Public description of the document
                  language: EnUs
                  markupOrientation: Bottom
      responses:
        '204':
          description: Envelope updated successfully
          headers: {}
          x-apidog-name: No Content
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
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252677-run
components:
  schemas:
    UpdateEnvelopeApiRequest:
      required:
        - vaultId
        - name
        - jurisdictionCountryCode
      type: object
      properties:
        vaultId:
          type: string
          description: Id of the vault where the envelope will be created
          format: uuid
        name:
          maxLength: 200
          minLength: 2
          type: string
          description: Name of the envelope
        privateDescription:
          maxLength: 4000
          type: string
          description: Private description of the envelope
          nullable: true
        publicDescription:
          maxLength: 4000
          type: string
          description: Public description of the envelope
          nullable: true
        message:
          maxLength: 4000
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
        toBePublishedAtUtc:
          type: string
          description: Date and time when the envelope will be published in UTC
          format: date-time
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
          description: Markup orientation of the envelope
        signInOrder:
          type: boolean
          description: Indicates if the signers should sign in order
        isAiChatForSignersEnabled:
          type: boolean
          description: Indicates if the AI chat for signers is enabled
        isLtvEnabled:
          type: boolean
          description: Indicates if long-term validation (LTV) is enabled for the envelope
        redirectUrl:
          maxLength: 2000
          type: string
          description: >-
            URL to redirect the signer after completing their action (approve,
            reject, cancel). The system will inject 'envelopeId' and 'signerId'
            as query params. If the URL already contains these params (in any
            casing, e.g. ENVELOPEID), they will be overwritten. Additional
            custom query params (e.g. envelope_id) are accepted and preserved.
          nullable: true
        jurisdictionCountryCode:
          maxLength: 2
          minLength: 2
          type: string
          description: >-
            ISO 3166-1 alpha-2 jurisdiction country code of the envelope (e.g.
            'BR', 'US'). Changes are only allowed while the envelope is in Draft
            status.
        signers:
          type: array
          items:
            $ref: '#/components/schemas/UpdateEnvelopeSignerApiRequest'
          description: List of signers
          nullable: true
        documents:
          type: array
          items:
            $ref: '#/components/schemas/UpdateEnvelopeDocumentApiRequest'
          description: List of documents
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - vaultId
        - name
        - privateDescription
        - publicDescription
        - message
        - reviewReminder
        - expiresAtUtc
        - expirationReminder
        - toBePublishedAtUtc
        - language
        - markupOrientation
        - signInOrder
        - isAiChatForSignersEnabled
        - isLtvEnabled
        - redirectUrl
        - jurisdictionCountryCode
        - signers
        - documents
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    UpdateEnvelopeDocumentApiRequest:
      required:
        - id
        - name
      type: object
      properties:
        id:
          type: string
          description: Id of the document
          format: uuid
        name:
          maxLength: 200
          minLength: 2
          type: string
          description: Name of the document
        privateDescription:
          maxLength: 4000
          type: string
          description: Private description of the document
          nullable: true
        publicDescription:
          maxLength: 4000
          type: string
          description: Public description of the document
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
          description: >-
            Language of the document, if provided, it will override the language
            of the envelope
          nullable: true
        markupOrientation:
          enum:
            - None
            - Bottom
            - Top
            - Left
            - Right
          type: string
          description: >-
            Markup orientation of the document, if provided, it will override
            the markup orientation of the envelope
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - name
        - privateDescription
        - publicDescription
        - language
        - markupOrientation
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    UpdateEnvelopeSignerApiRequest:
      required:
        - name
      type: object
      properties:
        id:
          type: string
          description: Id of the signer
          format: uuid
          nullable: true
        name:
          maxLength: 200
          minLength: 2
          type: string
          description: Name of the signer
        email:
          pattern: ^[\w\.\-\+]+@([\w\-]+\.)+[a-zA-Z]{2,}$
          type: string
          description: Email of the signer
          nullable: true
        emailCommunicationMode:
          enum:
            - None
            - InvitationOnly
            - Full
          type: string
          description: Indicates how the signer prefers to receive email communications
        title:
          maxLength: 200
          type: string
          description: Title of the signer, e.g. Lawyer, Company Director, etc.
          nullable: true
        shouldEnforceEmailValidation:
          type: boolean
          description: >-
            Indicates if the email OTP validation should be enforced for the
            signer
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
        shouldEnforceSimpleSelfieValidation:
          type: boolean
          description: >-
            Indicates if the simple selfie validation should be enforced for the
            signer
        phoneIdd:
          type: integer
          description: >-
            Phone IDD of the signer, e.g. '1' for the United States or '55' for
            Brazil
          format: int32
          nullable: true
        phoneNumber:
          maxLength: 200
          type: string
          description: Phone number of the signer excluding the IDD
          nullable: true
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
        documentType:
          enum:
            - GenericIdentification
            - BrazilianCpf
            - BrazilianIdentity
            - Passport
          type: string
          description: Type of the document to be used for the signer
          nullable: true
        documentValue:
          maxLength: 100
          type: string
          description: Value of the document
          nullable: true
        shouldEnforcePasscodeValidation:
          type: boolean
          description: >-
            Indicates if the passcode validation should be enforced for the
            signer
        passcode:
          maxLength: 100
          type: string
          description: Passcode of the signer
          nullable: true
        passcodeHint:
          maxLength: 200
          type: string
          description: >-
            Passcode hint of the signer, which will be shown to the signer
            during the signing process
          nullable: true
        signMarks:
          type: array
          items:
            $ref: '#/components/schemas/UpdateEnvelopeSignerSignMarkApiRequest'
          description: List of sign marks for the signer
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - name
        - email
        - emailCommunicationMode
        - title
        - shouldEnforceEmailValidation
        - shouldEnforceSmsValidation
        - shouldEnforceWhatsAppValidation
        - shouldEnforcePixValidation
        - shouldEnforceCustomDigitalCertificateValidation
        - shouldEnforceSimpleSelfieValidation
        - phoneIdd
        - phoneNumber
        - smsCommunicationMode
        - whatsAppCommunicationMode
        - documentType
        - documentValue
        - shouldEnforcePasscodeValidation
        - passcode
        - passcodeHint
        - signMarks
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    UpdateEnvelopeSignerSignMarkApiRequest:
      required:
        - documentId
      type: object
      properties:
        id:
          type: string
          description: Id of the sign mark
          format: uuid
          nullable: true
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
        documentId:
          type: string
          description: Id of the document
          format: uuid
        page:
          type: integer
          description: Page number of the document
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
          maxLength: 10000
          type: string
          description: Schema of the sign mark as JSON (for Text, Dropdown, etc.)
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - type
        - documentId
        - page
        - x
        - 'y'
        - rotation
        - width
        - height
        - isRequired
        - schemaJson
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
