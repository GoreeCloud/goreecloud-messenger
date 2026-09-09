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
- Native Android Development client foundations with explicit communication-provenance/readiness boundaries and a repository-local GLAZE UI V1.3 Stable source mapping.
- Android Data-send readiness keeps GoreeCloud Identity authentication, conversation authorization, Data transport availability, and verified active E2EE as independent authorities. Conversation authorization and E2EE scopes must each be canonical bounded opaque identifiers and must identify the same exact conversation before `Data · E2EE` readiness can be projected.

## Identity and discovery

Messenger must use GoreeCloud Identity for account/session authority and consumer username resolution. Username discovery must not create a Messenger-owned browsable account directory. The preferred contract is exact-handle resolution with explicit Identity-owned discoverability and per-service disclosure policy.

## Security and privacy requirements

- The service must not decrypt GoreeCloud Data message or attachment ciphertext.
- Encryption state must be represented only when verified by the applicable client/session protocol.
- Encrypted conversations must not silently downgrade to SMS/MMS.
- Attachment raw-byte transport must remain generic binary with no content sniffing and no server-side plaintext MIME interpretation.
- Wardveil Security, Privacy Shield, Everkeep, GoreeCloud Mesh, and GoreeCloud Identity integration are required where applicable.
- Client surfaces must track the **current approved Stable GLAZE UI release**. The current target is GLAZE UI V1.3 / `1.3.0` — Adaptive Resonance, at exact Stable integration revision `fc7cc91d2eace8da2371371c2855c24cbcb326a1`, with `1.2.0` retained as the rollback baseline.
- Client presentation must not manufacture security, privacy, encryption, delivery, authorization, or transport state through decorative material or semantic color. Visible state must remain derived from the responsible technical authority.
- Repository-local token/source tests are not rendered application acceptance. Exact-revision visual, accessibility, adaptive/form-factor, localization/directionality, representative-device, performance, rollback, and other applicable downstream evidence remain required before Stable promotion.
- Authority-owned conversation identifiers are opaque and exact. Leading/trailing whitespace, overlong values, C0 controls, and DEL are noncanonical and fail closed; valid internal whitespace and punctuation remain part of the identifier and are not normalized away.

## Current acceptance boundary

This is not production-ready. Production-grade identity/device keys, cryptographic session establishment, multi-device synchronization, distributed message/object persistence, push delivery, abuse controls, carrier adapters, calling media infrastructure, client packaging, complete GLAZE UI V1.3 downstream acceptance, and deployment acceptance remain incomplete.
