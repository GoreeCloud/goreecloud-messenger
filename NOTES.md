# GoreeCloud Messenger — Development Notes

## Current stabilization context

- Repository lifecycle remains Development and is not Stable or production accepted.
- Verified baseline for this notes change: `main` at `1273a57d08acb375a38127df2e2fee193d0dbe65`.
- Authoritative `main` is currently server-first and does not contain an accepted first-party Android Messenger client.
- The current platform declaration is stale relative to the governed Glaze UI 1.5.1 consumer authority and remains nonconformant.
- Draft PR #76 contains a disconnected Android Development shell with green historical exact-head Android/Foundation checks, but the PR explicitly records an integration blocker and targets Glaze UI 1.4.1 on a stacked parent. It must not be treated as current integration authority.

## Active stabilization observations

- A native Android client remains a required delivery target.
- Client integration must not activate message transport, Identity, E2EE, privacy, security, or synchronization authority merely because UI source exists.
- Production Identity/session/device integration, accepted cryptographic session/key lifecycle, Wardveil, Privacy Shield, Everkeep, Mesh, Manager, current Glaze acceptance, signing/deployment, and representative target-environment evidence remain incomplete.

## Maintenance notes

Resume Android client work from a current authoritative base and current governed platform/Glaze contracts rather than merging the stale stack unchanged. Keep Android client state explicitly disconnected/fail-closed until authenticated transport, ownership, cryptographic, privacy, and security boundaries are accepted.
