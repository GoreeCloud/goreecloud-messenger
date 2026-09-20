package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class GlazeClientTokensTest {
    @Test
    fun sourceMappingTargetsExactCurrentStableV16Authority() {
        assertEquals("1.6.0", GlazeClientTokens.Version)
        assertEquals(
            "a7180679ea851389e0f3004515f9a25f420e716d",
            GlazeClientTokens.StableReleaseRevision,
        )
        assertEquals("v1.6.0", GlazeClientTokens.StableReleaseTag)
        assertEquals("1.5.1", GlazeClientTokens.RollbackBaselineVersion)
        assertTrue(GlazeClientTokens.SharedStableConsumerEligible)
        assertFalse(GlazeClientTokens.ApplicationAcceptanceAutomatic)
        assertEquals(GlazeClientTokens.Version, GlazeMessengerOptics.Version)
        assertEquals(GlazeClientTokens.StableReleaseRevision, GlazeMessengerOptics.StableRevision)
        assertEquals(GlazeClientTokens.Version, GlazeMessengerPresentationPolicy.StableVersion)
        assertEquals(
            GlazeClientTokens.StableReleaseRevision,
            GlazeMessengerPresentationPolicy.StableSourceRevision,
        )
    }

    @Test
    fun messengerKeepsStricterTargetsThanInheritedCoarseFloor() {
        assertEquals(44, GlazeClientTokens.InheritedCoarseInteractionFloorDp)
        assertEquals(32, GlazeClientTokens.InheritedPointerCompactFloorDp)
        assertEquals(48, GlazeClientTokens.InteractionFloorDp)
        assertEquals(56, GlazeClientTokens.TouchAssistanceFloorDp)
        assertTrue(
            GlazeClientTokens.InteractionFloorDp >
                GlazeClientTokens.InheritedCoarseInteractionFloorDp,
        )
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
        assertFalse(GlazeMessengerOptics.OpticalEngineAdapterAccepted)
    }

    @Test
    fun V16AccessibilityAndPerformancePrecedenceFailsClosed() {
        assertTrue(GlazeMessengerOptics.ReducedTransparencyUsesSolidAccessibleTreatment)
        assertTrue(GlazeMessengerOptics.ReducedMotionUsesMinimalPresentation)
        assertTrue(GlazeMessengerOptics.EssentialPerformanceUsesSolidMinimalPresentation)
        assertTrue(GlazeMessengerOptics.IncreasedContrastSuppressesDecorativeTintAndWarmth)
        assertTrue(GlazeMessengerOptics.StrongFocusCannotBeSuppressed)
        assertFalse(GlazeMessengerOptics.AccessibilityMayBeOverriddenByOpticalContext)
        assertFalse(GlazeMessengerOptics.OpticalContextMayCarrySemanticAuthority)
        assertFalse(GlazeMessengerOptics.ContextCapabilityPresentationMayGrantAuthority)
        assertTrue(GlazeMessengerOptics.ProviderConflictFailsClosed)
        assertFalse(GlazeMessengerOptics.AutomaticConsequentialActionAllowed)
    }

    @Test
    fun privateMessagingAndPlatformStateCannotDrivePresentation() {
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
    fun sharedV16QualificationDoesNotFabricateMessengerAcceptance() {
        assertFalse(GlazeMessengerOptics.PhysicalDeviceAcceptanceEstablished)
        assertFalse(GlazeMessengerOptics.ManualAssistiveTechnologyAcceptanceEstablished)
        assertFalse(GlazeMessengerOptics.HumanVisualExcellenceAcceptanceEstablished)
        assertFalse(GlazeMessengerOptics.RepresentativeRealDevicePerformanceAccepted)
    }
}
