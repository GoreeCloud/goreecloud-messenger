package com.goreecloud.messenger.client

internal data class DataMessagingReadinessChecklistItem(
    val label: String,
    val verified: Boolean,
    val detail: String? = null,
)

internal data class DataMessagingReadinessPresentation(
    val summary: String,
    val checklist: List<DataMessagingReadinessChecklistItem>,
) {
    fun body(): String = buildString {
        append(summary)
        checklist.forEach { item ->
            append('\n')
            append(if (item.verified) "Verified — " else "Not verified — ")
            append(item.label)
            if (!item.verified && item.detail != null) {
                append('\n')
                append("Why — ")
                append(item.detail)
            }
        }
    }
}

/**
 * Presentation-only projection of the existing fail-closed Data messaging readiness result.
 *
 * It does not authenticate, authorize, connect transport, inspect keys, derive cryptographic state,
 * persist readiness, or create send authority. It only makes the already-evaluated prerequisite
 * result legible in a stable order for the disconnected Development shell.
 */
internal object DataMessagingReadinessPresentationPolicy {
    fun present(
        result: DataMessagingReadiness.Result,
        e2eeFailure: E2EEAcceptanceFailure? = null,
    ): DataMessagingReadinessPresentation {
        val missing = (result as? DataMessagingReadiness.Result.Blocked)?.reasons.orEmpty()
        val ready = result is DataMessagingReadiness.Result.Ready

        return DataMessagingReadinessPresentation(
            summary = if (ready) {
                val provenance = (result as DataMessagingReadiness.Result.Ready).provenance.displayLabel()
                "Data messaging prerequisites are independently verified ($provenance)."
            } else {
                "Send remains unavailable until every prerequisite is independently verified. No Send action is exposed."
            },
            checklist = listOf(
                item(
                    label = "Identity authentication",
                    reason = DataMessagingReadiness.BlockReason.IDENTITY_NOT_AUTHENTICATED,
                    missing = missing,
                    ready = ready,
                ),
                item(
                    label = "Conversation authorization",
                    reason = DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED,
                    missing = missing,
                    ready = ready,
                ),
                item(
                    label = "GoreeCloud Data transport",
                    reason = DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE,
                    missing = missing,
                    ready = ready,
                ),
                item(
                    label = "Verified active E2EE",
                    reason = DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE,
                    missing = missing,
                    ready = ready,
                    detail = e2eeFailure?.presentationReason(),
                ),
            ),
        )
    }

    private fun item(
        label: String,
        reason: DataMessagingReadiness.BlockReason,
        missing: Set<DataMessagingReadiness.BlockReason>,
        ready: Boolean,
        detail: String? = null,
    ): DataMessagingReadinessChecklistItem {
        val verified = ready || reason !in missing
        return DataMessagingReadinessChecklistItem(
            label = label,
            verified = verified,
            detail = detail.takeUnless { verified },
        )
    }
}
