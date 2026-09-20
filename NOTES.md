# GoreeCloud Messenger — Development Notes

## Current stabilization context

- Repository lifecycle remains **Development**, overall platform conformance remains **nonconformant**, and Messenger is not Stable or production accepted.
- Authoritative `main` is `edb69de4a9f7c71f298dd98f8e000d05881fe711`, including the disconnected native Android Development client, exact-source Android CI hardening from PR #82, accessibility-heading semantics from PR #83, immutable Messenger Foundation CI dependencies from PR #84, the non-interactive disconnected-shell runtime guard from PR #85, and the integrated PR #86 Identity session/device authority hardening.
- PR #79 exact source head `3cd0066c47d5c2779f4aaa231b1d220a19b0c337` passed the Messenger Foundation, Platform Contract 0.4, Android client build/evidence, and Android 16 emulator workflows before merge.
- The Android client is therefore authoritative Development source, but its disconnected/fail-closed behavior remains intentional and does not establish production messaging capability.
- Draft PR #76 and its stacked Glaze 1.4.1 lineage remain historical provenance only and must not be treated as current integration authority.

## Active stabilization observations

- The native Android client has no production account/session binding, connected message transport, durable client message storage, active production E2EE lifecycle, carrier authority, or live Send control.
- Client integration must not activate message transport, Identity, E2EE, privacy, security, recovery, policy, observability, synchronization, or delivery authority merely because UI source is present.
- The repository now declares accepted Platform Contract 0.4 structure and all nine Integral Platform Systems, while unresolved runtime systems remain explicitly blocked.
- Current Official Stable Glaze UI V1.5 / 1.5.1 is source-mapped, but Messenger-local rendered, accessibility, representative-device/form-factor, localization/RTL, performance, rollback, and Human Visual Excellence acceptance remain incomplete.
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
