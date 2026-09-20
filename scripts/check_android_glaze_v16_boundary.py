#!/usr/bin/env python3
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[1]
CLIENT = ROOT / "client" / "android" / "app" / "src" / "main" / "kotlin" / "com" / "goreecloud" / "messenger" / "client"
TOKENS = CLIENT / "GlazeClientTokens.kt"
POLICY = CLIENT / "GlazeMessengerPresentationPolicy.kt"
OPTICS = CLIENT / "GlazeMessengerOptics.kt"
ACTIVITY = CLIENT / "MessengerClientActivity.kt"
ANDROID_CONTEXT = CLIENT / "MessengerAndroidGlazeContext.kt"

errors: list[str] = []


def require(text: str, fragment: str, label: str) -> None:
    if fragment not in text:
        errors.append(f"{label}: required fragment missing: {fragment!r}")


def forbid(text: str, fragment: str, label: str) -> None:
    if fragment in text:
        errors.append(f"{label}: forbidden fragment present: {fragment!r}")


if (
    not TOKENS.is_file()
    or not POLICY.is_file()
    or not OPTICS.is_file()
    or not ACTIVITY.is_file()
    or not ANDROID_CONTEXT.is_file()
):
    errors.append("Messenger V1.6 Glaze source files are incomplete")
else:
    tokens = TOKENS.read_text(encoding="utf-8")
    policy = POLICY.read_text(encoding="utf-8")
    optics = OPTICS.read_text(encoding="utf-8")
    activity = ACTIVITY.read_text(encoding="utf-8")
    android_context = ANDROID_CONTEXT.read_text(encoding="utf-8")

    for marker in (
        'const val Version = "1.6.0"',
        'const val StableReleaseRevision = "a7180679ea851389e0f3004515f9a25f420e716d"',
        'const val StableReleaseTag = "v1.6.0"',
        'const val RollbackBaselineVersion = "1.5.1"',
        "const val SharedStableConsumerEligible = true",
        "const val ApplicationAcceptanceAutomatic = false",
        "const val InheritedCoarseInteractionFloorDp = 44",
        "const val InheritedPointerCompactFloorDp = 32",
        "const val InteractionFloorDp = 48",
        "const val TouchAssistanceFloorDp = 56",
        "const val LargeTextScreenGutterDp = 16",
    ):
        require(tokens, marker, "GlazeClientTokens")

    for marker in (
        'const val StableVersion = "1.6.0"',
        'const val StableSourceRevision = "a7180679ea851389e0f3004515f9a25f420e716d"',
        "GlazeMessengerMaterialRole.FUNCTIONAL_GLASS",
        "GlazeMessengerMaterialRole.CLEAR_GLASS",
        "GlazeMessengerMaterialRole.SOLID",
        "GlazeMessengerPerformanceLevel.ESSENTIAL",
        "GlazeMessengerPerformanceLevel.EFFICIENT",
        "GlazeMessengerMotionMode.MINIMAL",
        "context.reducedTransparency",
        "context.reducedMotion",
        "context.largeText || context.extraLargeText",
        "context.keyboardFirst",
    ):
        require(policy, marker, "GlazeMessengerPresentationPolicy")

    for marker in (
        'const val Version = "1.6.0"',
        'const val StableRevision = "a7180679ea851389e0f3004515f9a25f420e716d"',
        "const val OpticalEngineIsLocalAndDeterministic = true",
        "const val TelemetryRequired = false",
        "const val CameraRequired = false",
        "const val RemoteContextRequired = false",
        "const val EnvironmentalColorMemoryInfluence = 0.0f",
        "const val ReducedTransparencyUsesSolidAccessibleTreatment = true",
        "const val ReducedMotionUsesMinimalPresentation = true",
        "const val EssentialPerformanceUsesSolidMinimalPresentation = true",
        "const val StrongFocusCannotBeSuppressed = true",
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

    # The disconnected UI may consume only the pure presentation resolver and privacy-safe
    # platform configuration projection. The optical boundary remains inactive and cannot inspect
    # runtime messaging, Identity, authorization, transport, E2EE, or other authority state.
    for marker in (
        "object MessengerAndroidGlazeContext",
        "fun fromFontScale(fontScale: Float)",
        "largeText = normalized > DefaultFontScale",
        "extraLargeText = normalized >= ExtraLargeTextScale",
    ):
        require(android_context, marker, "MessengerAndroidGlazeContext")
    for forbidden in (
        "GoreeCloudIdentitySessionAuthority",
        "ConversationAuthorizationAuthority",
        "GoreeCloudDataTransportAuthority",
        "E2EESessionAuthority",
        "DataMessagingReadiness",
        "PreparedEncryptedDataMessage",
        "java.net",
        "android.net",
    ):
        forbid(android_context, forbidden, "MessengerAndroidGlazeContext")

    require(activity, "GlazeMessengerPresentationPolicy.resolve(", "MessengerClientActivity")
    require(activity, "MessengerAndroidGlazeContext.fromFontScale(", "MessengerClientActivity")
    require(activity, "resources.configuration.fontScale", "MessengerClientActivity")
    require(activity, "GlazeClientTokens.LargeTextScreenGutterDp", "MessengerClientActivity")
    forbid(activity, "GlazeMessengerOptics", "MessengerClientActivity")

if errors:
    print("Messenger Android GLAZE UI V1.6 boundary FAILED:", file=sys.stderr)
    for error in errors:
        print(f"- {error}", file=sys.stderr)
    raise SystemExit(1)

print("Messenger Android GLAZE UI V1.6 source boundary passed")
