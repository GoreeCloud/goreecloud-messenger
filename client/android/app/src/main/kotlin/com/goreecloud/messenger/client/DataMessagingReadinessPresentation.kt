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
 * persist readiness, or create send authority. Optional evidence is used only to explain why an
 * already-blocked prerequisite is not verified; the evaluated Result remains authoritative.
 */
internal object DataMessagingReadinessPresentationPolicy {
    fun present(
        result: DataMessagingReadiness.Result,
        evidence: DataMessagingReadiness.Evidence? = null,
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
                    detail = identityFailureDetail(evidence),
                ),
                item(
                    label = "Conversation authorization",
                    reason = DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED,
                    missing = missing,
                    ready = ready,
                    detail = conversationFailureDetail(evidence),
                ),
                item(
                    label = "GoreeCloud Data transport",
                    reason = DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE,
                    missing = missing,
                    ready = ready,
                    detail = transportFailureDetail(evidence),
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

    private fun identityFailureDetail(evidence: DataMessagingReadiness.Evidence?): String? =
        when (evidence?.identity) {
            null -> null
            DataMessagingReadiness.IdentityState.UNKNOWN ->
                "No verified GoreeCloud Identity session evidence is available."
            DataMessagingReadiness.IdentityState.UNAUTHENTICATED ->
                "GoreeCloud Identity reports that this client is not authenticated."
            DataMessagingReadiness.IdentityState.AUTHENTICATED ->
                "The evaluated readiness result does not verify GoreeCloud Identity authentication."
        }

    private fun conversationFailureDetail(evidence: DataMessagingReadiness.Evidence?): String? {
        evidence ?: return null
        return when (evidence.conversationAccess) {
            DataMessagingReadiness.ConversationAccessState.UNKNOWN ->
                "No verified conversation-participant authorization evidence is available."
            DataMessagingReadiness.ConversationAccessState.NOT_PARTICIPANT ->
                "This identity is not verified as a participant in the exact conversation."
            DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT -> {
                if (evidence.authorizedConversationId?.let(DataReceiptIdentifierPolicy::canonicalOrNull) == null) {
                    "Conversation authorization is not verified for an exact canonical conversation."
                } else {
                    "The evaluated readiness result does not verify conversation authorization."
                }
            }
        }
    }

    private fun transportFailureDetail(evidence: DataMessagingReadiness.Evidence?): String? =
        when (evidence?.transport) {
            null -> null
            DataMessagingReadiness.DataTransportState.UNKNOWN ->
                "No verified GoreeCloud Data transport availability evidence is available."
            DataMessagingReadiness.DataTransportState.UNAVAILABLE ->
                "GoreeCloud Data transport is reported unavailable."
            DataMessagingReadiness.DataTransportState.AVAILABLE ->
                "The evaluated readiness result does not verify GoreeCloud Data transport availability."
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
