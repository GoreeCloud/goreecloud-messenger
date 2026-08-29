# GoreeCloud Messenger Data HTTP API Foundation

## Purpose

This document defines the HTTP transport boundary for GoreeCloud Data messaging. The HTTP layer is an adapter around the Data, receipt, and encrypted-attachment services and does not replace their authorization, transport-provenance, or E2EE-only validation rules.

## Application-facing runtime composition

`DataRuntimeHandler` is the current application-facing composition boundary for the implemented Data HTTP surface. It requires the Data service, receipt service, attachment service, and one `Authenticator`, then registers all implemented routes directly on a shared `http.ServeMux`.

Direct route registration is intentional. Message listing and attachment listing both occupy the `/v1/data/conversations/{conversationID}/...` namespace; prefix-mounting separate child handlers could hide one surface behind another. The composed runtime keeps both route families reachable without duplicating authorization logic or inventing a second authentication boundary.

The runtime composition does not create credentials, establish cryptographic sessions, choose production persistence, terminate TLS, configure rate limits, or start a production listener. Those remain explicit outer application/deployment responsibilities.

## Current endpoints

- `POST /v1/data/messages` accepts an encrypted GoreeCloud Data envelope.
- `GET /v1/data/conversations/{conversationID}/messages` returns encrypted envelopes for one authorized conversation.
- `POST /v1/data/messages/{messageID}/receipts` records an authenticated recipient delivery or read acknowledgement.
- `GET /v1/data/messages/{messageID}/receipts` returns receipt state to an authorized conversation participant.
- `POST /v1/data/attachments` accepts already-encrypted opaque attachment ciphertext plus bounded metadata.
- `GET /v1/data/attachments/{attachmentID}` returns opaque ciphertext as base64 JSON to an authorized conversation participant.
- `GET /v1/data/attachments/{attachmentID}/ciphertext` returns the exact stored ciphertext bytes to an authorized conversation participant as `application/octet-stream`.
- `GET /v1/data/conversations/{conversationID}/attachments` returns a bounded metadata-only attachment projection.
- `DELETE /v1/data/attachments/{attachmentID}` removes retrievable ciphertext for an authorized participant and is idempotent for already-missing/deleted attachment identifiers.

Message and ordinary attachment JSON ciphertext is transported as standard base64 text. The raw ciphertext download endpoint transports encrypted bytes directly. The API does not accept plaintext message-body or plaintext attachment-content fields.

## Authentication boundary

The HTTP package requires an `Authenticator` implementation. The authenticator must resolve the request to an authenticated GoreeCloud user identifier before a service is called. The composed runtime deliberately reuses the same injected authenticator for message, receipt, and attachment routes.

Credential issuance, login, token creation, token storage, session management, device identity, and identity-provider integration remain outside this milestone. The API must not infer identity from client-supplied sender, attachment, or receipt metadata.

## Authorization boundary

The authenticated user is passed to service-layer authorization. The implementation:

- binds the authenticated user to message and attachment senders;
- verifies server-side conversation membership;
- binds a receipt to the authenticated recipient;
- rejects delivery/read acknowledgements from the message sender for that sender's own message;
- verifies that a receipt references an existing message in the stated conversation;
- prevents receipt state from moving backwards from `read` to `delivered`;
- rejects duplicate message and attachment identifiers;
- rejects message and attachment client-nonce reuse;
- authorizes attachment JSON fetch, raw ciphertext fetch, list, and delete operations against conversation membership; and
- enforces the E2EE-only Data-envelope contract.

The HTTP layer maps authorization failures to bounded non-success responses without returning internal error details.

## Raw ciphertext download semantics

`GET /v1/data/attachments/{attachmentID}/ciphertext` is a binary transport surface for E2EE clients that already possess the cryptographic context needed to decrypt the object. The server returns the stored ciphertext bytes unchanged and does not decode, decrypt, inspect, transcode, render, or content-sniff them.

The response is always labeled `application/octet-stream`, even when attachment metadata records a sender-declared plaintext MIME type. It also uses `Cache-Control: no-store`, `Pragma: no-cache`, and `X-Content-Type-Options: nosniff`. Filename and plaintext MIME metadata remain separate client metadata and are not allowed to change how the server or browser interprets encrypted bytes.

This endpoint is not an object-storage redirect, signed URL, CDN distribution contract, or evidence of production attachment storage. Those remain future deployment concerns.

## Attachment deletion semantics

Attachment deletion is designed to remove retrievable encrypted payload bytes without reopening replay state. The durable local store removes ciphertext first, then replaces user-facing attachment metadata with a minimal deletion tombstone that retains only the attachment identifier, client nonce, and deletion marker required to keep the identifier and nonce reserved.

Deleted attachments are not returned by fetch or metadata listing. Repeating `DELETE` is safe. A previously deleted attachment identifier or client nonce cannot be reused by a later submission. The tombstone is replay-prevention metadata; it is not evidence that a remote replica, backup, client cache, or production storage system has deleted corresponding data.

## Receipt semantics

Delivery receipts are GoreeCloud Data metadata, not proof of carrier delivery and not cryptographic proof that a human viewed content. `delivered` means an authorized recipient client reported delivery progress. `read` is a later recipient-observed state. Receipt state is monotonic for each message and recipient in the current service contract.

Receipts contain message, conversation, recipient, state, and observation timestamp only. They do not contain plaintext message bodies, encryption keys, device secrets, or carrier metadata.

## Privacy and security behavior

The HTTP adapter:

- accepts ciphertext rather than plaintext content;
- rejects unknown JSON fields;
- bounds request bodies, including attachment uploads;
- sets `Cache-Control: no-store` on JSON and raw ciphertext responses;
- sets `X-Content-Type-Options: nosniff`;
- keeps SMS, MMS, and RCS outside the Data API;
- authorizes receipt and attachment operations against server-side state;
- transports raw attachment ciphertext as generic binary rather than sender-declared plaintext media;
- removes attachment ciphertext before committing a privacy-minimized deletion tombstone;
- shares one required authentication boundary across the composed Data runtime; and
- does not claim that a cryptographic session exists merely because the envelope contract records `e2ee`.

## Current limitations

This remains a Development source foundation. It does not provide production credentials, production-grade distributed persistence, TLS termination, rate limiting, push delivery, cryptographic session establishment, key management, multi-device synchronization, carrier adapters, client packaging, production listener/bootstrap configuration, deployment, or production acceptance. The durable attachment file store is a single-node local implementation, not a distributed deletion or backup-erasure guarantee. Receipt storage remains an in-memory development implementation.