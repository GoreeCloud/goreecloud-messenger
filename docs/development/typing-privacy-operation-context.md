# Typing privacy operation context boundary

Status: Development

The single-node durable typing-privacy adapter treats its existing `context.Context` parameter as an actual fail-closed operation boundary instead of ignoring it.

Reads reject a missing or already-ended context before consulting durable state. They recheck cancellation after acquiring the adapter read lock, before returning the configured missing-record privacy default, before entering the durable file read, and again after that read completes before interpreting the stored record. A cancellation that arrives during either durable lookup path therefore cannot silently become a successful privacy decision.

Writes reject a missing or ended context before persistence work and recheck the context after acquiring the adapter lock and immediately before the atomic rename commit. A canceled write therefore does not replace an already-committed preference and temporary artifacts remain subject to the existing cleanup path.

The durable record parser is now bounded and strict as well. Preference records may not exceed 4096 bytes, JSON decoding rejects unknown fields, and any second/trailing JSON value is rejected. This keeps corrupted or substituted local state from expanding the accepted privacy schema or forcing an unbounded file read. The current schema remains version plus `publish_typing` and `observe_typing` booleans only.

This is cooperative cancellation and local record validation around the single-node persistence boundary. It does not claim that operating-system file calls are interruptible once entered and does not convert the adapter into a transactional database or distributed store.

## Privacy and authority boundary

Conversation and user identifiers remain hash inputs only and are not written into record bodies or filenames. No additional user preference, conversation metadata, message content, or network state is persisted by this change.

This change adds no cross-device convergence, multi-process locking, distributed consensus, production Privacy Shield persistence, backup/restore acceptance, network authority, deployment, release, or Stable qualification.
