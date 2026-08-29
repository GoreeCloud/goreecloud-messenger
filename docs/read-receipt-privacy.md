# Read-Receipt Privacy Boundary

## Development boundary

GoreeCloud Messenger delivery and read receipts are recipient-authenticated GoreeCloud Data metadata. This slice adds explicit privacy controls for `read` state without changing the existing delivery acknowledgement contract. It remains Development source and does not establish production Privacy Shield preference persistence, production receipt storage, multi-device convergence, release acceptance, or Stable status.

## State semantics

- `delivered` means an authorized recipient client reported application delivery progress.
- `read` means an authorized recipient client reported a later read observation.
- Neither state is carrier proof, cryptographic proof, or proof that a human understood message content.
- Receipt metadata contains no plaintext message body, encryption key, device secret, or carrier metadata.

## Privacy policy contract

`ReceiptPrivacyPolicy` has two independent decisions for one conversation and user:

- `CanPublishRead` controls whether that recipient may publish `read` state.
- `CanObserveRead` controls whether that observer may receive read projections for the conversation.

The default `NewReceiptService` compatibility constructor retains the previous allow-all Development behavior. Privacy-aware callers use `NewReceiptServiceWithPrivacy` and supply an explicit policy implementation. The in-memory policy implementation exists only for deterministic Development tests and is not a substitute for persistent Privacy Shield preferences.

## Publication behavior

Conversation membership, authenticated-user binding, message/conversation validation, self-receipt rejection, and monotonic delivery-to-read progression are still enforced before receipt state is accepted.

A `delivered` receipt is not blocked by read-receipt privacy settings. A `read` receipt is rejected when the authenticated recipient's current publish policy denies read sharing.

## Observation behavior

Authorized receipt reads still require current conversation membership. Stored `read` projections are returned only when both conditions are true at read time:

1. the observer's current policy permits observing read state; and
2. the receipt owner's current policy permits publishing read state.

This means a later privacy change can hide a previously stored read projection without rewriting receipt history. The service deliberately omits a hidden read projection rather than fabricating a `delivered` timestamp from the later read observation. The current one-state-per-recipient Development store therefore may not expose an earlier delivery projection after a stored read becomes hidden; production receipt persistence may preserve richer delivery history while maintaining the same privacy boundary.

## HTTP behavior

The existing receipt endpoints remain:

- `POST /v1/data/messages/{messageID}/receipts`
- `GET /v1/data/messages/{messageID}/receipts`

A privacy-denied read publication maps to HTTP 403. Authorized listing filters read projections in the service layer before the HTTP response is built. Shared JSON responses retain `Cache-Control: no-store` and `X-Content-Type-Options: nosniff`.

## Remaining gates

Persistent Privacy Shield preference storage and synchronization, production receipt persistence, cross-device read-policy convergence, push propagation, production GoreeCloud Identity sessions, abuse/rate controls, client settings and rendering, Wardveil Security/Privacy Shield/Everkeep acceptance, deployment, release, and Stable qualification remain separate work.
