# Get template

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/templates/{templateId}:
    get:
      summary: Get template
      deprecated: false
      description: >-
        Retrieves a template by its unique identifier. This endpoint allows you
        to access the template's content and fields configuration.
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
      responses:
        '200':
          description: Template retrieved successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GetTemplateApiResponse'
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
      x-apidog-folder: Template
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252694-run
components:
  schemas:
    GetTemplateApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the vault
          format: uuid
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
          description: Name of the template
          nullable: true
        description:
          type: string
          description: Description of the template
          nullable: true
        fields:
          type: array
          items:
            $ref: '#/components/schemas/GetTemplateFieldApiResponse'
          description: Fields of the template
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - vaultId
        - vaultName
        - name
        - description
        - fields
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    GetTemplateFieldApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the field
          format: uuid
        alias:
          type: string
          description: Alias of the field
          nullable: true
        isRequired:
          type: boolean
          description: Indicates if the field is required
        type:
          enum:
            - Text
            - LongText
            - Number
            - Currency
            - Date
            - Checkbox
            - RadioButton
            - Dropdown
          type: string
          description: Type of the field
        schema:
          allOf:
            - $ref: '#/components/schemas/TemplateFieldSchema'
          description: >-
            Schema of the field, defining its label, placeholder, options and
            constraints
          type: 'null'
        index:
          type: integer
          description: Index of the field, starting from 0
          format: int32
      additionalProperties: false
      x-apidog-orders:
        - id
        - alias
        - isRequired
        - type
        - schema
        - index
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    TemplateFieldSchema:
      type: object
      properties:
        placeholder:
          type: string
          nullable: true
        label:
          type: string
          nullable: true
        description:
          type: string
          nullable: true
        options:
          type: array
          items:
            type: string
          nullable: true
        currency:
          type: string
          nullable: true
        defaultValue:
          type: string
          nullable: true
        constraints:
          allOf:
            - $ref: '#/components/schemas/TemplateFieldConstraints'
          type: 'null'
        dateFormat:
          type: string
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - placeholder
        - label
        - description
        - options
        - currency
        - defaultValue
        - constraints
        - dateFormat
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    TemplateFieldConstraints:
      type: object
      properties:
        minLength:
          type: integer
          format: int32
          nullable: true
        maxLength:
          type: integer
          format: int32
          nullable: true
        pattern:
          type: string
          nullable: true
        min:
          type: number
          format: double
          nullable: true
        max:
          type: number
          format: double
          nullable: true
        decimalPlaces:
          type: integer
          format: int32
          nullable: true
        allowNegative:
          type: boolean
          nullable: true
        allowMultiple:
          type: boolean
          nullable: true
        minDate:
          type: string
          nullable: true
        maxDate:
          type: string
          nullable: true
        profile:
          enum:
            - None
            - Email
            - Phone
            - Cpf
            - Cnpj
            - Cep
          type: string
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - minLength
        - maxLength
        - pattern
        - min
        - max
        - decimalPlaces
        - allowNegative
        - allowMultiple
        - minDate
        - maxDate
        - profile
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
