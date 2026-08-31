# Typing Privacy Persistence Selection — Development

GoreeCloud Messenger now has an explicit Development composition contract for choosing the persistence implementation behind typing-presence privacy preferences.

Supported modes are deliberately bounded:

- `memory` uses the existing process-memory policy/preference implementation and rejects a configured file root.
- `file` requires a persistence root and uses the minimized file-backed adapter from the preceding durable-preference slice.

Missing, contradictory, or unsupported configuration fails closed. Both implementations satisfy the same `TypingPrivacyPolicy` and `TypingPrivacyPreferenceStore` boundaries, so the typing service and authenticated preference service can share one authoritative preference source within a composed runtime.

This is not production Privacy Shield integration. No environment-variable/process bootstrap, distributed storage, cross-device synchronization, deployment configuration, production Identity acceptance, presence fan-out, or Stable qualification is established by this source contract.
