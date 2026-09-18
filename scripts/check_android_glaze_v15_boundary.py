#!/usr/bin/env python3
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[1]
CLIENT = ROOT / "client" / "android" / "app" / "src" / "main" / "kotlin" / "com" / "goreecloud" / "messenger" / "client"
TOKENS = CLIENT / "GlazeClientTokens.kt"
OPTICS = CLIENT / "GlazeMessengerOptics.kt"
ACTIVITY = CLIENT / "MessengerClientActivity.kt"

errors: list[str] = []


def require(text: str, fragment: str, label: str) -> None:
    if fragment not in text:
        errors.append(f"{label}: required fragment missing: {fragment!r}")


def forbid(text: str, fragment: str, label: str) -> None:
    if fragment in text:
        errors.append(f"{label}: forbidden fragment present: {fragment!r}")


if not TOKENS.is_file() or not OPTICS.is_file() or not ACTIVITY.is_file():
    errors.append("Messenger V1.5.1 Glaze source files are incomplete")
else:
    tokens = TOKENS.read_text(encoding="utf-8")
    optics = OPTICS.read_text(encoding="utf-8")
    activity = ACTIVITY.read_text(encoding="utf-8")

    for marker in (
        'const val Version = "1.5.1"',
        'const val ReleaseTheme = "Contextual + Capability Awareness"',
        'const val StableReleaseRevision = "98da57064ede0f334627b632bc16801f580331af"',
        'const val StableWebEntrypoint = "css/glaze-v1.4.1.css"',
        'const val StableRuntimeEntrypoint = "js/glaze-v1.5.1.mjs"',
        'const val RollbackBaselineVersion = "1.5.0"',
        'const val ReviewedV15ImplementationAnchor = "ee1032a0822ab8e103f8afe48e5c1859fde65cc9"',
        'const val StableQualificationAnchor = "5b59d0e36950d737dba35b58ae58058684e0831b"',
        "const val InteractionFloorDp = 48",
        "const val TouchAssistanceFloorDp = 56",
    ):
        require(tokens, marker, "GlazeClientTokens")

    for marker in (
        'const val Version = "1.5.1"',
        'const val StableRevision = "98da57064ede0f334627b632bc16801f580331af"',
        "const val OpticalEngineIsLocalAndDeterministic = true",
        "const val TelemetryRequired = false",
        "const val CameraRequired = false",
        "const val RemoteContextRequired = false",
        "const val EnvironmentalColorMemoryInfluence = 0.0f",
        "const val SharedMaximumEnvironmentalMemoryInfluence = 0.08f",
        "const val ReducedTransparencyUsesSolidAccessibleTreatment = true",
        "const val ForcedColorsUsesSolidAccessibleTreatment = true",
        "const val IncreasedContrastSuppressesDecorativeTintAndWarmth = true",
        "const val AccessibilityMayBeOverriddenByOpticalContext = false",
        "const val MessageContentSamplingAllowed = false",
        "const val ComposerDraftSamplingAllowed = false",
        "const val ConversationIdentitySamplingAllowed = false",
        "const val ParticipantIdentitySamplingAllowed = false",
        "const val DeliveryReceiptSamplingAllowed = false",
        "const val TypingPresenceSamplingAllowed = false",
        "const val E2EEStateSamplingAllowed = false",
        "const val TransportStateSamplingAllowed = false",
        "const val IdentitySessionSamplingAllowed = false",
        "const val PrivacySecurityRecoverySyncStateSamplingAllowed = false",
        "const val RemoteArtworkOrMediaSamplingAllowed = false",
        "const val OpticalContextMayCarrySemanticAuthority = false",
        "const val ContextCapabilityPresentationMayGrantAuthority = false",
        "const val ProviderConflictFailsClosed = true",
        "const val AutomaticConsequentialActionAllowed = false",
        "const val OpticalEngineAdapterAccepted = false",
        "const val PhysicalDeviceAcceptanceEstablished = false",
        "const val ManualAssistiveTechnologyAcceptanceEstablished = false",
        "const val HumanVisualExcellenceAcceptanceEstablished = false",
        "const val RepresentativeRealDevicePerformanceAccepted = false",
    ):
        require(optics, marker, "GlazeMessengerOptics")

    # Source mapping must not silently activate contextual/optical authority in
    # the disconnected Development UI.
    forbid(activity, "GlazeMessengerOptics", "MessengerClientActivity")

if errors:
    print("Messenger Android GLAZE UI V1.5.1 boundary FAILED:", file=sys.stderr)
    for error in errors:
        print(f"- {error}", file=sys.stderr)
    raise SystemExit(1)

print("Messenger Android GLAZE UI V1.5.1 source boundary passed")
