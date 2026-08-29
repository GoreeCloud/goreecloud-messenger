# Durable GoreeCloud Data receipt storage

## Development milestone

This milestone adds a durable single-node implementation of the existing `ReceiptStore` contract for authenticated GoreeCloud Data delivery/read receipts.

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

## Explicit limitations

This is a Development single-node durability boundary. It is not production distributed persistence, replication, cross-region delivery state, backup/restore acceptance, multi-device convergence, push delivery, a Privacy Shield read-receipt preference store, or a universal proof that a human read or understood a message.

The HTTP runtime does not automatically select this store. Production composition, database selection, deployment, migration, monitoring, recovery, and acceptance remain explicit later work.
