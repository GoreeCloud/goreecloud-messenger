#!/usr/bin/env python3
"""Generate fail-closed, privacy-minimized evidence for the Development Android APK.

This script validates the exact checked-out source revision, APK package/version/API
identity, checksum, debuggable classification, APK signature state, and a bound
CycloneDX runtime SBOM. It does not create or assert release signing, Release
Candidate acceptance, production acceptance, or Stable qualification.
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
from uuid import NAMESPACE_URL, uuid5

ROOT = Path(__file__).resolve().parents[1]
EXPECTED_REPOSITORY = "GoreeCloud/goreecloud-messenger"
EXPECTED_PACKAGE = "com.goreecloud.messenger"
EXPECTED_VERSION_NAME = "0.1.0-dev"
EXPECTED_VERSION_CODE = "1"
EXPECTED_MIN_SDK = "26"
EXPECTED_TARGET_SDK = "35"
EXPECTED_SBOM_FORMAT = "CycloneDX"
EXPECTED_SBOM_SPEC = "1.6"
EXPECTED_SBOM_COMPONENT = "goreecloud-messenger-android"
EXPECTED_SBOM_CONFIGURATION = "debugRuntimeClasspath"
EXPECTED_SBOM_ROOT_REF = (
    f"pkg:generic/{EXPECTED_SBOM_COMPONENT}@{EXPECTED_VERSION_NAME}?type=android-debug"
)
SHA_RE = re.compile(r"^[0-9a-f]{40}$")
SHA256_RE = re.compile(r"^[0-9a-f]{64}$")


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
        if SHA256_RE.fullmatch(normalized):
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


def path_inside_root(value: str, label: str, *, must_exist: bool) -> Path:
    path = (ROOT / value).resolve()
    try:
        path.relative_to(ROOT)
    except ValueError:
        fail(f"{label} path must remain inside the repository workspace")
    if must_exist and not path.is_file():
        fail(f"{label} is missing: {path}")
    return path


def unique_properties(items: object, label: str) -> dict[str, str]:
    if not isinstance(items, list):
        fail(f"{label} properties must be a list")
    result: dict[str, str] = {}
    for item in items:
        if not isinstance(item, dict):
            fail(f"{label} property entry must be an object")
        name = item.get("name")
        value = item.get("value")
        if not isinstance(name, str) or not isinstance(value, str) or not name:
            fail(f"{label} property entries require nonempty string name/value")
        if name in result:
            fail(f"{label} contains duplicate property {name}")
        result[name] = value
    return result


def validate_sbom(
    path: Path,
    *,
    repository: str,
    source_revision: str,
    apk_name: str,
    apk_sha256: str,
) -> dict[str, object]:
    sbom_bytes = path.read_bytes()
    sbom_sha256 = hashlib.sha256(sbom_bytes).hexdigest()
    try:
        sbom = json.loads(sbom_bytes)
    except (UnicodeDecodeError, json.JSONDecodeError) as exc:
        fail(f"SBOM is not valid UTF-8 JSON: {exc}")

    if not isinstance(sbom, dict):
        fail("SBOM root must be a JSON object")
    if sbom.get("bomFormat") != EXPECTED_SBOM_FORMAT:
        fail(f"SBOM bomFormat must be {EXPECTED_SBOM_FORMAT}")
    if sbom.get("specVersion") != EXPECTED_SBOM_SPEC:
        fail(f"SBOM specVersion must be {EXPECTED_SBOM_SPEC}")
    if sbom.get("version") != 1:
        fail("SBOM document version must be 1")

    serial_number = sbom.get("serialNumber")
    serial_seed = (
        f"{repository}@{source_revision}:{apk_sha256}:{EXPECTED_SBOM_CONFIGURATION}"
    )
    expected_serial_number = f"urn:uuid:{uuid5(NAMESPACE_URL, serial_seed)}"
    if serial_number != expected_serial_number:
        fail("SBOM serialNumber does not match the exact source/APK/configuration identity")

    metadata = sbom.get("metadata")
    if not isinstance(metadata, dict):
        fail("SBOM metadata must be an object")
    metadata_properties = unique_properties(metadata.get("properties"), "SBOM metadata")
    expected_metadata = {
        "goreecloud:repository": repository,
        "goreecloud:source-revision": source_revision,
        "goreecloud:gradle-configuration": EXPECTED_SBOM_CONFIGURATION,
        "goreecloud:dependency-resolution": "resolved-transitive-runtime-closure",
        "goreecloud:release-candidate": "false",
        "goreecloud:production-accepted": "false",
        "goreecloud:stable": "false",
    }
    for name, expected_value in expected_metadata.items():
        if metadata_properties.get(name) != expected_value:
            fail(f"SBOM metadata property {name} does not match expected value")

    component = metadata.get("component")
    if not isinstance(component, dict):
        fail("SBOM metadata.component must be an object")
    if component.get("type") != "application":
        fail("SBOM root component type must be application")
    if component.get("name") != EXPECTED_SBOM_COMPONENT:
        fail(f"SBOM root component name must be {EXPECTED_SBOM_COMPONENT}")
    if component.get("version") != EXPECTED_VERSION_NAME:
        fail(f"SBOM root component version must be {EXPECTED_VERSION_NAME}")
    if component.get("bom-ref") != EXPECTED_SBOM_ROOT_REF:
        fail("SBOM root component bom-ref does not match the Android Development identity")
    if component.get("purl") != EXPECTED_SBOM_ROOT_REF:
        fail("SBOM root component purl does not match the Android Development identity")

    component_properties = unique_properties(component.get("properties"), "SBOM root component")
    if component_properties.get("goreecloud:android-package") != EXPECTED_PACKAGE:
        fail("SBOM Android package identity does not match the APK contract")
    if component_properties.get("goreecloud:lifecycle") != "Development":
        fail("SBOM lifecycle must remain Development")
    if component_properties.get("goreecloud:artifact-name") != apk_name:
        fail("SBOM artifact name does not match the expected Development APK")

    hashes = component.get("hashes")
    if not isinstance(hashes, list):
        fail("SBOM root component hashes must be a list")
    sha256_hashes = {
        item.get("content", "").lower()
        for item in hashes
        if isinstance(item, dict) and item.get("alg") == "SHA-256"
    }
    if sha256_hashes != {apk_sha256}:
        fail("SBOM root component SHA-256 does not exactly match the built APK")

    components = sbom.get("components")
    if not isinstance(components, list):
        fail("SBOM components must be a list")
    refs: set[str] = set()
    for entry in components:
        if not isinstance(entry, dict):
            fail("SBOM component entries must be objects")
        required_fields = ("type", "bom-ref", "group", "name", "version", "purl")
        for field in required_fields:
            if not isinstance(entry.get(field), str) or not entry.get(field):
                fail(f"SBOM component is missing nonempty string field {field}")
        if entry.get("type") != "library":
            fail("SBOM resolved runtime components must be library entries")
        ref = entry["bom-ref"]
        if entry.get("purl") != ref or not ref.startswith("pkg:maven/"):
            fail("SBOM runtime component purl/bom-ref must be the same Maven package URL")
        if entry.get("scope") != "required":
            fail("SBOM resolved runtime components must use required scope")
        if ref in refs:
            fail(f"SBOM contains duplicate component bom-ref {ref}")
        refs.add(ref)

    dependencies = sbom.get("dependencies")
    if not isinstance(dependencies, list) or len(dependencies) != 1:
        fail("SBOM must contain one root dependency-closure entry")
    root_dependency = dependencies[0]
    if not isinstance(root_dependency, dict):
        fail("SBOM root dependency entry must be an object")
    if root_dependency.get("ref") != component.get("bom-ref"):
        fail("SBOM dependency closure is not rooted at metadata.component")
    depends_on = root_dependency.get("dependsOn")
    if not isinstance(depends_on, list) or set(depends_on) != refs or len(depends_on) != len(refs):
        fail("SBOM root dependency closure does not exactly cover resolved runtime components")

    return {
        "name": path.name,
        "sha256": sbom_sha256,
        "size_bytes": len(sbom_bytes),
        "format": EXPECTED_SBOM_FORMAT,
        "spec_version": EXPECTED_SBOM_SPEC,
        "serial_number": serial_number,
        "runtime_component_count": len(components),
        "configuration": EXPECTED_SBOM_CONFIGURATION,
    }


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--apk", required=True)
    parser.add_argument("--sbom", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--checksum-output", required=True)
    parser.add_argument("--signature-output", required=True)
    args = parser.parse_args()

    apk = path_inside_root(args.apk, "APK", must_exist=True)
    sbom_path = path_inside_root(args.sbom, "SBOM", must_exist=True)

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
    sbom_evidence = validate_sbom(
        sbom_path,
        repository=repository,
        source_revision=source_revision,
        apk_name=apk.name,
        apk_sha256=apk_sha256,
    )

    output = path_inside_root(args.output, "build evidence output", must_exist=False)
    checksum_output = path_inside_root(
        args.checksum_output,
        "APK checksum output",
        must_exist=False,
    )
    signature_output = path_inside_root(
        args.signature_output,
        "signature evidence output",
        must_exist=False,
    )
    for path in (output, checksum_output, signature_output):
        path.parent.mkdir(parents=True, exist_ok=True)

    evidence = {
        "schema_version": 2,
        "evidence_type": "goreecloud.android-development-apk-build.v2",
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
        "sbom": sbom_evidence,
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
            "Build/package/SBOM evidence only. Debug APK signature verification and "
            "Development SBOM identity are not GoreeCloud production release acceptance "
            "or production release-signing authority."
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
        f"{apk.name} sha256={apk_sha256} sbom={sbom_evidence['sha256']} "
        f"source={source_revision}",
    )


if __name__ == "__main__":
    main()
