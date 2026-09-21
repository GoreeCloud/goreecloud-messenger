package com.goreecloud.messenger.client

import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class MessengerAndroidGlazeContextTest {
    @Test
    fun defaultFontScaleDoesNotInventLargeTextState() {
        val context = MessengerAndroidGlazeContext.fromFontScale(1.0f)

        assertFalse(context.largeText)
        assertFalse(context.extraLargeText)
    }

    @Test
    fun increasedFontScaleEnablesLargeTextReflow() {
        val context = MessengerAndroidGlazeContext.fromFontScale(1.3f)

        assertTrue(context.largeText)
        assertFalse(context.extraLargeText)
    }

    @Test
    fun twoHundredPercentClassScaleEnablesExtraLargeText() {
        val context = MessengerAndroidGlazeContext.fromFontScale(2.0f)

        assertTrue(context.largeText)
        assertTrue(context.extraLargeText)
    }

    @Test
    fun reducedMotionAndTouchExplorationProjectOnlyPresentationState() {
        val context = MessengerAndroidGlazeContext.fromSignals(
            MessengerAndroidGlazeSignals(
                fontScale = 1.0f,
                animationsEnabled = false,
                touchExplorationEnabled = true,
            ),
        )

        assertTrue(context.reducedMotion)
        assertTrue(context.touchAssistance)
        assertTrue(context.screenReaderOptimized)
        assertFalse(context.largeText)
        assertFalse(context.extraLargeText)
    }

    @Test
    fun ordinaryAndroidSignalsPreserveNeutralAccessibilityState() {
        val context = MessengerAndroidGlazeContext.fromSignals(
            MessengerAndroidGlazeSignals(
                fontScale = 1.0f,
                animationsEnabled = true,
                touchExplorationEnabled = false,
            ),
        )

        assertFalse(context.reducedMotion)
        assertFalse(context.touchAssistance)
        assertFalse(context.screenReaderOptimized)
    }

    @Test
    fun invalidFontScaleFailsClosedToNeutralContext() {
        listOf(Float.NaN, Float.POSITIVE_INFINITY, 0f, -1f).forEach { value ->
            val context = MessengerAndroidGlazeContext.fromFontScale(value)
            assertFalse(context.largeText)
            assertFalse(context.extraLargeText)
        }
    }
}
