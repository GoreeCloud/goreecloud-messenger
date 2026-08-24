# GoreeCloud Messenger

Native GoreeCloud messaging and calling with usernames, end-to-end encryption, Data messaging, SMS/RCS integration, groups, and video calls.

## Status

Active Development — initial native messaging foundation and Milestone 1 GoreeCloud Data messaging service foundation merged; Data HTTP API foundation in review.

The merged foundation establishes the transport-provenance domain model used to keep Data, SMS, MMS, and RCS communication technically distinct. Milestone 1 adds the first GoreeCloud-controlled Data service boundary for encrypted envelope validation, authenticated sender enforcement, conversation authorization, deterministic retry protection, and persistence abstraction. The current milestone adds an authenticated HTTP adapter for encrypted Data submission and authorized conversation history without introducing plaintext message handling or production credential issuance.

## Product principles

- Every message identifies its actual transport.
- End-to-end encryption is represented only when the application has verified that state.
- GoreeCloud Data messaging works independently of cellular service.
- Username identities are first-class and are not required to map one-to-one to phone numbers.
- Encrypted GoreeCloud conversations do not silently downgrade to SMS or MMS.
- RCS is integrated only where supported platform and carrier APIs legitimately allow it.
- Voice and video calling remain distinguishable from carrier calling.
- Glaze UI, Privacy Shield, Wardveil Security, and Everkeep are substantive platform integrations.

## Repository layout

- `cmd/messenger/` — development executable for exercising core contracts
- `internal/domain/` — transport, encryption, identity, conversation, message, call, and Data-envelope contracts
- `internal/service/` — GoreeCloud Data service and persistence/authorization boundaries
- `internal/api/` — authenticated HTTP transport boundary for encrypted GoreeCloud Data envelopes
- `docs/architecture.md` — product architecture and trust boundaries
- `docs/security.md` — encryption and security constraints
- `docs/data-messaging.md` — Data service authorization, storage, retry, and carrier-separation contract
- `docs/data-http-api.md` — HTTP API, authentication, authorization, and privacy boundary

## Planned clients

Native or platform-appropriate clients are planned for Android, tablets, desktop Linux, and other approved GoreeCloud client platforms. Client work will consume the shared transport and security contracts established here rather than redefining them independently.

## License

GNU Affero General Public License v3.0 only (`AGPL-3.0-only`).
