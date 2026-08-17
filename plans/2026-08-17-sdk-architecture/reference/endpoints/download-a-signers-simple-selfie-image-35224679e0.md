# Download a signer's simple-selfie image

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/envelopes/{envelopeId}/signers/{signerId}/simple-selfie:
    get:
      summary: Download a signer's simple-selfie image
      deprecated: false
      description: >-
        Returns a 302 redirect to a short-lived pre-signed URL for the
        simple-selfie image captured during signer validation. Only available
        for signers whose envelope enforced simple-selfie validation.
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
        - name: signerId
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
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-35224679-run
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
