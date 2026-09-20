# GoreeCloud Messenger — Development Notes

## Current stabilization context

- Repository lifecycle remains **Development**, overall platform conformance remains **nonconformant**, and Messenger is not Stable or production accepted.
- Current authoritative main before this candidate is `96aabdd7baf257821a3b9b3288c0e577b1caec88`, including the disconnected native Android Development client, exact-source Android CI hardening, accessibility-heading semantics, immutable Foundation CI, non-interactive runtime guard, PR #86 Identity session/device authority hardening, PR #89 exact-conversation authorization hardening, PR #91 current deployment/freshness-bound Data transport readiness, and PR #92 documentation reconciliation.
- PR #79 exact source head `3cd0066c47d5c2779f4aaa231b1d220a19b0c337` passed the Messenger Foundation, Platform Contract 0.4, Android client build/evidence, and Android 16 emulator workflows before merge.
- The Android client is therefore authoritative Development source, but its disconnected/fail-closed behavior remains intentional and does not establish production messaging capability.
- Draft PR #76 and its stacked Glaze 1.4.1 lineage remain historical provenance only and must not be treated as current integration authority.

## Active stabilization observations

- The native Android client has no production account/session binding, connected message transport, durable client message storage, active production E2EE lifecycle, carrier authority, or live Send control.
- Client integration must not activate message transport, Identity, E2EE, privacy, security, recovery, policy, observability, synchronization, or delivery authority merely because UI source is present.
- The repository now declares accepted Platform Contract 0.4 structure and all nine Integral Platform Systems, while unresolved runtime systems remain explicitly blocked.
- This V1.6 source-migration branch maps Android presentation source to exact Official Stable GLAZE UI V1.6 / 1.6.0 authority while Messenger-local rendered, accessibility, representative-device/form-factor, localization/RTL, performance, rollback, and Human Visual Excellence acceptance remain incomplete.
- Production Identity/session/device integration, accepted cryptographic session/key lifecycle, Wardveil, Privacy Shield, Everkeep, Mesh, Manager, GoreeCloud Policy, GoreeCloud Observability, protected signing/deployment, and representative target-environment evidence remain incomplete.

## Maintenance notes

Continue from authoritative `main`; do not revive or merge the stale stacked Android lineage unchanged. Preserve the disconnected/fail-closed client until authenticated transport, exact conversation ownership, cryptographic, privacy, security, policy, observability, recovery, and identity boundaries are accepted.

Any material Android source change requires fresh exact-head validation. Existing PR #79 evidence proves only the merged Development foundation at its exact source head; it does not establish future-revision, production, Release Candidate, or Stable acceptance.


## September 18 CI stabilization slice

- Android client CI now targets Ubuntu 24.04 rather than the moving `ubuntu-latest` label.
- Both Android jobs verify that the checked-out revision exactly matches the pull-request head or pushed `main` SHA before executing source.
- The build job also assembles the release variant as compilation/package-shape evidence without retaining or promoting it as a production artifact.
- These controls strengthen exact-source and build-variant evidence only. They do not activate network transport, Identity, E2EE, durable client state, Send authority, production signing, Release Candidate, or Stable status.


## Accessibility stabilization — September 18, 2026

PR #83 is integrated on `main`. The app title and major readiness/provenance/platform section headings use Android accessibility-heading semantics on API 28+, and Android 16 runtime acceptance verifies them. This does not widen messaging or platform authority.


## Foundation CI supply-chain stabilization — September 18, 2026

- PR #84 is integrated on `main`; Messenger Foundation uses Ubuntu 24.04 plus immutable checkout and setup-go commit SHAs while preserving the accepted action major versions.
- Foundation checkout now targets and verifies the exact pull-request head or pushed main revision before Go source is executed.
- This changes CI provenance only; it does not activate Data transport, Identity, E2EE, durable state, Send authority, production signing, RC, or Stable status.


## Disconnected runtime authority stabilization — September 19, 2026

- Android 16 runtime acceptance now treats the current disconnected shell as explicitly read-only: no descendant view may be clickable or long-clickable.
- The existing no-network/no-live-communication permission boundary and exact `Send`-control absence check remain in force.
- This protects the Development shell from silently acquiring user-triggered messaging authority before accepted Identity/session/device binding, exact conversation authorization, Data transport, E2EE lifecycle, durable state, Privacy Shield, Wardveil Security, Everkeep, Mesh, Manager, Policy, and Observability integration exist.
- The change does not connect Messenger, add a transport, create a Send action, or satisfy issue #78 production/release blockers.


## Identity session/device authority boundary — September 19, 2026

