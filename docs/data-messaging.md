# GoreeCloud Data Messaging

## Purpose

GoreeCloud Data is the first-party Internet messaging transport for GoreeCloud Messenger. This milestone establishes the service boundary used to accept, authorize, persist, and retrieve encrypted Data envelopes without introducing SMS, MMS, or RCS behavior into the Data path.

## Security boundary

The Data service does not accept plaintext message bodies. A submitted envelope contains ciphertext and must declare active E2EE. The current implementation validates the E2EE state as a contract; it does not yet implement cryptographic session establishment or key management.

The service receives an authenticated user identity separately from the envelope. The authenticated user must match the envelope sender, preventing a client from claiming another GoreeCloud identity through message metadata.

Conversation membership is verified server-side through a dedicated access interface before submission or history retrieval. Authorization is therefore not derived from a client-supplied participant list.

## Retry and idempotency foundation

Each envelope carries both a message identifier and a client nonce. The development store rejects duplicate message identifiers and nonce reuse. This creates a deterministic boundary for later retry and synchronization behavior without silently creating duplicate messages.

## Storage boundary

`DataStore` is the persistence interface. `MemoryDataStore` exists only for deterministic development and tests; it is not a production database implementation. Ciphertext byte slices are cloned when stored and returned so callers cannot mutate the stored representation through shared memory.

## Carrier separation

The Data service has no SMS, MMS, or RCS adapter entry point. Carrier fallback remains a separate, explicitly confirmed transport transition governed by the Messenger domain layer.

## Current limitations

This milestone does not implement a production database, network API, device authentication issuance, cryptographic protocol, key distribution, attachment transport, delivery receipts, push delivery, multi-device synchronization, carrier adapters, or production deployment.
