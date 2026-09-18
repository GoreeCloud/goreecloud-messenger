package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class GlazeClientTokensTest {
    @Test
    fun sourceMappingTargetsCurrentStableV151Authority() {
        assertEquals("1.5.1", GlazeClientTokens.Version)
        assertEquals("Contextual + Capability Awareness", GlazeClientTokens.ReleaseTheme)
        assertEquals(
            "98da57064ede0f334627b632bc16801f580331af",
            GlazeClientTokens.StableReleaseRevision,
        )
        assertEquals("css/glaze-v1.4.1.css", GlazeClientTokens.StableWebEntrypoint)
        assertEquals("js/glaze-v1.5.1.mjs", GlazeClientTokens.StableRuntimeEntrypoint)
        assertEquals("1.5.0", GlazeClientTokens.RollbackBaselineVersion)
        assertEquals("neutral-glass-is-material-color-is-accent", GlazeClientTokens.MaterialRule)
        assertEquals("ee1032a0822ab8e103f8afe48e5c1859fde65cc9", GlazeClientTokens.ReviewedV15ImplementationAnchor)
        assertEquals("5b59d0e36950d737dba35b58ae58058684e0831b", GlazeClientTokens.StableQualificationAnchor)
        assertEquals(GlazeClientTokens.Version, GlazeMessengerOptics.Version)
        assertEquals(GlazeClientTokens.StableReleaseRevision, GlazeMessengerOptics.StableRevision)
    }

    @Test
    fun interactionFloorsPreserveInheritedAccessibilityTargets() {
        assertEquals(48, GlazeClientTokens.InteractionFloorDp)
        assertEquals(56, GlazeClientTokens.TouchAssistanceFloorDp)
    }

    @Test
    fun currentLightAndDarkBaseMaterialsRemainNeutral() {
        assertTrue(GlazeClientTokens.isNeutralSubstrate(GlazeClientTokens.LightCanvas))
        assertTrue(GlazeClientTokens.isNeutralSubstrate(GlazeClientTokens.LightSurface))
        assertTrue(GlazeClientTokens.isNeutralSubstrate(GlazeClientTokens.DarkCanvas))
        assertTrue(GlazeClientTokens.isNeutralSubstrate(GlazeClientTokens.DarkSurface))
    }

    @Test
    fun chromaticBaseMaterialWouldFailNeutralityGuard() {
        assertFalse(GlazeClientTokens.isNeutralSubstrate(0xFF0F6B6F))
        assertFalse(GlazeClientTokens.isNeutralSubstrate(0xFFD9A35F))
    }

    @Test
    fun opticalEngineIsLocalBoundedAndInactive() {
        assertTrue(GlazeMessengerOptics.OpticalEngineIsLocalAndDeterministic)
        assertFalse(GlazeMessengerOptics.TelemetryRequired)
        assertFalse(GlazeMessengerOptics.CameraRequired)
        assertFalse(GlazeMessengerOptics.RemoteContextRequired)
        assertEquals(0.0f, GlazeMessengerOptics.EnvironmentalColorMemoryInfluence)
        assertEquals(0.08f, GlazeMessengerOptics.SharedMaximumEnvironmentalMemoryInfluence)
        assertFalse(GlazeMessengerOptics.OpticalEngineAdapterAccepted)
    }

    @Test
    fun opticalAccessibilityPrecedenceFailsClosed() {
        assertTrue(GlazeMessengerOptics.ReducedTransparencyUsesSolidAccessibleTreatment)
        assertTrue(GlazeMessengerOptics.ForcedColorsUsesSolidAccessibleTreatment)
        assertTrue(GlazeMessengerOptics.IncreasedContrastSuppressesDecorativeTintAndWarmth)
        assertFalse(GlazeMessengerOptics.AccessibilityMayBeOverriddenByOpticalContext)
        assertFalse(GlazeMessengerOptics.OpticalContextMayCarrySemanticAuthority)
        assertFalse(GlazeMessengerOptics.ContextCapabilityPresentationMayGrantAuthority)
        assertTrue(GlazeMessengerOptics.ProviderConflictFailsClosed)
        assertFalse(GlazeMessengerOptics.AutomaticConsequentialActionAllowed)
    }

    @Test
    fun privateMessagingAndPlatformStateCannotDriveOptics() {
        assertFalse(GlazeMessengerOptics.MessageContentSamplingAllowed)
        assertFalse(GlazeMessengerOptics.ComposerDraftSamplingAllowed)
        assertFalse(GlazeMessengerOptics.ConversationIdentitySamplingAllowed)
        assertFalse(GlazeMessengerOptics.ParticipantIdentitySamplingAllowed)
        assertFalse(GlazeMessengerOptics.DeliveryReceiptSamplingAllowed)
        assertFalse(GlazeMessengerOptics.TypingPresenceSamplingAllowed)
        assertFalse(GlazeMessengerOptics.E2EEStateSamplingAllowed)
        assertFalse(GlazeMessengerOptics.TransportStateSamplingAllowed)
        assertFalse(GlazeMessengerOptics.IdentitySessionSamplingAllowed)
        assertFalse(GlazeMessengerOptics.PrivacySecurityRecoverySyncStateSamplingAllowed)
        assertFalse(GlazeMessengerOptics.RemoteArtworkOrMediaSamplingAllowed)
    }

    @Test
    fun sharedV151QualificationDoesNotFabricateMessengerAcceptance() {
        assertFalse(GlazeMessengerOptics.PhysicalDeviceAcceptanceEstablished)
        assertFalse(GlazeMessengerOptics.ManualAssistiveTechnologyAcceptanceEstablished)
        assertFalse(GlazeMessengerOptics.HumanVisualExcellenceAcceptanceEstablished)
        assertFalse(GlazeMessengerOptics.RepresentativeRealDevicePerformanceAccepted)
    }
}
