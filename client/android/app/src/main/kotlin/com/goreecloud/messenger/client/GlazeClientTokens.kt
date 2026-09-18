package com.goreecloud.messenger.client

/**
 * Repository-local Android source mapping for GLAZE UI V1.5.1 Stable.
 *
 * V1.5.1 preserves the reviewed V1.5 Contextual + Capability Awareness behavior over the inherited V1.4.1 optical material baseline. These constants pin the disconnected Development
 * client to the exact current Stable design authority without claiming rendered
 * or native-device acceptance. Presentation must never manufacture data transport,
 * E2EE, Identity, privacy, security, recovery, synchronization, or delivery truth.
 */
object GlazeClientTokens {
    const val Version = "1.5.1"
    const val ReleaseTheme = "Contextual + Capability Awareness"
    const val StableReleaseRevision = "98da57064ede0f334627b632bc16801f580331af"
    const val StableWebEntrypoint = "css/glaze-v1.4.1.css"
    const val StableRuntimeEntrypoint = "js/glaze-v1.5.1.mjs"
    const val RollbackBaselineVersion = "1.5.0"
    const val ReviewedV15ImplementationAnchor = "ee1032a0822ab8e103f8afe48e5c1859fde65cc9"
    const val StableQualificationAnchor = "5b59d0e36950d737dba35b58ae58058684e0831b"
    const val MaterialRule = "neutral-glass-is-material-color-is-accent"

    const val InteractionFloorDp = 48
    const val TouchAssistanceFloorDp = 56
    const val ScreenGutterDp = 20
    const val SurfaceRadiusDp = 22f
    const val SectionSpacingDp = 14

    // This disconnected Development shell currently renders Android System Light/Dark only.
    // Deep Dark, adaptive accents, and optical context remain independently gated rather than
    // inferred or silently activated by the version bump.
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

    /** V1.5.1 preserves the inherited neutral material substrate; context/capability presentation remains non-authoritative. */
    fun isNeutralSubstrate(argb: Long, maximumChannelSpread: Int = 16): Boolean {
        require(maximumChannelSpread >= 0)
        val red = ((argb shr 16) and 0xff).toInt()
        val green = ((argb shr 8) and 0xff).toInt()
        val blue = (argb and 0xff).toInt()
        return maxOf(red, green, blue) - minOf(red, green, blue) <= maximumChannelSpread
    }
}
