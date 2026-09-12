# GoreeCloud Data transport acceptance boundary

## Status

Development acceptance evidence only. No live Android Data transport is implemented by this boundary.

## Purpose

The Android messaging authority resolver must not treat a bare transport `AVAILABLE` value as enough evidence to make GoreeCloud Data messaging ready.

The governing GoreeCloud production-readiness requirements require application/service integrations to use validated configuration, explicit authentication, protected network exposure, bounded timeouts and retries, and fail-closed degraded/error behavior. This Development contract carries only minimized acceptance facts for those requirements. It does not define an endpoint, credential, token, certificate, retry schedule, or network client.

## Required minimized evidence

A future responsible Data transport authority may project `DataTransportState.AVAILABLE` into Messenger readiness only when it independently supplies `ACCEPTED` for all of the following:

- **configuration** — required transport configuration has been validated and accepted by the responsible runtime;
- **authentication binding** — the transport operation is bound to the required authenticated authority rather than local UI state or caller-authored identity;
- **channel protection** — the responsible runtime accepts the network channel/protection boundary required for the operation;
- **failure policy** — timeout, retry, failure, and degraded-state behavior is bounded and accepted rather than unbounded or success-by-default.

If any fact is `UNKNOWN` or `NOT_ACCEPTED`, a bare `AVAILABLE` state is reduced to `UNKNOWN` for messaging readiness. `UNAVAILABLE` and other negative states are never upgraded by positive auxiliary facts.

## Authority separation

These acceptance projections do not create or replace:

- production GoreeCloud Identity sessions or credentials;
- conversation authorization;
- a GoreeCloud Data endpoint or service-discovery authority;
- TLS certificates or secret material;
- Android Internet permission or a networking library;
- E2EE implementation, device keys, sessions, or key lifecycle;
- message persistence, delivery, push, offline queues, retries, or synchronization;
- Wardveil Security, Privacy Shield, Everkeep, Mesh, or Manager production acceptance; or
- deployment, Release Candidate, production-release, or Stable authority.

The Android Development source guard continues to reject network implementation authority, Internet permission, local persistent storage authority, and local cryptographic implementation code.

## Testing boundary

Unit-test fixtures may construct fully accepted transport evidence so the existing coordinator can prove that all independently owned gates compose correctly before an injected mock transport is invoked. Those fixtures are test evidence only and do not establish a production provider.

Regression coverage must prove that bare `AVAILABLE`, rejected configuration, rejected authentication binding, rejected channel protection, rejected failure policy, provider exceptions, and unavailable transport all fail closed without calling the injected send transport.

## Current product state

Messenger remains Development. FR-006 is not complete: there is still no production Android GoreeCloud Data network adapter, accepted production endpoint/configuration, production credential binding, production transport security acceptance, delivery/push/offline implementation, or production persistence/synchronization authority.
