package com.goreecloud.messenger.client

/**
 * Privacy-safe Android runtime projection into the Messenger-local GLAZE UI V1.6 presentation
 * context.
 *
 * Only platform configuration supplied by Android is accepted here. No message, conversation,
 * Identity, transport, cryptographic, privacy, security, recovery, policy, observability, or
 * remote-content state may enter this projection.
 */
object MessengerAndroidGlazeContext {
    const val DefaultFontScale = 1.0f
    const val ExtraLargeTextScale = 2.0f

    fun fromFontScale(fontScale: Float): GlazeMessengerPresentationContext {
        val normalized = fontScale.takeIf { it.isFinite() && it > 0f } ?: DefaultFontScale
        return GlazeMessengerPresentationContext(
            largeText = normalized > DefaultFontScale,
            extraLargeText = normalized >= ExtraLargeTextScale,
        )
    }
}
