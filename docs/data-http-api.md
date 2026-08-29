# GoreeCloud Messenger Data HTTP API Foundation

## Purpose

This document defines the HTTP transport boundary for GoreeCloud Data messaging. The HTTP layer is an adapter around the Data, receipt, encrypted-reaction, typing-presence, and encrypted-attachment services and does not replace their authorization, transport-provenance, privacy-policy, or E2EE-only validation rules.

## Current endpoints

- `POST /v1/data/messages` accepts an encrypted GoreeCloud Data envelope, including optional direct-reply and threaded-reply metadata.
- `GET /v1/data/conversations/{conversationID}/messages` returns encrypted envelopes for one authorized conversation.
- `GET /v1/data/conversations/{conversationID}/threads/{rootMessageID}/messages` returns one authorized thread root plus its encrypted threaded replies.
- `POST /v1/data/messages/{messageID}/reactions` records an authenticated encrypted reaction set/clear event.
- `GET /v1/data/messages/{messageID}/reactions` returns the current active encrypted reaction projection to an authorized conversation participant.
- `POST /v1/data/conversations/{conversationID}/typing` publishes content-free authenticated `typing` or `idle` presence state.
- `GET /v1/data/conversations/{conversationID}/typing` returns currently active typing participants to an authorized observer whose privacy policy permits observation.
- `POST /v1/data/messages/{messageID}/receipts` records an authenticated recipient delivery or privacy-permitted read acknowledgement.
- `GET /v1/data/messages/{messageID}/receipts` returns authorized delivery state plus only read projections permitted by the current receipt privacy policy.
- `POST /v1/data/attachments` accepts already-encrypted opaque attachment ciphertext plus bounded metadata.
- `GET /v1/data/attachments/{attachmentID}` returns opaque ciphertext as base64 JSON to an authorized conversation participant.
- `GET /v1/data/attachments/{attachmentID}/ciphertext` returns the exact stored ciphertext bytes to an authorized conversation participant as `application/octet-stream`.
- `GET /v1/data/conversations/{conversationID}/attachments` returns a bounded metadata-only attachment projection.
- `DELETE /v1/data/attachments/{attachmentID}` removes retrievable ciphertext for an authorized participant and is idempotent for already-missing/deleted attachment identifiers.

Message, reaction, and ordinary attachment JSON ciphertext is transported as standard base64 text. Typing endpoints do not accept message or reaction ciphertext at all. Receipt endpoints carry delivery/read metadata only. The raw ciphertext download endpoint transports encrypted bytes directly. The API does not accept plaintext message-body, plaintext reaction-value, plaintext typing-draft, or plaintext attachment-content fields.

## Authentication boundary

The HTTP package requires an `Authenticator` implementation. The authenticator must resolve the request to an authenticated GoreeCloud user identifier before a service is called.

Credential issuance, login, token creation, token storage, session management, device identity, and identity-provider integration remain outside this milestone. The API must not infer identity from client-supplied sender, reactor, typing user, attachment, receipt, reply, or thread metadata.

## Authorization boundary

The authenticated user is passed to service-layer authorization. The implementation:

- binds the authenticated user to message and attachment senders, reaction reactors, typing-state publishers, and receipt recipients;
- verifies server-side conversation membership;
- validates direct-reply targets against same-conversation message state;
- validates thread roots and immediate thread parents against same-conversation message state;
- authorizes thread-history reads before returning root or reply metadata;
- validates reaction targets against same-conversation message state and returns one bounded target error for missing/cross-conversation reaction targets;
- rejects stale reaction events and reaction ID/client-nonce replay while maintaining one current active projection per reactor/message;
- applies an explicit typing privacy-policy boundary for both publishing and observing typing state;
- rejects stale or equal typing sequences so delayed state cannot overwrite a newer per-conversation/per-user typing decision;
- rejects delivery/read acknowledgements from the message sender for that sender's own message;
- verifies that a receipt references an existing message in the stated conversation;
- prevents receipt state from moving backwards from `read` to `delivered`;
- applies explicit read-receipt privacy policy to both read publication and read observation while leaving `delivered` state outside that privacy gate;
- re-evaluates the receipt owner's current read-publication policy when producing read projections so a later privacy change can hide stored read state;
- rejects duplicate message and attachment identifiers;
- rejects message and attachment client-nonce reuse;
- authorizes attachment JSON fetch, raw ciphertext fetch, list, and delete operations against conversation membership; and
- enforces the E2EE-only Data-envelope and encrypted-reaction contracts.

