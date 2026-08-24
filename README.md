# GoreeCloud Messenger

Native GoreeCloud messaging and calling with usernames, end-to-end encryption, Data messaging, SMS/RCS integration, groups, and video calls.

## Status

Development — initial native messaging foundation.

The first implementation slice establishes the transport-provenance domain model used to keep Data, SMS, MMS, and RCS communication technically distinct. The core deliberately treats carrier transports as adapters rather than as the foundation of GoreeCloud messaging.

## Product principles

- Every message identifies its actual transport.
- End-to-end encryption is represented only when the application has verified that state.
- GoreeCloud Data messaging works independently of cellular service.
- Username identities are first-class and are not required to map one-to-one to phone numbers.
- Encrypted GoreeCloud conversations do not silently downgrade to SMS or MMS.
- RCS is integrated only where supported platform and carrier APIs legitimately allow it.
- Voice and video calling remain distinguishable from carrier calling.
- Glaze UI, Privacy Shield, Wardveil Security, and Everkeep are substantive platform integrations.

## Initial repository layout

- `cmd/messenger/` — development executable for exercising core contracts
- `internal/domain/` — transport, encryption, identity, conversation, and message domain contracts
- `docs/architecture.md` — initial product architecture and trust boundaries
- `docs/security.md` — initial encryption and security constraints

## Planned clients

Native or platform-appropriate clients are planned for Android, tablets, desktop Linux, and other approved GoreeCloud client platforms. Client work will consume the shared transport and security contracts established here rather than redefining them independently.

## License

GNU Affero General Public License v3.0 only (`AGPL-3.0-only`).
