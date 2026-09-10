# Protocol-Neutral E2EE Authority Acceptance Boundary

Status: Development

This document describes the current Messenger FR-005 Android acceptance boundary. It is an evidence contract for a future reviewed cryptographic authority, not an end-to-end encryption implementation.

## Governing requirement

The canonical GoreeCloud Messenger specification requires an established, security-reviewed cryptographic protocol or implementation rather than an unreviewed custom cryptographic design. The eventual architecture must account for cryptographic device identity, authenticated session establishment, key lifecycle and change handling, verification, group behavior, multi-device enrollment/revocation, secure attachments, protected local state, and safe recovery behavior.

No authoritative GoreeCloud record currently selects Signal Protocol, MLS, Double Ratchet, libsignal, or another concrete Messenger protocol. Messenger therefore does not choose or imply one in this Development slice.

## Android acceptance projection

`E2EESessionAuthority` remains the future provider-owned cryptographic seam. A provider may project `E2EE_ACTIVE` for one exact conversation only when the minimized evidence also reports all of the following:

- the responsible cryptographic implementation is accepted under the applicable security review;
- the local device has an enrolled cryptographic identity;
- the conversation-scoped cryptographic session is established; and
- the provider accepts the current session/device key lifecycle state, including required key-change or rotation processing.

The evidence must also carry the exact bounded opaque conversation identifier for which protection is asserted.

`E2EESessionEvidence.readinessProjectionFor(...)` fails closed. A bare `E2EE_ACTIVE` enum, missing or rejected review state, missing device enrollment, unestablished session, non-current key lifecycle, noncanonical conversation scope, or conversation-scope mismatch becomes `UNKNOWN` before the evidence reaches `DataMessagingReadiness`.

Explicit negative cryptographic states are not upgraded by positive auxiliary evidence.

## Minimized boundary

The acceptance projection carries no:

- private key or public key material;
- session secret;
- reusable credential or token;
- algorithm identifier or cipher-suite claim;
- ciphertext transformation;
- device secret;
- participant key inventory;
- key fingerprint;
- recovery secret; or
- network endpoint.

Group and multi-device implementations remain responsible for evaluating their complete participant/device state before returning the minimized acceptance projection.

## Development client guard

The disconnected Android Development source guard continues to reject Internet, contacts, SMS, calling, microphone, camera, local persistence, and client networking authority. This FR-005 slice additionally rejects direct `java.security` and `javax.crypto` implementation markers in the current client source. That guard exists to prevent this acceptance-only branch from silently becoming a local unreviewed cryptographic implementation.

A future governed cryptographic implementation will require an explicit architectural change, its own authoritative protocol/implementation selection, review evidence, and corresponding update to this Development guard.

## Send-readiness relationship

The E2EE acceptance projection is still only one independent prerequisite. A future encrypted GoreeCloud Data send remains blocked unless all of the following independently succeed for the same exact conversation:

1. GoreeCloud Identity authentication;
2. conversation participation authorization;
3. GoreeCloud Data transport availability; and
4. the accepted active-E2EE projection described here.

No one authority can manufacture another authority's result, and transport availability does not imply encryption.

## Acceptance boundary

Tests may construct positive E2EE acceptance fixtures in order to exercise the policy contract. Those fixtures do not establish a real cryptographic implementation, real device enrollment, real session establishment, real key lifecycle, real security review, Wardveil Security acceptance, Privacy Shield acceptance, or production E2EE.

This Development slice does not implement encryption, select a protocol, establish keys or sessions, encrypt/decrypt messages, create production key storage, provide multi-device cryptographic synchronization, perform contact/device verification, or establish release/production/Stable acceptance.
