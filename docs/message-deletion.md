# Delete-for-everyone tombstone foundation

GoreeCloud Messenger models delete-for-everyone as an immutable control tombstone associated with an existing GoreeCloud Data message.

## Current Development contract

- Only the authenticated original sender may create the tombstone.
- The sender must still be an authorized participant in the message conversation.
- The request carries a unique deletion ID and client nonce for replay protection.
- A message accepts at most one delete-for-everyone tombstone in the current foundation.
- The deletion timestamp may not precede the original message timestamp.
- Authorized conversation participants may retrieve the tombstone so clients can synchronize deletion state.
- HTTP deletion responses deliberately omit the client nonce and never include message plaintext or ciphertext.

## HTTP surface

- `POST /v1/data/messages/{messageID}/deletions`
- `GET /v1/data/messages/{messageID}/deletions`

The POST body contains deletion metadata only. The GET response contains zero or one tombstone in the current model.

## Security, privacy, and continuity boundary

A delete-for-everyone tombstone is a synchronization instruction, not proof of universal physical erasure. This foundation does not claim that an offline device, independent export, backup, legal hold, or already-retained encrypted object has been securely erased. Production retention, encrypted-object purge, Everkeep backup semantics, GoreeCloud Sync propagation, and device-side suppression/erasure require separate implementation and acceptance.

The server continues to operate without decrypting GoreeCloud Data message content. Wardveil Security, Privacy Shield, Everkeep, GoreeCloud Identity, and GoreeCloud Sync integration remain evidence-gated production work.
