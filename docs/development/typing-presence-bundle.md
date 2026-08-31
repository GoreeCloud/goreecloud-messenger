# Typing Presence Runtime Bundle — Development

GoreeCloud Messenger now provides one explicit runtime-composition helper for enabling ephemeral typing events together with their authenticated privacy-control surface.

## Behavior

- `WithTypingPresence` composes both `TypingService` and `TypingPrivacyPreferenceService` under the runtime's existing Authenticator.
- The helper builds on the existing independently optional typing and preference composition methods rather than adding new authority.
- Typing publication continues to derive actor identity from authentication.
- Conversation membership and publish/observe privacy policy remain enforced by the existing services.
- The preference store remains Development memory state unless a deployment deliberately supplies a durable implementation.

## Boundary

Development only. This does not establish durable Privacy Shield preference backing, production fan-out, native Messenger typing UI, deployment acceptance, or Stable qualification.
