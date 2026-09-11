#!/usr/bin/env python3
"""Regression tests for Android Development SBOM dependency parsing."""

from __future__ import annotations

import importlib.util
from pathlib import Path
import unittest

SCRIPT = Path(__file__).with_name("generate_android_sbom.py")
SPEC = importlib.util.spec_from_file_location("generate_android_sbom", SCRIPT)
if SPEC is None or SPEC.loader is None:
    raise RuntimeError("unable to load generate_android_sbom.py")
MODULE = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(MODULE)


class ParseDependenciesTest(unittest.TestCase):
    def test_resolved_runtime_components_are_sorted_and_selected(self) -> None:
        report = """debugRuntimeClasspath - Runtime classpath of /debug.
+--- org.jetbrains.kotlin:kotlin-stdlib:2.0.0 -> 2.1.10
\\--- org.jetbrains:annotations:13.0
"""
        components = MODULE.parse_dependencies(report)
        self.assertEqual(
            [(c["group"], c["name"], c["version"]) for c in components],
            [
                ("org.jetbrains", "annotations", "13.0"),
                ("org.jetbrains.kotlin", "kotlin-stdlib", "2.1.10"),
            ],
        )

    def test_failed_resolution_fails_closed(self) -> None:
        report = """debugRuntimeClasspath - Runtime classpath of /debug.
\\--- example:missing:1.0 FAILED
"""
        with self.assertRaises(SystemExit):
            MODULE.parse_dependencies(report)

    def test_empty_resolution_fails_closed(self) -> None:
        report = "debugRuntimeClasspath - Runtime classpath of /debug.\nNo dependencies\n"
        with self.assertRaises(SystemExit):
            MODULE.parse_dependencies(report)

    def test_conflicting_selected_versions_fail_closed(self) -> None:
        report = """debugRuntimeClasspath - Runtime classpath of /debug.
+--- example:library:1.0
\\--- example:library:2.0
"""
        with self.assertRaises(SystemExit):
            MODULE.parse_dependencies(report)


if __name__ == "__main__":
    unittest.main()
