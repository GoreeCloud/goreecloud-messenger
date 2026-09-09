# GoreeCloud Messenger — Native Android Development Client

This directory contains the first original GoreeCloud-owned native Android client foundation for GoreeCloud Messenger.

## Current status

Lifecycle: **Development**

This branch is a stacked Development continuation of the current Messenger Android client line. Exact source revision and CI evidence must be read from the branch/PR being reviewed rather than copied forward from an older candidate.

Successful source/build evidence for earlier Development revisions does not establish Release Candidate, production, Stable, signing, deployment, representative-device, accessibility, or rendered GLAZE UI acceptance for this materially changed branch.

## Deliberate Development boundary

The current Android client is a disconnected presentation and provenance foundation. It has no production account, live conversation, message composer, Send action, network transport, production GoreeCloud Identity session, reviewed cryptographic session/key lifecycle, carrier transport, calling implementation, or durable message storage.

The Android manifest intentionally declares none of the following permissions in this slice:

- Internet
- Contacts
- SMS
- Phone calling
- Microphone
- Camera

Android application backup is disabled. The client source boundary also rejects network/storage implementation primitives and dependencies that would silently expand this slice beyond its reviewed authority.

## Communication provenance and messaging authority

`CommunicationProvenance` separates actual transport from verified protection state. The contract supports visible transport labels for GoreeCloud Data, SMS, MMS, and RCS while preventing carrier transports from being represented as GoreeCloud `E2EE_ACTIVE`.

`Data · E2EE` is valid only for an explicitly verified GoreeCloud Data protection state. Unknown or unavailable protection is disclosed rather than upgraded into an encryption claim. The sample cards rendered by the Development Activity are examples only; they are not conversations, user data, or runtime service state.

The future Data-send seam is intentionally split across independent authorities:

- `GoreeCloudIdentitySessionAuthority` supplies authentication state only;
- `ConversationAuthorizationAuthority` supplies participant state plus its exact authorized conversation scope;
- `GoreeCloudDataTransportAuthority` supplies Data transport availability only; and
- `E2EESessionAuthority` supplies reviewed cryptographic state plus the exact conversation scope for which active E2EE is verified.

`DataMessagingAuthorityResolver` queries those providers independently for the prepared encrypted message's exact canonical conversation. A provider exception fails closed to that provider's `UNKNOWN` state. Positive authorization cannot substitute for Identity, transport, or E2EE evidence, and positive E2EE cannot substitute for conversation authorization. Authorization and E2EE scopes must still identify the same exact conversation.

`DataMessageSendCoordinator` resolves the provider evidence itself rather than accepting caller-preassembled positive `DataMessagingReadiness.Evidence`. It invokes an injected encrypted Data transport only after all four authorities are positive and the resulting verified conversation exactly matches the prepared message target. This is a composition contract only; the current client still supplies no production providers or network transport and exposes no real Send control.

## GLAZE UI boundary

The current shared design-system consumer target is **GLAZE UI V1.3 / `1.3.0` Stable — Adaptive Resonance**.

Messenger's repository-local Android mapping records:

- exact Stable integration revision `fc7cc91d2eace8da2371371c2855c24cbcb326a1`;
- inherited optical foundation `tokens/glaze-v1.2-optical-foundation.candidate.json`;
- adaptive contract `contracts/v1.3/adaptive-resonance.plan.json`;
- Stable web entrypoint `css/glaze-v1.3.0.css`;
- Stable runtime entrypoint `js/glaze-v1.3.0.mjs`;
- rollback baseline `1.2.0`;
- the rule **Neutral glass is the material. Color is an accent.**;
- 48dp ordinary interaction floor and 56dp Touch Assistance floor; and
- neutral Light and Dark canvas/surface values, with tests that reject the historical Deep Teal/Soft Amber substrate pattern.

The current Development Activity only maps Android System Light and Dark. It does **not** silently alias Deep Dark to ordinary Dark or infer adaptive environmental accent behavior; those remain downstream application-specific acceptance items if exposed in supported release scope.

This source mapping does not establish rendered/native-device conformance. Messenger still requires exact-revision visual review, accessibility/assistive-technology evidence, large-text behavior, Reduced Motion, Reduced Transparency/effects-free behavior where applicable, contrast/high-contrast behavior, adaptive/form-factor behavior, RTL/localization, representative-device evidence, performance evidence, rollback acceptance, Human Visual Excellence review, and final product-specific acceptance before production promotion.

GLAZE UI presentation cannot manufacture Data transport availability, E2EE state, Identity authorization, Wardveil Security state, Privacy Shield privacy state, Everkeep recovery state, Mesh coordination, message delivery, or any other runtime authority.

## Remaining client gates

Before this Android client can approach Release Candidate status it still requires, as applicable:

- production GoreeCloud Identity account/session/device binding and exact-handle username resolution;
- accepted exact conversation authorization authority;
- reviewed cryptographic device/session/key lifecycle plus truthful exact conversation-scoped E2EE state derivation;
- a live GoreeCloud Data client transport and production delivery/synchronization boundary;
- Privacy Shield, Wardveil Security, Everkeep, GoreeCloud Mesh, and Manager integration/acceptance;
- complete application-specific GLAZE UI V1.3 rendered/accessibility/adaptive/device/performance/rollback acceptance;
- representative Android lifecycle, IME, Back, adaptive/form-factor, and accessibility acceptance; and
- protected release signing, artifact provenance, release/rollback/upgrade documentation, and final exact-candidate acceptance.

The server and Android client remain separate runtime compositions. This client foundation does not grant itself server, Identity, cryptographic, recovery, carrier, calling, deployment, or production authority.
