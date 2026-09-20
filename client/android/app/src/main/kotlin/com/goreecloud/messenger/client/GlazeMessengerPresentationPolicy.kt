package com.goreecloud.messenger.client

/**
 * V1.6 presentation-only policy for the disconnected Messenger Android Development shell.
 *
 * Inputs must come from the caller/platform. The resolver performs no private-content inspection,
 * network lookup, provider precedence, authorization inference, or consequential action.
 */
enum class GlazeMessengerMaterialRole {
    SOLID,
    RAISED,
    FUNCTIONAL_GLASS,
    CLEAR_GLASS,
}

enum class GlazeMessengerPerformanceLevel {
    FULL,
    BALANCED,
    EFFICIENT,
    ESSENTIAL,
}

enum class GlazeMessengerMotionMode {
    STANDARD,
    REDUCED,
    MINIMAL,
}

data class GlazeMessengerPresentationContext(
    val reducedMotion: Boolean = false,
    val reducedTransparency: Boolean = false,
    val increasedContrast: Boolean = false,
    val largeText: Boolean = false,
    val extraLargeText: Boolean = false,
    val touchAssistance: Boolean = false,
    val strongFocus: Boolean = false,
    val keyboardFirst: Boolean = false,
    val screenReaderOptimized: Boolean = false,
    val performanceLevel: GlazeMessengerPerformanceLevel = GlazeMessengerPerformanceLevel.FULL,
)

data class GlazeMessengerResolvedPresentation(
    val materialRole: GlazeMessengerMaterialRole,
    val motionMode: GlazeMessengerMotionMode,
    val minimumInteractionTargetDp: Int,
    val focusRingWidthDp: Int,
    val focusRingOffsetDp: Int,
    val densityMayYieldToReflow: Boolean,
    val strongVisibleFocusRequired: Boolean,
)

object GlazeMessengerPresentationPolicy {
    const val StableVersion = "1.6.0"
    const val StableSourceRevision = "a7180679ea851389e0f3004515f9a25f420e716d"
    const val InheritedFocusRingWidthDp = 2
    const val InheritedFocusRingOffsetDp = 2

    fun resolve(
        requestedMaterial: GlazeMessengerMaterialRole,
        context: GlazeMessengerPresentationContext,
    ): GlazeMessengerResolvedPresentation {
        val material = when {
            context.reducedTransparency &&
                requestedMaterial in setOf(
                    GlazeMessengerMaterialRole.FUNCTIONAL_GLASS,
                    GlazeMessengerMaterialRole.CLEAR_GLASS,
                ) -> GlazeMessengerMaterialRole.SOLID

            context.performanceLevel == GlazeMessengerPerformanceLevel.ESSENTIAL &&
                requestedMaterial != GlazeMessengerMaterialRole.SOLID ->
                GlazeMessengerMaterialRole.SOLID

            context.performanceLevel == GlazeMessengerPerformanceLevel.EFFICIENT &&
                requestedMaterial in setOf(
                    GlazeMessengerMaterialRole.FUNCTIONAL_GLASS,
                    GlazeMessengerMaterialRole.CLEAR_GLASS,
                ) -> GlazeMessengerMaterialRole.RAISED

            else -> requestedMaterial
        }

        val motion = when {
            context.reducedMotion ||
                context.performanceLevel == GlazeMessengerPerformanceLevel.ESSENTIAL ->
                GlazeMessengerMotionMode.MINIMAL

            context.performanceLevel == GlazeMessengerPerformanceLevel.EFFICIENT ->
                GlazeMessengerMotionMode.REDUCED

            else -> GlazeMessengerMotionMode.STANDARD
        }

        return GlazeMessengerResolvedPresentation(
            materialRole = material,
            motionMode = motion,
            minimumInteractionTargetDp = if (context.touchAssistance) {
                GlazeClientTokens.TouchAssistanceFloorDp
            } else {
                GlazeClientTokens.InteractionFloorDp
            },
            focusRingWidthDp = InheritedFocusRingWidthDp,
            focusRingOffsetDp = InheritedFocusRingOffsetDp,
            densityMayYieldToReflow = context.largeText || context.extraLargeText,
            strongVisibleFocusRequired =
                context.strongFocus ||
                    context.keyboardFirst ||
                    context.screenReaderOptimized ||
                    context.increasedContrast,
        )
    }
}
