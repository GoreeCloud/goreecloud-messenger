# Typing privacy file commit durability

Status: Development

The single-node file-backed typing privacy adapter now strengthens its local commit and read boundary after the existing temporary-file write and file sync.

The persistence root is resolved to its canonical filesystem location after creation/protection. Before reads and writes, Messenger verifies that the canonical root still exists as an owner-private non-symlink directory. After an atomic record rename, Messenger synchronizes the persistence directory before reporting the preference update as successful. This reduces the gap between a synced temporary record and durable directory-entry metadata after a process or host interruption.

Existing preference records are also required to be owner-private regular files rather than symlinks or other filesystem object types. Missing records still use the configured default policy, while substituted or broadly readable records fail closed instead of being interpreted as trusted typing authorization state.

The stored data remains minimized to schema version plus publish/observe booleans. Conversation and user identifiers remain outside record bodies and filenames.

## Authority boundary

This is still Development single-node durability. It does not provide cross-device convergence, distributed consensus, adversarial same-user filesystem isolation, production Privacy Shield persistence, backup/restore acceptance, deployment, or Stable qualification.
