# Encrypted Message Reactions

## Development boundary

GoreeCloud Messenger message reactions are represented as immutable GoreeCloud Data reaction-state events linked to an existing message. This slice is Development source and does not establish production persistence, push propagation, multi-device convergence, client rendering, release acceptance, or Stable status.

## Privacy and encryption model

A `set` event carries an opaque E2EE ciphertext value. The server does not decode that value into an emoji, label, sticker, or other plaintext reaction. A `clear` event carries no reaction ciphertext and removes that reactor's active projection for the message.

The server retains only the metadata needed to authorize and order the event: reaction ID, target message/conversation, authenticated reactor, replay nonce, operation, encryption state, and timestamp. The HTTP active-reaction projection deliberately omits the client nonce and operation because every returned item is an active `set` state.

## Authorization and consistency

- The authenticated user must match the reactor identity.
- The reactor must be a current participant in the stated conversation.
- The target message must already exist in that same conversation.
- Missing and cross-conversation target identifiers use the same bounded `reaction target unavailable` state.
- Each reaction event ID and per-reactor client nonce is replay-protected.
- A newer event from one reactor replaces that reactor's prior state for the message; stale or equal-timestamp events are rejected.
- A `clear` event removes that reactor's current active reaction without revealing the prior plaintext reaction value to the service.

## Current HTTP surface

- `POST /v1/data/messages/{messageID}/reactions` records a set or clear event.
- `GET /v1/data/messages/{messageID}/reactions` returns current active encrypted reaction projections to an authorized conversation participant.

Responses continue to use `Cache-Control: no-store` and `X-Content-Type-Options: nosniff`.

## Remaining gates

Production reaction persistence/indexing, offline and multi-device convergence, push delivery, client-side encrypt/decrypt and rendering behavior, reaction accessibility/UI treatment under the current Stable Glaze UI contract, production GoreeCloud Identity integration, complete E2EE session/key operation, Wardveil Security/Privacy Shield/Everkeep acceptance, deployment, release, and Stable qualification remain separate work.
