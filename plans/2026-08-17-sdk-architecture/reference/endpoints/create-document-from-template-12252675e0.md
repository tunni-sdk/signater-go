# Create document from template

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/documents/templates:
    post:
      summary: Create document from template
      deprecated: false
      description: >-
        Creates a new document by applying a predefined template, which can be
        customized with specific data. The endpoint will return a unique
        document ID, which must be used for further operations such as signing
        or editing. If the document is not used within 24 hours, it will be
        automatically removed from the system's storage.
      tags:
        - Document
        - Document
      parameters: []
      requestBody:
        content:
          application/json:
            schema:
              allOf:
                - $ref: '#/components/schemas/CreateDocumentFromTemplateApiRequest'
              description: Request to create document from template
            example:
              templateId: 3e533f62-031a-4ea6-8e13-5be21f6ded7b
              name: Example document
              privateDescription: Example private description
              publicDescription: Example public description
              language: EnUs
              markupOrientation: Bottom
              fields:
                - fieldId: 178712be-0599-4387-b68c-d8376cf9a5a3
                  value:
                    text: Example field value
                    date: null
                    checked: null
                    selected: null
                    number: null
                    amount: null
      responses:
        '201':
          description: Document created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CreateDocumentFromTemplateApiResponse'
              example:
                id: 11c5c486-4238-4dcd-bbd1-c5cdbe3321e8
                fileName: FileName.pdf
                name: Name
                privateDescription: Private Description
                publicDescription: Public Description
                originalFileSize: 0
                pageSizes:
                  - page: 1
                    width: 100
                    height: 200
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
      x-apidog-folder: Document
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252675-run
components:
  schemas:
    CreateDocumentFromTemplateApiRequest:
      required:
        - templateId
      type: object
      properties:
        templateId:
          type: string
          description: Template ID
          format: uuid
        name:
          maxLength: 200
          type: string
          description: >-
            Name of the document; if not provided, the original file name will
            be used
          nullable: true
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
          description: Language of the document, to be used in the markup if configured
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
            Markup orientation of the document; it stamps information about the
            document on all pages, including a link to access the document
            publicly
          nullable: true
        fields:
          type: array
          items:
            $ref: '#/components/schemas/CreateDocumentFromTemplateFieldApiRequest'
          description: Fields of the document
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - templateId
        - name
        - privateDescription
        - publicDescription
        - language
        - markupOrientation
        - fields
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    CreateDocumentFromTemplateFieldApiRequest:
      required:
        - fieldId
      type: object
      properties:
        fieldId:
          type: string
          description: Field ID
          format: uuid
        value:
          allOf:
            - $ref: '#/components/schemas/TemplateFieldValue'
          description: Value of the field, typed according to the field schema
          type: 'null'
      additionalProperties: false
      x-apidog-orders:
        - fieldId
        - value
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    TemplateFieldValue:
      type: object
      properties:
        text:
          type: string
          nullable: true
        date:
          type: string
          nullable: true
        checked:
          type: boolean
          nullable: true
        selected:
          type: array
          items:
            type: string
          nullable: true
        number:
          type: number
          format: double
          nullable: true
        amount:
          type: number
          format: double
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - text
        - date
        - checked
        - selected
        - number
        - amount
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    CreateDocumentFromTemplateApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the created document
          format: uuid
        fileName:
          type: string
          description: Original file name
          nullable: true
        name:
          type: string
          description: Name
          nullable: true
        privateDescription:
          type: string
          description: Private description
          nullable: true
        publicDescription:
          type: string
          description: Public description
          nullable: true
        originalFileSize:
          type: integer
          description: Size of the original file in bytes
          format: int64
        pageSizes:
          type: array
          items:
            $ref: '#/components/schemas/CreateDocumentFromTemplatePageSizeApiResponse'
          description: Page sizes
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - fileName
        - name
        - privateDescription
        - publicDescription
        - originalFileSize
        - pageSizes
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    CreateDocumentFromTemplatePageSizeApiResponse:
      type: object
      properties:
        page:
          type: integer
          description: Page number, starting from 1
          format: int32
        width:
          type: number
          description: Width in pixels
          format: double
        height:
          type: number
          description: Height in pixels
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
