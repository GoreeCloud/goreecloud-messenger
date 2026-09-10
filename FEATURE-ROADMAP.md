# GoreeCloud Messenger — Feature Roadmap

Status: Active roadmap control — as of September 10, 2026  
Drive authority: `GoreeCloud/Feature Roadmap/GoreeCloud Messenger/FEATURE-ROADMAP.docx`  
Authoritative project record: `GoreeCloud/Projects/Project Specification — Messenger`  
Canonical repository: `GoreeCloud/goreecloud-messenger`

## Purpose

This repository roadmap mirrors the Drive-side GoreeCloud Messenger feature-roadmap control. It records current planned and recommended feature obligations without replacing the authoritative project record, verified repository implementation evidence, release gates, or GoreeCloud Tasks Management.

This file and the corresponding Drive `FEATURE-ROADMAP.docx` must remain materially synchronized. A feature appearing here does not establish implementation, production acceptance, release acceptance, or Stable status by itself.

Current Development checkpoint: Draft PR #67 (`feature/e2ee-acceptance-explanations`) extended the protocol-neutral E2EE acceptance boundary from PR #66 with fail-closed explanation categories for inactive cryptography, wrong conversation scope, unaccepted implementation review, unenrolled cryptographic device identity, unestablished conversation-scoped session, and non-current key lifecycle. Exact PR #67 head `24f436d817e3d29671fbd7faccd87c45ad414834` passed Messenger Android Client run `34524705684` and Messenger Foundation run `34524705786`. Stacked Draft PR #68 (`feature/e2ee-acceptance-explanation-ui`) projects those existing failure states into bounded user-facing `Why — …` detail in the disconnected Android readiness UI without changing readiness evaluation. Exact PR #68 implementation head `76a7b6afe2a431f75c0b54eb16e4d724a008f95c` passed Messenger Foundation run `34542866611` and Messenger Android Client run `34542866700`. The shell still supplies UNKNOWN live authority, shows no Send control, and cannot create cryptographic, transport, key, account, or conversation authority. No authoritative GoreeCloud record currently selects or accepts a concrete Messenger E2EE protocol implementation; production Identity, Data transport, reviewed cryptography, Wardveil/Privacy Shield acceptance, deployment, release, and Stable authority remain pending.

## Roadmap

| ID | Feature / obligation | Priority | Current state |
| --- | --- | --- | --- |
| FR-001 | Reconcile and maintain every current planned or recommended GoreeCloud Messenger feature from the authoritative project record and verified repository evidence in this roadmap. | High | Ongoing control |
| FR-002 | Move actionable feature obligations into GoreeCloud Tasks Management when required, preserving priority, dependency, and lifecycle disposition. | High | Ongoing control |
| FR-003 | Do not mark features implemented, complete, cancelled, or superseded without authoritative evidence and synchronized repository/Drive roadmap updates. | High | Ongoing control |
| FR-004 | Integrate production GoreeCloud Identity sessions, device identity, authorization, and Identity-owned exact-handle username resolution for Messenger. | High | Development consumer boundary implemented; production authority and live integration pending |
| FR-005 | Establish and security-review the end-to-end cryptographic device, session, key, verification, rotation, and exact conversation-scoped E2EE authority required for production communication. | High | Development acceptance and explanation boundaries implemented through PR #67/#68; concrete reviewed cryptographic authority remains pending |
| FR-006 | Implement live GoreeCloud Data networking plus production message, receipt, attachment, delivery, push, offline, and distributed persistence foundations without weakening ciphertext opacity or authorization. | High | Planned — production transport/delivery pending |
| FR-007 | Advance the native Android client from disconnected readiness evidence to a real composer/send path only after Identity, exact conversation authorization, Data transport, and verified active E2EE authorities are independently accepted. | High | Development readiness UI now explains unverified E2EE without exposing Send; real composer/send remains blocked by FR-004 through FR-006 |
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
