# Shared typing privacy runtime composition — Development

GoreeCloud Messenger now has a bounded runtime composition helper that turns the explicit typing privacy persistence configuration into one shared `TypingPrivacyRuntimeStore` and injects that same store into both the ephemeral typing service and the authenticated typing privacy preference service.

Authority and privacy boundary:

- `TypingPresenceRuntimeFromEnvironment` first resolves the existing fail-closed `GOREECLOUD_MESSENGER_TYPING_PRIVACY_PERSISTENCE` configuration.
- Exactly one memory-backed or file-backed typing privacy runtime store is constructed from that accepted configuration.
- The same store instance is supplied to `TypingService` as its publish/observe privacy policy and to `TypingPrivacyPreferenceService` as its mutable preference store.
- Default policy remains deny. Tests prove a typing publish is denied before an authenticated preference update and becomes allowed only after the shared preference service changes the same backing policy state.
- Observation is likewise tested through the shared store: an authorized observer sees another active typing participant only after the relevant privacy preferences permit it.
- Memory mode now strictly requires `GOREECLOUD_MESSENGER_TYPING_PRIVACY_ROOT` to be absent; an explicitly set blank root is rejected rather than silently ignored.
- No message content, draft text, keystrokes, cursor data, ciphertext, credentials, identity secrets, or client timestamps are added to the typing contract.

This helper is a Development bootstrap boundary. It does not assemble the complete production Data HTTP process, establish production GoreeCloud Identity sessions, distributed or replicated preference persistence, cross-device synchronization, production presence fan-out, deployment configuration, release acceptance, or Stable qualification.
