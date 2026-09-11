# GoreeCloud Messenger Security Model

## Security principle

Messenger security indicators describe verified technical state. E2EE, Privacy Shield, Wardveil Security, and related protection labels must not be displayed as decorative claims.

## Cryptography boundary

Messenger will use an established, security-reviewed end-to-end encryption protocol or implementation. GoreeCloud will not invent a proprietary cryptographic algorithm for message confidentiality.

No authoritative GoreeCloud record currently selects a concrete Messenger E2EE protocol or implementation. Development code must therefore not claim or imply Signal Protocol, MLS, Double Ratchet, libsignal, or another concrete protocol unless a future authoritative decision and implementation evidence establish it.

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

## Protocol-neutral E2EE acceptance seam

The current native Android Development client does not implement cryptography. Its `E2EESessionAuthority` is a provider-owned seam that may expose only a minimized readiness projection.

A provider-supplied `E2EE_ACTIVE` claim is accepted for Data-send readiness only when the same evidence also establishes all of the following for the exact requested conversation:

- `E2EEImplementationReviewState.ACCEPTED`
- `E2EEDeviceIdentityState.ENROLLED`
- `E2EESessionEstablishmentState.ESTABLISHED`
- `E2EEKeyLifecycleState.CURRENT`
- an exact canonical bounded opaque E2EE conversation identifier matching the requested conversation

`E2EESessionEvidence.readinessProjectionFor(...)` downgrades a contradictory active claim to `UNKNOWN` when any required acceptance fact is missing, negative, noncanonical, or scoped to another conversation. Explicit negative cryptographic states are never upgraded by positive auxiliary facts.

The projection carries no key material, session secret, algorithm claim, cipher-suite identifier, device secret, fingerprint, ciphertext transformation, reusable credential, or network endpoint. Group and multi-device implementations remain responsible for resolving their complete participant/device state before exposing this minimized projection.

Positive unit-test fixtures exercise only this policy contract. They are not a real security review, cryptographic session, key lifecycle, or production E2EE implementation.

The disconnected Android Development source guard also rejects direct `java.security` and `javax.crypto` implementation markers. That guard must be deliberately replaced under a future governed cryptographic implementation rather than silently bypassed.

See [`development/e2ee-authority-acceptance.md`](development/e2ee-authority-acceptance.md) for the current FR-005 Development boundary.

## Initial enforcement

The current domain layer rejects an `e2ee` assertion on SMS, MMS, or RCS. This prevents the application core from presenting GoreeCloud E2EE on a carrier transport without a separately implemented and verified cryptographic envelope.

## Authority identifier boundary

GoreeCloud authority-owned identifiers are opaque scopes, not user-entered text to normalize. Current message, conversation, participant/sender, authenticated-user, receipt, attachment, typing-presence, typing-preference, and replay-nonce boundaries require valid Unicode, nonblank values that are already canonical, no C0 or DEL control characters, and a maximum of 512 UTF-16 code units. Internal spaces and punctuation remain part of the exact identifier.

Data HTTP handlers pass path identifiers through without trimming or rewriting them. The message, receipt, attachment, typing-presence, typing-preference, and authenticated runtime-diagnostic surfaces reject noncanonical authenticated identities. Domain and service layers revalidate the same boundary before authorization or persistence so a caller cannot bypass the rule by invoking a service directly. A value such as ` conversation-1` must therefore fail closed rather than being transformed into the authority scope `conversation-1`.

The Go boundary also rejects malformed UTF-8 before UTF-16 length evaluation so byte sequences that cannot exist as a native Android `String` cannot become a different replacement-character scope on the server.

This Development rule prevents identifier aliasing; it is not an authentication, authorization, cryptographic-session, or production Identity implementation by itself.

## Metadata minimization

Message content, attachment content, encryption key material, device identity secrets, contact data, and communication metadata must be collected and retained only when required for the application to function.

## Fallback behavior

Encrypted Data communication must not silently downgrade. If Data messaging is unavailable and a carrier transport is offered, the user must be told that the transport and protection state are changing before the message is sent when practical and required by the final interaction contract.

## Security identities

Wardveil Security is the platform security and protection identity. Privacy Shield is the platform privacy identity. Both integrations must be backed by actual application state and evidence.
