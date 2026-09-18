package com.goreecloud.messenger.client

/**
 * Source-level GLAZE UI V1.5.1 Contextual + Capability Awareness boundary for Messenger's
 * disconnected Android Development shell.
 *
 * This policy intentionally accepts no conversation/message/editor/identity/security context.
 * Optical presentation remains local, deterministic, accessibility-subordinate, and unable to
 * become messaging or platform authority. Shared Glaze V1.5.1 qualification does not establish
 * Messenger-local rendered, accessibility, device, performance, release, or production acceptance.
 */
object GlazeMessengerOptics {
    const val Version = "1.5.1"
    const val StableRevision = "98da57064ede0f334627b632bc16801f580331af"

    const val OpticalEngineIsLocalAndDeterministic = true
    const val TelemetryRequired = false
    const val CameraRequired = false
    const val RemoteContextRequired = false

    // Messenger uses no environmental memory in this disconnected Development shell.
    const val EnvironmentalColorMemoryInfluence = 0.0f
    const val SharedMaximumEnvironmentalMemoryInfluence = 0.08f

    const val ReducedTransparencyUsesSolidAccessibleTreatment = true
    const val ForcedColorsUsesSolidAccessibleTreatment = true
    const val IncreasedContrastSuppressesDecorativeTintAndWarmth = true
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
