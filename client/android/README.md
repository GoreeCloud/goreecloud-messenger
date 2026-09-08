# GoreeCloud Messenger — Native Android Development Client

This directory contains the first original GoreeCloud-owned native Android client foundation for GoreeCloud Messenger.

## Current status

Lifecycle: **Development**

Current exact candidate revision: `d67dd1f38eb68e86859af4f56c38009f164fad6e`.

Exact-head validation completed successfully on September 4, 2026:

- Platform Contract #15 — success.
- Messenger Foundation #434 — success.
- Messenger Android Client #4 — success, including the Development authority boundary, client unit tests, and debug APK build.

This establishes source/build evidence for this Development client foundation only. It does not establish Release Candidate, production, Stable, signing, deployment, representative-device, accessibility, or rendered GLAZE UI acceptance.

## Deliberate Development boundary

The current Android client is a disconnected presentation and provenance foundation. It has no account, live conversation, message composer, send action, network transport, production GoreeCloud Identity session, cryptographic session/key lifecycle, carrier transport, calling implementation, or durable message storage.

The Android manifest intentionally declares none of the following permissions in this slice:

- Internet
- Contacts
- SMS
- Phone calling
- Microphone
- Camera

Android application backup is disabled. The client source boundary also rejects network/storage implementation primitives and dependencies that would silently expand this slice beyond its reviewed authority.

## Communication provenance

`CommunicationProvenance` separates actual transport from verified protection state. The contract supports visible transport labels for GoreeCloud Data, SMS, MMS, and RCS while preventing carrier transports from being represented as GoreeCloud `E2EE_ACTIVE`.

`Data · E2EE` is valid only for an explicitly verified GoreeCloud Data protection state. Unknown or unavailable protection is disclosed rather than upgraded into an encryption claim.

The sample cards rendered by the Development Activity are examples only. They are not conversations, user data, or runtime service state.

## GLAZE UI boundary

The current native shell is a Development source/design mapping, not accepted application conformance. The GoreeCloud Platform Contract requires GLAZE UI V1.1 / `1.1.0`, but that immutable Stable release has a known import-closure defect. A corrected immutable Stable release, explicit Messenger re-pin, and fresh application-specific rendered, accessibility, resilience, adaptive, device, and Human Visual Excellence acceptance remain required.

## Remaining client gates

Before this Android client can approach Release Candidate status it still requires, as applicable:

- production GoreeCloud Identity account/session/device binding;
- reviewed cryptographic session and key lifecycle plus truthful E2EE state derivation;
- a real GoreeCloud Data client transport and synchronization boundary;
- Privacy Shield, Wardveil Security, Everkeep, GoreeCloud Mesh, and Manager integration/acceptance;
- current corrected-Stable GLAZE UI re-pin and rendered/accessibility/device acceptance;
- representative Android lifecycle, IME, Back, adaptive/form-factor, and accessibility acceptance;
- protected release signing, artifact provenance, release/rollback/upgrade documentation, and final exact-candidate acceptance.

The server and Android client remain separate runtime compositions. This client foundation does not grant itself server, Identity, cryptographic, recovery, carrier, calling, deployment, or production authority.
