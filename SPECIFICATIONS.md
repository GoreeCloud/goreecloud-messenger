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
- Authenticated delivery/read receipt contracts with monotonic state.
- Opaque encrypted attachment submission, authorized JSON/base64 fetch, metadata listing, replay-safe deletion, and raw ciphertext download.
- Local Development persistence abstractions and focused tests.
- Native Android Development client foundations with explicit communication-provenance/readiness boundaries and a repository-local GLAZE UI V1.2 Stable source mapping.

## Identity and discovery

Messenger must use GoreeCloud Identity for account/session authority and consumer username resolution. Username discovery must not create a Messenger-owned browsable account directory. The preferred contract is exact-handle resolution with explicit Identity-owned discoverability and per-service disclosure policy.

## Security and privacy requirements

- The service must not decrypt GoreeCloud Data message or attachment ciphertext.
- Encryption state must be represented only when verified by the applicable client/session protocol.
- Encrypted conversations must not silently downgrade to SMS/MMS.
- Attachment raw-byte transport must remain generic binary with no content sniffing and no server-side plaintext MIME interpretation.
- Wardveil Security, Privacy Shield, Everkeep, GoreeCloud Mesh, and GoreeCloud Identity integration are required where applicable.
- Client surfaces must track the **current approved Stable GLAZE UI release**. The current target is V1.2 / `1.2.0`, Stable promotion merge revision `f285b9145e27e6e7027b075c37299d101945c272`, and V1.2 source-qualification anchor `b0eadf9a60f73d45caffb62ffc7e9e0334cddc97`.
- The V1.2 material rule is **Neutral glass is the material. Color is an accent.** Client base material must not use security, privacy, delivery, encryption, or other semantic colors as substrate authority.
- Repository-local token/source tests are not rendered application acceptance. Exact-revision visual, accessibility, adaptive/form-factor, localization/directionality, representative-device, performance, and other applicable downstream evidence remain required before Stable promotion.

## Current acceptance boundary

This is not production-ready. Production-grade identity/device keys, cryptographic session establishment, multi-device synchronization, distributed message/object persistence, push delivery, abuse controls, carrier adapters, calling media infrastructure, client packaging, complete GLAZE UI V1.2 downstream acceptance, and deployment acceptance remain incomplete.
