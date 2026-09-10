# GoreeCloud Messenger — Feature Roadmap

Status: Active roadmap control — as of September 9, 2026  
Drive authority: `GoreeCloud/Feature Roadmap/GoreeCloud Messenger/FEATURE-ROADMAP.docx`  
Authoritative project record: `GoreeCloud/Projects/Project Specification — Messenger`  
Canonical repository: `GoreeCloud/goreecloud-messenger`

## Purpose

This repository roadmap mirrors the Drive-side GoreeCloud Messenger feature-roadmap control. It records current planned and recommended feature obligations without replacing the authoritative project record, verified repository implementation evidence, release gates, or GoreeCloud Tasks Management.

This file and the corresponding Drive `FEATURE-ROADMAP.docx` must remain materially synchronized. A feature appearing here does not establish implementation, production acceptance, release acceptance, or Stable status by itself.

Current Development checkpoint: Draft PR #63 (`feature/android-messaging-authority-providers`) remains the validated Android authority-provider parent at final exact head `5c03b8d198015b33dc89cea94c6679635700dd7d`. Draft PR #64 (`feature/identity-handle-resolution-boundary`) is stacked on that exact parent and advances FR-004 with an optional, authenticated Messenger consumer boundary for GoreeCloud Identity exact-handle resolution. It preserves Identity-owned canonicalization, discoverability, and service-disclosure policy; exposes only the minimized opaque subject, canonical handle, and optional display name projection; preserves one unresolved result for nonexistent/private/disclosure-unauthorized identities; rejects client-supplied requester-service authority; redacts resolver failures; and keeps the route absent unless explicitly composed. PR #64 is pinned to GoreeCloud Identity Draft PR #5 contract head `5904a44997b90cc44f5d196620bb7e189fea0eeb`, which remains Development and unmerged. Reconciled PR #64 exact head `eb5ca79f9a5f7f053aaaeaa1ac4404cb66570d97` passed Messenger Foundation run `34418188655`, Messenger Android Client run `34418188656` including Android 15 disconnected-shell emulator acceptance, and Platform Contract run `34418189166`. PR #64 remains Draft, open, mergeable, and unmerged. This is Development source/CI evidence only: it adds no live Identity endpoint/client, trusted production service-authentication profile, production session/device authority, production username resolver, Privacy Shield acceptance, Wardveil abuse/rate-limit acceptance, deployment, release, or Stable authority.

## Roadmap

| ID | Feature / obligation | Priority | Current state |
| --- | --- | --- | --- |
| FR-001 | Reconcile and maintain every current planned or recommended GoreeCloud Messenger feature from the authoritative project record and verified repository evidence in this roadmap. | High | Ongoing control |
| FR-002 | Move actionable feature obligations into GoreeCloud Tasks Management when required, preserving priority, dependency, and lifecycle disposition. | High | Ongoing control |
| FR-003 | Do not mark features implemented, complete, cancelled, or superseded without authoritative evidence and synchronized repository/Drive roadmap updates. | High | Ongoing control |
| FR-004 | Integrate production GoreeCloud Identity sessions, device identity, authorization, and Identity-owned exact-handle username resolution for Messenger. | High | Development consumer boundary implemented; production authority and live integration pending |
| FR-005 | Establish and security-review the end-to-end cryptographic device, session, key, verification, rotation, and exact conversation-scoped E2EE authority required for production communication. | High | Planned — cryptographic authority pending |
| FR-006 | Implement live GoreeCloud Data networking plus production message, receipt, attachment, delivery, push, offline, and distributed persistence foundations without weakening ciphertext opacity or authorization. | High | Planned — production transport/delivery pending |
| FR-007 | Advance the native Android client from disconnected readiness evidence to a real composer/send path only after Identity, exact conversation authorization, Data transport, and verified active E2EE authorities are independently accepted. | High | Planned — blocked by FR-004 through FR-006 |
| FR-008 | Complete GLAZE UI V1.3 / 1.3.0 Adaptive Resonance rendered, accessibility, adaptive/form-factor, representative-device, performance, rollback, and release acceptance for Messenger clients. | High | Planned — application acceptance pending |
| FR-009 | Complete applicable Wardveil Security, Privacy Shield, Everkeep, GoreeCloud Mesh, and GoreeCloud Manager integrations with evidence-backed production acceptance. | High | Planned — platform-system acceptance pending |
| FR-010 | Add SMS, MMS, and RCS adapters only through legitimate supported platform/carrier interfaces while preserving explicit transport provenance and no silent encrypted-Data downgrade. | Medium | Planned — platform/carrier capability dependent |
| FR-011 | Implement GoreeCloud Data voice and video calling signaling and media transport with applicable identity, E2EE, privacy, security, and provenance controls. | Medium | Planned |
| FR-012 | Implement authorized multi-device synchronization and encrypted backup/restore/recovery behavior without weakening E2EE or bypassing Everkeep and Privacy Shield requirements. | Medium | Planned |
| FR-013 | Continue evidence-gated modern conversation capabilities including replies, threads, reactions, edits, deletion, search, disappearing communication, media/document exchange, and approved group/community controls. | Medium | Planned / Development slices as separately evidenced |

## Maintenance and synchronization

Update both roadmap copies whenever feature scope, priority, dependency, implementation status, cancellation, supersession, recommendation, or verification state materially changes. No feature may be represented as complete or Stable solely because it appears in this roadmap; completion and lifecycle claims require the applicable authoritative implementation, validation, review, release, and production evidence.

## Reconciliation rule

At each material feature change, reconcile this roadmap against the current authoritative project record, repository implementation state, applicable platform-system requirements, the Drive-side roadmap, and GoreeCloud Tasks Management. Missing obligations, stale status, duplicated work, roadmap drift, or undocumented disposition changes are defects to correct.
