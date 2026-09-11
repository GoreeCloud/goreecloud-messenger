# Android Development build provenance

## Status

Development build/package evidence only.

This document describes the bounded provenance evidence generated for the GoreeCloud Messenger Android Development debug APK. It does not establish Release Candidate status, production acceptance, production signing authority, a production release, or Stable qualification.

## Exact-source build boundary

The Messenger Android workflow checks out the exact pull-request head revision for pull-request runs and the exact event revision for main-branch push runs. The evidence generator compares the checked-out `HEAD` to that declared exact source revision and fails closed if they differ.

The generated evidence records:

- repository identity;
- exact source commit SHA;
- workflow, run, attempt, job, and event identity;
- APK filename, SHA-256, and byte size;
- Android package name;
- version name and version code;
- minimum and target Android API levels;
- debuggable Development-build state;
- Android APK signature verification result;
- signer-certificate SHA-256 digest;
- Android build-tools version used for metadata and signature verification; and
- explicit false states for Release Candidate, release acceptance, production acceptance, Stable, and production release-signing authority.

The evidence generator currently requires the Development APK to remain:

- package `com.goreecloud.messenger`;
- version name `0.1.0-dev`;
- version code `1`;
- minimum API 26;
- target API 35; and
- debuggable.

A deliberate package/version/API transition must update the evidence contract and then pass fresh exact-head validation. The workflow must not silently accept an unexpected APK identity.

## Signature boundary

`apksigner verify --verbose --print-certs` is executed against the built APK. The evidence retains only the verification state and signer-certificate SHA-256 digest needed to identify the Development build signature.

A verified debug APK signature proves package integrity under that Development signer. It does **not** establish GoreeCloud production release-signing authority. No production signing key, private key, keystore password, reusable credential, or equivalent release secret belongs in this repository, generated evidence, ordinary logs, or documentation.

A future Release Candidate or production release-signing workflow must use the approved GoreeCloud secret-storage and release-governance boundary, bind the signed artifact to an exact accepted candidate, preserve artifact checksum/provenance evidence, and complete all separately applicable release and production-acceptance gates.

## CI artifact evidence

The Android workflow retains two bounded artifact families for 30 days:

- `goreecloud-messenger-android-debug` — the Development debug APK; and
- `goreecloud-messenger-android-debug-evidence` — privacy-minimized build-evidence JSON, APK SHA-256 file, and signature-evidence text.

Artifact retention is build evidence, not a release publication mechanism. A successful APK build or artifact upload does not advance Messenger beyond its verified Development lifecycle.

## Supply-chain boundary

The Android workflow pins its GitHub Actions dependencies to reviewed immutable commit SHAs rather than relying on moving major-version tags. The Android emulator runner was already commit-pinned. Pinning automation improves reproducibility and provenance but does not by itself establish application security, release acceptance, or Stable status.

## Current product boundary

The retained APK inherits the Messenger Android Development client's fail-closed communication boundary. It does not gain production GoreeCloud Identity sessions, live conversation authorization, GoreeCloud Data networking, a reviewed concrete E2EE implementation, message persistence/delivery/synchronization, carrier messaging authority, calling authority, production Wardveil Security/Privacy Shield/Everkeep/Mesh/Manager acceptance, production signing, or a Send action merely because provenance evidence exists.
