# GoreeCloud Messenger

Native GoreeCloud messaging and calling with usernames, end-to-end encryption, Data messaging, SMS/RCS integration, groups, and video calls.

## Status

Active Development — native messaging, GoreeCloud Data messaging, authenticated HTTP transport, delivery/read receipt, encrypted-attachment, and unified Data HTTP runtime-composition foundations are implemented in source; production identity, cryptographic session establishment, distributed delivery/storage, and client acceptance remain incomplete.

The current foundation establishes the transport-provenance domain model used to keep Data, SMS, MMS, and RCS communication technically distinct. GoreeCloud Data adds encrypted-envelope validation, authenticated sender enforcement, conversation authorization, deterministic retry protection, persistence abstraction, authenticated delivery/read receipts, and opaque encrypted-attachment transport. The attachment surface can submit, fetch as JSON/base64, list metadata, delete with replay-safe tombstones, and download exact ciphertext bytes without asking the server to interpret plaintext media.

The HTTP layer now has one application-facing composition boundary that registers message, receipt, and attachment routes onto the same mux. This preserves both message and attachment listing routes under the shared `/v1/data/conversations/...` namespace and requires one injected authenticator across the composed surface. It does not supply production credentials, persistence, cryptographic sessions, TLS, or deployment configuration.

## Product principles

- Every message identifies its actual transport.
- End-to-end encryption is represented only when the application has verified that state.
- GoreeCloud Data messaging works independently of cellular service.
- Username identities are first-class and are not required to map one-to-one to phone numbers.
- Encrypted GoreeCloud conversations do not silently downgrade to SMS or MMS.
- RCS is integrated only where supported platform and carrier APIs legitimately allow it.
- Voice and video calling remain distinguishable from carrier calling.
- Delivery/read state is recipient-authenticated and is not presented as carrier or cryptographic proof.
- Attachment services transport opaque ciphertext and do not decrypt user content.
- Glaze UI 2.0 or newer, Privacy Shield, Wardveil Security, Everkeep, GoreeCloud Mesh, and GoreeCloud Identity are substantive platform integration requirements for applicable surfaces.

## Repository layout

- `cmd/messenger/` — development executable for exercising core contracts
- `internal/domain/` — transport, encryption, identity, conversation, message, call, Data-envelope, receipt, and attachment contracts
- `internal/service/` — GoreeCloud Data, receipt, attachment, persistence, and authorization boundaries
- `internal/api/` — authenticated HTTP transport plus the unified application-facing Data route-composition boundary
- `docs/architecture.md` — product architecture and trust boundaries
- `docs/security.md` — encryption and security constraints
- `docs/data-messaging.md` — Data service authorization, storage, retry, and carrier-separation contract
- `docs/data-http-api.md` — HTTP API, authentication, authorization, receipt, attachment, runtime composition, and privacy boundary

## Planned clients

Native or platform-appropriate clients are planned for Android, tablets, desktop Linux, and other approved GoreeCloud client platforms. Client work will consume the shared transport and security contracts established here rather than redefining them independently. Consumer username resolution is expected to use a GoreeCloud Identity-owned exact-handle disclosure contract rather than a Messenger-owned account directory.

## Current limitations

This repository remains Development. It does not yet establish production-grade Identity sessions, device/key lifecycle, end-to-end cryptographic session establishment, distributed message delivery, production object storage, push notification delivery, anti-abuse/rate-limit acceptance, carrier adapters, calling media transport, Glaze UI client acceptance, or production deployment evidence. The unified Data handler is a composition boundary only; it is not a production server bootstrap or authentication implementation.

## License

GNU Affero General Public License v3.0 only (`AGPL-3.0-only`).