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
- models future conversation-access readiness as an exact verified conversation identifier rather than an unscoped participant flag;
- refuses the injected encrypted Data transport when the verified conversation scope does not exactly match the prepared message target;
- disables Android application backup for this Development shell; and
- declares no Internet, contacts, SMS, phone, microphone, or camera permission.

The source boundary check also rejects client-side network-library/import markers and local persistence primitives in the current shell. It additionally requires conversation-scoped authorization evidence at the future pre-transport seam. This is intentionally stricter than the eventual product because the first slice has no accepted GoreeCloud Identity, cryptographic-session, Data-client, carrier-adapter, call-media, or local encrypted-state authority yet.

## GLAZE UI boundary

The client contains a bounded repository-local source mapping for neutral Light/Dark surfaces and the current 48 dp general interaction floor. This is Development presentation source only. It does **not** establish current GLAZE UI application conformance.

The platform requires the GLAZE UI V1.1 line. Immutable 1.1.0 remains the formal Stable source promotion but has a known import-closure defect, and the corrective 1.1.1 work remains a Draft Release Candidate. Messenger must explicitly re-pin the corrected immutable Stable release after governed promotion and repeat applicable rendered, accessibility, adaptive/form-factor, representative-device, and Human Visual Excellence acceptance before current Glaze conformance can be claimed.

## Security and privacy boundary

The client currently has no account/session handling, credentials, keys, cryptographic protocol implementation, network transport, message plaintext/ciphertext persistence, contacts access, carrier messaging, calling, telemetry, or background synchronization. It must not display a real E2EE/security/privacy state until the responsible platform and cryptographic authorities supply verifiable state.

`DataMessagingReadiness` does not itself authorize a conversation. A future authorization authority must supply both a positive participant decision and the exact verified conversation identifier. Missing/blank scope fails closed, and `DataMessageSendCoordinator` compares that verified scope with the prepared encrypted message's target before the transport seam can be invoked. This prevents authorization evidence for one conversation from being reused for another while adding no sender identity, credential, key, plaintext, endpoint, or transport authority to the client model.

Future connected work must preserve the product specification's separation between GoreeCloud Identity, GoreeCloud Data transport, cryptographic sessions/key lifecycle, optional SMS/MMS/RCS adapters, calling, notifications, attachment handling, local protected state, multi-device synchronization, Privacy Shield, Wardveil Security, Everkeep, Mesh, and Manager responsibilities.

## Acceptance boundary

A green Android build, provenance/readiness/coordinator unit tests, or source-boundary check proves only this Development slice. It does not establish a working Messenger client, production Identity, conversation authorization, E2EE, Data messaging, SMS/MMS/RCS support, calling, push delivery, multi-device operation, current-Stable Glaze acceptance, platform-system acceptance, protected signing, physical-device acceptance, Release Candidate status, production deployment, or Stable qualification.
