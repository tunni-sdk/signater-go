# Authentication

To use the Signater API, you must give an `x-api-token` header with each request. The value of this header should be the token generated from the API Tokens section in your profile at Signater App. Failing to send a token, using an incorrect or inactive token, or using a token associated with an inactive user will result in a `401 Unauthorized` response.

## Sandbox Environment

Signater provides a sandbox environment for testing without needing a paid plan. You can generate tokens specifically for sandbox use. When making API calls with a sandbox token, the API automatically switches to sandbox mode to access resources accordingly.

If you encounter any issues accessing the sandbox environment, please get in touch with our support team for assistance and instructions on enabling sandbox mode.

## Rate Limiting

The Signater API enforces a rate limit of 1,000 requests per minute. Exceeding this limit may result in throttling or temporary restrictions on further requests.
