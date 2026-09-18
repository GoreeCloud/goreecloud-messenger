# GoreeCloud Messenger — Development Notes

## Current stabilization context

- Repository lifecycle remains **Development**, overall platform conformance remains **nonconformant**, and Messenger is not Stable or production accepted.
- Authoritative `main` includes the disconnected native Android Development client through merge `08ea079527960e6c3ac7e5f1f4234236aef6d508` (PR #79).
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
