# List available envelope owners

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/envelopes/{envelopeId}/owners:
    get:
      summary: List available envelope owners
      deprecated: false
      description: >-
        Retrieves a list of available envelope owners who are eligible to manage
        or own the specified envelope.
      tags:
        - Envelope
        - Envelope
      parameters:
        - name: envelopeId
          in: path
          description: Envelope ID
          required: true
          example: ''
          schema:
            type: string
            description: Envelope ID
            format: uuid
      responses:
        '200':
          description: Available envelope owners listed successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ListAvailableEnvelopeOwnersApiResponse'
              example:
                items:
                  - id: 1bea60bc-03fc-4124-ba35-2afda8a5b3d9
                    name: John Doe
                    isActive: true
                    isCurrent: false
                    role: Administrator
                    avatar: https://example.com/avatar.jpg
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
      x-apidog-folder: Envelope
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252681-run
components:
  schemas:
    ListAvailableEnvelopeOwnersApiResponse:
      type: object
      properties:
        items:
          type: array
          items:
            $ref: '#/components/schemas/ListAvailableEnvelopeOwnersApiResponseItem'
          description: List of available envelope owners
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - items
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ListAvailableEnvelopeOwnersApiResponseItem:
      type: object
      properties:
        id:
          type: string
          description: Id of the owner
          format: uuid
        name:
          type: string
          description: Name of the owner
          nullable: true
        isActive:
          type: boolean
          description: Indicates if the owner is active
        isCurrent:
          type: boolean
          description: Indicates if is the current owner
        role:
          enum:
            - Administrator
            - User
          type: string
          description: Role of the owner
        avatar:
          type: string
          description: Url avatar of the owner
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - name
        - isActive
        - isCurrent
        - role
        - avatar
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
