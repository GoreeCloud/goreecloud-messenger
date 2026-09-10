# GoreeCloud Messenger — Repository Specifications

Status: Development  
Canonical project record: `GoreeCloud/Projects/Project Specification — Messenger`  
Repository: `GoreeCloud/goreecloud-messenger`

## Product boundary

GoreeCloud Messenger is the native GoreeCloud messaging and calling application/service. It owns GoreeCloud Data messaging, conversation/message/receipt/attachment contracts, and approved client experiences. SMS, MMS, RCS, and carrier calling remain technically distinct transports and must never be mislabeled as GoreeCloud Data E2EE.

## Current implemented foundation

- Native Go domain/service/API layers for Data conversations and encrypted message envelopes.
- Authenticated sender and conversation-membership enforcement.
- Deterministic duplicate/idempotency and nonce-reuse protections.
- Authenticated delivery/read receipt contracts with monotonic state.
- Opaque encrypted attachment submission, authorized JSON/base64 fetch, metadata listing, replay-safe deletion, and raw ciphertext download.
- Local Development persistence abstractions and focused tests.
- Native Android Development client foundations with explicit communication-provenance/readiness boundaries and a repository-local GLAZE UI V1.3 Stable source mapping.
- Android Data-send readiness keeps GoreeCloud Identity authentication, conversation authorization, Data transport availability, and verified active E2EE as independent authorities. Conversation authorization and E2EE scopes must each be canonical bounded opaque identifiers and must identify the same exact conversation before `Data · E2EE` readiness can be projected.
- A protocol-neutral Android E2EE acceptance projection now prevents a provider's bare `E2EE_ACTIVE` enum from reaching readiness unless the same provider evidence reports accepted implementation security review, enrolled local cryptographic device identity, established conversation-scoped session, current key lifecycle, and the exact canonical conversation scope. Missing, negative, contradictory, noncanonical, or mismatched evidence fails closed to an unverified cryptographic state.
- The current Android source guard prohibits network/local-persistence authority and additionally prohibits direct `java.security`/`javax.crypto` implementation markers for this acceptance-only Development slice.
- A Development-only, explicitly composed Messenger consumer boundary for GoreeCloud Identity exact-handle resolution. It accepts only an exact handle from an authenticated Messenger user, preserves one uniform unresolved state, returns only opaque Identity subject/canonical handle/optional display name on success, and contains no production Identity network client or service credential.

## Identity and discovery

Messenger must use GoreeCloud Identity for account/session authority and consumer username resolution. Username discovery must not create a Messenger-owned browsable account directory.

The current Development consumer boundary is pinned to GoreeCloud Identity Draft PR #5 (`agent/native-directory-contract`) at exact head `5904a44997b90cc44f5d196620bb7e189fea0eeb`, contract `goreecloud-identity.consumer-directory.v1`. That upstream contract remains Draft/unmerged and does not establish production authority.

Messenger's current optional `POST /v1/identity/resolve` application endpoint accepts only `{handle}` after Messenger-user authentication and delegates the supplied handle unchanged to an injected `IdentityDirectoryResolver`. Messenger does not accept a client-supplied requester-service identity and does not own Identity canonicalization, discoverability, or per-service disclosure policy. The future resolver must derive Messenger's verified service principal from an accepted trusted runtime/service-authentication configuration.

Resolved provider output is limited to opaque Identity subject, canonical lowercase handle, and optional display name. Nonexistent, private, and service-disclosure-unauthorized accounts remain one indistinguishable unresolved result. Prefix search, fuzzy search, directory browsing, account enumeration, phone/email lookup, and administrative listing are outside this contract.

Invitation, conversation membership, contact, messaging, block/report, and other Messenger-specific authorization remain Messenger responsibilities and are not granted by a successful Identity resolution.

## End-to-end cryptography authority

Messenger must use an established, security-reviewed cryptographic protocol or implementation. It must not invent an unreviewed proprietary cryptographic algorithm or treat an application enum, network transport, or test fixture as cryptographic proof.