The HTTP layer maps authorization, privacy-policy, and relationship-validation failures to bounded non-success responses without returning internal error details.

## Reply and thread semantics

`reply_to_message_id` identifies an existing message that the new encrypted message directly references. When `thread_root_message_id` is absent, the relationship is a direct reply.

When `thread_root_message_id` is present, the message is part of a focused thread. It must also carry `reply_to_message_id`. The declared root must be a stable same-conversation root rather than an existing threaded child, and the immediate parent must be either that root or a message already assigned to that root.

`GET /v1/data/conversations/{conversationID}/threads/{rootMessageID}/messages` returns the root followed by messages assigned to that thread in deterministic Development-store order. Missing, cross-conversation, invalid-root, and unrelated-parent thread identifiers use the same bounded `thread target unavailable` state. Direct replies use the corresponding bounded `reply target unavailable` state.

The server never decrypts original messages, replies, thread roots, quoted text, or previews. Any plaintext quote, preview, summary, or thread title is a client concern and must be derived only from content the authorized client can already decrypt.

## Reaction semantics

A reaction is represented as an immutable `set` or `clear` state event for one message and one authenticated reactor. A `set` event carries opaque E2EE reaction ciphertext. The service does not decode the encrypted value into an emoji, text label, sticker, or other plaintext reaction. A `clear` event carries no reaction ciphertext and removes that reactor's active projection.

Each reaction event has a unique event ID and per-reactor client nonce. The in-memory Development store retains the latest event per message/reactor and rejects events whose timestamp does not advance that reactor's current state, preventing a delayed stale event from overwriting a newer reaction. The active-reaction GET response returns only current `set` projections and deliberately omits replay nonces and the set/clear operation field.

Missing and cross-conversation reaction targets share the bounded `reaction target unavailable` state during submission. Reaction value interpretation and rendering remain client responsibilities after authorized decryption.

## Typing indicator semantics

Typing indicators are ephemeral, content-free conversation-presence state. `POST /v1/data/conversations/{conversationID}/typing` accepts only `user_id`, a monotonically increasing `sequence`, and `state` (`typing` or `idle`). The authenticated user must match the published typing user, must still be a conversation participant, and must pass the configured publish privacy policy.

A `typing` update creates a server-expiring active projection for 10 seconds. An `idle` update clears that active projection immediately. The Development store retains only the latest sequence needed to reject stale updates and the currently active projection. Equal or older sequences are rejected, including delayed `typing` signals that arrive after a newer `idle` signal.

`GET /v1/data/conversations/{conversationID}/typing` requires conversation membership plus observe permission from the typing privacy policy. It returns only other active participants, sequence numbers, and server-assigned expiry times. It does not return the observer's own state, a client timestamp, draft content, cursor location, keystrokes, message ciphertext, device secret, or a free-form presence payload.

## Raw ciphertext download semantics

`GET /v1/data/attachments/{attachmentID}/ciphertext` is a binary transport surface for E2EE clients that already possess the cryptographic context needed to decrypt the object. The server returns the stored ciphertext bytes unchanged and does not decode, decrypt, inspect, transcode, render, or content-sniff them.

The response is always labeled `application/octet-stream`, even when attachment metadata records a sender-declared plaintext MIME type. It also uses `Cache-Control: no-store`, `Pragma: no-cache`, and `X-Content-Type-Options: nosniff`. Filename and plaintext MIME metadata remain separate client metadata and are not allowed to change how the server or browser interprets encrypted bytes.

