# GoreeCloud Messenger Security Model

## Security principle

Messenger security indicators describe verified technical state. E2EE, Privacy Shield, Wardveil Security, and related protection labels must not be displayed as decorative claims.

## Cryptography boundary

Messenger will use an established, security-reviewed end-to-end encryption protocol or implementation. GoreeCloud will not invent a proprietary cryptographic algorithm for message confidentiality.

The eventual implementation must cover:

- device identity keys
- authenticated session establishment
- forward secrecy where supported
- key rotation and key-change handling
- multi-device enrollment and revocation
- group key management
- device verification
- secure attachment encryption
- encrypted local state where appropriate
- safe backup, migration, and recovery behavior

## Initial enforcement

The current domain layer rejects an `e2ee` assertion on SMS, MMS, or RCS. This prevents the application core from presenting GoreeCloud E2EE on a carrier transport without a separately implemented and verified cryptographic envelope.

## Authority identifier boundary

GoreeCloud Data authority-owned identifiers are opaque scopes, not user-entered text to normalize. Message, conversation, participant/sender, receipt, attachment, and replay-nonce identifiers must be nonblank, must already be canonical at the boundary, must not contain C0 or DEL control characters, and are bounded to 512 UTF-16 code units. Internal spaces and punctuation remain part of the exact identifier.

HTTP handlers must pass path identifiers through without trimming or rewriting them. Domain and service layers revalidate the same boundary before authorization or persistence so a caller cannot bypass the rule by invoking a service directly. A value such as ` conversation-1` must therefore fail closed rather than being transformed into the authority scope `conversation-1`.

This Development rule prevents identifier aliasing; it is not an authentication, authorization, cryptographic-session, or production Identity implementation by itself.

## Metadata minimization

Message content, attachment content, encryption key material, device identity secrets, contact data, and communication metadata must be collected and retained only when required for the application to function.

## Fallback behavior

Encrypted Data communication must not silently downgrade. If Data messaging is unavailable and a carrier transport is offered, the user must be told that the transport and protection state are changing before the message is sent when practical and required by the final interaction contract.

## Security identities

Wardveil Security is the platform security and protection identity. Privacy Shield is the platform privacy identity. Both integrations must be backed by actual application state and evidence.
