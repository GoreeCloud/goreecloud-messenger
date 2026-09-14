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
- `GoreeCloudDataTransportAuthority` supplies minimized transport evidence whose configuration, authentication binding, protected channel, and bounded failure policy must all be independently accepted before availability may participate in readiness; and
- `E2EESessionAuthority` supplies reviewed cryptographic state plus the exact conversation scope for which active E2EE is verified.

`DataMessagingAuthorityResolver` queries those providers independently for the prepared encrypted message's exact canonical conversation. A provider exception fails closed to that provider's `UNKNOWN` state. Positive authorization cannot substitute for Identity, accepted transport, or E2EE evidence, and positive E2EE cannot substitute for conversation authorization. Authorization and E2EE scopes must still identify the same exact conversation.

`DataMessageSendCoordinator` resolves the provider evidence itself rather than accepting caller-preassembled positive `DataMessagingReadiness.Evidence`. It invokes an injected encrypted Data transport only after all four authorities are positive and the resulting verified conversation exactly matches the prepared message target. This is a composition contract only; the current client still supplies no production providers or network transport and exposes no real Send control.

## GLAZE UI V1.4 boundary

The shared design-system source target is **GLAZE UI V1.4 / `1.4.0` Stable — Optical Intelligence** at exact Stable revision `84cb3db4884042f0fa25ed6d475a127fb110f596`.

V1.4 inherits the V1.3 Adaptive Resonance token/component baseline. Messenger preserves the neutral material foundation, 48dp ordinary interaction floor, 56dp Touch Assistance floor, System Light/Dark Development behavior, and V1.3 as the immediate rollback baseline.

`GlazeMessengerOptics` adds a Messenger-specific fail-closed V1.4 source policy:

- the shared Optical Engine is treated as local and deterministic;
- telemetry, camera access, and remote context are not required;
- Messenger environmental color-memory influence is fixed at `0.0`, below the shared V1.4 maximum of `0.08`;
- Reduced Transparency and Forced Colors require solid-accessible treatment;
- Increased Contrast suppresses decorative tint and warmth;
- accessibility cannot be overridden by optical context;
- optical presentation cannot become semantic/product authority; and
- message content, composer drafts, conversation/participant identity, delivery receipts, typing presence, E2EE state, Data-transport state, Identity-session state, Privacy Shield/Wardveil/Everkeep/Sync state, and remote artwork/media are prohibited optical inputs in this disconnected Development shell.

The optical adapter remains inactive. `MessengerClientActivity` does not consume `GlazeMessengerOptics`, so migrating source authority does not silently activate context-aware rendering.

This source mapping does not establish rendered/native-device conformance. GLAZE UI V1.4.0 also does not claim the human/manual/physical-device/subjective qualification assigned to V1.4.1. Messenger still requires fresh exact-revision visual, accessibility/assistive-technology, large-text, Reduced Motion, Reduced Transparency, contrast/high-contrast, adaptive/form-factor, RTL/localization, representative-device, performance, rollback, Human Visual Excellence, and final product-specific acceptance before production promotion.

GLAZE UI presentation cannot manufacture Data transport availability, E2EE state, Identity authorization, Wardveil Security state, Privacy Shield privacy state, Everkeep recovery state, GoreeCloud Sync state, Mesh coordination, message delivery, or any other runtime authority.

## Eight-system Platform Contract boundary

Messenger declares Platform Contract `0.3`, which evaluates all eight Integral Platform Systems: Manager, Privacy Shield, Wardveil Security, Everkeep, GLAZE UI, Mesh, Identity, and Sync.

GoreeCloud Sync remains `applicable-blocked`. Messenger requires future authorized cross-device conversation state, device/key lifecycle coordination, offline continuation, receipts, edits/deletions, and preference continuity, but the disconnected Android shell does not persist messages and no accepted Sync change tracking, version reconciliation, conflict handling, authorized replication, offline-resume, or cross-device runtime exists. Sync is not Everkeep backup/recovery and does not create Identity or E2EE authority.

## Remaining client gates

Before this Android client can approach Release Candidate status it still requires, as applicable:

- production GoreeCloud Identity account/session/device binding and exact-handle username resolution;
- accepted exact conversation authorization authority;
- reviewed cryptographic device/session/key lifecycle plus truthful exact conversation-scoped E2EE state derivation;
- a live GoreeCloud Data client transport and production delivery/synchronization boundary;
- Privacy Shield, Wardveil Security, Everkeep, GoreeCloud Mesh, GoreeCloud Sync, and Manager integration/acceptance;
- complete application-specific GLAZE UI V1.4 rendered/accessibility/adaptive/device/performance/rollback acceptance plus applicable V1.4.1 human/manual/device validation;
- representative Android lifecycle, IME, Back, adaptive/form-factor, and accessibility acceptance; and
- protected release signing, artifact provenance, release/rollback/upgrade documentation, and final exact-candidate acceptance.

The server and Android client remain separate runtime compositions. This client foundation does not grant itself server, Identity, cryptographic, recovery, synchronization, carrier, calling, deployment, or production authority.