This endpoint is not an object-storage redirect, signed URL, CDN distribution contract, or evidence of production attachment storage. Those remain future deployment concerns.

## Attachment deletion semantics

Attachment deletion is designed to remove retrievable encrypted payload bytes without reopening replay state. The durable local store removes ciphertext first, then replaces user-facing attachment metadata with a minimal deletion tombstone that retains only the attachment identifier, client nonce, and deletion marker required to keep the identifier and nonce reserved.

Deleted attachments are not returned by fetch or metadata listing. Repeating `DELETE` is safe. A previously deleted attachment identifier or client nonce cannot be reused by a later submission. The tombstone is replay-prevention metadata; it is not evidence that a remote replica, backup, client cache, or production storage system has deleted corresponding data.

## Receipt semantics

Delivery receipts are GoreeCloud Data metadata, not proof of carrier delivery and not cryptographic proof that a human viewed content. `delivered` means an authorized recipient client reported delivery progress. `read` is a later recipient-observed state. Receipt state remains monotonic for each message and recipient in the current service contract.

Read state now has an explicit privacy boundary. A recipient may publish `read` only when `ReceiptPrivacyPolicy.CanPublishRead` permits that conversation/user pair. Authorized receipt reads include a stored `read` projection only when the observer's `CanObserveRead` policy permits read state and the receipt owner's current `CanPublishRead` policy still permits sharing it. A later privacy change can therefore hide an already-stored read projection without rewriting stored receipt state.

`delivered` is deliberately not blocked by read-receipt privacy settings. When a stored read projection is privacy-hidden, the Development service omits that read projection rather than fabricating or backdating a delivery acknowledgement from the later read timestamp. Because the current in-memory store keeps one latest receipt state per recipient/message, an earlier delivered projection may no longer be available after it was replaced by read; production receipt persistence may preserve richer delivery history while keeping the same privacy contract.

Receipts contain message, conversation, recipient, state, and observation timestamp only. They do not contain plaintext message bodies, encryption keys, device secrets, or carrier metadata.

## Privacy and security behavior

The HTTP adapter:

- accepts ciphertext rather than plaintext message or reaction content;
- keeps typing state content-free and separate from drafts or message bodies;
- keeps read-receipt state subject to explicit publish/observe privacy policy while preserving delivery-state semantics;
- rejects unknown JSON fields;
- bounds request bodies, including attachment uploads;
- sets `Cache-Control: no-store` on JSON and raw ciphertext responses;
- sets `X-Content-Type-Options: nosniff`;
- keeps SMS, MMS, and RCS outside the Data API;
- authorizes reply, thread, reaction, typing, receipt, and attachment operations against server-side state;
- keeps reply/thread metadata limited to message relationships rather than plaintext previews or summaries;
- keeps active reaction values as ciphertext and omits replay/control metadata not needed by the current projection;
- keeps active typing projections time-bounded, privacy-gated, and free of draft/keypress data;
- filters privacy-disallowed read projections before HTTP response construction;
- transports raw attachment ciphertext as generic binary rather than sender-declared plaintext media;
- removes attachment ciphertext before committing a privacy-minimized deletion tombstone; and
- does not claim that a cryptographic session exists merely because an envelope or reaction contract records `e2ee`.

## Current limitations

This remains a Development source foundation. It does not provide production credentials, production-grade distributed persistence, TLS termination, rate limiting, push delivery, production presence fan-out, durable Privacy Shield preference storage for typing/read-receipt controls, cryptographic session establishment, key management, production thread or reaction indexing, offline or multi-device thread/reaction/typing/receipt-policy synchronization, cross-device typing sequence allocation, carrier adapters, client reaction encrypt/decrypt/render behavior, client typing/read-receipt settings and rendering, client packaging, deployment, or production acceptance. The durable attachment file store is a single-node local implementation, not a distributed deletion or backup-erasure guarantee. Receipt and reaction storage remain in-memory Development implementations; typing state is intentionally ephemeral in-memory Development state.
