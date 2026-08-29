# GoreeCloud Messenger Typing Indicator Boundary

## Purpose

Typing indicators are short-lived conversation-presence metadata for GoreeCloud Data conversations. They are not messages, drafts, delivery receipts, or proof that a message will be sent.

## Current Development contract

- An authenticated conversation participant may publish `typing` or `idle` state for their own GoreeCloud identity.
- The service verifies current conversation membership before accepting a signal.
- A separate `TypingPrivacyPolicy` boundary controls whether a participant may publish typing state and whether a participant may observe active typing state.
- Signals carry a monotonically increasing per-conversation/per-user sequence. Equal or older sequences are rejected so delayed device traffic cannot resurrect stale typing state after a newer update.
- `typing` creates a short-lived active projection with a server-assigned 10-second expiration.
- `idle` removes the active projection immediately while preserving the latest sequence needed to reject stale updates.
- Authorized reads return only other currently active participants, their sequence, and expiration time. The observer's own indicator is omitted.

## Privacy boundary

The typing surface does not accept or retain draft text, keystrokes, cursor positions, message ciphertext, attachment content, client nonces, device secrets, or user-supplied timestamps. It does not attempt to infer what a participant is composing.

The Development memory store keeps only the latest sequence required for ordering and the currently active short-lived projection. Expired active entries are removed during reads. This is a source-level minimization contract, not evidence of production telemetry, logging, retention, or distributed-presence behavior.

## HTTP boundary

- `POST /v1/data/conversations/{conversationID}/typing` accepts `user_id`, `sequence`, and `state`.
- `GET /v1/data/conversations/{conversationID}/typing` returns active typing projections for authorized observers.

Requests use strict JSON decoding and bounded bodies. Authorization and privacy failures are bounded non-success responses. Active responses are non-cacheable through the shared JSON response protections.

## Production boundary

This foundation does not establish production presence infrastructure, WebSocket or push fan-out, distributed expiry coordination, offline/multi-device convergence, cross-device sequence allocation, production GoreeCloud Identity sessions, production Privacy Shield preference storage, abuse/rate controls, client rendering, deployment, release, or Stable acceptance. Typing indicators remain Development source until those separate gates are satisfied.
