# Error handling

The Signater API provides detailed error-handling mechanisms to ensure clarity and ease of troubleshooting. Each endpoint can respond with different HTTP status codes documented per endpoint in the API reference. For example, some endpoints may return a `404 Not Found` response when a requested resource is unavailable, while others do not. The potential for a `404` response is explicitly noted in the documentation for each endpoint.

## Common HTTP Status Codes

Here is a list of the common HTTP status codes returned by the API:

### 400 Bad Request

This status is used when the input data is invalid or improperly formatted. The response body includes an `errors` array with detailed messages:

```json
{
  "errors": [
    {
      "message": "The length of 'Name' must be at least 2 characters. You entered 1 characters.",
      "metadata": null
    }
  ]
}
````

### 401 Unauthorized

This status code indicates that the provided token is invalid, or inactive, or the user associated with the token is no longer active for the corresponding account. In such cases, the response body includes a JSON object:

```json
{
  "message": "Authentication failed."
}
```

### 402 Payment Required

This status code is returned when the account does not have sufficient funds to perform the requested operation. Currently, the only endpoint returning this status is the envelope publication endpoint. The response includes a `ShouldBuy` property in the body, which indicates the product needed to proceed:

```json
{
  "ShouldBuy": "ApiEnvelopes"
}
````

Possible values for `ShouldBuy`:

- **Envelopes**: Indicates that the envelope was created through the Signater web client, requiring an active subscription plan for unlimited envelopes.
- **ApiEnvelopes**: Indicates that the envelope was created via the API, requiring a one-time envelope credit purchase.
- **SignerMfaCredits**: Indicates insufficiently advanced authentication credits for one or more signers.

For all purchases, users can visit https://app.signater.com/billing/plans. Note that users must have administrator privileges to access this page.

### 403 Forbidden

This status code occurs when the user lacks the necessary permissions to access the resource. For example, if administrator rights are required:

```json
{
  "message": "User should be an administrator."
}
````

### 429 Too Many Requests

This status code indicates that the rate limit for the API has been exceeded. The Signater API enforces a limit of 1,000 requests per minute. When this limit is reached, the client must wait before making additional requests. A response with a `Retry-After` header may be included to specify the recommended wait time in seconds:

```makefile
HTTP/1.1 429 Too Many Requests
Retry-After: 60
```

To prevent exceeding rate limits, consider implementing client-side request throttling or retry strategies with exponential backoff.

## Request Telemetry


Every response from the Signater API includes a unique telemetry ID in the X-`Signater-Telemetry-Operation-Id` response header. This ID can be provided to our support team to facilitate the investigation of specific requests:

```makefile
X-Signater-Telemetry-Operation-Id: 123e4567-e89b-12d3-a456-426614174000
```
