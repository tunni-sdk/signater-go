# Get document original file URL

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/documents/{documentId}/original:
    get:
      summary: Get document original file URL
      deprecated: false
      description: Returns the URL to download the original file of a document.
      tags:
        - Document
        - Document
      parameters:
        - name: documentId
          in: path
          description: Document ID
          required: true
          example: ''
          schema:
            type: string
            description: Document ID
            format: uuid
      responses:
        '200':
          description: Document found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/GetDocumentOriginalFileUrlApiResponse'
              example:
                url: https://example.com/document.pdf
          headers: {}
          x-apidog-name: OK
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
          description: Document not found
          headers: {}
          x-apidog-name: Record Not Found
      security:
        - apikey-header-x-api-token: []
      x-apidog-folder: Document
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-13038397-run
components:
  schemas:
    GetDocumentOriginalFileUrlApiResponse:
      type: object
      properties:
        url:
          type: string
          description: The document original file URL
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - url
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
