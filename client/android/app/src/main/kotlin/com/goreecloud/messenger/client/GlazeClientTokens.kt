package com.goreecloud.messenger.client

/**
 * Repository-local Android source mapping for the exact GLAZE UI V1.6 Stable release source.
 *
 * Shared V1.6 consumer eligibility does not establish Messenger-local rendered, accessibility,
 * device, performance, production, release, or Stable acceptance. These values are presentation
 * inputs only and cannot manufacture Data transport, E2EE, Identity, privacy, security, recovery,
 * synchronization, delivery, policy, observability, or other runtime authority.
 */
object GlazeClientTokens {
    const val Version = "1.6.0"
    const val StableReleaseRevision = "a7180679ea851389e0f3004515f9a25f420e716d"
    const val StableReleaseTag = "v1.6.0"
    const val RollbackBaselineVersion = "1.5.1"

    const val SharedStableConsumerEligible = true
    const val ApplicationAcceptanceAutomatic = false

    // Inherited Stable V1.6 layout source floors.
    const val InheritedCoarseInteractionFloorDp = 44
    const val InheritedPointerCompactFloorDp = 32

    // Messenger deliberately keeps stricter touch-oriented product targets.
    const val InteractionFloorDp = 48
    const val TouchAssistanceFloorDp = 56

    // Messenger-owned composition values; not claimed as canonical V1.6 tokens.
    const val ScreenGutterDp = 20
    const val LargeTextScreenGutterDp = 16
    const val SurfaceRadiusDp = 22f
    const val SectionSpacingDp = 14

    // Messenger-owned neutral Development palettes. V1.6 semantics govern presentation behavior;
    // these pigments are not represented as exact shared token values.
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

    fun isNeutralSubstrate(argb: Long, maximumChannelSpread: Int = 16): Boolean {
        require(maximumChannelSpread >= 0)
        val red = ((argb shr 16) and 0xff).toInt()
        val green = ((argb shr 8) and 0xff).toInt()
        val blue = (argb and 0xff).toInt()
        return maxOf(red, green, blue) - minOf(red, green, blue) <= maximumChannelSpread
    }
}
