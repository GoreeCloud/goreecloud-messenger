package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class GlazeClientTokensTest {
    @Test
    fun sourceMappingTargetsCurrentStableV13Authority() {
        assertEquals("1.3.0", GlazeClientTokens.Version)
        assertEquals("Adaptive Resonance", GlazeClientTokens.ReleaseTheme)
        assertEquals(
            "fc7cc91d2eace8da2371371c2855c24cbcb326a1",
            GlazeClientTokens.StableReleaseRevision,
        )
        assertEquals(
            "tokens/glaze-v1.2-optical-foundation.candidate.json",
            GlazeClientTokens.OpticalContract,
        )
        assertEquals(
            "contracts/v1.3/adaptive-resonance.plan.json",
            GlazeClientTokens.AdaptiveContract,
        )
        assertEquals("css/glaze-v1.3.0.css", GlazeClientTokens.StableWebEntrypoint)
        assertEquals("js/glaze-v1.3.0.mjs", GlazeClientTokens.StableRuntimeEntrypoint)
        assertEquals("1.2.0", GlazeClientTokens.RollbackBaselineVersion)
        assertEquals(
            "neutral-glass-is-material-color-is-accent",
            GlazeClientTokens.MaterialRule,
        )
    }

    @Test
    fun interactionFloorsPreserveV13AccessibilityTargets() {
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
