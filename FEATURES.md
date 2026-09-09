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
- Authenticated per-conversation typing privacy preferences for the current Development memory policy, allowing a participant to independently disable publishing or observing typing presence.
- Strict typing-preference HTTP input that derives user identity from authentication and rejects request-body identity fields.
- Explicit optional runtime composition of the mutable typing-preference route, separate from the base message runtime and separate from typing-signal composition.
- Native Android Development readiness policy and presentation that report GoreeCloud Identity authentication, exact conversation authorization, GoreeCloud Data transport, and exact conversation-scoped verified active E2EE independently before any future `Data · E2EE` send operation can be considered ready.
- Exact bounded opaque conversation-scope validation for Android Data readiness: noncanonical leading/trailing whitespace, overlong values, C0 controls, and DEL fail closed rather than being normalized into another authority scope; valid internal spaces and punctuation are preserved exactly.
- GLAZE UI V1.3 / `1.3.0` Adaptive Resonance source mapping for the current native Android Development client, with rendered/application acceptance remaining separate.

## Planned / incomplete

- Production GoreeCloud Identity sessions and device/key lifecycle.
- Identity-owned exact-handle username resolution integration.
- Production E2EE session establishment, verification, rotation, and multi-device state.
- Distributed message and attachment persistence/object storage.
- Push delivery, production presence fan-out/offline synchronization, and production rate limiting.
- **Durable Privacy Shield-backed** typing/presence preference persistence and native client typing/privacy presentation. The current preference store is Development memory state only.
- SMS/MMS/RCS carrier/platform adapters where legitimate APIs permit.
- Voice/video call signaling and media transport.
- Native client packaging and complete GLAZE UI V1.3 rendered, accessibility, adaptive/form-factor, representative-device, performance, rollback, release, and production acceptance.
- Wardveil encrypted-object/security acceptance, Privacy Shield controls, Everkeep continuity, GoreeCloud Mesh and GoreeCloud Manager production integration, and production deployment evidence.
