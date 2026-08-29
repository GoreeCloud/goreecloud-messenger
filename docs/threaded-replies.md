# GoreeCloud Messenger Threaded Reply Foundation

## Purpose

This Development slice adds focused threaded discussions to GoreeCloud Data without creating a server-side plaintext conversation layer. Thread structure is represented by message identifiers only; message bodies remain opaque E2EE ciphertext at the Data service boundary.

## Thread metadata

A threaded message carries both:

- `reply_to_message_id` — the immediate parent message; and
- `thread_root_message_id` — the stable top-level message that owns the focused discussion.

A direct reply that is not part of a focused thread may continue to carry only `reply_to_message_id`.

The first message in a thread replies directly to the root and sets both identifiers to the root. A nested reply may point `reply_to_message_id` at another message already in that same thread while preserving the original `thread_root_message_id`.

## Validation rules

The Data service:

- authenticates the sender and verifies current conversation membership before thread validation;
- normalizes reply and thread-root identifiers before persistence;
- requires the thread root to already exist in the same conversation;
- requires the thread root itself not to be another threaded child;
- requires every threaded message to have an immediate reply target;
- permits the immediate target to be either the root or a message already assigned to that root;
- rejects self-references and unrelated same-conversation parents; and
- uses one bounded `thread target unavailable` state for missing, cross-conversation, invalid-root, and unrelated-parent cases.

The bounded failure behavior prevents the HTTP surface from becoming a message-existence oracle across conversations.

## Thread history

`GET /v1/data/conversations/{conversationID}/threads/{rootMessageID}/messages` returns the root message followed by messages whose `thread_root_message_id` matches that root, in the current deterministic conversation-store order.

The caller must be an authenticated participant in the conversation. A missing root, a root from another conversation, or an attempt to use a threaded child as a new thread root is rejected through the same bounded thread-target state.

## Privacy and encryption boundary

The server does not decrypt thread roots, replies, quoted text, or previews. Thread metadata identifies relationships between encrypted messages only. Any plaintext thread title, quote, summary, preview, or navigation label must be derived by an authorized client from content that client can already decrypt.

This slice does not establish production thread indexing, push propagation, offline synchronization, multi-device convergence, moderation workflows, server-side semantic summaries, production persistence, complete cryptographic session/key operation, deployment, release, or Stable acceptance.
