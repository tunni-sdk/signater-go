# Create vault

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/vaults:
    post:
      summary: Create vault
      deprecated: false
      description: >-
        Creates a new vault in the system. A vault is a secure storage container
        used to organize and store templates and envelopes. Vaults can have
        different access levels: personal vaults, which are accessible only to
        the user who created them; account vaults, which are accessible to all
        users within the account; and member based vaults, where access can be
        granted to specific members, allowing for granular control over who can
        view or modify the contents.
      tags:
        - Vault
        - Vault
      parameters: []
      requestBody:
        content:
          application/json:
            schema:
              allOf:
                - $ref: '#/components/schemas/CreateVaultApiRequest'
            example:
              name: My Vault
              description: My Vault Description
              type: UserAccountGroup
              memberIds:
                - 136ad90e-2f5e-46d1-85cf-b0b94457f334
                - bfdf91ca-093b-4c15-968b-e602375d8152
      responses:
        '201':
          description: Vault created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CreateVaultApiResponse'
              example:
                id: b01eabec-dd0f-4c23-b537-65f40a1e216d
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
      x-apidog-folder: Vault
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252701-run
components:
  schemas:
    CreateVaultApiRequest:
      required:
        - name
      type: object
      properties:
        name:
          maxLength: 200
          minLength: 2
          type: string
          description: Name of the vault
        description:
          maxLength: 400
          type: string
          description: Description of the vault
          nullable: true
        type:
          enum:
            - Account
            - UserAccount
            - UserAccountGroup
          type: string
          description: Type of the vault
        memberIds:
          type: array
          items:
            type: string
            format: uuid
          description: >-
            User member IDs of the vault (required for UserAccountGroup type,
            you also must explicitly add yourself if you want to be a member)
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - name
        - description
        - type
        - memberIds
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    CreateVaultApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the vault
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
