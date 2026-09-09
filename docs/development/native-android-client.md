# Native Android Client Foundation

Status: Development source candidate.

GoreeCloud Messenger now has a first original GoreeCloud-owned native Android client shell under `client/android`. This is intentionally a client-structure and communication-provenance slice, not a connected messaging application.

## Current capability

The client:

- builds as package `com.goreecloud.messenger`;
- uses Android platform UI controls without a web wrapper;
- presents an explicit disconnected **Native Android Development preview** boundary;
- presents transport/protection examples only, not live conversations or user data;
- defines a typed `CommunicationProvenance` contract for Data, SMS, MMS, and RCS presentation;
- permits `Data · E2EE` only when the modeled protection state is explicitly `E2EE_ACTIVE`;
- rejects `E2EE_ACTIVE` construction for SMS, MMS, and RCS in the current client foundation;
- renders unverified Data protection as `Data · Protection not verified`, `Data · E2EE unavailable`, or `Data · Not end-to-end encrypted` rather than implying encryption;
- requires future conversation authorization and future active-E2EE evidence to each carry their own exact conversation identifier rather than relying on unscoped positive enums;
- declares readiness only when those independent authorization and cryptographic scopes identify the same conversation;
- defines independent provider seams for GoreeCloud Identity session state, exact conversation authorization, GoreeCloud Data transport availability, and exact conversation-scoped E2EE state;
- resolves those four authorities independently for the exact prepared-message conversation and maps a failed provider to `UNKNOWN` rather than borrowing readiness from another authority;
- prevents the final Data send coordinator from accepting a caller-preassembled positive readiness object;
- refuses the injected encrypted Data transport unless independently resolved readiness is positive and its unified verified conversation scope exactly matches the prepared message target;
- disables Android application backup for this Development shell; and
- declares no Internet, contacts, SMS, phone, microphone, or camera permission.

The source boundary check also rejects client-side network-library/import markers and local persistence primitives in the current shell. It requires independently conversation-scoped authorization and E2EE evidence, explicit provider-owned readiness resolution at the future pre-transport seam, and exact prepared-target binding before an injected transport can be called. This is intentionally stricter than the eventual product because the current client has no accepted production GoreeCloud Identity, cryptographic-session, Data-client, carrier-adapter, call-media, or local encrypted-state authority yet.

## GLAZE UI boundary

The native Android Development client targets current Stable **GLAZE UI V1.3 / 1.3.0 Adaptive Resonance** at exact Stable integration revision `fc7cc91d2eace8da2371371c2855c24cbcb326a1`, with `1.2.0` retained as the recorded rollback baseline. The repository-local client token mapping retains the inherited neutral Light/Dark material foundation and the current 48/56 dp interaction floors.

This is source-level Development integration evidence only. Complete rendered, interaction, accessibility, adaptive/form-factor, representative-device, workflow, performance, rollback, Human Visual Excellence, release, and production acceptance remain independent gates before Messenger can claim complete current GLAZE UI application acceptance.

## Security and privacy boundary

The client currently has no production account/session handling, credentials, keys, cryptographic protocol implementation, network transport, message plaintext/ciphertext persistence, contacts access, carrier messaging, calling, telemetry, or background synchronization. It must not display a real E2EE/security/privacy state until the responsible platform and cryptographic authorities supply verifiable state.

`DataMessagingReadiness` does not itself authenticate, authorize a conversation, establish transport, or establish E2EE. `DataMessagingAuthorityResolver` instead provides a narrow integration seam for four separately responsible future authorities:

- `GoreeCloudIdentitySessionAuthority` supplies authentication state only, with no credentials or tokens exposed to this projection;
- `ConversationAuthorizationAuthority` supplies participant state plus the exact authorized conversation scope;
- `GoreeCloudDataTransportAuthority` supplies Data transport availability only and cannot imply E2EE or authorization; and
- `E2EESessionAuthority` supplies cryptographic state plus the exact conversation for which active E2EE has been verified, without exposing keys or session secrets.

Each authority is queried independently for the prepared message's exact canonical conversation. Provider exceptions fail closed to the corresponding `UNKNOWN` state without upgrading any other authority. Positive conversation authorization and E2EE still have to name the same exact bounded opaque conversation identifier. `DataMessageSendCoordinator` resolves this evidence itself and no longer accepts a caller-assembled `DataMessagingReadiness.Evidence` object at the send seam. Only after the independent authorities agree does the coordinator compare the resulting verified conversation with the prepared encrypted message target before an injected transport can be invoked.

These interfaces are integration contracts, not authority implementations. The disconnected Development Activity still supplies no production providers and exposes no real Send action. Future connected work must preserve the product specification's separation between GoreeCloud Identity, GoreeCloud Data transport, cryptographic sessions/key lifecycle, optional SMS/MMS/RCS adapters, calling, notifications, attachment handling, local protected state, multi-device synchronization, Privacy Shield, Wardveil Security, Everkeep, Mesh, and Manager responsibilities.

## Acceptance boundary

A green Android build, provider/readiness/coordinator unit tests, emulator acceptance, or source-boundary check proves only this Development slice. It does not establish a working Messenger client, production Identity, device identity, username resolution, conversation authorization, reviewed E2EE device/session/key authority, live Data networking, SMS/MMS/RCS support, calling, push delivery, multi-device operation, complete current-Stable GLAZE UI acceptance, platform-system acceptance, protected signing, production deployment, Release Candidate approval, or Stable qualification.
