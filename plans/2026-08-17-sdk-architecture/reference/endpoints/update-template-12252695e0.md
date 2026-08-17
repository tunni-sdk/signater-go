# Update template

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/templates/{templateId}:
    put:
      summary: Update template
      deprecated: false
      description: >-
        Updates an existing template with new content, structure, or fields.
        This allows you to modify predefined elements of the template to reflect
        changes in document requirements. The updated template can then be used
        to create new documents and envelopes, ensuring that future documents
        align with the latest configuration and structure.
      tags:
        - Template
        - Template
      parameters:
        - name: templateId
          in: path
          description: Template ID
          required: true
          example: ''
          schema:
            type: string
            description: Template ID
            format: uuid
      requestBody:
        content:
          application/json:
            schema:
              allOf:
                - $ref: '#/components/schemas/UpdateTemplateApiRequest'
            example:
              vaultId: c982b250-70a0-4293-b5cc-25a6a97e65fc
              name: Contract Template
              description: A template description for contracts between parties.
              payload: '{ "key": "value" }'
              fields:
                - id: 84819214-2d3e-4416-9b78-8400dc2b59c7
                  name: First Name
                  alias: first_name
                  description: The first name of the individual.
                  isRequired: true
                  type: Text
                  customRegex: null
                - id: null
                  name: Age
                  alias: age
                  description: The age of the individual.
                  isRequired: false
                  type: Number
                  customRegex: null
                - id: null
                  name: Score
                  alias: score
                  description: Score number.
                  isRequired: false
                  type: Custom
                  customRegex: ^[0-9]+$
      responses:
        '204':
          description: Template updated successfully
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
      x-apidog-folder: Template
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252695-run
components:
  schemas:
    UpdateTemplateApiRequest:
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
            $ref: '#/components/schemas/UpdateTemplateFieldApiRequest'
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
    UpdateTemplateFieldApiRequest:
      required:
        - name
        - alias
      type: object
      properties:
        id:
          type: string
          description: Id of the field
          format: uuid
          nullable: true
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
        - id
        - name
        - alias
        - description
        - isRequired
        - type
        - customRegex
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
