# GoreeCloud Messenger User Manual

## Current availability

GoreeCloud Messenger is in **Active Development**. The repository currently provides messaging-domain, service, authenticated HTTP transport, receipt, persistence, and encrypted-attachment foundations. A production-ready end-user Messenger client is **not yet available** from this repository.

This manual therefore explains what current Development builds/services do and which product behaviors must not yet be assumed.

## GoreeCloud Data messaging

GoreeCloud Data is the first-party GoreeCloud messaging transport for authorized GoreeCloud conversations. It is technically distinct from SMS, MMS, and RCS and must not silently downgrade an encrypted GoreeCloud conversation to a carrier transport.

Current source validates encrypted Data envelopes, enforces authenticated sender/conversation authorization, provides deterministic retry protection, and exposes authenticated message transport through the Development HTTP service.

End-to-end encryption is represented only when the relevant cryptographic state is actually verified. The current repository does not yet establish production device/key lifecycle or complete production cryptographic session establishment.

## Delivery and read state

The current service foundation supports authenticated delivery and read receipts. Receipt state represents authorized application-level state; it is not carrier proof, cryptographic proof that a human read content, or evidence that every planned client has rendered a message.

Later privacy controls for receipt behavior must remain explicit and evidence-backed before user-facing production claims.

## Encrypted attachments

Current GoreeCloud Data attachment transport treats attachment content as opaque ciphertext.

The service can currently:

- accept encrypted attachment bytes through authorized conversation boundaries;
- return bounded attachment metadata;
- fetch ciphertext through JSON/base64 Development transport;
- delete attachments with replay-safe tombstone semantics; and
- download the exact stored ciphertext bytes through a participant-authorized raw transport.

The server does not decrypt attachment media in order to provide this transport. Plaintext filenames, MIME interpretation, client-side decrypt/render behavior, production object storage, and any plaintext security-scanning workflow require separate contracts and acceptance.

## Identity and usernames

GoreeCloud usernames are intended to be first-class Messenger identities and do not have to map one-to-one to phone numbers. Production username resolution and account disclosure must use GoreeCloud Identity-owned contracts rather than turning Messenger into an unrestricted account directory.

The current service does not establish production GoreeCloud Identity sessions or a production consumer username-resolution experience.

## SMS, MMS, RCS, and calls

The domain model keeps Data, SMS, MMS, and RCS transports distinct. Carrier integration is permitted only through legitimate supported platform/carrier APIs. Voice/video calling and carrier calling must remain distinguishable.

Current source does not establish production carrier adapters, voice/video media transport, or complete calling clients.

## Privacy and security expectations

- Message transport must preserve actual transport provenance.
- Encrypted conversations must not claim verified end-to-end protection unless the client has verified that state.
- Attachment services transport ciphertext and must not invent plaintext interpretation.
- Privacy Shield governs privacy-sensitive messaging behavior and user controls.
- Wardveil Security governs applicable protection, detection, trust, verification, and response boundaries.
- GoreeCloud Identity governs account, device, credential, session, and authorization authority.
- Everkeep governs accepted backup, recovery, preservation, portability, and succession behavior.
- GoreeCloud Mesh governs authenticated cross-service coordination.
- Glaze UI governs approved client-interface design once client surfaces are implemented and accepted.

## Current limitations

The Development repository does not yet establish production identity sessions, complete cryptographic session/key management, distributed delivery/storage, production object storage, push delivery, carrier adapters, calling media transport, anti-abuse/rate-limit acceptance, complete native clients, Glaze UI client acceptance, deployment, signed release, or Stable qualification.

Until a later release record states otherwise, there is no supported production end-user installation workflow from this repository.

Refer to `README.md`, `SPECIFICATIONS.md`, `FEATURES.md`, and `docs/` for technical and acceptance details.
