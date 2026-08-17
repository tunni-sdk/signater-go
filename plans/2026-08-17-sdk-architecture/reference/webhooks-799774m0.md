# Webhooks

Signater's webhooks allow you to receive real-time events related to envelope activities on the platform. Webhooks are configured at the account level and can be created in either production or sandbox mode. The behavior of webhooks depends on the selected environment filter, meaning they can be configured to receive events from either production or sandbox.

## Environment Modes
- Sandbox Mode: Webhooks created in sandbox mode do not require a paid plan or the purchase of the API addon. They are ideal for testing and development.
- Production Mode: For production webhooks, a paid plan and the purchase of the API addon are required.

## Accessing and Managing Webhooks

To `view`, `create`, `edit`, `activate`or `deactivate` your configured webhooks, visit your account's webhook page:
https://app.signater.com/account/webhooks

## Creating Webhooks
Each webhook registered in Signater must include the following information:

- **Name** (required)
- **Description** (optional)
- **Endpoint Method** (required): The HTTP method to be used (e.g., `POST`, `GET`, etc.)
- **Endpoint URL** (required): The URL of the endpoint to be called when the event occurs
- **Events** (required): The types of events the webhook will receive
- **Authentication Type** (optional): The authentication method used, with `HookdeckSignature` being the default

## Supported Events

The events that can be configured for the webhook are:

- `envelope.created`: This event is triggered when a new envelope is created in the system. An envelope is a container for documents that need to be signed, and this event signifies that a new envelope has been initiated and is ready for further processing.
- `envelope.publish_scheduled`: This event occurs when an envelope's publication is scheduled. It indicates that the envelope has been planned for publishing at a later time, but hasn't yet been made available for signing.
- `envelope.published`: Triggered when an envelope is officially published. This means that the envelope is now available for signing by the designated parties. It signifies the start of the signing process.
- `envelope.published_by_schedule`: This event happens when an envelope is published automatically according to a pre-scheduled time. It's similar to the envelope.published event, but specifically indicates that the publication was triggered by a scheduling rule.
- `envelope.unpublished`: This event is sent when an envelope is unpublished, meaning it has been removed or made unavailable for signing. It could happen if the envelope was canceled or revoked before being signed.
- `envelope.updated`: This event occurs when an envelope is updated. Updates may include changes to the envelope's details, documents, recipients, or other attributes that are necessary for the workflow.
- `envelope.cancelled`: Triggered when an envelope is canceled, meaning the envelope and its associated signing process are no longer valid. This could happen for various reasons, such as an error or a change in requirements.
- `envelope.expired`: This event happens when an envelope expires, meaning the envelope is no longer valid for signing due to the expiration of its designated timeframe.
- `envelope.signed`: This event is triggered when an envelope has been fully signed by all the intended signers. It signifies that the signing process is complete, and the envelope is now finalized.
- `envelope.viewed_by_signer`: This event occurs when a signer views the envelope, but has not necessarily signed it yet. It indicates that the signer has accessed the envelope for review.
- `envelope.cancelled_due_to_mfa_error_by_signer`: Triggered when an envelope is canceled because the signer encountered an error related to Multi-Factor Authentication (MFA). This event signifies that the signing process was stopped due to a failure in the authentication step.
- `envelope.rejected_by_signer`: This event happens when a signer rejects an envelope, meaning they decided not to sign or agree to the contents of the document. This can happen at any stage of the signing process.
- `envelope.approved_by_signer`: This event occurs when a signer approves the envelope, indicating their consent to the contents, but this approval might not necessarily be a final signature. It could be part of an approval workflow or a pre-signing step.


## Data Sent

For all events, the envelope ID will be sent. For the events with the `_by_signer` suffix (the last four types of events), the signer ID will also be included.

## Important Notes

- **Event Order**: Signater does not guarantee the order of event delivery. It is the client's responsibility to manage events correctly. We always recommend querying the API for the most up-to-date resource information.
- **Retry Mechanism**: Signater's webhook system uses Hookdeck, which provides a robust retry mechanism. The retry rule is **Retry exponentially every 30 seconds up to 5 times for any error**.
- **Asynchronous Processing**: We recommend processing received webhooks asynchronously to improve the reliability of the webhook ingestion system.

## Example Payload

The payload sent by a webhook might look like this:

