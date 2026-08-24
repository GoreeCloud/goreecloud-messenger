# GoreeCloud Messenger Data HTTP API Foundation

## Purpose

This document defines the first HTTP transport boundary for GoreeCloud Data messaging. The HTTP layer is an adapter around the existing Data service and does not replace its authorization, transport-provenance, or E2EE-only validation rules.

## Current endpoints

- `POST /v1/data/messages` accepts an encrypted GoreeCloud Data envelope.
- `GET /v1/data/conversations/{conversationID}/messages` returns encrypted envelopes for one authorized conversation.

The API transports ciphertext as standard base64 text inside JSON. It does not accept a plaintext message-body field.

## Authentication boundary

The HTTP package requires an `Authenticator` implementation. The authenticator must resolve the request to an authenticated GoreeCloud user identifier before the Data service is called.

Credential issuance, login, token creation, token storage, session management, device identity, and identity-provider integration are intentionally outside this milestone. The API must not infer identity from client-supplied message metadata.

## Authorization boundary

The authenticated user is passed to the Data service, which remains responsible for:

- binding the authenticated user to the envelope sender;
- verifying server-side conversation membership;
- rejecting duplicate message identifiers;
- rejecting client-nonce reuse; and
- enforcing the E2EE-only Data-envelope contract.

The HTTP layer maps authorization failures to non-success responses without returning internal error details.

## Privacy and security behavior

The HTTP adapter:

- accepts ciphertext rather than plaintext message bodies;
- rejects unknown JSON fields;
- bounds request bodies to 1 MiB;
- sets `Cache-Control: no-store` on JSON responses;
- sets `X-Content-Type-Options: nosniff`;
- keeps SMS, MMS, and RCS outside the Data API; and
- does not claim that a cryptographic session exists merely because the envelope contract records `e2ee`.

## Current limitations

This is a Development source foundation. It does not provide production credentials, production persistence, TLS termination, rate limiting, push delivery, attachment transfer, delivery receipts, cryptographic session establishment, key management, multi-device synchronization, carrier adapters, client packaging, deployment, or production acceptance.
