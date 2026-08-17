# Upload document

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/documents:
    post:
      summary: Upload document
      deprecated: false
      description: >-
        Uploads a document to the system and returns a unique document ID. This
        ID must be used in subsequent operations for creating or editing an
        envelope. If the document is not used within 24 hours, it will be
        automatically removed from the system's storage.
      tags:
        - Document
        - Document
      parameters: []
      requestBody:
        content:
          multipart/form-data:
            schema:
              type: object
              properties:
                File:
                  type: string
                  description: >-
                    File to upload. Accepted formats: PDF, DOCX, DOC, DOCM, RTF,
                    HTML, TXT, ODT, EPUB, MHT, MHTML, XLSX, PNG, JPG, JPEG, GIF,
                    BMP, and WEBP. Non-PDF files are automatically converted to
                    PDF.
                  format: binary
                  example: ''
              required:
                - File
      responses:
        '201':
          description: Document uploaded successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UploadDocumentApiResponse'
              example:
                id: 851f8457-733d-4af3-81fa-29d86274a0a4
                fileName: FileName.pdf
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
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252674-run
components:
  schemas:
    UploadDocumentApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the created document
          format: uuid
        fileName:
          type: string
          description: Original document file name
          nullable: true
        originalFileSize:
          type: integer
          description: Size of the original file in bytes
          format: int64
        pageSizes:
          type: array
          items:
            $ref: '#/components/schemas/UploadDocumentPageSizeApiResponse'
          description: Page sizes
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - fileName
        - originalFileSize
        - pageSizes
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    UploadDocumentPageSizeApiResponse:
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
