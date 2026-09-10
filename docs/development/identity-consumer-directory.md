# GoreeCloud Messenger — Identity Consumer Directory Development Boundary

Status: Development contract-consumer candidate  
Messenger branch: `feature/identity-handle-resolution-boundary`  
Upstream authority candidate: GoreeCloud Identity Draft PR #5 (`agent/native-directory-contract`)  
Pinned Identity contract revision: `5904a44997b90cc44f5d196620bb7e189fea0eeb`  
Pinned contract: `contracts/consumer-directory.v1.json` (`goreecloud-identity.consumer-directory.v1`)

## Purpose

Messenger requires GoreeCloud Identity-owned username resolution without creating a Messenger-owned account directory. This Development slice defines the Messenger-side consumer boundary for that future integration while deliberately stopping before production service-to-service authentication or live Identity networking.

The upstream Identity contract is itself Draft/Development and unmerged. Its current semantics are therefore a pinned Development dependency, not production authority.

## Responsibility split

GoreeCloud Identity remains responsible for:

- canonical handle rules and normalization;
- account discoverability;
- disclosure authorization for the verified requesting service;
- mapping an Identity account to the minimized consumer projection; and
- making nonexistent, private, and service-unauthorized accounts indistinguishable through one unresolved result.

Messenger remains responsible for:

- authenticating the Messenger user before exposing the application-facing resolution endpoint;
- forwarding only the supplied exact handle to the injected Identity resolver;
- validating the minimized provider projection before returning it to a Messenger client;
- invitation, conversation membership, contact, messaging, blocking, and other Messenger-specific authorization; and
- refusing to infer account existence or disclosure policy from an unresolved result.

## Messenger-facing Development endpoint

When and only when `DataRuntimeHandler.WithIdentityDirectory(...)` is explicitly composed, Messenger exposes:

`POST /v1/identity/resolve`

The request accepts exactly one field:

```json
{"handle":"@example"}
```

The request is bounded to 4096 bytes and unknown JSON fields are rejected. A client cannot provide a requester-service name, service principal, service credential, disclosure policy, user identifier, email address, or phone number.

The existing Messenger `Authenticator` must authenticate the requesting Messenger user before the resolver is called. This user authentication is separate from the future service-to-service authentication used by the injected resolver when it eventually communicates with GoreeCloud Identity.

## Resolver boundary

`IdentityDirectoryResolver.ResolveExact(ctx, handle)` is an injected service boundary. A future implementation must derive and authenticate Messenger's service identity from trusted runtime configuration or an approved service-authentication mechanism. The requester service is never accepted from client input.

Messenger deliberately does not normalize the supplied handle before calling the resolver. The Identity authority owns leading-`@`, lowercase, and canonical-handle behavior. A successful provider result is validated as a minimized projection containing only:

- opaque Identity subject;
- canonical lowercase handle; and
- optional display name.

No email, phone number, discoverability flag, service-disclosure policy, credentials, roles, session information, or Identity administration state belongs in this projection.

## Privacy-preserving negative result

The resolver returns one unresolved state for all privacy-sensitive negative outcomes. Messenger maps that state to the same HTTP 404 body:

```json
{"error":{"code":"not_resolved","message":"The requested identity could not be resolved."}}
```

Messenger must not distinguish nonexistent accounts from private accounts or accounts that are not disclosable to Messenger. Malformed client requests may still return a bounded 400 because the requester already knows the malformed input it submitted.

Unexpected resolver failures and invalid provider projections fail closed to a generic 503 response. Upstream topology, credentials, policy details, or backend exceptions are not returned to clients.

All JSON responses use the repository's private-response hardening, including `Cache-Control: no-store`.

## No directory browsing

This slice provides exact-handle resolution only. It does not add:

- prefix search;
- fuzzy search;
- directory browsing;
- account enumeration;
- administrative user listing;
- contact graph discovery; or
- phone/email lookup.

Those capabilities must not be inferred from the existence of the resolver interface.

## Runtime composition

The base `DataRuntimeHandler` does not expose the Identity route. The route exists only after an `IdentityDirectoryService` is deliberately supplied through `WithIdentityDirectory(...)`.

The repository currently supplies no production `IdentityDirectoryResolver`, HTTP client, service credential, OAuth/OIDC service flow, endpoint URL, retry policy, circuit breaker, deployment configuration, or production service principal. This prevents source-level contract work from masquerading as a connected Identity integration.

## Security and abuse boundary

The endpoint is authenticated and privacy-minimized, but this Development slice does not establish production Wardveil Security abuse protection or distributed rate limiting. Before production-family acceptance, the connected implementation requires approved abuse controls, request-volume policy, telemetry minimization, operational monitoring, and service-authentication handling appropriate to account-discovery risk.

## Acceptance boundary

Passing Messenger unit/API tests or repository CI proves only the Development consumer boundary. Production acceptance still requires, at minimum:

- an accepted GoreeCloud Identity consumer-directory contract and deployment;
- an approved Messenger service-principal authentication profile;
- a production resolver/transport implementation;
- production Messenger Identity session/device integration where required;
- Privacy Shield review of discovery/disclosure behavior;
- Wardveil Security abuse/rate-limit acceptance;
- application-specific invitation/membership authorization;
- deployment, rollback, monitoring, and recovery evidence; and
- exact-candidate validation under the applicable release process.

Until those gates are satisfied, FR-004 remains incomplete and GoreeCloud Messenger remains Development.
