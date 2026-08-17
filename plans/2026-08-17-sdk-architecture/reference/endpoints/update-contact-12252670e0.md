# Update contact

## OpenAPI Specification

```yaml
openapi: 3.0.1
info:
  title: ''
  description: ''
  version: 1.0.0
paths:
  /v1/ecm/contacts/{contactId}:
    put:
      summary: Update contact
      deprecated: false
      description: >-
        Updates the details of an existing contact in the system using its
        unique identifier. This endpoint allows modifications to fields such as
        name, email, or other associated metadata. Ensure the contact ID
        provided exists and that the payload contains the updated information.
      tags:
        - Contact
        - Contact
      parameters:
        - name: contactId
          in: path
          description: Contact ID
          required: true
          example: ''
          schema:
            type: string
            description: Contact ID
            format: uuid
      requestBody:
        content:
          application/json:
            schema:
              allOf:
                - $ref: '#/components/schemas/UpdateContactApiRequest'
            example:
              name: John Doe
              title: Project Manager
              email: john.doe@example.com
              phoneIdd: 1
              phoneNumber: '1234567890'
              documentType: GenericIdentification
              documentValue: A1234567
              description: A contact person for project management.
      responses:
        '204':
          description: Contact updated successfully
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
      x-apidog-folder: Contact
      x-apidog-status: released
      x-run-in-apidog: https://app.apidog.com/web/project/755605/apis/api-12252670-run
components:
  schemas:
    UpdateContactApiRequest:
      required:
        - name
      type: object
      properties:
        name:
          maxLength: 200
          minLength: 2
          type: string
          description: Name of the contact
        title:
          maxLength: 200
          type: string
          description: Title of the contact
          nullable: true
        email:
          pattern: ^[\w\.\-\+]+@([\w\-]+\.)+[a-zA-Z]{2,}$
          type: string
          description: Email address of the contact
          nullable: true
        phoneIdd:
          type: integer
          description: >-
            Phone IDD of the contact, e.g. '1' for the United States or '55' for
            Brazil
          format: int32
          nullable: true
        phoneNumber:
          maxLength: 200
          type: string
          description: Phone number of the contact excluding the IDD
          nullable: true
        documentType:
          enum:
            - GenericIdentification
            - BrazilianCpf
            - BrazilianIdentity
            - Passport
          type: string
          description: Document type of the contact
          nullable: true
        documentValue:
          maxLength: 100
          type: string
          description: Document value of the contact
          nullable: true
        description:
          maxLength: 2000
          type: string
          description: Description of the contact
          nullable: true
      additionalProperties: false
      x-apidog-orders:
        - name
        - title
        - email
        - phoneIdd
        - phoneNumber
        - documentType
        - documentValue
        - description
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
