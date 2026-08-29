# Durable GoreeCloud Data receipt storage

## Development milestone

This milestone provides a durable single-node implementation of the existing `ReceiptStore` contract for authenticated GoreeCloud Data delivery/read receipts, an explicit Data-runtime composition selector for choosing receipt persistence, a hardened local filesystem durability boundary, and an explicit executable environment-configuration contract.

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

This closes the earlier composition gap where `FileReceiptStore` existed but the Data runtime could only receive a receipt service that another caller had already assembled. The selection path is tested to reopen file-backed receipt state across store instances.

## Executable environment configuration

The development executable now derives receipt persistence from an explicit process-environment contract before it declares the Messenger development contract active.

Required configuration:

- `GOREECLOUD_MESSENGER_RECEIPT_PERSISTENCE=memory` selects the deterministic non-durable store and requires `GOREECLOUD_MESSENGER_RECEIPT_ROOT` to be unset or empty.
- `GOREECLOUD_MESSENGER_RECEIPT_PERSISTENCE=file` requires `GOREECLOUD_MESSENGER_RECEIPT_ROOT` to be an absolute, non-root filesystem path.

Missing mode, unsupported mode, file mode without a root, a relative file root, filesystem-root persistence, or an ignored durable root supplied alongside memory mode all fail closed. Mode parsing is case-insensitive and surrounding whitespace is removed. The returned file root is normalized with the platform path cleaner before it reaches runtime composition.

This configuration parser is deliberately testable through an injected environment lookup and does not infer a fallback from host state, current working directory, or an automatically created default path.

## Explicit limitations

This remains a Development single-node durability boundary. Parent-directory fsync, private-path checks, ambiguity poisoning, and explicit environment parsing strengthen local persistence semantics but do not make the file store a production distributed database.

The command now validates and derives the receipt-persistence selection from its process environment, but the current `cmd/messenger` executable is still a development contract exerciser rather than the full Data HTTP server bootstrap. It does not yet construct the Data service, authentication, conversation access, attachment service, and `NewConfiguredDataRuntimeHandler` from one production-accepted process configuration.

This work is not replication, cross-region delivery state, multi-writer coordination, backup/restore acceptance, multi-device convergence, push delivery, a Privacy Shield read-receipt preference store, or a universal proof that a human read or understood a message. Production server bootstrap, secrets/Identity configuration, migration, monitoring, recovery, backup/restore evidence, distributed persistence, and deployment acceptance remain explicit later work.
