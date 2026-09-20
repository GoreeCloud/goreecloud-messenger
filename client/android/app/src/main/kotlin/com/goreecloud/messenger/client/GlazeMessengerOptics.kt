package com.goreecloud.messenger.client

/**
 * Fail-closed GLAZE UI V1.6 presentation boundary for Messenger's disconnected Android shell.
 *
 * The current shell accepts no conversation/message/editor/identity/security content as an
 * appearance input. Presentation remains local and deterministic and cannot become messaging or
 * platform authority. Shared V1.6 Stable acceptance does not establish Messenger-local acceptance.
 */
object GlazeMessengerOptics {
    const val Version = "1.6.0"
    const val StableRevision = "a7180679ea851389e0f3004515f9a25f420e716d"

    const val OpticalEngineIsLocalAndDeterministic = true
    const val TelemetryRequired = false
    const val CameraRequired = false
    const val RemoteContextRequired = false

    // The disconnected shell uses no environmental-memory input.
    const val EnvironmentalColorMemoryInfluence = 0.0f

    const val ReducedTransparencyUsesSolidAccessibleTreatment = true
    const val ReducedMotionUsesMinimalPresentation = true
    const val EssentialPerformanceUsesSolidMinimalPresentation = true
    const val IncreasedContrastSuppressesDecorativeTintAndWarmth = true
    const val StrongFocusCannotBeSuppressed = true
    const val AccessibilityMayBeOverriddenByOpticalContext = false

    const val MessageContentSamplingAllowed = false
    const val ComposerDraftSamplingAllowed = false
    const val ConversationIdentitySamplingAllowed = false
    const val ParticipantIdentitySamplingAllowed = false
    const val DeliveryReceiptSamplingAllowed = false
    const val TypingPresenceSamplingAllowed = false
    const val E2EEStateSamplingAllowed = false
    const val TransportStateSamplingAllowed = false
    const val IdentitySessionSamplingAllowed = false
    const val PrivacySecurityRecoverySyncStateSamplingAllowed = false
    const val RemoteArtworkOrMediaSamplingAllowed = false

    const val OpticalContextMayCarrySemanticAuthority = false
    const val ContextCapabilityPresentationMayGrantAuthority = false
    const val ProviderConflictFailsClosed = true
    const val AutomaticConsequentialActionAllowed = false
    const val OpticalEngineAdapterAccepted = false
    const val PhysicalDeviceAcceptanceEstablished = false
    const val ManualAssistiveTechnologyAcceptanceEstablished = false
    const val HumanVisualExcellenceAcceptanceEstablished = false
    const val RepresentativeRealDevicePerformanceAccepted = false
}
