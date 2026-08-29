# Messenger runtime diagnostics

GoreeCloud Messenger exposes a deliberately minimized Development startup diagnostic for the accepted receipt-persistence configuration boundary.

## Current projection

The diagnostic may report only categorical runtime configuration facts:

- `receipt_persistence=memory|file`
- `receipt_durability=process-local|single-node-durable`
- `receipt_config_source=environment`

The projection is derived only after the existing fail-closed environment configuration parser accepts the receipt persistence selection.

## Sensitive-information boundary

The diagnostic must not contain:

- `GOREECLOUD_MESSENGER_RECEIPT_ROOT` or its configured filesystem path;
- filenames or directory components from durable receipt storage;
- message plaintext or ciphertext;
- receipt identifiers or recipient identities;
- conversation identifiers;
- credentials, authentication state, cryptographic material, or environment-secret values.

The executable's existing Development provenance label is separate from this configuration projection and does not contain message content.

## Failure behavior

Malformed or unsupported receipt persistence configuration continues to fail closed before a diagnostic is emitted. The diagnostic projection independently rejects inconsistent memory/file configuration rather than attempting to repair or downgrade it.

This diagnostic does not probe storage health, filesystem capacity, replication state, backup state, process readiness, or production service health. `single-node-durable` describes the selected implemented persistence class only; it is not a distributed-durability or production-acceptance claim.

## Acceptance boundary

This milestone is Development observability hygiene. It does not establish a production process supervisor, metrics exporter, centralized logging pipeline, health/readiness endpoint, alerting, distributed receipt persistence, Everkeep backup/restore acceptance, production Privacy Shield/Wardveil acceptance, deployment, release, or Stable qualification.
