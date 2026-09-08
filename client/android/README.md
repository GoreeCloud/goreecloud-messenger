# GoreeCloud Messenger — Native Android Development Client

This directory contains the first original GoreeCloud-owned native Android client foundation for GoreeCloud Messenger.

## Current status

Lifecycle: **Development**

This branch is a stacked Development continuation of the current Messenger Android client line. Exact source revision and CI evidence must be read from the branch/PR being reviewed rather than copied forward from an older candidate.

Successful source/build evidence for earlier Development revisions does not establish Release Candidate, production, Stable, signing, deployment, representative-device, accessibility, or rendered GLAZE UI acceptance for this materially changed branch.

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

The current shared design-system consumer target is **GLAZE UI V1.2 / `1.2.0` Stable**.

Messenger's repository-local Android mapping records:

- Stable promotion merge revision `f285b9145e27e6e7027b075c37299d101945c272`;
- V1.2 source-qualification anchor `b0eadf9a60f73d45caffb62ffc7e9e0334cddc97`;
- optical foundation `tokens/glaze-v1.2-optical-foundation.candidate.json` (the filename preserves upstream source lineage; V1.2 lifecycle is Stable);
- the rule **Neutral glass is the material. Color is an accent.**;
- 48dp ordinary interaction floor and 56dp Touch Assistance floor;
- neutral Light and Dark canvas/surface values, with tests that reject the historical Deep Teal/Soft Amber substrate pattern.

The current Development Activity only maps Android System Light and Dark. It does **not** silently alias Deep Dark to ordinary Dark; a distinct Deep Dark runtime mapping remains an application-specific downstream acceptance item if that appearance is exposed in the supported release scope.

This source mapping does not establish rendered/native-device conformance. Messenger still requires exact-revision visual review, accessibility/assistive-technology evidence, large-text behavior, Reduced Motion, Reduced Transparency/effects-free behavior where applicable, contrast/high-contrast behavior, adaptive/form-factor behavior, RTL/localization, representative-device evidence, performance evidence, and final product-specific acceptance before production promotion.

GLAZE UI presentation cannot manufacture Data transport availability, E2EE state, Identity authorization, Wardveil security state, Privacy Shield privacy state, Everkeep recovery state, Mesh coordination, message delivery, or any other runtime authority.

## Remaining client gates

Before this Android client can approach Release Candidate status it still requires, as applicable:

- production GoreeCloud Identity account/session/device binding;
- reviewed cryptographic session and key lifecycle plus truthful E2EE state derivation;
- a real GoreeCloud Data client transport and synchronization boundary;
- Privacy Shield, Wardveil Security, Everkeep, GoreeCloud Mesh, and Manager integration/acceptance;
- complete application-specific GLAZE UI V1.2 rendered/accessibility/adaptive/device acceptance, including any Deep Dark scope actually exposed;
- representative Android lifecycle, IME, Back, adaptive/form-factor, and accessibility acceptance;
- protected release signing, artifact provenance, release/rollback/upgrade documentation, and final exact-candidate acceptance.

The server and Android client remain separate runtime compositions. This client foundation does not grant itself server, Identity, cryptographic, recovery, carrier, calling, deployment, or production authority.
