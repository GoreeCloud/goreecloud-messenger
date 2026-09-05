package com.goreecloud.messenger.client

/**
 * A delivery/read receipt event that has already crossed the future transport/authentication
 * boundary.
 *
 * This value does not authenticate a recipient and must not be constructed from untrusted network
 * input directly. Its only responsibility is carrying the already-derived scope needed by the
 * local projection so state for one message/conversation cannot be reused for another.
 */
data class DataReceiptEvent private constructor(
    val messageId: String,
    val conversationId: String,
    val recipientId: String,
    val stage: Stage,
) {
    enum class Stage(val rank: Int) {
        DELIVERED(1),
        READ(2),
    }

    companion object {
        fun create(
            messageId: String,
            conversationId: String,
            recipientId: String,
            stage: Stage,
        ): DataReceiptEvent {
            require(messageId.isNotBlank()) { "messageId must not be blank" }
            require(conversationId.isNotBlank()) { "conversationId must not be blank" }
            require(recipientId.isNotBlank()) { "recipientId must not be blank" }

            return DataReceiptEvent(
                messageId = messageId.trim(),
                conversationId = conversationId.trim(),
                recipientId = recipientId.trim(),
                stage = stage,
            )
        }
    }
}

/**
 * Immutable, message-scoped local projection of recipient delivery/read progress.
 *
 * The projection cannot authenticate receipts, persist them, fetch them, or publish them. A future
 * reviewed Data adapter remains responsible for authentication and authorization. This class only
 * enforces exact message/conversation scope and monotonic per-recipient progression once an event
 * reaches the client boundary.
 */
class DataMessageReceiptProjection private constructor(
    val messageId: String,
    val conversationId: String,
    private val recipientStages: Map<String, DataReceiptEvent.Stage>,
) {
    sealed interface ApplyResult {
        data class Applied(val projection: DataMessageReceiptProjection) : ApplyResult

        data class Ignored(
            val projection: DataMessageReceiptProjection,
            val reason: IgnoreReason,
        ) : ApplyResult

        data class Rejected(
            val projection: DataMessageReceiptProjection,
            val reason: RejectReason,
        ) : ApplyResult
    }

    enum class IgnoreReason {
        DUPLICATE,
        STALE,
    }

    enum class RejectReason {
        MESSAGE_SCOPE_MISMATCH,
        CONVERSATION_SCOPE_MISMATCH,
    }

    fun stageFor(recipientId: String): DataReceiptEvent.Stage? =
        recipientStages[recipientId.trim()]

    fun snapshot(): Map<String, DataReceiptEvent.Stage> = recipientStages.toMap()

    fun apply(event: DataReceiptEvent): ApplyResult {
        if (event.messageId != messageId) {
            return ApplyResult.Rejected(this, RejectReason.MESSAGE_SCOPE_MISMATCH)
        }
        if (event.conversationId != conversationId) {
            return ApplyResult.Rejected(this, RejectReason.CONVERSATION_SCOPE_MISMATCH)
        }

        val current = recipientStages[event.recipientId]
        if (current != null) {
            if (event.stage == current) {
                return ApplyResult.Ignored(this, IgnoreReason.DUPLICATE)
            }
            if (event.stage.rank < current.rank) {
                return ApplyResult.Ignored(this, IgnoreReason.STALE)
            }
        }

        return ApplyResult.Applied(
            DataMessageReceiptProjection(
                messageId = messageId,
                conversationId = conversationId,
                recipientStages = recipientStages + (event.recipientId to event.stage),
            ),
        )
    }

    companion object {
        fun empty(
            messageId: String,
            conversationId: String,
        ): DataMessageReceiptProjection {
            require(messageId.isNotBlank()) { "messageId must not be blank" }
            require(conversationId.isNotBlank()) { "conversationId must not be blank" }

            return DataMessageReceiptProjection(
                messageId = messageId.trim(),
                conversationId = conversationId.trim(),
                recipientStages = emptyMap(),
            )
        }
    }
}
