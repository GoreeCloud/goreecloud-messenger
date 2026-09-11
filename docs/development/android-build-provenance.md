# Android Development build provenance

## Status

Development build/package/SBOM evidence only.

This document describes the bounded provenance evidence generated for the GoreeCloud Messenger Android Development debug APK. It does not establish Release Candidate status, production acceptance, production signing authority, a production release, or Stable qualification.

## Exact-source build boundary

The Messenger Android workflow checks out the exact pull-request head revision for pull-request runs and the exact event revision for main-branch push runs. The evidence generators compare the checked-out `HEAD` to that declared exact source revision and fail closed if they differ.

The generated build evidence records:

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
- Android build-tools version used for metadata and signature verification;
- exact Development SBOM identity, SHA-256, byte size, CycloneDX specification version, deterministic serial number, resolved runtime-component count, and Gradle configuration; and
- explicit false states for Release Candidate, release acceptance, production acceptance, Stable, and production release-signing authority.

The evidence generator currently requires the Development APK to remain:

- package `com.goreecloud.messenger`;
- version name `0.1.0-dev`;
- version code `1`;
- minimum API 26;
- target API 35; and
- debuggable.

A deliberate package/version/API transition must update the evidence contract and then pass fresh exact-head validation. The workflow must not silently accept an unexpected APK identity.

## Development SBOM boundary

The Android workflow resolves the exact `debugRuntimeClasspath` dependency closure after the same exact-head debug APK build succeeds. The raw Gradle dependency report is a temporary generator input and is removed before CI evidence is retained.

`scripts/generate_android_sbom.py` converts that resolved runtime closure into deterministic CycloneDX 1.6 JSON without introducing a new runtime library or SBOM-generation plugin into the Android application. The SBOM is bound to:

- repository `GoreeCloud/goreecloud-messenger`;
- the exact checked-out source revision;
- the built `app-debug.apk` SHA-256;
- Android package `com.goreecloud.messenger`;
- version `0.1.0-dev`;
- lifecycle `Development`; and
- Gradle configuration `debugRuntimeClasspath`.

The generated SBOM uses Maven package URLs for the resolved transitive runtime-component inventory. Its document serial number is deterministically derived from the exact repository, source revision, APK SHA-256, and runtime configuration. The SBOM has no wall-clock timestamp, so identical accepted inputs produce a stable evidence identity.

`scripts/generate_android_build_evidence.py` independently revalidates the SBOM before build evidence can pass. It requires the expected CycloneDX 1.6 document/root identity, recomputes the deterministic serial number, requires the SBOM root SHA-256 to equal the built APK SHA-256, verifies exact repository/source/configuration metadata, rejects duplicate or ambiguous runtime components, requires the root dependency closure to cover exactly the listed runtime components, and records the SBOM checksum and identity in build evidence.

This SBOM is a Development packaged-runtime inventory. It is not by itself a vulnerability scan, license approval, dependency-security review, source-composition proof for every build tool, production dependency acceptance, or Stable qualification. Android Gradle/Kotlin build plugins and GitHub Actions automation are governed separately through exact version or immutable action provenance.

## Signature boundary

`apksigner verify --verbose --print-certs` is executed against the built APK. The evidence retains only the verification state and signer-certificate SHA-256 digest needed to identify the Development build signature.

A verified debug APK signature proves package integrity under that Development signer. It does **not** establish GoreeCloud production release-signing authority. No production signing key, private key, keystore password, reusable credential, or equivalent release secret belongs in this repository, generated evidence, ordinary logs, or documentation.

Clean CI runners may create a fresh Android debug keystore. The Development signer-certificate identity and therefore the signed APK checksum may differ across exact-head workflow runs even when Messenger runtime source and resolved runtime dependencies are unchanged. Each retained debug APK is therefore treated as its own Development artifact and must be bound to its own exact source revision, APK checksum, signer-certificate digest, and SBOM identity. This is exact per-artifact traceability; it is not reproducible production signing or signer continuity.

A future Release Candidate or production release-signing workflow must use the approved GoreeCloud secret-storage and release-governance boundary, bind the signed artifact to an exact accepted candidate, preserve artifact checksum/provenance and required SBOM evidence, and complete all separately applicable release and production-acceptance gates.

## CI artifact evidence

The Android workflow retains two bounded artifact families for 30 days:

- `goreecloud-messenger-android-debug` — the Development debug APK; and
- `goreecloud-messenger-android-debug-evidence` — privacy-minimized build-evidence JSON, APK SHA-256 file, signature-evidence text, CycloneDX Development SBOM, and SBOM SHA-256 file.

Artifact retention is build evidence, not a release publication mechanism. A successful APK build, SBOM generation, or artifact upload does not advance Messenger beyond its verified Development lifecycle.

## Supply-chain boundary

The Android workflow pins its GitHub Actions dependencies to reviewed immutable commit SHAs rather than relying on moving major-version tags. The Android emulator runner is also commit-pinned. Pinning automation and retaining exact runtime-dependency inventory improve reproducibility and provenance but do not by themselves establish application security, dependency vulnerability/license acceptance, release acceptance, or Stable status.

The governing GoreeCloud Stable release standard separately requires dependency maintenance/review and exact-candidate SBOM identity where applicable. This Development evidence advances traceability toward that gate without claiming that the complete Stable supply-chain qualification has been satisfied.

## Current product boundary

The retained APK inherits the Messenger Android Development client's fail-closed communication boundary. It does not gain production GoreeCloud Identity sessions, live conversation authorization, GoreeCloud Data networking, a reviewed concrete E2EE implementation, message persistence/delivery/synchronization, carrier messaging authority, calling authority, production Wardveil Security/Privacy Shield/Everkeep/Mesh/Manager acceptance, production signing, or a Send action merely because build/SBOM provenance evidence exists.
