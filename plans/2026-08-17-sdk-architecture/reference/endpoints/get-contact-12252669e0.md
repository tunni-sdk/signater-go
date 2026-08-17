# Get contact

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/contacts/{contactId}:
    get:
      summary: Get contact
      deprecated: false
      description: >-
        Retrieves detailed information about a specific contact by its unique
        identifier. This endpoint is useful for viewing or validating contact
        data before creating or updating related entities. Ensure the provided
        contact ID exists in the system.
      tags:
        - Contact
        - Contact
      parameters:
        - name: contactId
          in: path
          description: Contact ID
          required: true
          example: ''
          schema:
            type: string
            description: Contact ID
            format: uuid
      responses:
        '200':
          description: Contact retrieved successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GetContactApiResponse'
              example:
                id: bd9c82af-4549-4ad8-bf3c-148cb3858ac1
                name: John Doe
                title: Company Director
                email: john.doe@example.com
                phoneIdd: 1
                phoneNumber: '5555555555'
                documentType: GenericIdentification
                documentValue: '123456789'
                description: This is a description
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
      x-apidog-folder: Contact
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252669-run
components:
  schemas:
    GetContactApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the contact
          format: uuid
        name:
          type: string
          description: Name of the contact
          nullable: true
        title:
          type: string
          description: Title of the contact
          nullable: true
        email:
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
          type: string
          description: Document value of the contact
          nullable: true
        description:
          type: string
          description: Description of the contact
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
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
