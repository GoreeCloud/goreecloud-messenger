# GoreeCloud Messenger — Native Android Development Client

This directory contains the first original GoreeCloud-owned native Android client foundation for GoreeCloud Messenger.

## Current status

Lifecycle: **Development**

This client has been restacked onto current authoritative Messenger `main` as a bounded migration candidate. Historical stacked branches remain provenance only; exact source revision and CI evidence must be read from the current branch/PR.

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

## GLAZE UI V1.5.1 boundary

The shared design-system source target is **GLAZE UI V1.5 / `1.5.1` Stable — Contextual + Capability Awareness**. Stable lifecycle authority is `98da57064ede0f334627b632bc16801f580331af`; the reviewed V1.5 implementation anchor is `ee1032a0822ab8e103f8afe48e5c1859fde65cc9`; V1.5.1 qualification remains bound to `5b59d0e36950d737dba35b58ae58058684e0831b`; V1.5.0 is the immediate shared rollback baseline.

V1.5.1 preserves the reviewed V1.5 contextual/capability presentation over the inherited V1.4.1 optical material baseline. Messenger preserves the neutral material foundation, 48dp ordinary interaction floor, 56dp Touch Assistance floor, System Light/Dark Development behavior, and V1.5.0 as the immediate rollback baseline.

`GlazeMessengerOptics` keeps a Messenger-specific fail-closed source policy:

- the shared Optical Engine is treated as local and deterministic;
- telemetry, camera access, and remote context are not required;
- Messenger environmental color-memory influence is fixed at `0.0`, below the shared maximum of `0.08`;
- Reduced Transparency and Forced Colors require solid-accessible treatment;
- Increased Contrast suppresses decorative tint and warmth;
- accessibility cannot be overridden by optical context;
- optical presentation cannot become semantic/product authority; and
- message content, composer drafts, conversation/participant identity, delivery receipts, typing presence, E2EE state, Data-transport state, Identity-session state, Privacy Shield/Wardveil/Everkeep/Sync state, and remote artwork/media are prohibited optical inputs in this disconnected Development shell.

The optical adapter remains inactive. `MessengerClientActivity` does not consume `GlazeMessengerOptics`, so migrating source authority does not silently activate context-aware rendering.

Shared GLAZE UI V1.5.1 human/manual/device qualification does not auto-certify Messenger. Messenger still requires fresh exact-revision visual, accessibility/assistive-technology, large-text, Reduced Motion, Reduced Transparency, contrast/high-contrast, adaptive/form-factor, RTL/localization, representative-device, performance, rollback, Human Visual Excellence, and final product-specific acceptance before production promotion.

GLAZE UI presentation cannot manufacture Data transport availability, E2EE state, Identity authorization, Wardveil Security state, Privacy Shield privacy state, Everkeep recovery state, GoreeCloud Sync state, Mesh coordination, message delivery, or any other runtime authority.

## Platform Contract boundary

Messenger targets Platform Contract `0.4`, which evaluates all nine Integral Platform Systems: Manager, Privacy Shield, Wardveil Security, Everkeep, GLAZE UI, Mesh, Identity, GoreeCloud Policy, and GoreeCloud Observability. GoreeCloud Sync remains separately governed and is not a tenth Integral Platform System.

Messenger still requires future authorized cross-device conversation state, device/key lifecycle coordination, offline continuation, receipts, edits/deletions, and preference continuity, but the disconnected Android shell does not persist messages and no accepted Sync change tracking, version reconciliation, conflict handling, authorized replication, offline-resume, or cross-device runtime exists. Sync is not Everkeep backup/recovery and does not create Identity or E2EE authority.

The accepted central Platform Contract `0.4` implementation is pinned from `GoreeCloud/GoreeCloud` merge `6cb150d512647a0401b4da9e4741d7591693dee0`. Component-specific Policy and Observability runtime integration remains blocked; the declaration does not manufacture acceptance.

## Remaining client gates

Before this Android client can approach Release Candidate status it still requires, as applicable:

- production GoreeCloud Identity account/session/device binding and exact-handle username resolution;
- accepted exact conversation authorization authority;
- reviewed cryptographic device/session/key lifecycle plus truthful exact conversation-scoped E2EE state derivation;
- a live GoreeCloud Data client transport and production delivery/synchronization boundary;
- Privacy Shield, Wardveil Security, Everkeep, GoreeCloud Mesh, GoreeCloud Sync, and Manager integration/acceptance;
- complete application-specific GLAZE UI V1.5.1 rendered/accessibility/adaptive/device/performance/rollback/Human Visual Excellence acceptance;
- representative Android lifecycle, IME, Back, adaptive/form-factor, and accessibility acceptance; and
- protected release signing, artifact provenance, release/rollback/upgrade documentation, and final exact-candidate acceptance.

The server and Android client remain separate runtime compositions. This client foundation does not grant itself server, Identity, cryptographic, recovery, synchronization, carrier, calling, deployment, or production authority.
