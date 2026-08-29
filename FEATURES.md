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
- Immutable encrypted sent-message revision records on the current stacked Development branch.
- Delete-for-everyone synchronization tombstones on the current stacked Development branch.
- Direct reply references that preserve opaque message ciphertext and verify same-conversation targets on the current stacked Development branch.
- Threaded reply metadata and authorized thread-history reads with stable same-conversation roots on the current stacked Development branch.

## Planned / incomplete

- Production GoreeCloud Identity sessions and device/key lifecycle.
- Identity-owned exact-handle username resolution integration.
- Production E2EE session establishment, verification, rotation, and multi-device state.
- Distributed message and attachment persistence/object storage.
- Push delivery, presence/typing policy, offline synchronization, and production rate limiting.
- SMS/MMS/RCS carrier/platform adapters where legitimate APIs permit.
- Voice/video call signaling and media transport.
- Native client packaging and Glaze UI 2.0+ rendered acceptance.
- Wardveil encrypted-object/security acceptance, Privacy Shield controls, Everkeep continuity, and production deployment evidence.
