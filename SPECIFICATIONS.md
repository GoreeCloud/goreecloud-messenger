# GoreeCloud Messenger — Repository Specifications

Status: Development  
Canonical project record: `GoreeCloud/Projects/Project Specification — Messenger`  
Repository: `GoreeCloud/goreecloud-messenger`

## Product boundary

GoreeCloud Messenger is the native GoreeCloud messaging and calling application/service. It owns GoreeCloud Data messaging, conversation/message/receipt/attachment contracts, and approved client experiences. SMS, MMS, RCS, and carrier calling remain technically distinct transports and must never be mislabeled as GoreeCloud Data E2EE.

## Current implemented foundation

- Native Go domain/service/API layers for Data conversations and encrypted message envelopes.
- Authenticated sender and conversation-membership enforcement.
- Deterministic duplicate/idempotency and nonce-reuse protections.
- Authenticated delivery/read receipt contracts with monotonic state; the current stacked Development source adds conversation-scoped read publication/observation privacy-policy gates while keeping delivered state operational.
- Opaque encrypted attachment submission, authorized JSON/base64 fetch, metadata listing, replay-safe deletion, and raw ciphertext download.
- Immutable encrypted sent-message revisions and delete-for-everyone synchronization tombstones on the current stacked Development source.
- Direct replies with same-conversation target validation and threaded replies with stable root/parent relationships plus authorized thread-history reads on the current stacked Development source.
- Encrypted message-reaction set/clear events with authenticated per-reactor current-state projection, replay protection, and stale-event rejection on the current stacked Development source.
- Privacy-controlled ephemeral typing indicators with authenticated participant binding, monotonic per-user sequence ordering, server-assigned expiry, publish/observe policy gates, and content-free active projections on the current stacked Development source.
- Local Development persistence abstractions and focused tests.

## Identity and discovery

Messenger must use GoreeCloud Identity for account/session authority and consumer username resolution. Username discovery must not create a Messenger-owned browsable account directory. The preferred contract is exact-handle resolution with explicit Identity-owned discoverability and per-service disclosure policy.

## Security and privacy requirements

- The service must not decrypt GoreeCloud Data message, reaction, or attachment ciphertext.
- Encryption state must be represented only when verified by the applicable client/session protocol.
- Encrypted conversations must not silently downgrade to SMS/MMS.
- Reply and thread metadata must not become a server-side plaintext quote, preview, or conversation-summary channel.
- Reaction values must remain opaque ciphertext at the service boundary; server-visible set/clear state is control metadata and must not be presented as plaintext reaction content.
- Typing indicators must remain content-free ephemeral presence metadata. Draft text, keystrokes, cursor positions, message ciphertext, and device secrets must not enter the typing-indicator service.
- Typing publication and observation must remain subject to explicit privacy-policy checks in addition to conversation membership.
- Read-receipt publication and observation must remain subject to explicit privacy-policy checks. Current policy changes must be able to hide stored read projections without changing the underlying delivery acknowledgement contract.
- Delivered state remains operational GoreeCloud Data metadata and must not be mislabeled as read, human-view proof, carrier delivery proof, or cryptographic proof.
- Thread, reply, and reaction target errors must remain bounded so cross-conversation message identifiers are not exposed through existence-oracle behavior.
- Attachment raw-byte transport must remain generic binary with no content sniffing and no server-side plaintext MIME interpretation.
- Wardveil Security, Privacy Shield, Everkeep, GoreeCloud Mesh, and GoreeCloud Identity integration are required where applicable.
- Client surfaces must target Glaze UI 2.0 or newer and pass rendered acceptance before Stable promotion.

## Current acceptance boundary

This is not production-ready. Production-grade identity/device keys, cryptographic session establishment, multi-device synchronization, distributed message/reaction/object persistence, push delivery, production presence fan-out and durable typing/read-receipt preference persistence, thread and reaction synchronization/indexing at production scale, abuse controls, carrier adapters, calling media infrastructure, client packaging, and deployment acceptance remain incomplete.