No authoritative GoreeCloud record currently selects Signal Protocol, MLS, Double Ratchet, libsignal, or another concrete Messenger E2EE protocol/implementation. The current FR-005 Android slice therefore remains deliberately protocol-neutral.

The future responsible `E2EESessionAuthority` must internally own or consume the real cryptographic evidence necessary to determine protection state. Before it may expose `E2EE_ACTIVE` into Messenger send readiness, the minimized Development projection currently requires:

- accepted implementation review state;
- enrolled local cryptographic device identity;
- established conversation-scoped cryptographic session;
- current key lifecycle state, including required rotation/key-change handling; and
- the exact canonical bounded opaque conversation identifier.

The minimized projection contains no private/public keys, session secrets, device secrets, reusable credentials, key fingerprints, algorithm/cipher-suite identifiers, ciphertext transformation, or network endpoint. Group and multi-device implementations must resolve their complete participant/device state before returning this projection.

Positive test fixtures establish policy behavior only. They do not establish an actual protocol implementation, device enrollment, session establishment, key lifecycle, security review, Wardveil Security acceptance, or production E2EE.

The complete production architecture still must address forward secrecy where supported, key rotation/change handling, device verification, group security, attachment protection, multi-device enrollment/revocation, protected local state, and backup/recovery without weakening E2EE.

## Security and privacy requirements

- The service must not decrypt GoreeCloud Data message or attachment ciphertext.
- Encryption state must be represented only when verified by the applicable client/session protocol and accepted cryptographic authority.
- Encrypted conversations must not silently downgrade to SMS/MMS.
- Attachment raw-byte transport must remain generic binary with no content sniffing and no server-side plaintext MIME interpretation.
- Identity resolution must preserve uniform privacy-sensitive negative results and minimized successful disclosure; upstream provider errors and policy details must not be exposed to clients.
- Identity service identity must be derived from a trusted service-authentication authority, never from client-supplied request fields.
- Wardveil Security, Privacy Shield, Everkeep, GoreeCloud Mesh, and GoreeCloud Identity integration are required where applicable.
- Client surfaces must track the **current approved Stable GLAZE UI release**. The current target is GLAZE UI V1.3 / `1.3.0` — Adaptive Resonance, at exact Stable integration revision `fc7cc91d2eace8da2371371c2855c24cbcb326a1`, with `1.2.0` retained as the rollback baseline.
- Client presentation must not manufacture security, privacy, encryption, delivery, authorization, or transport state through decorative material or semantic color. Visible state must remain derived from the responsible technical authority.
- Repository-local token/source tests are not rendered application acceptance. Exact-revision visual, accessibility, adaptive/form-factor, localization/directionality, representative-device, performance, rollback, and other applicable downstream evidence remain required before Stable promotion.
- Authority-owned conversation identifiers are opaque and exact. Leading/trailing whitespace, overlong values, C0 controls, and DEL are noncanonical and fail closed; valid internal whitespace and punctuation remain part of the identifier and are not normalized away.

## Current acceptance boundary

This is not production-ready. Production-grade Identity sessions/device identity, accepted service-principal authentication for the consumer directory, a live Identity resolver/transport, selection and security review of a concrete E2EE implementation, real cryptographic device/session/key lifecycle, multi-device synchronization, distributed message/object persistence, push delivery, abuse controls/rate limiting, carrier adapters, calling media infrastructure, client packaging, complete GLAZE UI V1.3 downstream acceptance, and deployment acceptance remain incomplete.

A Development exact-handle resolution boundary does not complete FR-004. Production Identity credential/session/device authority, accepted consumer-directory deployment and service authentication, production username-resolution transport, Privacy Shield/Wardveil acceptance, and application-specific authorization still require separate evidence.

The protocol-neutral E2EE acceptance projection does not complete FR-005. It strengthens what a future provider must prove before Messenger can accept an active-E2EE claim, but there is still no selected/implemented security-reviewed protocol, production device identity/key storage, real session establishment, verification/rotation/group/multi-device operation, or production cryptographic authority.
