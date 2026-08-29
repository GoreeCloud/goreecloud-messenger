# Durable GoreeCloud Data receipt storage

## Development milestone

This milestone adds a durable single-node implementation of the existing `ReceiptStore` contract for authenticated GoreeCloud Data delivery/read receipts and an explicit Data-runtime composition selector for choosing receipt persistence.

The store preserves the current receipt semantics:

- one latest receipt projection per message and recipient;
- monotonic `delivered` to `read` state progression;
- no plaintext message body, encryption key, device secret, or carrier state in the receipt record; and
- service-layer conversation/message/user authorization remains authoritative before a receipt reaches the store.

## Persistence behavior

`FileReceiptStore` keeps a private versioned JSON document under a caller-supplied local directory. Updates are built from the current in-memory snapshot, written to a private temporary file, synchronized, and atomically renamed over the prior document before the new in-memory state becomes authoritative.

The store:

- creates its directory with owner-only permissions;
- writes the receipt document with owner-only permissions;
- fails closed on corrupt or unsupported persisted state;
- refuses receipt-state regression without modifying durable state;
- returns deterministic per-message receipt ordering; and
- honors canceled request contexts before reads or writes.

## Runtime selection

`NewConfiguredDataRuntimeHandler` now owns an explicit receipt-persistence selection step while composing the existing Data HTTP runtime. The caller must select one of the implemented modes:

- `memory` for the deterministic non-durable development store; or
- `file` for the current durable single-node file store, with an explicit root directory.

An omitted mode, unknown mode, or file mode without a usable root fails closed. There is deliberately no silent fallback from durable to memory persistence.

This closes the previous composition gap where `FileReceiptStore` existed but the Data runtime could only receive a receipt service that another caller had already assembled. The selection path is tested to reopen file-backed receipt state across store instances.

## Explicit limitations

This remains a Development single-node durability boundary. It is not production distributed persistence, replication, cross-region delivery state, backup/restore acceptance, multi-device convergence, push delivery, a Privacy Shield read-receipt preference store, or a universal proof that a human read or understood a message.

The repository command entrypoint does not yet derive this configuration from an accepted production configuration source, and no deployed process is claimed to be using file receipt persistence. Production process wiring, migration, monitoring, recovery, backup/restore evidence, distributed persistence, and deployment acceptance remain explicit later work.
