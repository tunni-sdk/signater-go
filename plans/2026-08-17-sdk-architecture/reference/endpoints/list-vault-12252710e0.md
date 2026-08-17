# List vault

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/vaults/{vaultId}/list:
    get:
      summary: List vault
      deprecated: false
      description: >-
        Retrieves the contents of a specific vault, identified by its unique ID.
        This includes a list of envelopes and templates stored within the vault,
        allowing for easy management and access to the resources.
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
        - name: ExcludeEnvelopes
          in: query
          description: If true, excludes envelopes from the response
          required: false
          schema:
            type: boolean
            description: If true, excludes envelopes from the response
        - name: ExcludeTemplates
          in: query
          description: If true, excludes templates from the response
          required: false
          schema:
            type: boolean
            description: If true, excludes templates from the response
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
          description: Vault listed successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ListVaultApiResponse'
              example:
                vault:
                  id: a5105f68-a502-4c99-8ca0-3a1be6331531
                  isSandbox: false
                  name: Example Vault
                  type: UserAccount
                  canBeRemoved: true
                  userAccountMembers:
                    - id: 29c00159-c5cc-4748-9740-5baf1cf88c1a
                      name: John Doe
                      avatar: null
                    - id: 357e9756-e9fe-4cca-92e2-6ddb86c8f502
                      name: Jane Smith
                      avatar: https://example.com/avatars/janesmith.png
                items:
                  - type: Envelope
                    envelope:
                      id: b76cc7e2-2ba6-450c-98f8-e57cb812935c
                      createdAtUtc: '2026-08-17T00:52:54.0210429Z'
                      createdById: f1811b7d-6841-44c7-9ed0-8a48f7227498
                      createdByName: John Doe
                      createdByAvatar: null
                      updatedAtUtc: '2026-08-17T00:52:54.0210444Z'
                      lastUpdateById: f7e3d9eb-ca18-44e7-85dc-c0a356fd223c
                      updatedByName: John Doe
                      updatedByAvatar: null
                      status: Draft
                      name: Example Envelope
                      privateDescription: >-
                        This is an example of an envelope used for demonstration
                        purposes.
                      canBeRemoved: true
                      signers:
                        - name: Jane Smith
                          email: jane.smith@example.com
                    template: null
                  - type: Template
                    envelope: null
                    template:
                      id: fe4276b7-03ce-420c-939c-748c9fb5ba6f
                      createdAtUtc: '2026-08-17T00:52:54.0210477Z'
                      createdById: 42e6dc9c-f97b-41be-a623-9e936122afd1
                      createdByName: John Doe
                      createdByAvatar: null
                      updatedAtUtc: '2026-08-17T00:52:54.0210491Z'
                      lastUpdateById: d12337d6-235b-43da-bb98-86d881b4d566
                      updatedByName: John Doe
                      updatedByAvatar: null
                      name: Example Template
                pagination:
                  totalItems: 2
                  pageSize: 10
                  pageNumber: 1
                  pageItems: 2
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
      x-apidog-folder: Vault
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252710-run
components:
  schemas:
    ListVaultApiResponse:
      type: object
      properties:
        vault:
          allOf:
            - $ref: '#/components/schemas/ListVaultVaultApiResponse'
          description: Vault object
          type: 'null'
        items:
          type: array
          items:
            $ref: '#/components/schemas/ListVaultApiResponseItem'
          description: Items of the vault
          nullable: true
        pagination:
          allOf:
            - $ref: '#/components/schemas/Pagination'
          description: Pagination metadata
          type: 'null'
      additionalProperties: false
      x-apidog-orders:
        - vault
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
    ListVaultApiResponseItem:
      type: object
      properties:
        type:
          enum:
            - Envelope
            - Template
          type: string
          description: Type of the item
        envelope:
          allOf:
            - $ref: '#/components/schemas/ListVaultApiResponseItemEnvelope'
          description: Envelope object
          type: 'null'
        template:
          allOf:
            - $ref: '#/components/schemas/ListVaultApiResponseItemTemplate'
          description: Template object
          type: 'null'
      additionalProperties: false
      x-apidog-orders:
        - type
        - envelope
        - template
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ListVaultApiResponseItemTemplate:
      type: object
      properties:
        id:
          type: string
          description: Id of the template
          format: uuid
        createdAtUtc:
          type: string
          description: Date and time when the template was created in UTC
          format: date-time
        createdById:
          type: string
          description: Id of the that created the template
          format: uuid
        createdByName:
          type: string
          description: Name of the that created the template
          nullable: true
        createdByAvatar:
          type: string
          description: Url of the avatar that created the template
          nullable: true
        updatedAtUtc:
          type: string
          description: Date and time when the template was last updated in UTC
          format: date-time
        lastUpdateById:
          type: string
          description: Id of the that last updated the template
          format: uuid
        updatedByName:
          type: string
          description: Name of the that last updated the template
          nullable: true
        updatedByAvatar:
          type: string
          description: Url of the avatar that last updated the template
          nullable: true
        name:
          type: string
          description: Name of the template
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - createdAtUtc
        - createdById
        - createdByName
        - createdByAvatar
        - updatedAtUtc
        - lastUpdateById
        - updatedByName
        - updatedByAvatar
        - name
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ListVaultApiResponseItemEnvelope:
      type: object
      properties:
        id:
          type: string
          description: Id of the envelope
          format: uuid
        createdAtUtc:
          type: string
          description: Date and time when the envelope was created in UTC
          format: date-time
        createdById:
          type: string
          description: Id of the that created the envelope
          format: uuid
        createdByName:
          type: string
          description: Name of the that created the envelope
          nullable: true
        createdByAvatar:
          type: string
          description: Url of the avatar that created the envelope
          nullable: true
        updatedAtUtc:
          type: string
          description: Date and time when the envelope was last updated in UTC
          format: date-time
        lastUpdateById:
          type: string
          description: Id of the that last updated the envelope
          format: uuid
        updatedByName:
          type: string
          description: Name of the that last updated the envelope
          nullable: true
        updatedByAvatar:
          type: string
          description: Url of the avatar that last updated the envelope
          nullable: true
        status:
          enum:
            - Draft
            - PublishScheduled
            - Published
            - Hold
            - Cancelled
            - CancelledBySignerMfaError
            - Rejected
            - Signed
            - Expired
          type: string
          description: Status of the envelope
        name:
          type: string
          description: Name of the envelope
          nullable: true
        privateDescription:
          type: string
          description: Public description of the envelope
          nullable: true
        canBeRemoved:
          type: boolean
          description: >-
            Indicates if the envelope can be removed, envelopes with status
            'Published' cannot be removed, in this case first you should
            unpublish (hold) the envelope
        signers:
          type: array
          items:
            $ref: '#/components/schemas/ListVaultApiResponseItemEnvelopeSigner'
          description: Signers of the envelope
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - createdAtUtc
        - createdById
        - createdByName
        - createdByAvatar
        - updatedAtUtc
        - lastUpdateById
        - updatedByName
        - updatedByAvatar
        - status
        - name
        - privateDescription
        - canBeRemoved
        - signers
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ListVaultApiResponseItemEnvelopeSigner:
      type: object
      properties:
        name:
          type: string
          description: Id of the signer
          nullable: true
        email:
          type: string
          description: Email of the signer
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - name
        - email
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ListVaultVaultApiResponse:
      type: object
      properties:
        id:
          type: string
          description: Id of the vault
          format: uuid
        isSandbox:
          type: boolean
          description: Indicates if the vault is in the sandbox mode
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
        canBeRemoved:
          type: boolean
          description: >-
            Indicates if the vault can be removed, only vaults with no items
            (envelopes and templates) can be removed
        userAccountMembers:
          type: array
          items:
            $ref: '#/components/schemas/ListVaultUserAccountMemberApiResponse'
          description: Members of the vault
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - isSandbox
        - name
        - type
        - canBeRemoved
        - userAccountMembers
      x-apidog-ignore-properties: []
      x-apidog-folder: ''
    ListVaultUserAccountMemberApiResponse:
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
        avatar:
          type: string
          description: Url avatar of the member
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - id
        - name
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
