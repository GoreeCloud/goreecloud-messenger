# GoreeCloud Messenger

Native GoreeCloud messaging and calling with usernames, end-to-end encryption, Data messaging, SMS/RCS integration, groups, and video calls.

## Status

Active Development — native messaging, GoreeCloud Data messaging, authenticated HTTP transport, delivery/read receipt, encrypted-attachment, optional content-free typing presence, authenticated Development typing-privacy preferences, unified Data HTTP runtime composition, a Development-only opt-in GoreeCloud Identity exact-handle consumer boundary, hardened single-node receipt persistence, explicit receipt-persistence environment configuration, minimized startup-diagnostic foundations, and a disconnected native Android readiness/presentation foundation are implemented in source. Production Identity sessions/device authority, a live authenticated Identity directory transport, cryptographic session establishment, distributed delivery/storage, durable Privacy Shield-backed preference storage, live Data transport, complete client acceptance, and release acceptance remain incomplete.

The current foundation establishes the transport-provenance domain model used to keep Data, SMS, MMS, and RCS communication technically distinct. GoreeCloud Data adds encrypted-envelope validation, authenticated sender enforcement, conversation authorization, deterministic retry protection, persistence abstraction, authenticated delivery/read receipts, opaque encrypted-attachment transport, and short-lived content-free typing state. The attachment surface can submit, fetch as JSON/base64, list metadata, delete with replay-safe tombstones, and download exact ciphertext bytes without asking the server to interpret plaintext media.

The HTTP layer has one application-facing composition boundary that registers message, receipt, and attachment routes onto the same mux. Typing presence, typing-privacy preferences, and the Development Identity exact-handle resolution route are independently optional compositions under that same Messenger-user authenticator; none of those capabilities is enabled merely by constructing the base Data runtime. This preserves shared runtime composition without creating a second client-authentication boundary.

The optional Identity consumer boundary exposes authenticated `POST /v1/identity/resolve` only when an `IdentityDirectoryService` is explicitly composed. It accepts only an exact handle from the Messenger client, rejects client-supplied requester-service identity, forwards the handle to an injected resolver without Messenger-owned canonicalization, preserves one uniform unresolved response for nonexistent/private/service-unauthorized identities, and returns only opaque Identity subject, canonical handle, and optional display name on success. The repository contains no production GoreeCloud Identity HTTP client, service credential, service-principal authentication profile, or deployed directory integration. The Development semantics are pinned to GoreeCloud Identity Draft PR #5 at exact head `5904a44997b90cc44f5d196620bb7e189fea0eeb`.

The Development typing policy exposes independent per-conversation choices for publishing and observing typing state. The preference API derives the acting user from authentication, checks conversation membership, rejects request-body identity fields, and returns only the two minimized boolean choices. Its current storage implementation is memory-backed for deterministic Development validation and is not durable Privacy Shield preference persistence.

Receipt persistence can be explicitly selected as memory or hardened file-backed storage; there is no implicit durable-to-memory fallback. The development executable requires an explicit receipt-persistence environment selection. `memory` must not carry an ignored durable root. `file` requires an explicit absolute non-root persistence directory. Missing, unsupported, relative, root-level, or contradictory settings fail closed before the executable reports its development contract active.

After configuration is accepted, the executable can report a minimized categorical diagnostic containing only receipt persistence mode, implemented durability class, and configuration source. File-mode diagnostics deliberately omit the configured receipt root and any message, receipt, conversation, credential, or cryptographic content. `single-node-durable` describes the selected implementation class only; it is not a distributed-durability or production-readiness claim.

The native Android Development client remains deliberately disconnected from production communication authorities. Its readiness policy keeps GoreeCloud Identity authentication, conversation authorization, GoreeCloud Data transport, and verified active E2EE independent. Conversation authorization and cryptographic scopes must each be canonical bounded opaque identifiers and must identify the same exact conversation before a future operation can project `Data · E2EE`. Noncanonical scopes fail closed; the client does not trim or otherwise normalize authority-owned identifiers into a different conversation scope, and blocked Data readiness does not create SMS/MMS/RCS fallback authority.

## Product principles

