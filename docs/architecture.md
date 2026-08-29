# GoreeCloud Messenger Architecture

## Role

GoreeCloud Messenger is a native GoreeCloud Suite communication application. Its core must remain functional through GoreeCloud-controlled Data messaging without depending on carrier messaging services.

## Primary boundaries

### Identity

GoreeCloud usernames are first-class identifiers. A user may have verified phone numbers or other approved identifiers, but Data messaging must not require a telephone number.

### Data transport

The Data transport is the GoreeCloud-controlled Internet messaging path and is the primary transport for usernames, end-to-end encryption, groups, attachments, voice calls, and video calls.

The current HTTP implementation has an application-facing `DataRuntimeHandler` composition boundary that joins message, receipt, and encrypted-attachment routes under one required authenticator. The runtime handler composes established services; it does not become the Identity authority, cryptographic-session authority, persistence authority, or deployment configuration.

### Carrier adapters

SMS, MMS, and RCS are adapters. Their availability depends on operating-system, device-role, carrier, and API support. A carrier adapter cannot redefine the security properties of the Data transport.

### Transport provenance

Each message stores its transport independently. A conversation can therefore contain messages with different transports without losing their provenance.

The initial contract recognizes:

- `data`
- `sms`
- `mms`
- `rcs`

### Encryption provenance

Encryption state is explicit and evidence-backed. The initial domain model recognizes `none`, `e2ee`, and `unknown` states. GoreeCloud E2EE may only be asserted on the GoreeCloud Data transport.

### No silent downgrade

An encrypted Data conversation must not silently turn into SMS or MMS. Any fallback or transport transition must be represented explicitly by the client.

## Planned logical components

1. Identity and device service
2. Conversation and message service
3. Data transport service
4. Cryptographic session service
5. SMS/MMS adapter
6. RCS adapter where supported
7. Attachment and media service
8. Voice/video calling service
9. Multi-device synchronization
10. Notification delivery
11. Local client storage
12. Glaze UI presentation layer
13. Privacy Shield integration
14. Wardveil Security integration
15. Everkeep portability and recovery integration

## Client direction

Clients must share the same message, transport, encryption, and identity contracts. Platform-specific implementations may differ, but no client may present a transport or encryption state that conflicts with the underlying verified state.
