package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class GlazeClientTokensTest {
    @Test
    fun sourceMappingTargetsCurrentStableV12Authority() {
        assertEquals("1.2.0", GlazeClientTokens.Version)
        assertEquals(
            "f285b9145e27e6e7027b075c37299d101945c272",
            GlazeClientTokens.StableReleaseRevision,
        )
        assertEquals(
            "b0eadf9a60f73d45caffb62ffc7e9e0334cddc97",
            GlazeClientTokens.SourceQualificationAnchor,
        )
        assertEquals(
            "tokens/glaze-v1.2-optical-foundation.candidate.json",
            GlazeClientTokens.OpticalContract,
        )
        assertEquals(
            "neutral-glass-is-material-color-is-accent",
            GlazeClientTokens.MaterialRule,
        )
    }

    @Test
    fun interactionFloorsPreserveV12AccessibilityTargets() {
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
}
