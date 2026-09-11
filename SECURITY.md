# GoreeCloud Messenger Security

Status: Development  
Detailed security architecture: [`docs/security.md`](docs/security.md)  
Current FR-005 boundary: [`docs/development/e2ee-authority-acceptance.md`](docs/development/e2ee-authority-acceptance.md)  
Canonical project record: `GoreeCloud/Projects/Project Specification — Messenger`

## Security policy

GoreeCloud Messenger must represent security state only when that state is supported by the responsible implementation and evidence. End-to-end encryption, Wardveil Security, Privacy Shield, authentication, authorization, transport protection, and related indicators are not decorative claims.

The application must preserve explicit transport provenance between GoreeCloud Data, SMS, MMS, RCS, and calling paths. Encrypted GoreeCloud Data communication must not silently become a carrier transport or inherit an E2EE label from transport availability alone.

## Current Development boundaries

- Production GoreeCloud Identity credentials, sessions, device identity, and username-resolution authority are not yet accepted.
- Complete security-reviewed E2EE device/session/key establishment, verification, rotation, group behavior, and multi-device lifecycle are not yet accepted.
- No authoritative GoreeCloud record currently selects a concrete Messenger E2EE protocol or implementation; this repository must not invent or imply one.
- The native Android Development client remains disconnected from live communication authority.
- Android Data readiness requires independent positive evidence for Identity authentication, exact conversation authorization, GoreeCloud Data transport, and exact conversation-scoped active E2EE before a future operation may project `Data · E2EE`.
- A future `E2EESessionAuthority` may project active E2EE only when the minimized Development acceptance seam also reports accepted implementation review, enrolled local cryptographic device identity, established conversation session, current key lifecycle, and the exact canonical conversation scope. Missing, rejected, contradictory, or mismatched evidence fails closed.
- Test fixtures that construct positive E2EE acceptance evidence prove only policy behavior; they do not create a security review, cryptographic implementation, key/session state, or production acceptance.
- The current disconnected Android source guard rejects direct `java.security` and `javax.crypto` implementation markers so this acceptance-only slice cannot silently become an unreviewed client cryptographic implementation.
- Authority-owned conversation identifiers are bounded opaque values. Leading/trailing whitespace, overlong identifiers, C0 controls, and DEL are rejected rather than normalized into another authorization or cryptographic scope.
- Server message and attachment services operate on encrypted/opaque payload boundaries and must not claim plaintext inspection where the architecture does not provide it.
- Current local persistence, diagnostics, typing anti-flood controls, and runtime health/readiness foundations remain Development evidence rather than production Wardveil Security acceptance.

## Cryptography requirements

Messenger will use established, security-reviewed cryptographic protocols and implementations rather than inventing proprietary cryptographic primitives. Cryptographic architecture must address device identity, authenticated session establishment, forward secrecy where supported, key rotation/change handling, device verification, group security, attachment protection, multi-device enrollment/revocation, local protected state, and safe backup/recovery behavior.

The protocol-neutral FR-005 acceptance seam is deliberately narrower than that complete architecture. It is intended to prevent a bare application enum from claiming active E2EE before the eventual responsible cryptographic authority has independently accepted its implementation, device identity, session, key lifecycle, and conversation scope.

## Vulnerability handling

Do not place credentials, reusable tokens, private keys, signing material, message plaintext, attachment plaintext, device secrets, or other sensitive user data in public issue content, logs, test fixtures, or repository commits. Security reports should use an authorized private GoreeCloud/GitHub reporting path when available; if no private reporting path is available, do not publish sensitive exploit details merely to create a report.

## Acceptance rule

A successful unit test, CI workflow, source review, Development acceptance projection, or runtime check establishes only the evidence that it actually performs. Production security acceptance, platform-system acceptance, release approval, signing/provenance acceptance, deployment approval, and Stable qualification remain separate governed gates.
