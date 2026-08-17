# Create template

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/templates:
    post:
      summary: Create template
      deprecated: false
      description: >-
        Creates a new template that can be used to generate documents and
        envelopes. Templates allow you to define predefined structures and
        fields, which can be reused across multiple documents, saving time and
        ensuring consistency in document creation. Once created, the template
        can be accessed and used to quickly generate new documents with the same
        configuration.
      tags:
        - Template
        - Template
      parameters: []
      requestBody:
        content:
          application/json:
            schema:
              allOf:
                - $ref: '#/components/schemas/CreateTemplateApiRequest'
            example:
              vaultId: 1900bc26-7e01-40a1-aa83-d2abc112401c
              name: Contract Template
              description: A template description for contracts between parties.
              payload: '{ "key": "value" }'
              fields:
                - name: First Name
                  alias: first_name
                  description: The first name of the individual.
                  isRequired: true
                  type: Text
                  customRegex: null
                - name: Age
                  alias: age
                  description: The age of the individual.
                  isRequired: false
                  type: Number
                  customRegex: null
                - name: Score
                  alias: score
                  description: Score number.
                  isRequired: false
                  type: Custom
                  customRegex: ^[0-9]+$
      responses:
        '201':
          description: Template created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CreateTemplateApiResponse'
              example:
                id: 8a936624-9998-4c30-8486-62fd059017a7
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
      x-apidog-folder: Template
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252697-run
components:
  schemas:
    CreateTemplateApiRequest:
      required:
        - vaultId
        - name
      type: object
      properties:
        vaultId:
          type: string
          description: Id of the vault
          format: uuid
        name:
          maxLength: 200
          minLength: 2
          type: string
          description: Name of the template
        description:
          maxLength: 4000
          type: string
          description: Description of the template
          nullable: true
        payload:
          type: string
          description: Payload of the template
          nullable: true
        fields:
          type: array
          items:
            $ref: '#/components/schemas/CreateTemplateFieldApiRequest'
          description: Fields of the template
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - vaultId
        - name
        - description
        - payload
        - fields
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    CreateTemplateFieldApiRequest:
      required:
        - name
        - alias
      type: object
      properties:
        name:
          maxLength: 200
          minLength: 2
          type: string
          description: Name of the field
        alias:
          maxLength: 200
          minLength: 2
          pattern: ^[a-zA-Z0-9_-]*$
          type: string
          description: >-
            Alias of the field, must be unique and only contain alphanumeric
            characters, dashes and underscores
        description:
          maxLength: 4000
          type: string
          description: Description of the field
          nullable: true
        isRequired:
          type: boolean
          description: Indicates if the field is required
        type:
          enum:
            - BrCep
            - BrCpf
            - BrCnpj
            - DateOnly
            - TimeOnly
            - DateTime
            - Email
            - Number
            - Text
            - Custom
          type: string
          description: Type of the field
        customRegex:
          maxLength: 200
          type: string
          description: Custom regex when Type is Custom
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - name
        - alias
        - description
        - isRequired
        - type
        - customRegex
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    CreateTemplateApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the template
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
