# GoreeCloud Messenger — Development Notes

## Current stabilization context

- Repository lifecycle remains Development and is not Stable or production accepted.
- Verified baseline for this notes change: `main` at `1273a57d08acb375a38127df2e2fee193d0dbe65`.
- Authoritative `main` at the migration start is server-first. This topic branch cleanly restacks the historical disconnected Android shell onto current `main`; it is not authoritative client state until governed integration completes.
- This migration branch updates the platform declaration to accepted Contract 0.4/nine-system semantics and current Glaze UI 1.5.1 source authority while keeping overall conformance nonconformant.
- Draft PR #76 remains historical stacked provenance only. The current migration branch copies the reviewed client subtree onto current `main`, updates the Android 16 toolchain and current authorities, and must obtain fresh exact-head validation rather than reusing #76's green results.

## Active stabilization observations

- A native Android Development candidate is now present on this migration branch; connected messaging, representative-device acceptance, production signing/deployment, and release acceptance remain separate required states.
- Client integration must not activate message transport, Identity, E2EE, privacy, security, or synchronization authority merely because UI source exists.
- Production Identity/session/device integration, accepted cryptographic session/key lifecycle, Wardveil, Privacy Shield, Everkeep, Mesh, Manager, GoreeCloud Policy, GoreeCloud Observability, application-specific Glaze acceptance, signing/deployment, and representative target-environment evidence remain incomplete.

## Maintenance notes

Keep this clean restack on current authoritative contracts rather than merging the stale stack unchanged. Preserve the disconnected/fail-closed client until authenticated transport, exact conversation ownership, cryptographic, privacy, security, policy, observability, recovery, and identity boundaries are accepted. Fresh exact-head Android/emulator/contract evidence is required before integration; production and Stable claims remain blocked.
