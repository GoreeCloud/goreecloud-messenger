# GoreeCloud Messenger — Repository Specifications

Status: Development  
Canonical project record: `GoreeCloud/Projects/Project Specification — Messenger`  
Repository: `GoreeCloud/goreecloud-messenger`

## Product boundary

GoreeCloud Messenger is the native GoreeCloud messaging and calling application/service. It owns GoreeCloud Data messaging, conversation/message/receipt/attachment contracts, and approved client experiences. SMS, MMS, RCS, and carrier calling remain technically distinct transports and must never be mislabeled as GoreeCloud Data E2EE.

Messenger is intended to be GoreeCloud's primary communications layer and a complete standalone messaging/calling application. Other GoreeCloud applications may consume approved Messenger capabilities through controlled first-party integration interfaces rather than duplicating Messenger's communication domain.

## Required built-in product scope

The canonical Drive project specification contains the complete feature inventory. Repository implementation must converge on that inventory while keeping each capability's implementation state explicit.

Required capability families include:

- private one-to-one, group, Internet/Data, SMS, MMS, and legitimately available RCS/rich-carrier messaging;
- cross-platform rich messaging interoperability only through supported standards, platform APIs, carrier APIs, or other legitimate integration paths;
- verified E2EE for eligible private and group GoreeCloud Data conversations, plus distinct secret/private conversation modes only when their additional security properties are actually implemented;
- message editing, delete-for-everyone, local deletion, direct replies, threaded replies, typing state, sent/delivered/read state, scheduling, recurring scheduling, silent delivery, synchronized drafts, forwarding, multi-message forwarding, copying, pinned messages, pinned conversations, saved messages, and private notes/saved-message conversations;
- unread state, archive, custom organization, global and in-conversation search, sender/date/media/file/link/document search, mentions, group mentions, rich-text formatting, timestamps, delivery details, and authorized history synchronization;
- voice messages with recording, playback-speed, and trimming controls; photos, high-resolution photos, videos, high-resolution videos, motion/live photos, GIFs, general files, large files, documents/PDFs, contact cards, static location, live location, and view-once media;
- one-to-one and group voice/video calling, invitation links, and supported screen sharing;
- multi-device synchronization, desktop access, web access, conversation continuity, and enrolled-device history synchronization;
- native GoreeCloud sharing, contextual first-party content previews, application-to-conversation and conversation-to-application handoff, integrated identity/contacts, notifications, search/discovery, permissions, contextual communication surfaces, and an extensible first-party Messenger integration framework;
- disappearing messages, expiration timers, self-destructing media, username and phone-number privacy controls, contact verification, registration protection, application passcodes, biometric locking, screenshot protection where the operating environment can enforce it, blocking, abuse reporting, spam/scam protection, encrypted local storage, metadata minimization, privacy-preserving routing where technically supported, and user-controlled key options for supported security configurations.

## First-party GoreeCloud integrations

Messenger must use the appropriate GoreeCloud platform systems rather than recreate isolated equivalents:

- **Wardveil Security:** protection, trust, detection, verification, response, attachment/abuse safety, and evidence-backed security state while respecting E2EE boundaries.
- **Privacy Shield:** consent, data minimization, metadata governance, retention/sharing controls, permission transparency, and user control.
- **Everkeep:** continuity, recovery, preservation, portability, succession, and resilience without weakening E2EE.
- **Glaze UI:** current Stable design, accessibility, responsiveness, interaction, and component contracts.
- **GoreeCloud Identity:** accounts, usernames, authentication, authorization, devices, credentials, sessions, verification, and delegated authority.
- **GoreeCloud Sync:** authorized encrypted synchronization of conversations, message state, drafts, settings, media state, and other eligible Messenger state.
- **GoreeCloud Backups:** user-controlled protected backup/restore for eligible Messenger state without silently weakening confidentiality or key protection.
- **GoreeCloud Drive:** explicit attachment selection, saving, export, sharing, large-file/document handoff, and first-party file integration without bypassing conversation authorization.
- **GoreeCloud Location:** consent-driven places, static location, maps, and live-location sharing with audience, precision, expiration, revocation, and privacy controls.

## Current implemented foundation

- Native Go domain/service/API layers for Data conversations and encrypted message envelopes.
- Authenticated sender and conversation-membership enforcement.
- Deterministic duplicate/idempotency and nonce-reuse protections.
- Authenticated delivery/read receipt contracts with monotonic state.
- Opaque encrypted attachment submission, authorized JSON/base64 fetch, metadata listing, replay-safe deletion, and raw ciphertext download.
- Local Development persistence abstractions and focused tests.

## Identity and discovery

Messenger must use GoreeCloud Identity for account/session authority and consumer username resolution. Username discovery must not create a Messenger-owned browsable account directory. The preferred contract is exact-handle resolution with explicit Identity-owned discoverability and per-service disclosure policy.

## Security and privacy requirements

- The service must not decrypt GoreeCloud Data message or attachment ciphertext merely to provide transport or storage.
- Encryption state must be represented only when verified by the applicable client/session protocol.
- Encrypted conversations must not silently downgrade to SMS/MMS.
- Attachment raw-byte transport must remain generic binary with no content sniffing and no server-side plaintext MIME interpretation.
- Wardveil Security, Privacy Shield, Everkeep, GoreeCloud Mesh, and GoreeCloud Identity integration are required where applicable.
- Client surfaces must target Glaze UI 2.0 or newer and pass rendered acceptance before Stable promotion.
- Screenshot protection, view-once behavior, private routing, carrier interoperability, and similar platform-dependent properties may be claimed only to the extent the actual client/platform/runtime can enforce them.

## Current acceptance boundary

The built-in inventory is product scope, not a statement that every item is implemented. This repository remains Development. Production-grade identity/device keys, cryptographic session establishment, multi-device synchronization, distributed message/object persistence, push delivery, abuse controls, carrier adapters, calling media infrastructure, complete native/web/desktop clients, full GoreeCloud integration acceptance, and deployment acceptance remain incomplete.