```json
{
  "envelope_id": "696d0b7f-8e2f-4fb1-9605-72b197ee154d",
  "event_type": "envelope.created",
  "account_id": "00cae3c2-8b7c-41b6-b934-6724808fe7f7",
  "env": "production"
}
```
In this example, we can see:

- `envelope_id`: The ID of the envelope that triggered the event
- `event_type`: The type of event that was triggered
- `account_id`: The ID of the account that triggered the event
- `env`: The environment, which can be either `production` or `sandbox`

This is a basic example of how the webhook data is sent to the configured endpoint.

## Webhooks CLI Debugging

For local debugging of webhooks, Signater provides an integration with the Hookdeck CLI, allowing you to receive webhook events locally on your machine. This integration is particularly useful for testing and debugging, especially when working with production or sandbox webhooks.

### Installing Hookdeck CLI

To use the Hookdeck CLI for debugging, first, you need to install it. Hookdeck CLI is available for macOS, Windows, Linux, and Docker. You can install it using various methods depending on your platform.

#### Installation Methods

- NPM & Yarn (Cross-platform)
```bash
$ npm install hookdeck-cli -g
$ yarn global add hookdeck-cli
```
- macOS (Using Homebrew)
```bash
$ brew install hookdeck/hookdeck/hookdeck
```
- Windows (Using Scoop package manager)
```bash
$ scoop bucket add hookdeck https://github.com/hookdeck/scoop-hookdeck-cli.git
$ scoop install hookdeck
```
- Linux (Without package manager)
    - Download the latest release's tar.gz file.
    - Unzip the file:
    - ```bash
        $ tar -xvf hookdeck_X.X.X_linux_x86_64.tar.gz
        ```
    - Run the executable:
    - ```bash
        $ ./hookdeck
        ```

For detailed instructions, refer to the [official Hookdeck CLI documentation](https://hookdeck.com/docs/cli).

### Setting Up CLI for Debugging in Signater

In Signater, you can enable the Hookdeck CLI for debugging webhook events directly on your local machine. To do this:

1. Go to the **Webhooks** section of your Signater account.
2. Click the button next to the "Create Webhook" button.
3. Select CLI Settings from the dropdown menu.

This will open a modal where you can configure the CLI URL generated by Hookdeck. Additionally, you'll find a toggle to activate or deactivate the CLI integration.

### Important Notes:
- **Environment-Specific**: The CLI is environment-dependent. To receive events from the sandbox environment, ensure you are viewing the sandbox webhook configuration page. This allows you to receive webhook events locally for both production and sandbox environments, without affecting previously registered webhooks.
- **Temporary Connection**: When you activate the CLI, the system will forward all webhook events (from production or sandbox) to the locally running Hookdeck CLI. This temporary connection allows you to redirect the events to a third-party endpoint, typically a locally running API, for full debugging of your solution.
- **Automatic Deactivation**: The CLI will be automatically deactivated after 2 hours unless explicitly deactivated by the user.

With this setup, you can efficiently debug webhook events in a local development environment, simulating the behavior of webhooks in production or sandbox without modifying your existing webhooks configuration.

### Webhook Request Headers

Each webhook request sent by the Signater system to the user's registered endpoint contains specific headers that provide important information about the request.

Example:

```json
{
  "content-length": "157",
  "content-type": "application/json; charset=utf-8",
  "request-id": "|6b1bc96c7cee9179e177e8ce95b18530.09a8251c01c86fdb.",
  "x-signater-apikey": "35198a2a282a41e5b1c64d55bc55c81c"
}
```
Here's a breakdown of the headers included in the request:

- `content-length`: This header indicates the size of the body content in bytes. For example, a value of 157 means that the body of the request is 157 bytes long.
- `content-type`: Specifies the media type of the request body. In this case, it is set to application/json; charset=utf-8, which means the body contains JSON data encoded in UTF-8.
- `request-id`: A unique identifier for the webhook request, useful for tracing and debugging. This helps correlate logs and track the request through the system. For example, the value |6b1bc96c7cee9179e177e8ce95b18530.09a8251c01c86fdb. is used to uniquely identify the request.
- `x-signater-apikey`: This header contains the API key associated with the Signater account that triggered the webhook. It is used for authentication and authorization purposes. For instance, the value 35198a2a282a41e5b1c64d55bc55c81c represents the API key used in the request.

These headers ensure that the webhook request is properly processed by the user's endpoint, allowing for secure and accurate handling of events triggered by the Signater system.

