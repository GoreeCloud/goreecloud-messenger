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
- Content-free privacy-controlled typing/idle signals with authenticated self-publication, conversation-membership checks, independent publish/observe policy gates, monotonic sequencing, stale-signal rejection, and a 10-second server expiry.
- Explicit optional composition of typing routes into the application-facing Data runtime under the same Authenticator boundary; typing is not enabled merely by constructing the base runtime.

## Planned / incomplete

- Production GoreeCloud Identity sessions and device/key lifecycle.
- Identity-owned exact-handle username resolution integration.
- Production E2EE session establishment, verification, rotation, and multi-device state.
- Distributed message and attachment persistence/object storage.
- Push delivery, production presence fan-out/offline synchronization, and production rate limiting.
- Durable Privacy Shield-backed typing/presence preferences and native client typing presentation.
- SMS/MMS/RCS carrier/platform adapters where legitimate APIs permit.
- Voice/video call signaling and media transport.
- Native client packaging and Glaze UI 2.0+ rendered acceptance.
- Wardveil encrypted-object/security acceptance, Privacy Shield controls, Everkeep continuity, and production deployment evidence.