- PR #86 is integrated on `main` and strengthens the future GoreeCloud Identity seam so a bare authenticated state is no longer sufficient for Data messaging readiness.
- A positive Identity projection now requires both session binding and device binding to be independently reported as bound by the responsible Identity authority.
- The minimized boundary carries no principal, session, device, token, credential, or secret identifiers.
- Unit tests verify that missing session binding, missing device binding, provider failure, and explicit unauthenticated state all fail closed before the send coordinator may invoke an injected transport.
- This is Development authority-contract hardening only. The disconnected Android shell still has no production Identity adapter, network authority, Send control, durable client state, or accepted production messaging capability.

## Exact-conversation authorization stabilization — September 19, 2026

- PR #89 is integrated on `main` and strengthens the future conversation authorization seam used by Android Data messaging readiness.
- A bare VERIFIED_PARTICIPANT claim is no longer sufficient. Positive readiness also requires the responsible authorization authority to bind its decision to the current accepted Identity session/device authority, report the decision as current, and identify the exact canonical conversation scope requested by the client.
- The minimized boundary carries no participant list, principal identifier, ACL, token, credential, or reusable authorization material.
- Unit coverage fails closed for bare participant claims, missing Identity binding, stale authorization decisions, and mismatched conversation scope.
- This does not create a production authorization adapter, authenticate a user, connect Data transport, expose Send, establish E2EE, persist messages, or establish production, Release Candidate, or Stable authority.

- Superseded candidate head `f73139adf5d5d7eb0a6e67cdaee8fe23a24d34e6` failed the Android client authority-boundary guard because the required literal conversation-authority call was split across lines. The authorization logic itself had not executed; the corrected candidate preserves the strengthened fail-closed projection while restoring the repository guard-visible call shape.

- Superseded candidate head `09663a64bd9aaf66f09bfe731ac424b4189742f2` cleared the authority-boundary guard but failed unit tests because existing send-coordinator fixtures still modeled a bare positive conversation claim. The corrected candidate updates only accepted test fixtures to carry the new authorization acceptance facts and drops unaccepted authorization scope before readiness evaluation so rejected scope cannot manufacture an unrelated E2EE mismatch reason.

- Superseded candidate head `e340cbdbf30548a06e1075e569050845795d7fb6` reduced the remaining failures to two stale expectations: one transport-isolation test still used a bare participant claim, and one mismatched-conversation test expected only the cryptographic gate to reject a scope mismatch. The corrected candidate uses accepted authorization evidence when testing transport isolation and expects both independent conversation-scoped gates to fail closed when both scopes mismatch.


## Data transport current-authority hardening — September 19, 2026

- PR #91 is integrated on `main` and strengthens the existing accepted-transport projection without adding a live network adapter or widening the disconnected Development shell.
- A positive Data transport state now requires accepted configuration, Identity/authentication binding, protected-channel behavior, bounded failure policy, **current deployment/environment binding**, and **fresh acceptance evidence**.
- Missing or rejected deployment binding or stale transport acceptance fails closed to `UNKNOWN`; coordinator tests verify that neither condition can invoke the injected send transport.
- No endpoint, credential, token, certificate, Android INTERNET permission, durable queue, retry engine, Send control, production adapter, deployment, Release Candidate, or Stable authority is created.
- Current Official Stable GLAZE UI V1.6 / 1.6.0 remains the required authority; the new source-migration candidate maps Android presentation source to that exact release while keeping application acceptance blocked.

- Accepted exact-head evidence for PR #91: Platform Contract run `35490212282`, Messenger Foundation run `35490212048`, and Messenger Android Client run `35490212049` all succeeded before squash merge to authoritative main commit `fa0417fcf7274fb611981af768d35b8555908c53`.


## Android GLAZE UI V1.6 source migration candidate — September 20, 2026

- This Development branch maps the disconnected Android presentation source from V1.5.1 to exact Official Stable GLAZE UI V1.6 / 1.6.0 release source `a7180679ea851389e0f3004515f9a25f420e716d`.
- Shared V1.6 consumer eligibility is not rebound as Messenger acceptance. The Platform Contract Glaze result becomes applicable-blocked, overall conformance stays nonconformant, and all Messenger-local rendered/device/release gates remain open.
- A pure presentation resolver adds Reduced Transparency, Reduced Motion, performance-cost, large-text, focus, and conservative target semantics without accepting message, conversation, Identity, E2EE, transport, privacy, security, recovery, policy, observability, or remote-content inputs.
- The Activity consumes only a neutral resolver result for its existing minimum interaction target. No INTERNET permission, connected transport, Send control, account/session provider, durable message storage, cryptographic session, telemetry, or optical adapter is activated.
- Retained neutral pigments and layout conveniences are explicitly Messenger-owned rather than represented as exact shared tokens.
- Issue #78 remains open for production communication authorities and full product acceptance.
