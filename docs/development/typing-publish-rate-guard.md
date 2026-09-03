# Typing Publish Rate Guard — Development Boundary

## Status

Development source hardening only. This checkpoint does **not** establish production abuse prevention, Wardveil Security acceptance, Privacy Shield acceptance, distributed rate limiting, deployment, or Stable qualification.

## Purpose

GoreeCloud Messenger typing presence is intentionally content-free and ephemeral, but an authenticated participant could still generate unnecessary server work and observer churn by publishing repeated `typing` signals at excessive frequency.

The current guard adds a bounded per-conversation/per-participant publication interval to the existing `TypingService` while preserving the ability to stop typing immediately.

Current Development constants:

- typing indicator TTL: 10 seconds;
- minimum interval between accepted repeated `typing` publications for the same conversation/participant: 250 milliseconds.

These are source-level Development defaults, not a release-level service guarantee.

## Authorization order

The service continues to validate the signal and authenticate/authorize the acting participant before the rate guard is relevant. Publication still requires:

1. a nonempty authenticated user identity supplied by the runtime;
2. a valid content-free `TypingSignal`;
3. exact authenticated-user/signal-user equality;
4. conversation participation; and
5. an allowed typing-publication privacy policy.

The guard does not replace any of those checks and does not turn rejected outsiders or privacy-denied users into rate-limiter state.

## Behavior

For `typing` state:

- the first authorized publication reserves the participant/conversation rate slot;
- another authorized `typing` signal before the 250 ms interval completes fails with `ErrTypingRateLimited`;
- the HTTP adapter maps that sentinel to HTTP 429 Too Many Requests;
- a signal accepted at or after the interval is allowed to proceed to the existing store/sequence boundary.

For `idle` state:

- the rate guard never delays the stop signal;
- after the store successfully accepts the idle sequence, the participant/conversation reservation is cleared;
- a genuine subsequent typing restart can therefore be published immediately.

## Store-failure rollback

A typing publication reserves its rate slot before the store mutation so concurrent publications cannot race through the interval check. If the underlying store then rejects the signal—for example because its sequence is stale—the service removes that reservation when it still corresponds to the failed attempt.

This prevents a rejected signal from incorrectly throttling the next valid signal.

## Memory boundary

The Development guard retains only the last accepted/reserved publication timestamp keyed by conversation and participant. It does not retain typed characters, message content, composing text, editor state, or typing duration history.

Entries older than the existing 10-second typing TTL are opportunistically removed when another typing publication is evaluated. Successful idle also removes the matching entry.

This is in-process ephemeral state. It is not persisted, exported, backed up, synchronized, or included in Messenger data portability.

## Privacy boundary

The guard observes only metadata the service already requires for authorized typing presence:

- conversation ID;
- authenticated participant ID;
- typing-versus-idle state; and
- service clock time.

It does not add message-body inspection, keystroke observation, client editor context, device telemetry, IP-derived identity, or cross-conversation behavioral profiling.

Typing presence remains excluded from durable user-data exports.

## Distributed-runtime limitation

The current guard is process-local. Multiple Messenger service instances would each maintain independent limiter state unless a future accepted distributed abuse-control design provides shared authority.

Therefore this checkpoint must not be described as complete production flood protection or a distributed rate-limit guarantee.

## HTTP boundary

`ErrTypingRateLimited` maps to HTTP 429. The response remains a minimized generic error and does not reveal other participant state, privacy settings, server topology, or limiter contents.

This checkpoint does not add a `Retry-After` contract. A future public API compatibility decision may add one only after rate-policy ownership and release behavior are defined.

## Validation

Source tests cover:

- rapid repeated `typing` publication rejection;
- immediate `idle` acceptance despite a recent typing signal;
- immediate genuine restart after a successful idle;
- acceptance exactly at the configured minimum interval;
- rollback of a reservation when the store rejects a stale typing sequence; and
- HTTP 429 mapping for the rate-limit sentinel.

## Acceptance boundary

Passing these tests proves only the in-process Development behavior at the tested revision. It does not establish production load testing, distributed enforcement, Wardveil Security acceptance, Privacy Shield acceptance, service deployment, release acceptance, or Stable qualification.