- Every message identifies its actual transport.
- End-to-end encryption is represented only when the application has verified that state.
- GoreeCloud Data messaging works independently of cellular service.
- Username identities are first-class and are not required to map one-to-one to phone numbers.
- Username resolution belongs to GoreeCloud Identity; Messenger must not create a browsable account directory or infer private-account existence from unresolved results.
- Encrypted GoreeCloud conversations do not silently downgrade to SMS or MMS.
- RCS is integrated only where supported platform and carrier APIs legitimately allow it.
- Voice and video calling remain distinguishable from carrier calling.
- Delivery/read state is recipient-authenticated and is not presented as carrier or cryptographic proof.
- Attachment services transport opaque ciphertext and do not decrypt user content.
- Typing presence remains content-free, short-lived, participant-authorized, and independently privacy-gated for publish and observe behavior.
- Operational diagnostics must minimize sensitive configuration and communication data.
- GLAZE UI V1.3 / 1.3.0 Adaptive Resonance, Privacy Shield, Wardveil Security, Everkeep, GoreeCloud Mesh, GoreeCloud Manager, and GoreeCloud Identity are substantive platform integration requirements for applicable surfaces.

## Repository layout

- `cmd/messenger/` — development executable for exercising core contracts and validating explicit receipt-persistence process configuration
- `internal/domain/` — transport, encryption, identity, conversation, message, call, Data-envelope, receipt, attachment, typing, and minimized Identity-directory projection contracts
- `internal/service/` — GoreeCloud Data, receipt, attachment, typing, typing-privacy, persistence, authorization, and Development Identity-directory resolver boundaries
- `internal/api/` — authenticated HTTP transport plus the unified application-facing Data route-composition boundary and optional typing/privacy/Identity-resolution routes
- `internal/runtimeconfig/` — fail-closed process configuration derivation and minimized diagnostic projection for currently implemented runtime options
- `client/android/` — disconnected native Android Development client, readiness/provenance policy, GLAZE UI source mapping, tests, and rendered shell acceptance
- `docs/architecture.md` — product architecture and trust boundaries
- `docs/security.md` — encryption and security constraints
- `docs/data-messaging.md` — Data service authorization, storage, retry, and carrier-separation contract
- `docs/data-http-api.md` — HTTP API, authentication, authorization, receipt, attachment, optional Identity resolution, runtime composition, and privacy boundary
- `docs/development/identity-consumer-directory.md` — pinned GoreeCloud Identity Development contract, privacy boundary, responsibility split, and production gates for exact-handle resolution
- `docs/durable-receipt-store.md` — file durability, runtime selection, and environment-configuration boundaries
- `docs/runtime-diagnostics.md` — minimized runtime configuration diagnostic and non-disclosure boundary

## Documentation

- [USER-MANUAL.md](USER-MANUAL.md)
- [SPECIFICATIONS.md](SPECIFICATIONS.md)
- [FEATURES.md](FEATURES.md)
- [FEATURE-ROADMAP.md](FEATURE-ROADMAP.md)
- [BENEFITS.md](BENEFITS.md)
- [COMPETITIVE-OBJECTIVES.md](COMPETITIVE-OBJECTIVES.md)
- [BRANDING.md](BRANDING.md)
- [SECURITY.md](SECURITY.md)

## Planned clients

Native or platform-appropriate clients are planned for Android, tablets, desktop Linux, and other approved GoreeCloud client platforms. Client work will consume the shared transport and security contracts established here rather than redefining them independently. Consumer username resolution is expected to use the GoreeCloud Identity-owned exact-handle disclosure contract; the current repository implements only the bounded Messenger-side Development consumer seam, not a live production directory client.

## Current limitations

This repository remains Development. It does not yet establish production-grade Identity sessions/device authority, accepted service-principal authentication for the Identity consumer directory, a live Identity resolver/network transport, application-specific post-resolution invitation/membership authorization, device/key lifecycle, end-to-end cryptographic session establishment, distributed message delivery, production object storage, push notification delivery, durable Privacy Shield-backed typing preference storage, production presence fan-out/offline synchronization, distributed anti-abuse/rate-limit acceptance, carrier adapters, calling media transport, complete GLAZE UI V1.3 client acceptance, GoreeCloud Mesh/Manager production integration, or production deployment evidence.

The unified Data handler is a composition boundary, and the command-level environment parser supplies a strict receipt-persistence selection contract with a minimized categorical diagnostic. The current executable still does not assemble the complete Data runtime dependencies, credentials/Identity boundaries, TLS, Identity service-authentication client, service lifecycle, production health/readiness monitoring, migration, and deployment configuration needed for a production server. The Android client still has no production Identity session, live conversation-authorization provider, verified cryptographic session/key authority, live Data networking, message persistence/delivery, composer/Send control, carrier fallback authority, calling authority, protected release signing, or Stable qualification.

## License

GNU Affero General Public License v3.0 only (`AGPL-3.0-only`).
