#!/usr/bin/env python3
"""Generate fail-closed, privacy-minimized evidence for the Development Android APK.

This script validates the exact checked-out source revision, APK package/version/API
identity, checksum, debuggable classification, and APK signature state. It does not
create or assert release signing, Release Candidate acceptance, production
acceptance, or Stable qualification.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
EXPECTED_REPOSITORY = "GoreeCloud/goreecloud-messenger"
EXPECTED_PACKAGE = "com.goreecloud.messenger"
EXPECTED_VERSION_NAME = "0.1.0-dev"
EXPECTED_VERSION_CODE = "1"
EXPECTED_MIN_SDK = "26"
EXPECTED_TARGET_SDK = "35"
SHA_RE = re.compile(r"^[0-9a-f]{40}$")
CERT_SHA256_RE = re.compile(r"^[0-9a-f]{64}$")


def fail(message: str) -> None:
    print(f"Android build evidence FAILED: {message}", file=sys.stderr)
    raise SystemExit(1)


def run_checked(command: list[str]) -> str:
    completed = subprocess.run(
        command,
        cwd=ROOT,
        check=False,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
    )
    if completed.returncode != 0:
        fail(f"command failed ({completed.returncode}): {' '.join(command)}\n{completed.stdout}")
    return completed.stdout


def require_env(name: str) -> str:
    value = os.environ.get(name, "").strip()
    if not value:
        fail(f"required environment variable {name} is missing")
    return value


def build_tools_version(path: Path) -> tuple[int, ...]:
    try:
        return tuple(int(part) for part in path.name.split("."))
    except ValueError:
        return (0,)


def find_build_tool(name: str) -> Path:
    sdk_root_value = os.environ.get("ANDROID_HOME") or os.environ.get("ANDROID_SDK_ROOT")
    if not sdk_root_value:
        fail("ANDROID_HOME/ANDROID_SDK_ROOT is unavailable")
    build_tools_root = Path(sdk_root_value) / "build-tools"
    if not build_tools_root.is_dir():
        fail(f"Android build-tools directory is missing: {build_tools_root}")
    candidates = sorted(
        (directory for directory in build_tools_root.iterdir() if directory.is_dir()),
        key=build_tools_version,
        reverse=True,
    )
    for directory in candidates:
        tool = directory / name
        if tool.is_file():
            return tool
    fail(f"Android build tool {name} was not found")
    raise AssertionError("unreachable")


def extract_single(pattern: str, text: str, label: str) -> str:
    match = re.search(pattern, text, re.MULTILINE)
    if not match:
        fail(f"unable to read {label} from APK metadata")
    return match.group(1)


def extract_signer_certificate_sha256(signature_report: str) -> str:
    candidate_lines: list[str] = []
    digests: set[str] = set()

    for raw_line in signature_report.splitlines():
        line = raw_line.strip()
        lowered = line.lower()
        if "certificate" not in lowered or not re.search(r"sha[ -]?256", lowered):
            continue
        candidate_lines.append(line)

        digest_marker = re.search(r"digest\s*[:=]?\s*(.*)$", line, re.IGNORECASE)
        if not digest_marker:
            continue
        suffix = digest_marker.group(1).strip()
        normalized = re.sub(r"[^0-9a-fA-F]", "", suffix).lower()
        if CERT_SHA256_RE.fullmatch(normalized):
            digests.add(normalized)

    if len(digests) == 1:
        return next(iter(digests))

    safe_lines = " | ".join(candidate_lines[:8]) if candidate_lines else "<none>"
    if len(digests) > 1:
        fail(
            "APK signature verification exposed multiple certificate SHA-256 digests; "
            f"refusing ambiguous signer identity. Candidate metadata: {safe_lines}",
        )
    fail(
        "APK signature verification did not expose one parseable signer-certificate "
        f"SHA-256 digest. Candidate metadata: {safe_lines}",
    )
    raise AssertionError("unreachable")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--apk", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--checksum-output", required=True)
    parser.add_argument("--signature-output", required=True)
    args = parser.parse_args()

    apk = (ROOT / args.apk).resolve()
    try:
        apk.relative_to(ROOT)
    except ValueError:
        fail("APK path must remain inside the repository workspace")
    if not apk.is_file():
        fail(f"APK is missing: {apk}")

    repository = require_env("GITHUB_REPOSITORY")
    if repository != EXPECTED_REPOSITORY:
        fail(f"unexpected repository identity: {repository}")

    source_revision = require_env("GOREECLOUD_SOURCE_REVISION").lower()
    if not SHA_RE.fullmatch(source_revision):
        fail("GOREECLOUD_SOURCE_REVISION must be an exact 40-character Git SHA")

    checked_out_revision = run_checked(["git", "rev-parse", "HEAD"]).strip().lower()
    if checked_out_revision != source_revision:
        fail(
            "checked-out source revision does not match the declared exact source "
            f"revision: {checked_out_revision} != {source_revision}",
        )

    aapt = find_build_tool("aapt")
    apksigner = find_build_tool("apksigner")
    badging = run_checked([str(aapt), "dump", "badging", str(apk)])

    package_name = extract_single(r"^package: name='([^']+)'", badging, "package name")
    version_code = extract_single(r"^package: .* versionCode='([^']+)'", badging, "version code")
    version_name = extract_single(r"^package: .* versionName='([^']+)'", badging, "version name")
    min_sdk = extract_single(r"^sdkVersion:'([^']+)'", badging, "minimum SDK")
    target_sdk = extract_single(r"^targetSdkVersion:'([^']+)'", badging, "target SDK")

    expected = {
        "package_name": EXPECTED_PACKAGE,
        "version_code": EXPECTED_VERSION_CODE,
        "version_name": EXPECTED_VERSION_NAME,
        "min_sdk": EXPECTED_MIN_SDK,
        "target_sdk": EXPECTED_TARGET_SDK,
    }
    observed = {
        "package_name": package_name,
        "version_code": version_code,
        "version_name": version_name,
        "min_sdk": min_sdk,
        "target_sdk": target_sdk,
    }
    if observed != expected:
        fail(f"APK identity mismatch: expected {expected}, observed {observed}")
    if "application-debuggable" not in badging:
        fail("APK is not identified as a debuggable Development build")

    signature_report = run_checked(
        [str(apksigner), "verify", "--verbose", "--print-certs", str(apk)],
    )
    signer_digest = extract_signer_certificate_sha256(signature_report)
    verified_scheme = bool(
        re.search(r"Verified using v[1-9][^:]*:\s*true", signature_report, re.IGNORECASE),
    )
    if not verified_scheme:
        fail("APK signature verification did not confirm any Android signature scheme")

    apk_bytes = apk.read_bytes()
    apk_sha256 = hashlib.sha256(apk_bytes).hexdigest()
    apk_size = len(apk_bytes)

    output = (ROOT / args.output).resolve()
    checksum_output = (ROOT / args.checksum_output).resolve()
    signature_output = (ROOT / args.signature_output).resolve()
    for path in (output, checksum_output, signature_output):
        try:
            path.relative_to(ROOT)
        except ValueError:
            fail("evidence output paths must remain inside the repository workspace")
        path.parent.mkdir(parents=True, exist_ok=True)

    evidence = {
        "schema_version": 1,
        "evidence_type": "goreecloud.android-development-apk-build.v1",
        "repository": repository,
        "source_revision": source_revision,
        "lifecycle": "Development",
        "artifact": {
            "name": apk.name,
            "sha256": apk_sha256,
            "size_bytes": apk_size,
            "package_name": package_name,
            "version_name": version_name,
            "version_code": version_code,
            "min_sdk": int(min_sdk),
            "target_sdk": int(target_sdk),
            "debuggable": True,
        },
        "signature": {
            "verified": True,
            "signer_certificate_sha256": signer_digest,
            "classification": "development-debug",
            "release_signing_authority": False,
        },
        "build": {
            "workflow": os.environ.get("GITHUB_WORKFLOW", ""),
            "run_id": os.environ.get("GITHUB_RUN_ID", ""),
            "run_attempt": os.environ.get("GITHUB_RUN_ATTEMPT", ""),
            "job": os.environ.get("GITHUB_JOB", ""),
            "event": os.environ.get("GITHUB_EVENT_NAME", ""),
            "android_build_tools": aapt.parent.name,
        },
        "acceptance": {
            "release_candidate": False,
            "release_accepted": False,
            "production_accepted": False,
            "stable": False,
        },
        "authority_boundary": (
            "Build/package evidence only. Debug APK signature verification is not "
            "GoreeCloud production release-signing authority."
        ),
    }

    output.write_text(json.dumps(evidence, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    checksum_output.write_text(f"{apk_sha256}  {apk.name}\n", encoding="utf-8")
    signature_output.write_text(
        "apk_signature_verified=true\n"
        f"signer_certificate_sha256={signer_digest}\n"
        "signature_classification=development-debug\n"
        "release_signing_authority=false\n",
        encoding="utf-8",
    )
    print(
        "Android Development APK build evidence passed: "
        f"{apk.name} sha256={apk_sha256} source={source_revision}",
    )


if __name__ == "__main__":
    main()
