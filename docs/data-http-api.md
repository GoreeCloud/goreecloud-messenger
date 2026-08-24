# GoreeCloud Messenger Data HTTP API Foundation

## Purpose

This document defines the HTTP transport boundary for GoreeCloud Data messaging. The HTTP layer is an adapter around the Data and delivery-receipt services and does not replace their authorization, transport-provenance, or E2EE-only validation rules.

## Current endpoints

- `POST /v1/data/messages` accepts an encrypted GoreeCloud Data envelope.
- `GET /v1/data/conversations/{conversationID}/messages` returns encrypted envelopes for one authorized conversation.
- `POST /v1/data/messages/{messageID}/receipts` records an authenticated recipient delivery or read acknowledgement.
- `GET /v1/data/messages/{messageID}/receipts` returns receipt state to an authorized conversation participant.

The API transports ciphertext as standard base64 text inside JSON. It does not accept a plaintext message-body field.

## Authentication boundary

The HTTP package requires an `Authenticator` implementation. The authenticator must resolve the request to an authenticated GoreeCloud user identifier before a service is called.

Credential issuance, login, token creation, token storage, session management, device identity, and identity-provider integration remain outside this milestone. The API must not infer identity from client-supplied sender or receipt metadata.

## Authorization boundary

The authenticated user is passed to service-layer authorization. The implementation:

- binds the authenticated user to the envelope sender;
- verifies server-side conversation membership;
- binds a receipt to the authenticated recipient;
- rejects delivery/read acknowledgements from the message sender for that sender's own message;
- verifies that a receipt references an existing message in the stated conversation;
- prevents receipt state from moving backwards from `read` to `delivered`;
- rejects duplicate message identifiers;
- rejects client-nonce reuse; and
- enforces the E2EE-only Data-envelope contract.

The HTTP layer maps authorization failures to bounded non-success responses without returning internal error details.

## Receipt semantics

Delivery receipts are GoreeCloud Data metadata, not proof of carrier delivery and not cryptographic proof that a human viewed content. `delivered` means an authorized recipient client reported delivery progress. `read` is a later recipient-observed state. Receipt state is monotonic for each message and recipient in the current service contract.

Receipts contain message, conversation, recipient, state, and observation timestamp only. They do not contain plaintext message bodies, encryption keys, device secrets, or carrier metadata.

## Privacy and security behavior

The HTTP adapter:

- accepts ciphertext rather than plaintext message bodies;
- rejects unknown JSON fields;
- bounds request bodies to 1 MiB;
- sets `Cache-Control: no-store` on JSON responses;
- sets `X-Content-Type-Options: nosniff`;
- keeps SMS, MMS, and RCS outside the Data API;
- authorizes receipt operations against server-side message and conversation state; and
- does not claim that a cryptographic session exists merely because the envelope contract records `e2ee`.

## Current limitations

This is a Development source foundation. It does not provide production credentials, production persistence, TLS termination, rate limiting, push delivery, attachment transfer, cryptographic session establishment, key management, multi-device synchronization, carrier adapters, client packaging, deployment, or production acceptance. Receipt storage is currently an in-memory development implementation and is not production persistence.
