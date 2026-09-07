package com.goreecloud.messenger.client

/**
 * Repository-local Android source mapping for GLAZE UI V1.2 Stable.
 *
 * These constants keep the disconnected Development client on the current Stable design-system
 * authority without claiming rendered/native-device acceptance. Presentation must not manufacture
 * transport, E2EE, Identity, privacy, security, recovery, or delivery truth.
 */
object GlazeClientTokens {
    const val Version = "1.2.0"
    const val StableReleaseRevision = "f285b9145e27e6e7027b075c37299d101945c272"
    const val SourceQualificationAnchor = "b0eadf9a60f73d45caffb62ffc7e9e0334cddc97"
    const val OpticalContract = "tokens/glaze-v1.2-optical-foundation.candidate.json"
    const val MaterialRule = "neutral-glass-is-material-color-is-accent"

    const val InteractionFloorDp = 48
    const val TouchAssistanceFloorDp = 56
    const val ScreenGutterDp = 20
    const val SurfaceRadiusDp = 22f
    const val SectionSpacingDp = 14

    // This Development shell currently renders only Android System Light/Dark. Deep Dark remains a
    // downstream acceptance/mapping gate rather than being silently aliased to ordinary Dark.
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

    /** V1.2 base material should remain neutral rather than becoming semantic/accent color. */
    fun isNeutralSubstrate(argb: Long, maximumChannelSpread: Int = 16): Boolean {
        require(maximumChannelSpread >= 0)
        val red = ((argb shr 16) and 0xff).toInt()
        val green = ((argb shr 8) and 0xff).toInt()
        val blue = (argb and 0xff).toInt()
        return maxOf(red, green, blue) - minOf(red, green, blue) <= maximumChannelSpread
    }
}
