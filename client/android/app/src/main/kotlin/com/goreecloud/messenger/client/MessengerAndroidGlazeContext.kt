package com.goreecloud.messenger.client

/**
 * Privacy-safe Android runtime projection into the Messenger-local GLAZE UI V1.6 presentation
 * context.
 *
 * Only platform configuration supplied by Android is accepted here. No message, conversation,
 * Identity, transport, cryptographic, privacy, security, recovery, policy, observability, or
 * remote-content state may enter this projection.
 */
data class MessengerAndroidGlazeSignals(
    val fontScale: Float,
    val animationsEnabled: Boolean,
    val touchExplorationEnabled: Boolean,
)

object MessengerAndroidGlazeContext {
    const val DefaultFontScale = 1.0f
    const val ExtraLargeTextScale = 2.0f

    fun fromSignals(signals: MessengerAndroidGlazeSignals): GlazeMessengerPresentationContext {
        val normalized = signals.fontScale
            .takeIf { it.isFinite() && it > 0f }
            ?: DefaultFontScale

        return GlazeMessengerPresentationContext(
            reducedMotion = !signals.animationsEnabled,
            largeText = normalized > DefaultFontScale,
            extraLargeText = normalized >= ExtraLargeTextScale,
            touchAssistance = signals.touchExplorationEnabled,
            screenReaderOptimized = signals.touchExplorationEnabled,
        )
    }

    fun fromFontScale(fontScale: Float): GlazeMessengerPresentationContext =
        fromSignals(
            MessengerAndroidGlazeSignals(
                fontScale = fontScale,
                animationsEnabled = true,
                touchExplorationEnabled = false,
            ),
        )
}
