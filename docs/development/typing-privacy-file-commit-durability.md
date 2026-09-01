# Typing privacy file commit durability

Status: Development

The single-node file-backed typing privacy adapter now strengthens its local commit boundary after the existing temporary-file write and file sync.

The persistence root is resolved to its canonical filesystem location after creation/protection. After an atomic record rename, Messenger synchronizes the persistence directory before reporting the preference update as successful. This reduces the gap between a synced temporary record and durable directory-entry metadata after a process or host interruption.

The stored data remains minimized to schema version plus publish/observe booleans. Conversation and user identifiers remain outside record bodies and filenames.

## Authority boundary

This is still Development single-node durability. It does not provide cross-device convergence, distributed consensus, production Privacy Shield persistence, backup/restore acceptance, deployment, or Stable qualification.
