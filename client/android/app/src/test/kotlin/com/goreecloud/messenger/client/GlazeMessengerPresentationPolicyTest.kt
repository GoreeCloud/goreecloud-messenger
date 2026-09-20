package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class GlazeMessengerPresentationPolicyTest {
    @Test
    fun exactStableSourceAndFocusGeometryArePinned() {
        assertEquals("1.6.0", GlazeMessengerPresentationPolicy.StableVersion)
        assertEquals(
            "a7180679ea851389e0f3004515f9a25f420e716d",
            GlazeMessengerPresentationPolicy.StableSourceRevision,
        )
        assertEquals(2, GlazeMessengerPresentationPolicy.InheritedFocusRingWidthDp)
        assertEquals(2, GlazeMessengerPresentationPolicy.InheritedFocusRingOffsetDp)
    }

    @Test
    fun reducedTransparencyFailsGlassDownToSolid() {
        val resolved = GlazeMessengerPresentationPolicy.resolve(
            requestedMaterial = GlazeMessengerMaterialRole.FUNCTIONAL_GLASS,
            context = GlazeMessengerPresentationContext(reducedTransparency = true),
        )

        assertEquals(GlazeMessengerMaterialRole.SOLID, resolved.materialRole)
        assertEquals(GlazeMessengerMotionMode.STANDARD, resolved.motionMode)
    }

    @Test
    fun essentialPerformanceUsesSolidMinimalPresentation() {
        val resolved = GlazeMessengerPresentationPolicy.resolve(
            requestedMaterial = GlazeMessengerMaterialRole.RAISED,
            context = GlazeMessengerPresentationContext(
                performanceLevel = GlazeMessengerPerformanceLevel.ESSENTIAL,
            ),
        )

        assertEquals(GlazeMessengerMaterialRole.SOLID, resolved.materialRole)
        assertEquals(GlazeMessengerMotionMode.MINIMAL, resolved.motionMode)
    }

    @Test
    fun reducedMotionPreservesConservativeTouchTarget() {
        val resolved = GlazeMessengerPresentationPolicy.resolve(
            requestedMaterial = GlazeMessengerMaterialRole.SOLID,
            context = GlazeMessengerPresentationContext(
                reducedMotion = true,
                touchAssistance = true,
            ),
        )

        assertEquals(GlazeMessengerMotionMode.MINIMAL, resolved.motionMode)
        assertEquals(56, resolved.minimumInteractionTargetDp)
    }

    @Test
    fun largeTextYieldsDensityAndKeyboardRequiresStrongFocus() {
        val resolved = GlazeMessengerPresentationPolicy.resolve(
            requestedMaterial = GlazeMessengerMaterialRole.SOLID,
            context = GlazeMessengerPresentationContext(
                largeText = true,
                keyboardFirst = true,
            ),
        )

        assertTrue(resolved.densityMayYieldToReflow)
        assertTrue(resolved.strongVisibleFocusRequired)
    }

    @Test
    fun neutralContextDoesNotInventAccessibilityOrPerformanceState() {
        val context = GlazeMessengerPresentationContext()
        assertFalse(context.reducedMotion)
        assertFalse(context.reducedTransparency)
        assertFalse(context.touchAssistance)
        assertEquals(GlazeMessengerPerformanceLevel.FULL, context.performanceLevel)
    }
}
