# List current user vaults

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/vaults/user-accounts:
    get:
      summary: List current user vaults
      deprecated: false
      description: >-
        Retrieves a list of vaults that are accessible to the currently
        authenticated user. This includes personal vaults, as well as shared or
        restricted vaults the user has access to, depending on their
        permissions.
      tags:
        - Vault
        - Vault
      parameters:
        - name: IsFavorite
          in: query
          description: ' Filter vaults by favorite status'
          required: false
          schema:
            type: boolean
            description: ' Filter vaults by favorite status'
        - name: Search
          in: query
          description: Search vaults by id or name
          required: false
          schema:
            maxLength: 200
            type: string
            description: Search vaults by id or name
        - name: PageSize
          in: query
          description: Number of items per page
          required: true
          schema:
            maximum: 100
            minimum: 1
            type: integer
            description: Number of items per page
            format: int32
            default: 10
        - name: PageNumber
          in: query
          description: Current page number
          required: true
          schema:
            maximum: 4294967295
            minimum: 1
            type: integer
            description: Current page number
            format: int32
            default: 1
        - name: OrderByDirection
          in: query
          description: Order by direction
          required: true
          schema:
            enum:
              - ASC
              - DESC
            type: string
            description: Order by direction
            default: DESC
      responses:
        '200':
          description: User account vaults listed successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ListUserAccountVaultsApiResponse'
              example:
                items:
                  - id: 500d5305-a6cc-44d3-8fee-f6c2f3c94327
                    createdAtUtc: '2026-08-17T00:52:53.9901844Z'
                    createdByName: John Doe
                    createdByAvatar: null
                    updatedAtUtc: '2026-08-17T00:52:53.9901845Z'
                    updatedByName: John Doe
                    updatedByAvatar: null
                    name: Example Vault
                    type: UserAccount
                    isFavorite: false
                    canBeRemoved: true
                    userAccountMembers:
                      - id: 2f4256cc-3db6-404a-9d30-a0c29957fccb
                        name: John Doe
                        email: john.doe@example.com
                        avatar: null
                        isActive: true
                pagination:
                  totalItems: 1
                  pageSize: 10
                  pageNumber: 1
                  pageItems: 1
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
      security:
        - apikey-header-x-api-token: []
      x-apidog-folder: Vault
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252708-run
components:
  schemas:
    ListUserAccountVaultsApiResponse:
      type: object
      properties:
        items:
          type: array
          items:
            $ref: '#/components/schemas/ListUserAccountVaultsApiResponseItem'
          description: List of vaults
          nullable: true
        pagination:
          allOf:
            - $ref: '#/components/schemas/Pagination'
          description: Pagination metadata
          type: 'null'
      additionalProperties: false
      x-apidog-orders:
        - items
        - pagination
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    Pagination:
      type: object
      properties:
        totalItems:
          type: integer
          description: Total number of items
          format: int32
        pageSize:
          type: integer
          description: Total number of pages
          format: int32
        pageNumber:
          type: integer
          description: Current page number
          format: int32
        pageItems:
          type: integer
          description: Number of items per page
          format: int32
      additionalProperties: false
      description: Pagination metadata
      x-apidog-orders:
        - totalItems
        - pageSize
        - pageNumber
        - pageItems
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ListUserAccountVaultsApiResponseItem:
      type: object
      properties:
        id:
          type: string
          description: Id of the vault
          format: uuid
        createdAtUtc:
          type: string
          description: Date and time when the vault was created in UTC
          format: date-time
        createdByName:
          type: string
          description: Name of the user who created the vault
          nullable: true
        createdByAvatar:
          type: string
          description: Url of the avatar of the user who created the vault
          nullable: true
        updatedAtUtc:
          type: string
          description: Date and time when the vault was last updated in UTC
          format: date-time
        updatedByName:
          type: string
          description: Name of the user who last updated the vault
          nullable: true
        updatedByAvatar:
          type: string
          description: Url of the avatar of the user who last updated the vault
          nullable: true
        name:
          type: string
          description: Name of the vault
          nullable: true
        type:
          enum:
            - Account
            - UserAccount
            - UserAccountGroup
          type: string
          description: Type of the vault
        isFavorite:
          type: boolean
          description: Indicates if the vault is marked as favorite
        canBeRemoved:
          type: boolean
          description: >-
            Indicates if the vault can be removed, only vaults with no items
            (envelopes and templates) can be removed
        userAccountMembers:
          type: array
          items:
            $ref: >-
              #/components/schemas/ListUserAccountVaultsApiResponseItemUserAccountMember
          description: Members of the vault
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - createdAtUtc
        - createdByName
        - createdByAvatar
        - updatedAtUtc
        - updatedByName
        - updatedByAvatar
        - name
        - type
        - isFavorite
        - canBeRemoved
        - userAccountMembers
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ListUserAccountVaultsApiResponseItemUserAccountMember:
      type: object
      properties:
        id:
          type: string
          description: Id of the member
          format: uuid
        name:
          type: string
          description: Name of the member
          nullable: true
        email:
          type: string
          description: Email of the member
          nullable: true
        avatar:
          type: string
          description: Url avatar of the member
          nullable: true
        isActive:
          type: boolean
          description: Indicates if the member is active
      additionalProperties: false
      x-apidog-orders:
        - id
        - name
        - email
        - avatar
        - isActive
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
