#!/usr/bin/env python3
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[1]
ANDROID = ROOT / "client" / "android"
CLIENT = ANDROID / "app" / "src" / "main"
MANIFEST = CLIENT / "AndroidManifest.xml"
STRINGS = CLIENT / "res" / "values" / "strings.xml"
APP_BUILD = ANDROID / "app" / "build.gradle.kts"

errors: list[str] = []

if not MANIFEST.is_file():
    errors.append("Messenger Android client manifest is missing")
else:
    manifest = MANIFEST.read_text(encoding="utf-8")
    forbidden_permissions = (
        "android.permission.INTERNET",
        "android.permission.READ_CONTACTS",
        "android.permission.READ_SMS",
        "android.permission.SEND_SMS",
        "android.permission.RECEIVE_SMS",
        "android.permission.CALL_PHONE",
        "android.permission.RECORD_AUDIO",
        "android.permission.CAMERA",
    )
    for permission in forbidden_permissions:
        if permission in manifest:
            errors.append(f"Development client unexpectedly declares {permission}")
    if 'android:allowBackup="false"' not in manifest:
        errors.append("Development client must keep Android backup disabled")

# Network/storage implementation authority must not enter Kotlin source in this
# disconnected Development shell. Android XML namespace URIs are intentionally
# not treated as network authority.
for path in CLIENT.rglob("*.kt"):
    text = path.read_text(encoding="utf-8")
    relative = path.relative_to(ROOT)
    forbidden_fragments = (
        "java.net.",
        "javax.net.",
        "okhttp",
        "retrofit",
        "http://",
        "https://",
        "SharedPreferences",
        "getSharedPreferences(",
        "RoomDatabase",
        "SQLiteDatabase",
    )
    for fragment in forbidden_fragments:
        if fragment in text:
            errors.append(f"{relative}: forbidden Development client authority fragment {fragment!r}")

if APP_BUILD.is_file():
    build_text = APP_BUILD.read_text(encoding="utf-8").lower()
    for dependency_fragment in ("okhttp", "retrofit", "ktor-client", "room-runtime"):
        if dependency_fragment in build_text:
            errors.append(
                f"client/android/app/build.gradle.kts: forbidden Development client dependency {dependency_fragment!r}",
            )
else:
    errors.append("Messenger Android client app build file is missing")

activity = CLIENT / "kotlin" / "com" / "goreecloud" / "messenger" / "client" / "MessengerClientActivity.kt"
if not activity.is_file():
    errors.append("Messenger native Android Development Activity is missing")
else:
    activity_text = activity.read_text(encoding="utf-8")
    for required in (
        "Development boundary",
        "Not Release Candidate",
        "R.string.provenance_heading",
    ):
        if required not in activity_text:
            errors.append(f"Messenger client is missing required visible boundary reference {required!r}")

if not STRINGS.is_file():
    errors.append("Messenger Android client strings resource is missing")
else:
    strings_text = STRINGS.read_text(encoding="utf-8")
    required_labels = (
        "Native Android Development preview",
        "Disconnected shell · No account · No network · No message storage",
        "Provenance examples",
    )
    for required in required_labels:
        if required not in strings_text:
            errors.append(f"Messenger client is missing required visible boundary label {required!r}")

provenance = CLIENT / "kotlin" / "com" / "goreecloud" / "messenger" / "client" / "CommunicationProvenance.kt"
if not provenance.is_file():
    errors.append("Messenger communication provenance contract is missing")
else:
    provenance_text = provenance.read_text(encoding="utf-8")
    if "E2EE_ACTIVE || transport == CommunicationTransport.DATA" not in provenance_text:
        errors.append("Messenger provenance contract does not fail closed on carrier E2EE claims")

if errors:
    print("Messenger Android Development client boundary FAILED:", file=sys.stderr)
    for error in errors:
        print(f"- {error}", file=sys.stderr)
    raise SystemExit(1)

print("Messenger Android Development client boundary passed")
