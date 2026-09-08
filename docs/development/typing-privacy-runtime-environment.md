# Typing privacy runtime environment — Development

The Messenger Development command now requires an explicit process-level typing privacy persistence selection.

Environment contract:

- `GOREECLOUD_MESSENGER_TYPING_PRIVACY_PERSISTENCE=memory` requires the root variable to be unset.
- `GOREECLOUD_MESSENGER_TYPING_PRIVACY_PERSISTENCE=file` requires `GOREECLOUD_MESSENGER_TYPING_PRIVACY_ROOT` to be an explicit absolute non-root filesystem path.
- Missing, blank, unsupported, relative, contradictory, or filesystem-root configuration fails closed.
- The process configuration fixes the default typing publication/observation policy to deny (`DefaultAllowed=false`).
- Startup status exposes only the categorical persistence mode and never logs the private durable root.

This is a source-level process configuration contract. It does not itself assemble the full production Data HTTP process, establish distributed persistence, synchronize preferences across devices, or constitute production Privacy Shield acceptance.
