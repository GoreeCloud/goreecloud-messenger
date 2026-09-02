# Typing privacy operation context boundary

Status: Development

The single-node durable typing-privacy adapter treats its existing `context.Context` parameter as an actual fail-closed operation boundary instead of ignoring it.

Reads reject a missing or already-ended context before consulting durable state. They recheck cancellation after acquiring the adapter read lock, before returning the configured missing-record privacy default, before entering the durable file read, and again after that read completes before interpreting the stored record. A cancellation that arrives during either durable lookup path therefore cannot silently become a successful privacy decision.

Writes reject a missing or ended context before persistence work and recheck the context after acquiring the adapter lock and immediately before the atomic rename commit. A canceled write therefore does not replace an already-committed preference and temporary artifacts remain subject to the existing cleanup path.

The durable record parser is bounded and strict. Preference records may not exceed 4096 bytes, JSON decoding rejects unknown fields, and any second/trailing JSON value is rejected. This keeps corrupted or substituted local state from expanding the accepted privacy schema or forcing an unbounded file read. The current schema remains version plus `publish_typing` and `observe_typing` booleans only.

The durable read boundary validates both the path-visible file and the file descriptor actually opened for decoding. `openValidatedTypingPrivacyRecord(...)` performs the pre-open `Lstat`, requires a regular owner-only bounded record, opens it, validates the opened file metadata again, and requires `os.SameFile(...)` between the prevalidated object and opened descriptor. A symlink or file replacement that does not preserve that identity fails closed instead of becoming trusted privacy state. This narrows the validate-then-open replacement window without adding a database or network dependency.

## Cooperative multi-process write serialization

The adapter now adds a per-record cooperative write lock using atomic owner-only directory creation at the hashed record path plus `.lock`. Two compliant policy instances or processes using the same persistence root therefore cannot silently race through the atomic rename and produce last-writer-wins behavior for the same preference record. The lock path contains only the existing SHA-256-derived record key and does not add conversation or user identifiers to durable state.

Lock acquisition is context-aware and bounded. A caller waiting behind another writer continues only if the lock becomes available while its operation context remains active; otherwise the write fails closed. Even with a non-expiring context, an already-held or crash-stale lock is rejected after a bounded wait rather than waiting forever or automatically breaking a lock whose owner cannot be proven inactive. A successful write keeps the cooperative lock through temporary-file synchronization, atomic rename, and persistence-root synchronization, then removes it before returning success.

Automatic stale-lock reclamation is deliberately not implemented. A lock left behind by process termination therefore blocks later compliant writes until an operator or future recovery mechanism can establish that no active writer owns it. This favors privacy-state integrity over speculative recovery. Reads remain lock-free because committed records are replaced atomically and the read boundary independently validates the exact opened file; a hostile or non-cooperating process can still bypass this cooperative protocol.

This remains a local single-node file persistence design. It does not claim operating-system file calls are interruptible once entered, provide hostile-process exclusion, distributed consensus, cross-host locking, or transactional database semantics.

## Privacy and authority boundary

Conversation and user identifiers remain hash inputs only and are not written into record bodies or filenames. The transient cooperative lock uses the same hash-derived key and stores no record body. No additional user preference, conversation metadata, message content, process identifier, hostname, or network state is persisted by this change.

This change adds no cross-device convergence, distributed authority, production Privacy Shield persistence, backup/restore acceptance, automatic crash-lock recovery, network authority, deployment, release, or Stable qualification.
