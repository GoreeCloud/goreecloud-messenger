# Typing publish limiter capacity

Status: Development source only.

The content-free typing publish guard now bounds the number of distinct active per-conversation/per-participant rate reservations retained inside one Messenger process.

## Capacity rule

- active reservation ceiling: 4096 distinct keys;
- reservations older than the existing 10-second typing indicator TTL are removed before the ceiling is evaluated;
- refreshing an already-present key after the 250 ms minimum interval remains allowed at capacity because it does not increase retained state;
- a new distinct key is rejected when the active ceiling is full;
- the rejection reuses the existing `ErrTypingRateLimited` boundary and therefore the existing low-detail HTTP 429 response;
- successful idle still removes its matching reservation immediately; and
- a failed typing-store mutation still rolls back its matching reservation when that reservation is current.

## Privacy boundary

The limiter retains only the same ephemeral conversation/participant key material and timestamps already required by the Development anti-flood guard. It does not retain draft text, characters, cursor state, message ciphertext, client timestamps, device telemetry, clipboard data, or editor context.

The capacity rejection does not reveal the current reservation count or identify another conversation/participant responsible for capacity pressure.

## Distributed-runtime boundary

This remains process-local abuse protection. It is not a distributed quota, cross-instance rate limiter, production denial-of-service defense, Wardveil Security acceptance, or Privacy Shield acceptance. A multi-instance production topology requires a separately reviewed shared abuse-control design that preserves privacy and availability semantics.

## Acceptance boundary

This hardening does not establish production GoreeCloud Identity, cryptographic session/key lifecycle, distributed persistence/delivery, backup/restore, Mesh integration, Glaze UI client acceptance, deployment, release, or Stable qualification.
