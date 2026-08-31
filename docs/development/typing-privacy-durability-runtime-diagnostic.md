# Typing privacy durability runtime diagnostic

Status: Development

Messenger startup now projects the validated typing-privacy persistence mode into a minimized durability diagnostic.

The diagnostic reports only:

- persistence category (`memory` or `file`),
- durability level (`transient` or `single-node-durable`),
- restart durability,
- cross-device durability,
- production-readiness truth.

## Privacy and authority boundary

The diagnostic deliberately excludes storage roots, filenames, preference values, identities, conversation data, and user-derived content. File persistence is reported only as single-node restart durability; it explicitly remains non-cross-device and non-production-ready.

The startup command consumes the same persistence mode selected by `TypingPrivacyPersistenceFromEnvironment`, so the diagnostic describes the configured runtime rather than a parallel or guessed state.

This is Development evidence only and does not establish Privacy Shield production acceptance.
