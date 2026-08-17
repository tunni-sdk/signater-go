# Create contact

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/contacts:
    post:
      summary: Create contact
      deprecated: false
      description: >-
        Creates a new contact in the system by accepting essential details such
        as name, email, and other optional metadata. This endpoint is designed
        to store reusable information that can streamline the process of
        assigning contacts as signers. Ensure all required fields are provided
        in the request payload.
      tags:
        - Contact
        - Contact
      parameters: []
      requestBody:
        content:
          application/json:
            schema:
              allOf:
                - $ref: '#/components/schemas/CreateContactApiRequest'
            example:
              name: John Doe
              title: Project Manager
              email: john.doe@example.com
              phoneIdd: 1
              phoneNumber: '1234567890'
              documentType: GenericIdentification
              documentValue: A1234567
              description: A contact person for project management.
      responses:
        '201':
          description: Contact created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CreateContactApiResponse'
              example:
                id: e303178d-de2c-42de-8f89-4f7dcac95beb
          headers: {}
          x-apidog-name: Created
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
      security:
        - apikey-header-x-api-token: []
      x-apidog-folder: Contact
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252668-run
components:
  schemas:
    CreateContactApiRequest:
      required:
        - name
      type: object
      properties:
        name:
          maxLength: 200
          minLength: 2
          type: string
          description: Name of the contact
        title:
          maxLength: 200
          type: string
          description: Title of the contact
          nullable: true
        email:
          pattern: ^[\w\.\-\+]+@([\w\-]+\.)+[a-zA-Z]{2,}$
          type: string
          description: Email address of the contact
          nullable: true
        phoneIdd:
          type: integer
          description: >-
            Phone IDD of the contact, e.g. '1' for the United States or '55' for
            Brazil
          format: int32
          nullable: true
        phoneNumber:
          maxLength: 200
          type: string
          description: Phone number of the contact excluding the IDD
          nullable: true
        documentType:
          enum:
            - GenericIdentification
            - BrazilianCpf
            - BrazilianIdentity
            - Passport
          type: string
          description: Document type of the contact
          nullable: true
        documentValue:
          maxLength: 100
          type: string
          description: Document value of the contact
          nullable: true
        description:
          maxLength: 2000
          type: string
          description: Description of the contact
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - name
        - title
        - email
        - phoneIdd
        - phoneNumber
        - documentType
        - documentValue
        - description
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    CreateContactApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the created contact
          format: uuid
      additionalProperties: false
      x-apidog-orders:
        - id
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
