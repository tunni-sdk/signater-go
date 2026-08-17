# List contacts

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/contacts:
    get:
      summary: List contacts
      deprecated: false
      description: >-
        Retrieves a list of contacts stored in the system. Contacts serve as a
        tool to streamline the creation of signers, allowing you to reuse
        previously registered information without needing to input it from
        scratch.
      tags:
        - Contact
        - Contact
      parameters:
        - name: IsFavorite
          in: query
          description: Filter contacts by favorite status
          required: false
          schema:
            type: boolean
            description: Filter contacts by favorite status
        - name: Search
          in: query
          description: Search contacts by id, name, email, or description
          required: false
          schema:
            type: string
            description: Search contacts by id, name, email, or description
        - name: PageSize
          in: query
          description: Number of items per page
          required: true
          schema:
            maximum: 100
            minimum: 1
            type: integer
            description: Number of items per page
            format: int32
            default: 10
        - name: PageNumber
          in: query
          description: Current page number
          required: true
          schema:
            maximum: 4294967295
            minimum: 1
            type: integer
            description: Current page number
            format: int32
            default: 1
        - name: OrderByDirection
          in: query
          description: Order by direction
          required: true
          schema:
            enum:
              - ASC
              - DESC
            type: string
            description: Order by direction
            default: DESC
      responses:
        '200':
          description: Contacts listed successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ListContactsApiResponse'
              example:
                items:
                  - id: 99c8a78c-71eb-4e20-a758-99082ee63ef9
                    createdAtUtc: '2026-08-17T00:52:53.7112783Z'
                    createdByName: John Doe
                    createdByAvatar: https://example.com/avatar.jpg
                    updatedAtUtc: '2026-08-17T00:52:53.7112785Z'
                    updatedByName: John Doe
                    updatedByAvatar: https://example.com/avatar.jpg
                    name: John Doe
                    title: Company Director
                    email: john.doe@example.com
                    phoneIdd: 1
                    phoneNumber: '5555555555'
                    documentType: GenericIdentification
                    documentValue: '123456789'
                    isFavorite: true
                pagination:
                  totalItems: 1
                  pageSize: 10
                  pageNumber: 1
                  pageItems: 1
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
      security:
        - apikey-header-x-api-token: []
      x-apidog-folder: Contact
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252667-run
components:
  schemas:
    ListContactsApiResponse:
      type: object
      properties:
        items:
          type: array
          items:
            $ref: '#/components/schemas/ListContactsApiResponseItem'
          description: List of contacts
          nullable: true
        pagination:
          allOf:
            - $ref: '#/components/schemas/Pagination'
          description: Pagination metadata
          type: 'null'
      additionalProperties: false
      x-apidog-orders:
        - items
        - pagination
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    Pagination:
      type: object
      properties:
        totalItems:
          type: integer
          description: Total number of items
          format: int32
        pageSize:
          type: integer
          description: Total number of pages
          format: int32
        pageNumber:
          type: integer
          description: Current page number
          format: int32
        pageItems:
          type: integer
          description: Number of items per page
          format: int32
      additionalProperties: false
      description: Pagination metadata
      x-apidog-orders:
        - totalItems
        - pageSize
        - pageNumber
        - pageItems
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ListContactsApiResponseItem:
      type: object
      properties:
        id:
          type: string
          description: Id of the contact
          format: uuid
        createdAtUtc:
          type: string
          description: Date and time when the contact was created (UTC)
          format: date-time
        createdByName:
          type: string
          description: Name of the user who created the contact
          nullable: true
        createdByAvatar:
          type: string
          description: Url Avatar of the user who created the contact
          nullable: true
        updatedAtUtc:
          type: string
          description: Date and time when the contact was last updated (UTC)
          format: date-time
        updatedByName:
          type: string
          description: Name of the user who last updated the contact
          nullable: true
        updatedByAvatar:
          type: string
          description: Url Avatar of the user who last updated the contact
          nullable: true
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
        isFavorite:
          type: boolean
          description: Description of the contact
      additionalProperties: false
      x-apidog-orders:
        - id
        - createdAtUtc
        - createdByName
        - createdByAvatar
        - updatedAtUtc
        - updatedByName
        - updatedByAvatar
        - name
        - title
        - email
        - phoneIdd
        - phoneNumber
        - documentType
        - documentValue
        - isFavorite
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
