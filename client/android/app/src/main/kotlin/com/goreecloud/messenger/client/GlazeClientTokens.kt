package com.goreecloud.messenger.client

/**
 * Repository-local Android source mapping for GLAZE UI V1.3 Stable.
 *
 * Adaptive Resonance extends the inherited V1.2 neutral-glass foundation with
 * bounded adaptive expression, ergonomic composition, and resilience behavior.
 * These constants pin the Development client to the exact current Stable design
 * authority without claiming rendered/native-device acceptance. Presentation
 * must not manufacture transport, E2EE, Identity, privacy, security, recovery,
 * or delivery truth.
 */
object GlazeClientTokens {
    const val Version = "1.3.0"
    const val ReleaseTheme = "Adaptive Resonance"
    const val StableReleaseRevision = "fc7cc91d2eace8da2371371c2855c24cbcb326a1"
    const val OpticalContract = "tokens/glaze-v1.2-optical-foundation.candidate.json"
    const val AdaptiveContract = "contracts/v1.3/adaptive-resonance.plan.json"
    const val StableWebEntrypoint = "css/glaze-v1.3.0.css"
    const val StableRuntimeEntrypoint = "js/glaze-v1.3.0.mjs"
    const val RollbackBaselineVersion = "1.2.0"
    const val MaterialRule = "neutral-glass-is-material-color-is-accent"

    const val InteractionFloorDp = 48
    const val TouchAssistanceFloorDp = 56
    const val ScreenGutterDp = 20
    const val SurfaceRadiusDp = 22f
    const val SectionSpacingDp = 14

    // This Development shell currently renders Android System Light/Dark only. Deep Dark and
    // adaptive environmental accent behavior remain independently acceptance-gated rather than
    // being inferred or silently activated.
    const val LightCanvas = 0xFFF4F6F7
    const val LightSurface = 0xEFFFFFFF
    const val LightText = 0xFF182126
    const val LightMutedText = 0xFF5B666C
    const val LightBorder = 0x33727E84
    const val DarkCanvas = 0xFF101619
    const val DarkSurface = 0xE61B2327
    const val DarkText = 0xFFF3F6F7
    const val DarkMutedText = 0xFFB6C0C5
    const val DarkBorder = 0x337E8A90

    /** The V1.3 material substrate remains neutral; adaptive color stays bounded to accent. */
    fun isNeutralSubstrate(argb: Long, maximumChannelSpread: Int = 16): Boolean {
        require(maximumChannelSpread >= 0)
        val red = ((argb shr 16) and 0xff).toInt()
        val green = ((argb shr 8) and 0xff).toInt()
        val blue = (argb and 0xff).toInt()
        return maxOf(red, green, blue) - minOf(red, green, blue) <= maximumChannelSpread
    }
}
