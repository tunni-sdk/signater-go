# Update vault

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/vaults/{vaultId}:
    patch:
      summary: Update vault
      deprecated: false
      description: >-
        Updates the properties of an existing vault. Note that the access type
        of a vault cannot be changed once it is created. However, for vaults
        with 'member' access type, you can add or remove members, modifying who
        has access to the vault.
      tags:
        - Vault
        - Vault
      parameters:
        - name: vaultId
          in: path
          description: Vault ID
          required: true
          example: ''
          schema:
            type: string
            description: Vault ID
            format: uuid
      requestBody:
        content:
          application/json:
            schema:
              allOf:
                - $ref: '#/components/schemas/UpdateVaultApiRequest'
            example:
              name: New Vault name
              description: New Vault Description
              memberIds:
                - 878341a6-a942-408e-880e-32910af99227
                - 90496500-64de-4225-9409-16b29de1bbc8
      responses:
        '204':
          description: Vault updated successfully
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
      x-apidog-folder: Vault
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252702-run
components:
  schemas:
    UpdateVaultApiRequest:
      required:
        - name
        - type
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
          description: >-
            Type of the vault — can be changed across UserAccount / Account /
            UserAccountGroup
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
