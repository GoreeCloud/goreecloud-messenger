#!/usr/bin/env python3
"""Generate exact-source Android Development dependency/SBOM evidence.

The input is the Gradle `:app:dependencies --configuration debugRuntimeClasspath`
report captured by CI. The output is a deterministic CycloneDX 1.6 JSON document
for the resolved Android debug runtime classpath. This is inventory evidence only;
it does not claim vulnerability coverage, license review, Wardveil acceptance,
release acceptance, production readiness, or Stable qualification.
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

ROOT = Path(__file__).resolve().parents[1]
EXPECTED_REPOSITORY = "GoreeCloud/goreecloud-messenger"
EXPECTED_CONFIGURATION = "debugRuntimeClasspath"
EXPECTED_APP_NAME = "GoreeCloud Messenger Android"
EXPECTED_APP_VERSION = "0.1.0-dev"
SHA_RE = re.compile(r"^[0-9a-f]{40}$")
COORDINATE_RE = re.compile(
    r"(?<![A-Za-z0-9_.-])"
    r"([A-Za-z0-9_.-]+):([A-Za-z0-9_.-]+):([^\s()]+)"
)


def fail(message: str) -> None:
    print(f"Android dependency SBOM FAILED: {message}", file=sys.stderr)
    raise SystemExit(1)


def require_env(name: str) -> str:
    value = os.environ.get(name, "").strip()
    if not value:
        fail(f"required environment variable {name} is missing")
    return value


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


def workspace_path(raw: str, label: str) -> Path:
    path = (ROOT / raw).resolve()
    try:
        path.relative_to(ROOT)
    except ValueError:
        fail(f"{label} path must remain inside the repository workspace")
    return path


def normalize_version(raw: str) -> str:
    version = raw.rstrip(",")
    version = version.removesuffix("(*)")
    version = version.removesuffix("(c)")
    return version.strip()


def parse_components(report: str) -> list[tuple[str, str, str]]:
    if EXPECTED_CONFIGURATION not in report:
        fail(f"dependency report does not identify {EXPECTED_CONFIGURATION}")
    if re.search(r"\bFAILED\b", report, re.IGNORECASE):
        fail("dependency report contains FAILED resolution state")

    components: set[tuple[str, str, str]] = set()
    unresolved_lines: list[str] = []

    for raw_line in report.splitlines():
        line = raw_line.strip()
        if not line:
            continue
        if "FAILED" in line.upper():
            unresolved_lines.append(line)
            continue

        matches = list(COORDINATE_RE.finditer(line))
        if not matches:
            continue

        # If Gradle prints `requested -> selected` with two full coordinates,
        # the final coordinate is authoritative. If it prints only a selected
        # version token after `->`, use that selected version instead.
        selected = matches[-1]
        group, name, version = selected.groups()
        tail = line[selected.end():]
        arrow = re.search(r"->\s*([^\s()]+)", tail)
        if arrow and len(matches) == 1:
            version = arrow.group(1)

        version = normalize_version(version)
        if not group or not name or not version or version in {"unspecified", "FAILED"}:
            unresolved_lines.append(line)
            continue
        components.add((group, name, version))

    if unresolved_lines:
        safe = " | ".join(unresolved_lines[:8])
        fail(f"dependency report contains unresolved coordinates: {safe}")
    if not components:
        fail("dependency report did not contain any resolved runtime components")

    return sorted(components)


def component_entry(group: str, name: str, version: str) -> dict[str, object]:
    encoded_group = quote(group, safe=".")
    encoded_name = quote(name, safe=".-_")
    encoded_version = quote(version, safe=".-_+")
    purl = f"pkg:maven/{encoded_group}/{encoded_name}@{encoded_version}"
    return {
        "type": "library",
        "group": group,
        "name": name,
        "version": version,
        "bom-ref": purl,
        "purl": purl,
    }


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--report", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--checksum-output", required=True)
    args = parser.parse_args()

    report_path = workspace_path(args.report, "dependency report")
    output_path = workspace_path(args.output, "SBOM output")
    checksum_path = workspace_path(args.checksum_output, "SBOM checksum output")
    if not report_path.is_file():
        fail(f"dependency report is missing: {report_path}")

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

    report = report_path.read_text(encoding="utf-8")
    components = parse_components(report)

    bom: dict[str, object] = {
        "bomFormat": "CycloneDX",
        "specVersion": "1.6",
        "version": 1,
        "metadata": {
            "component": {
                "type": "application",
                "name": EXPECTED_APP_NAME,
                "version": EXPECTED_APP_VERSION,
            },
            "properties": [
                {"name": "goreecloud.repository", "value": repository},
                {"name": "goreecloud.source_revision", "value": source_revision},
                {"name": "goreecloud.lifecycle", "value": "Development"},
                {"name": "goreecloud.gradle_configuration", "value": EXPECTED_CONFIGURATION},
                {"name": "goreecloud.vulnerability_coverage", "value": "not-assessed"},
                {"name": "goreecloud.license_review", "value": "not-assessed"},
                {"name": "goreecloud.release_candidate", "value": "false"},
                {"name": "goreecloud.production_accepted", "value": "false"},
                {"name": "goreecloud.stable", "value": "false"},
            ],
        },
        "components": [component_entry(*component) for component in components],
    }

    output_path.parent.mkdir(parents=True, exist_ok=True)
    checksum_path.parent.mkdir(parents=True, exist_ok=True)
    serialized = json.dumps(bom, indent=2, sort_keys=True) + "\n"
    output_path.write_text(serialized, encoding="utf-8")
    digest = hashlib.sha256(serialized.encode("utf-8")).hexdigest()
    checksum_path.write_text(f"{digest}  {output_path.name}\n", encoding="utf-8")

    print(
        "Android Development dependency SBOM passed: "
        f"components={len(components)} sha256={digest} source={source_revision}",
    )


if __name__ == "__main__":
    main()
