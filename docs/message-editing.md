# GoreeCloud Data Message Editing Foundation

GoreeCloud Messenger models sent-message editing as an immutable encrypted revision stream rather than server-side plaintext mutation.

## Boundary

- Only the authenticated original sender may append an edit revision.
- The original message and the edit must belong to the same conversation.
- Conversation membership is verified server-side.
- Every edit carries opaque ciphertext and requires the `e2ee` state at the service boundary.
- The server does not decrypt or interpret edited message content.
- Edit IDs are unique and client nonces are replay-protected per editor.
- An edit timestamp may not precede the original message timestamp.
- Authorized conversation participants may retrieve the encrypted revision history for a message.

## HTTP surface

`POST /v1/data/messages/{messageID}/edits`

Appends one immutable encrypted edit revision. The JSON body contains `edit_id`, `conversation_id`, `editor_id`, `client_nonce`, base64 `ciphertext`, `encryption`, and `edited_at`. Unknown fields are rejected and request bodies are bounded by the existing Data API request limit.

`GET /v1/data/messages/{messageID}/edits`

Returns the authorized encrypted revision history. Ciphertext remains base64-encoded JSON data and responses use the existing no-store and nosniff JSON protections.

## Current acceptance boundary

This is a Development service and HTTP-contract foundation. It does not provide production persistence, edit-window policy, client-side plaintext re-encryption, edit rendering, device synchronization, push propagation, moderation retention policy, cryptographic session establishment, deployment, or Stable acceptance.
