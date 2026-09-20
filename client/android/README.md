# GoreeCloud Messenger — Native Android Development Client

This directory contains the first original GoreeCloud-owned native Android client foundation for GoreeCloud Messenger.

## Current status

Lifecycle: **Development**

Authoritative Messenger source-bearing main is `b27114a816d4c480a6b4160b5d59d43ec78860c3` after integrated source PR #94. The disconnected native Android Development client remains on that line; historical stacked branches remain provenance only, and exact source revision plus CI evidence must always be read from the revision being evaluated.

Successful source/build evidence for the merged Development foundation does not establish Release Candidate, production, Stable, protected signing, deployment, representative-device, accessibility, or rendered GLAZE UI acceptance for future materially changed revisions.

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
- `GoreeCloudDataTransportAuthority` supplies minimized transport evidence whose configuration, authentication binding, protected channel, bounded failure policy, current deployment/environment binding, and freshness must all be independently accepted before availability may participate in readiness; and
- `E2EESessionAuthority` supplies reviewed cryptographic state plus the exact conversation scope for which active E2EE is verified.

`DataMessagingAuthorityResolver` queries those providers independently for the prepared encrypted message's exact canonical conversation. A provider exception fails closed to that provider's `UNKNOWN` state. Positive authorization cannot substitute for Identity, accepted transport, or E2EE evidence, and positive E2EE cannot substitute for conversation authorization. Authorization and E2EE scopes must still identify the same exact conversation.

`DataMessageSendCoordinator` resolves the provider evidence itself rather than accepting caller-preassembled positive `DataMessagingReadiness.Evidence`. It invokes an injected encrypted Data transport only after all four authorities are positive and the resulting verified conversation exactly matches the prepared message target. This is a composition contract only; the current client still supplies no production providers or network transport and exposes no real Send control.

## GLAZE UI V1.6 source boundary

Authoritative Android Development source maps presentation to **Official Stable GLAZE UI V1.6 / `1.6.0`** at exact release source `a7180679ea851389e0f3004515f9a25f420e716d`. Shared V1.6 consumer eligibility does not auto-certify Messenger, so application acceptance remains explicitly false. The immediate shared rollback Stable is V1.5.1.

`GlazeClientTokens` records exact V1.6 provenance plus the inherited Stable 44dp-equivalent coarse and 32dp pointer-compact target floors. Messenger deliberately keeps a stricter 48dp ordinary interaction target and 56dp accessibility-oriented target. The shell's 20dp gutter, 22dp radius, 14dp section spacing, and neutral pigments remain Messenger-owned product values rather than relabeled canonical V1.6 tokens.

`GlazeMessengerPresentationPolicy` maps presentation-only V1.6 behavior:

- Reduced Transparency fails functional/clear glass down to solid presentation;
- Essential performance reduces costly non-solid material to solid and uses minimal motion;
- Efficient performance reduces glass to raised presentation;
- Reduced Motion uses minimal motion;
- large/extra-large text may make density yield to reflow;
- keyboard-first, screen-reader-optimized, strong-focus, and increased-contrast contexts require strong visible focus;
- caller/platform context defaults remain neutral and no accessibility/performance state is inferred.

This Development candidate replaces the prior always-neutral text-scale input with a narrow Android runtime projection: `resources.configuration.fontScale` is normalized locally, values above the default enable the resolver's large-text reflow signal, and the shared 200%-class boundary enables the extra-large-text signal. Invalid/non-positive/non-finite values fail back to the neutral default. The Activity uses only that resolved presentation state for conservative target sizing and application-owned large-text gutter yielding. It does not activate network, messaging, Identity, E2EE, environmental/private-content sampling, or consequential behavior.

`GlazeMessengerOptics` remains a fail-closed authority boundary:

- presentation stays local and deterministic;
- telemetry, camera access, and remote context are not required;
- environmental-memory influence remains `0.0`;
- Reduced Transparency, Reduced Motion, constrained performance, increased contrast, and visible-focus semantics cannot be overridden by optical context;
- optical/context presentation cannot become semantic or product authority; and
- message content, drafts, conversation/participant identity, receipts, typing, E2EE, Data transport, Identity session, Privacy Shield/Wardveil/Everkeep/Sync state, and remote artwork/media remain prohibited appearance inputs.

The optical adapter remains inactive. Shared V1.6 qualification does not establish Messenger-local visual, accessibility/assistive-technology, large-text, Reduced Motion, Reduced Transparency, contrast/high-contrast, adaptive/form-factor, RTL/localization, representative-device, performance, rollback, Human Visual Excellence, release, or production acceptance.

GLAZE UI presentation cannot manufacture Data transport availability, E2EE state, Identity authorization, Wardveil Security state, Privacy Shield privacy state, Everkeep recovery state, GoreeCloud Sync state, Mesh coordination, Policy decisions, Observability health, message delivery, or any other runtime authority.

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
- complete application-specific GLAZE UI V1.6 rendered/accessibility/adaptive/device/performance/rollback/Human Visual Excellence acceptance;
- representative Android lifecycle, IME, Back, adaptive/form-factor, and accessibility acceptance; and
- protected release signing, artifact provenance, release/rollback/upgrade documentation, and final exact-candidate acceptance.

The server and Android client remain separate runtime compositions. This client foundation does not grant itself server, Identity, cryptographic, recovery, synchronization, carrier, calling, deployment, or production authority.
