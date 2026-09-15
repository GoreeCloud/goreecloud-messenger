package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class GlazeClientTokensTest {
    @Test
    fun sourceMappingTargetsCurrentStableV141Authority() {
        assertEquals("1.4.1", GlazeClientTokens.Version)
        assertEquals("Optical Hardening", GlazeClientTokens.ReleaseTheme)
        assertEquals(
            "4fab9da0fad2e5c974e0e66ec88632c61745751c",
            GlazeClientTokens.StableReleaseRevision,
        )
        assertEquals("css/glaze-v1.4.1.css", GlazeClientTokens.StableWebEntrypoint)
        assertEquals("js/glaze-v1.4.1.mjs", GlazeClientTokens.StableRuntimeEntrypoint)
        assertEquals("1.4.0", GlazeClientTokens.RollbackBaselineVersion)
        assertEquals("neutral-glass-is-material-color-is-accent", GlazeClientTokens.MaterialRule)
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
    fun sharedV141QualificationDoesNotFabricateMessengerAcceptance() {
        assertFalse(GlazeMessengerOptics.PhysicalDeviceAcceptanceEstablished)
        assertFalse(GlazeMessengerOptics.ManualAssistiveTechnologyAcceptanceEstablished)
        assertFalse(GlazeMessengerOptics.HumanVisualExcellenceAcceptanceEstablished)
        assertFalse(GlazeMessengerOptics.RepresentativeRealDevicePerformanceAccepted)
    }
}
