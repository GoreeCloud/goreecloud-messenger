# GoreeCloud Messenger Data HTTP API Foundation

## Purpose

This document defines the HTTP transport boundary for GoreeCloud Data messaging and the optional application-facing Development adapters composed beside it. The HTTP layer is an adapter around the Data, receipt, encrypted-attachment, typing/privacy, and optional Identity consumer services and does not replace their authorization, transport-provenance, privacy, or E2EE-only validation rules.

## Application-facing runtime composition

`DataRuntimeHandler` is the current application-facing composition boundary for the implemented Data HTTP surface. It requires the Data service, receipt service, attachment service, and one `Authenticator`, then registers all implemented base routes directly on a shared `http.ServeMux`.

Direct route registration is intentional. Message listing and attachment listing both occupy the `/v1/data/conversations/{conversationID}/...` namespace; prefix-mounting separate child handlers could hide one surface behind another. The composed runtime keeps those route families reachable without duplicating authorization logic or inventing a second authentication boundary.

Typing presence, typing privacy preferences, and the Messenger-facing Identity exact-handle resolution endpoint are independently optional compositions. The Identity route is absent unless an `IdentityDirectoryService` is deliberately installed with `WithIdentityDirectory(...)`.

The runtime composition does not create credentials, establish cryptographic sessions, authenticate a Messenger service principal to GoreeCloud Identity, choose production persistence, terminate TLS, configure production rate limits, or start a production listener. Those remain explicit outer application/deployment responsibilities.

## Current endpoints

Base Data routes:

- `POST /v1/data/messages` accepts an encrypted GoreeCloud Data envelope.
- `GET /v1/data/conversations/{conversationID}/messages` returns encrypted envelopes for one authorized conversation.
- `POST /v1/data/messages/{messageID}/receipts` records an authenticated recipient delivery or read acknowledgement.
- `GET /v1/data/messages/{messageID}/receipts` returns receipt state to an authorized conversation participant.
- `POST /v1/data/attachments` accepts already-encrypted opaque attachment ciphertext plus bounded metadata.
- `GET /v1/data/attachments/{attachmentID}` returns opaque ciphertext as base64 JSON to an authorized conversation participant.
- `GET /v1/data/attachments/{attachmentID}/ciphertext` returns the exact stored ciphertext bytes to an authorized conversation participant as `application/octet-stream`.
- `GET /v1/data/conversations/{conversationID}/attachments` returns a bounded metadata-only attachment projection.
- `DELETE /v1/data/attachments/{attachmentID}` removes retrievable ciphertext for an authorized participant and is idempotent for already-missing/deleted attachment identifiers.

Optional Development Identity consumer route:

- `POST /v1/identity/resolve` accepts one exact handle from an authenticated Messenger user and, when an `IdentityDirectoryService` is explicitly composed, delegates to the injected GoreeCloud Identity resolver boundary. Success returns only opaque Identity subject, canonical handle, and optional display name. Unresolved privacy-sensitive outcomes return one uniform 404.

Message and ordinary attachment JSON ciphertext is transported as standard base64 text. The raw ciphertext download endpoint transports encrypted bytes directly. The API does not accept plaintext message-body or plaintext attachment-content fields.

## Authentication boundary

The HTTP package requires an `Authenticator` implementation. The authenticator must resolve the request to an authenticated GoreeCloud user identifier before a protected service is called. The composed runtime deliberately reuses the same injected Messenger-user authenticator for message, receipt, attachment, typing/privacy, and optional Identity-resolution routes.

Credential issuance, login, token creation, token storage, session management, device identity, and production identity-provider integration remain outside this milestone. The API must not infer identity from client-supplied sender, attachment, receipt, or service-principal metadata.

The optional Identity consumer route has two intentionally separate authentication layers:

1. the Messenger-facing route requires the existing authenticated Messenger user; and
2. the future `IdentityDirectoryResolver` implementation must independently authenticate Messenger as a service to GoreeCloud Identity using an approved service-authentication profile.

The requester service is never accepted from the Messenger client. This repository currently implements no live service-to-service Identity client, credential, endpoint configuration, OAuth/OIDC service flow, or production service principal.

## GoreeCloud Identity consumer-directory boundary

The Development consumer boundary is pinned to GoreeCloud Identity Draft PR #5 (`agent/native-directory-contract`) at exact head `5904a44997b90cc44f5d196620bb7e189fea0eeb`, contract `goreecloud-identity.consumer-directory.v1`. That upstream contract remains Draft/unmerged and is not production authority.

The Messenger request schema accepts only:

```json
{"handle":"@example"}
```

Unknown JSON fields are rejected and the request body is limited to 4096 bytes. Messenger forwards the supplied handle unchanged to the injected resolver; GoreeCloud Identity owns leading-`@` handling, canonical lowercase normalization, discoverability, and requesting-service disclosure policy.

A successful provider projection is validated before it reaches the client and may contain only:

- opaque Identity subject;
- canonical lowercase handle; and
- optional display name.

Messenger does not expose email addresses, phone numbers, discoverability flags, allowed-service policy, service credentials, Identity roles, or administrative state through this route.

Nonexistent, private, and service-disclosure-unauthorized accounts intentionally map to the same response:

```json
{"error":{"code":"not_resolved","message":"The requested identity could not be resolved."}}
```

Unexpected resolver failures and invalid provider projections fail closed to a generic `503 identity directory unavailable` response. Backend error details are not returned.

This is exact-handle resolution, not search. Prefix search, fuzzy search, directory browsing, account enumeration, phone/email lookup, and administrative listing are not implemented.

A successful Identity resolution also does not grant a conversation invitation, membership, contact relationship, message authorization, or other Messenger-specific permission. Those remain application-owned authorization decisions.

## Authorization boundary

The authenticated user is passed to service-layer authorization for communication state. The implementation:

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

The Identity consumer adapter authenticates the Messenger user before resolution but does not convert a resolved identity into communication authorization. Application-specific invitation/membership rules remain separate.

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

- accepts ciphertext rather than plaintext communication content;
- rejects unknown JSON fields;
- bounds request bodies, including attachment uploads and Identity resolution;
- sets `Cache-Control: no-store` on JSON and raw ciphertext responses;
- sets `X-Content-Type-Options: nosniff`;
- keeps SMS, MMS, and RCS outside the Data API;
- authorizes receipt and attachment operations against server-side state;
- transports raw attachment ciphertext as generic binary rather than sender-declared plaintext media;
- removes attachment ciphertext before committing a privacy-minimized deletion tombstone;
- shares one required Messenger-user authentication boundary across the composed runtime;
- keeps Identity service authentication outside client input and outside the current Development adapter;
- preserves uniform privacy-sensitive Identity resolution negatives and minimized successful projection;
- keeps exact-handle Identity resolution disabled unless explicitly composed; and
- does not claim that a cryptographic session, production Identity integration, or directory deployment exists merely because the corresponding contract boundary is present.

## Current limitations

This remains a Development source foundation. It does not provide production credentials, a live GoreeCloud Identity consumer-directory client, accepted service-principal authentication, production account-discovery abuse/rate-limit controls, production-grade distributed persistence, TLS termination, push delivery, cryptographic session establishment, key management, multi-device synchronization, carrier adapters, client packaging, production listener/bootstrap configuration, deployment, or production acceptance. The durable attachment file store is a single-node local implementation, not a distributed deletion or backup-erasure guarantee. Receipt storage remains an in-memory development implementation.

See `docs/development/identity-consumer-directory.md` for the pinned upstream Development contract, responsibility split, privacy boundary, and production gates for exact-handle resolution.
