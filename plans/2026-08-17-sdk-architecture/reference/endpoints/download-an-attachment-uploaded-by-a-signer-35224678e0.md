# Download an attachment uploaded by a signer

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/envelopes/{envelopeId}/sign-marks/{signMarkId}/attachments/{attachmentExternalId}/download:
    get:
      summary: Download an attachment uploaded by a signer
      deprecated: false
      description: >-
        Returns a 302 redirect to a pre-signed URL for the attachment file
        uploaded by a signer to an Attachment SignMark.
      tags:
        - Envelope
        - Envelope
      parameters:
        - name: envelopeId
          in: path
          description: ''
          required: true
          example: ''
          schema:
            type: string
            format: uuid
        - name: signMarkId
          in: path
          description: ''
          required: true
          example: ''
          schema:
            type: string
            format: uuid
        - name: attachmentExternalId
          in: path
          description: ''
          required: true
          example: ''
          schema:
            type: string
            format: uuid
      responses:
        '302':
          description: Found
          headers: {}
          x-apidog-name: ''
        '400':
          description: Bad Request
          headers: {}
          x-apidog-name: Bad Request
        '401':
          description: Unauthorized
          headers: {}
          x-apidog-name: Unauthorized
        '403':
          description: Forbidden
          headers: {}
          x-apidog-name: Forbidden
        '404':
          description: Not Found
          headers: {}
          x-apidog-name: Record Not Found
      security:
        - apikey-header-x-api-token: []
      x-apidog-folder: Envelope
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-35224678-run
components:
  schemas: {}
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
