# GoreeCloud Data direct-reply foundation

GoreeCloud Messenger models a direct reply as encrypted GoreeCloud Data message content plus a minimal server-visible reference to the target message identifier.

## Current Development contract

- `reply_to_message_id` is optional. Messages without it remain ordinary GoreeCloud Data messages.
- A reply target must already exist and belong to the same conversation as the new message.
- The authenticated sender must still pass the existing conversation-membership and sender-binding checks.
- The reply target identifier is normalized before persistence.
- Missing targets, self-references, and cross-conversation targets are rejected through the same bounded `reply target unavailable` service state so the HTTP boundary does not distinguish whether a message identifier exists in another conversation.
- Conversation history returns the reply target identifier alongside the existing opaque ciphertext so authorized clients can render reply context after client-side decryption.

## E2EE and privacy boundary

The server does not decrypt the original message, the reply, or any reply preview. The reply relationship is routing/conversation metadata only. Any plaintext quote, sender display text, or preview shown in a client must be derived from content the authorized client can already decrypt and must not be added as a server-side plaintext field.

This slice does not implement threaded discussions, reply-preview caching, push propagation, offline synchronization, production persistence, production GoreeCloud Identity, or complete cryptographic session/key operation. Those remain separate implementation and acceptance work.
