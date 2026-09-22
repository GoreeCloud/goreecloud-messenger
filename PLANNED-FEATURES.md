# GoreeCloud Messenger — Planned Features and Open Obligations

**Record type:** Repository planned/incomplete-feature inventory  
**Repository:** `GoreeCloud/messenger`  
**Lifecycle:** Development / nonconformant  
**Authority:** Current `main` source, accepted repository evidence, and active GoreeCloud Tasks Management obligations  
**Governing standard:** Standard — Repository Feature Tracking and Changelog Governance v1.0, effective September 22, 2026.

## Interpretation

Items here are planned, incomplete, blocked, or acceptance-gated. Their presence does not imply implementation or release readiness. Partial foundations that already exist are also described in `IMPLEMENTED-FEATURES.md` for the verified portion only.

## Current stabilization obligations

- Connect a production GoreeCloud Identity adapter and verify exact account/session/device authority.
- Establish producer-authoritative exact conversation authorization in the connected runtime.
- Connect and accept the GoreeCloud Data transport without weakening the current fail-closed no-network/no-Send boundary before readiness is proven.
- Establish accepted production E2EE session creation, verification, rotation, recovery, and multi-device lifecycle.
- Add durable message and encrypted-attachment persistence/object storage with verified authorization and recovery boundaries.
- Complete production push delivery, offline synchronization, presence fan-out, and production rate limiting.
- Replace Development-memory typing/privacy preference state with durable Privacy Shield-governed persistence.
- Add connected native client messaging persistence and production packaging only after the prerequisite authority chain is accepted.
- Complete GLAZE UI V1.6 rendered, accessibility, adaptive-device, performance, rollback, and Human Visual Excellence acceptance.
- Complete representative-device, recovery, signing/distribution, Release Candidate, production, and Stable qualification.

## Planned product capabilities

- Identity-owned exact-handle username resolution integration.
- SMS/MMS/RCS carrier/platform adapters where legitimate platform APIs permit.
- Voice/video call signaling and media transport.
- Native client typing/privacy presentation backed by accepted persistence.
- Production-grade notification and offline-delivery behavior.

## Integral Platform System obligations

Evaluate and satisfy or explicitly justify the applicable state of all nine Integral Platform Systems:

- GoreeCloud Manager
- Privacy Shield
- Wardveil Security
- Everkeep
- GLAZE UI
- GoreeCloud Mesh
- GoreeCloud Identity
- GoreeCloud Policy
- GoreeCloud Observability

GoreeCloud Sync remains separately governed and must not be represented as a tenth Integral Platform System.

## Explicit non-claims

Until accepted evidence exists, this file does not claim production messaging, production E2EE, live connected Identity/Data transport, production deployment, Release Candidate, or Stable status.

## Maintenance rule

Move an item to `IMPLEMENTED-FEATURES.md` only after the authoritative implementation and required verification are integrated. Record material lifecycle changes in `CHANGELOGS.md`. Keep actionable execution work in GoreeCloud Tasks Management without creating duplicate task authority.