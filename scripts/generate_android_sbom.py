#!/usr/bin/env python3
"""Generate a deterministic CycloneDX SBOM for the Android Development debug APK.

The SBOM is bound to the exact checked-out source revision and exact APK bytes. It
represents the resolved debug runtime dependency closure only; it is Development
build evidence and does not assert release, production, or Stable acceptance.
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
from urllib.parse import quote
from uuid import NAMESPACE_URL, uuid5

ROOT = Path(__file__).resolve().parents[1]
EXPECTED_REPOSITORY = "GoreeCloud/goreecloud-messenger"
EXPECTED_PACKAGE = "com.goreecloud.messenger"
EXPECTED_VERSION_NAME = "0.1.0-dev"
EXPECTED_CONFIGURATION = "debugRuntimeClasspath"
SHA_RE = re.compile(r"^[0-9a-f]{40}$")
COORDINATE_RE = re.compile(
    r"(?P<group>[A-Za-z0-9_.-]+):"
    r"(?P<name>[A-Za-z0-9_.-]+):"
    r"(?P<declared>[A-Za-z0-9_.+\-]+)"
    r"(?:\s+->\s+(?P<selected>[A-Za-z0-9_.+\-]+))?"
)
MAX_REPORT_BYTES = 2 * 1024 * 1024


def fail(message: str) -> None:
    print(f"Android SBOM generation FAILED: {message}", file=sys.stderr)
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


def path_inside_root(value: str, label: str, *, must_exist: bool) -> Path:
    path = (ROOT / value).resolve()
    try:
        path.relative_to(ROOT)
    except ValueError:
        fail(f"{label} path must remain inside the repository workspace")
    if must_exist and not path.is_file():
        fail(f"{label} is missing: {path}")
    return path


def maven_purl(group: str, name: str, version: str) -> str:
    namespace = quote(group, safe=".-_")
    package_name = quote(name, safe=".-_")
    package_version = quote(version, safe=".-_+")
    return f"pkg:maven/{namespace}/{package_name}@{package_version}"


def parse_dependencies(report: str) -> list[dict[str, str]]:
    if f"{EXPECTED_CONFIGURATION} -" not in report:
        fail(f"dependency report does not identify {EXPECTED_CONFIGURATION}")

    versions: dict[tuple[str, str], str] = {}
    for raw_line in report.splitlines():
        match = COORDINATE_RE.search(raw_line)
        if not match:
            continue
        group = match.group("group")
        name = match.group("name")
        version = match.group("selected") or match.group("declared")
        key = (group, name)
        previous = versions.get(key)
        if previous is not None and previous != version:
            fail(
                "resolved dependency report contains conflicting selected versions for "
                f"{group}:{name}: {previous} vs {version}",
            )
        versions[key] = version

    components: list[dict[str, str]] = []
    for (group, name), version in sorted(versions.items()):
        purl = maven_purl(group, name, version)
        components.append(
            {
                "type": "library",
                "bom-ref": purl,
                "group": group,
                "name": name,
                "version": version,
                "purl": purl,
                "scope": "required",
            },
        )
    return components


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--dependency-report", required=True)
    parser.add_argument("--apk", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--checksum-output", required=True)
    args = parser.parse_args()

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

    dependency_report = path_inside_root(
        args.dependency_report,
        "dependency report",
        must_exist=True,
    )
    if dependency_report.stat().st_size > MAX_REPORT_BYTES:
        fail(f"dependency report exceeds {MAX_REPORT_BYTES} bytes")
    report_text = dependency_report.read_text(encoding="utf-8")
    components = parse_dependencies(report_text)

    apk = path_inside_root(args.apk, "APK", must_exist=True)
    apk_bytes = apk.read_bytes()
    apk_sha256 = hashlib.sha256(apk_bytes).hexdigest()

    root_ref = f"pkg:generic/goreecloud-messenger-android@{EXPECTED_VERSION_NAME}?type=android-debug"
    serial_seed = f"{repository}@{source_revision}:{apk_sha256}:{EXPECTED_CONFIGURATION}"
    serial_number = f"urn:uuid:{uuid5(NAMESPACE_URL, serial_seed)}"
    component_refs = [component["bom-ref"] for component in components]

    sbom = {
        "bomFormat": "CycloneDX",
        "specVersion": "1.6",
        "serialNumber": serial_number,
        "version": 1,
        "metadata": {
            "component": {
                "type": "application",
                "bom-ref": root_ref,
                "name": "goreecloud-messenger-android",
                "version": EXPECTED_VERSION_NAME,
                "purl": root_ref,
                "hashes": [
                    {
                        "alg": "SHA-256",
                        "content": apk_sha256,
                    },
                ],
                "properties": [
                    {"name": "goreecloud:android-package", "value": EXPECTED_PACKAGE},
                    {"name": "goreecloud:artifact-name", "value": apk.name},
                    {"name": "goreecloud:lifecycle", "value": "Development"},
                ],
            },
            "properties": [
                {"name": "goreecloud:repository", "value": repository},
                {"name": "goreecloud:source-revision", "value": source_revision},
                {"name": "goreecloud:gradle-configuration", "value": EXPECTED_CONFIGURATION},
                {
                    "name": "goreecloud:dependency-resolution",
                    "value": "resolved-transitive-runtime-closure",
                },
                {"name": "goreecloud:release-candidate", "value": "false"},
                {"name": "goreecloud:production-accepted", "value": "false"},
                {"name": "goreecloud:stable", "value": "false"},
            ],
        },
        "components": components,
        "dependencies": [
            {
                "ref": root_ref,
                "dependsOn": component_refs,
            },
        ],
    }

    output = path_inside_root(args.output, "SBOM output", must_exist=False)
    checksum_output = path_inside_root(
        args.checksum_output,
        "SBOM checksum output",
        must_exist=False,
    )
    output.parent.mkdir(parents=True, exist_ok=True)
    checksum_output.parent.mkdir(parents=True, exist_ok=True)

    serialized = (json.dumps(sbom, indent=2, sort_keys=True) + "\n").encode("utf-8")
    output.write_bytes(serialized)
    sbom_sha256 = hashlib.sha256(serialized).hexdigest()
    checksum_output.write_text(f"{sbom_sha256}  {output.name}\n", encoding="utf-8")

    print(
        "Android Development SBOM generation passed: "
        f"components={len(components)} sha256={sbom_sha256} source={source_revision}",
    )


if __name__ == "__main__":
    main()
