# GoreeCloud Messenger — Features

## Implemented in Development source

- GoreeCloud Data conversation/message domain contracts.
- Transport provenance separating Data, SMS, MMS, and RCS semantics.
- Authenticated sender and conversation authorization.
- Encrypted envelope submission and authorized history reads.
- Replay/duplicate and client-nonce protections.
- Recipient-authenticated delivery/read receipts with monotonic state progression.
- Opaque encrypted attachment upload and authorized fetch.
- Metadata-only attachment listing.
- Replay-safe encrypted attachment deletion.
- Exact raw ciphertext-byte download with generic binary transport and no content sniffing.

## Required built-in product capabilities — planned / incomplete unless listed above

### Messaging and conversation behavior

- Private one-to-one and group messaging.
- Internet/Data messaging independent of cellular service.
- SMS, MMS, and legitimately supported RCS/rich-carrier adapters.
- Cross-platform rich-messaging interoperability where legitimate interfaces exist.
- Verified E2EE for eligible private and group Data conversations.
- Distinct secret/private conversation modes only when additional privacy/security properties are actually implemented.
- Sent-message editing, delete for everyone, local deletion, direct replies, threaded replies, copying, forwarding, and multi-message forwarding.
- Typing indicators and configurable sent/delivered/read state.
- Scheduled and repeating scheduled messages.
- Silent delivery where supported.
- Synchronized message drafts.
- Pinned messages, pinned conversations, saved messages, and personal notes/saved-message conversations.

### Organization, search, and discovery

- Mark as unread, archive, and custom conversation organization.
- Global Messenger search and in-conversation search.
- Sender-, date-, media-, file-, link-, and document/PDF-based search.
- User mentions and supported group-wide mentions.
- Rich-text message formatting.
- Timestamps and expanded delivery details.
- Authorized conversation-history synchronization.

### Media, documents, files, and location

- Voice messages with recording controls, playback speed, and trimming.
- Photo and high-resolution photo sharing.
- Video and high-resolution video sharing.
- Motion/live-photo and GIF/animated-media sharing.
- General files, large files, documents/PDFs, and contact cards.
- Static location/place sharing through GoreeCloud Location.
- Live location with explicit audience, precision, duration, revocation, and privacy controls.
- View-once media where client/platform enforcement can support the promised behavior.

### Calling

- One-to-one and group voice calling.
- One-to-one and group video calling.
- Call invitation links.
- Supported screen sharing.

### Multi-device and clients

- Multi-device synchronization of eligible messages, conversations, drafts, media state, settings, and communication state.
- Desktop clients/access.
- Web access.
- Conversation continuity between authorized devices.
- Synchronized history subject to security, privacy, identity, and E2EE constraints.

### GoreeCloud integration

- Native first-party sharing into Messenger.
- Context-aware GoreeCloud links, files, media, locations, records, and object previews.
- Messenger entry points inside supported GoreeCloud applications.
- Cross-application conversation continuity.
- GoreeCloud Identity and contacts integration.
- Coordinated GoreeCloud notifications.
- Privacy-preserving integrated search/discovery.
- Application-to-conversation and conversation-to-application handoff.
- Shared-service use of Wardveil Security, Privacy Shield, Everkeep, Glaze UI, GoreeCloud Identity, GoreeCloud Sync, GoreeCloud Backups, GoreeCloud Drive, and GoreeCloud Location.
- Consistent authorization and permission state across integrations.
- Contextual messaging/calling/collaboration surfaces for approved first-party applications.
- Extensible first-party Messenger integration framework.

### Privacy and security

- End-to-end encryption for eligible private and group communications.
- Protected message backup/restore without unnecessary weakening of E2EE.
- Disappearing messages, expiration timers, and self-destructing media.
- Username-first communication and phone-number privacy controls.
- Contact/cryptographic verification.
- Registration/takeover protection.
- Application passcodes and supported biometric locking.
- Screenshot protection where the operating environment can enforce it.
- Blocking and abuse/account reporting.
- Spam and scam protection with evidence-backed detection behavior.
- Encrypted local storage.
- Metadata minimization and privacy protection.
- Privacy-preserving routing where technically supported.
- User-controlled encryption-key options for supported configurations.

## Current major implementation gaps

- Production GoreeCloud Identity sessions and device/key lifecycle.
- Identity-owned exact-handle username resolution integration.
- Production E2EE session establishment, verification, rotation, groups, and multi-device state.
- Distributed message and attachment persistence/object storage.
- Push delivery, presence/typing policy, offline synchronization, and production rate limiting.
- SMS/MMS/RCS carrier/platform adapters where legitimate APIs permit.
- Voice/video call signaling and media transport.
- Complete native, desktop, and web client surfaces.
- Glaze UI 2.0+ rendered acceptance.
- Wardveil Security, Privacy Shield, Everkeep, Sync, Backups, Drive, Location, and Identity production integration acceptance.
- Deployment and Stable acceptance evidence.

The complete product inventory is maintained in `GoreeCloud/Projects/Project Specification — Messenger`. This file must not represent planned scope as shipped functionality.