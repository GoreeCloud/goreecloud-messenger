# Durable GoreeCloud Data receipt storage

## Development milestone

This milestone provides a durable single-node implementation of the existing `ReceiptStore` contract for authenticated GoreeCloud Data delivery/read receipts, an explicit Data-runtime composition selector for choosing receipt persistence, and a hardened local filesystem durability boundary.

The store preserves the current receipt semantics:

- one latest receipt projection per message and recipient;
- monotonic `delivered` to `read` state progression;
- no plaintext message body, encryption key, device secret, or carrier state in the receipt record; and
- service-layer conversation/message/user authorization remains authoritative before a receipt reaches the store.

## Persistence behavior

`FileReceiptStore` keeps a private versioned JSON document under a caller-supplied local directory. Updates are built from the current in-memory snapshot, written to a private temporary file, synchronized, atomically renamed over the prior document, protected as owner-only, and followed by parent-directory synchronization before the new in-memory state becomes authoritative.

The store:

- creates or tightens its persistence directory to owner-only `0700` permissions;
- rejects a persistence root that is a symbolic link or not a directory;
- requires an existing `receipts.json` to be a regular non-symlink file with owner-only `0600` permissions;
- writes temporary and final receipt documents with owner-only permissions;
- synchronizes the temporary file before rename and the parent directory after rename;
- fails closed on corrupt, unsupported, permissively exposed, or structurally unsafe persisted state;
- refuses receipt-state regression without modifying durable state;
- returns deterministic per-message receipt ordering; and
- honors canceled request contexts before reads or writes.

## Post-rename ambiguity containment

A rename can become visible before a parent-directory synchronization failure establishes whether that directory entry is crash-durable. Continuing with the prior in-memory snapshot after such a failure would risk process/disk disagreement.

For that reason, a post-rename permission or directory-sync failure is surfaced as `ErrReceiptDurabilityUnknown`. The current `FileReceiptStore` instance then poisons itself and refuses subsequent receipt reads or writes. Reopening the store is the explicit reconciliation boundary: it revalidates the filesystem object and loads the versioned state that is actually present on disk.

This is fail-closed ambiguity containment, not a claim that a failed fsync magically proves durability.

## Runtime selection

`NewConfiguredDataRuntimeHandler` owns an explicit receipt-persistence selection step while composing the existing Data HTTP runtime. The caller must select one of the implemented modes:

- `memory` for the deterministic non-durable development store; or
- `file` for the current durable single-node file store, with an explicit root directory.

An omitted mode, unknown mode, or file mode without a usable root fails closed. There is deliberately no silent fallback from durable to memory persistence.

This closes the previous composition gap where `FileReceiptStore` existed but the Data runtime could only receive a receipt service that another caller had already assembled. The selection path is tested to reopen file-backed receipt state across store instances.

## Explicit limitations

This remains a Development single-node durability boundary. Parent-directory fsync, private-path checks, and ambiguity poisoning strengthen local persistence semantics but do not make the file store a production distributed database.

It is not replication, cross-region delivery state, multi-writer coordination, backup/restore acceptance, multi-device convergence, push delivery, a Privacy Shield read-receipt preference store, or a universal proof that a human read or understood a message.

The repository command entrypoint does not yet derive this configuration from an accepted production configuration source, and no deployed process is claimed to be using file receipt persistence. Production process wiring, migration, monitoring, recovery, backup/restore evidence, distributed persistence, and deployment acceptance remain explicit later work.
