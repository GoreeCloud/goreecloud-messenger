package com.goreecloud.messenger.client

internal object DataReceiptIdentifierPolicy {
    const val MAX_IDENTIFIER_LENGTH = 512

    fun requireCanonical(value: String, label: String): String {
        require(value.isNotEmpty()) { "$label must not be blank" }
        require(value == value.trim()) { "$label must already be canonical" }
        require(value.isNotBlank()) { "$label must not be blank" }
        require(value.length <= MAX_IDENTIFIER_LENGTH) {
            "$label must not exceed $MAX_IDENTIFIER_LENGTH UTF-16 code units"
        }
        require(value.none { character ->
            character.code in 0x00..0x1f || character.code == 0x7f
        }) { "$label must not contain control characters" }
        return value
    }

    fun canonicalOrNull(value: String): String? {
        if (value.isEmpty() || value != value.trim() || value.isBlank()) return null
        if (value.length > MAX_IDENTIFIER_LENGTH) return null
        if (value.any { character -> character.code in 0x00..0x1f || character.code == 0x7f }) {
            return null
        }
        return value
    }
}

/**
 * A delivery/read receipt event that has already crossed the future transport/authentication
 * boundary.
 *
 * This is deliberately not a data class: a generated copy() method would allow a caller to replace
 * one of the factory-validated identifiers without re-running the canonical identifier boundary.
 *
 * This value does not authenticate a recipient and must not be constructed from untrusted network
 * input directly. Its only responsibility is carrying the already-derived exact bounded scope
 * needed by the local projection so state for one message/conversation cannot be reused for
 * another. Identifiers remain opaque: internal spaces and punctuation are preserved exactly.
 */
class DataReceiptEvent private constructor(
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
        ): DataReceiptEvent = DataReceiptEvent(
            messageId = DataReceiptIdentifierPolicy.requireCanonical(messageId, "messageId"),
            conversationId = DataReceiptIdentifierPolicy.requireCanonical(
                conversationId,
                "conversationId",
            ),
            recipientId = DataReceiptIdentifierPolicy.requireCanonical(recipientId, "recipientId"),
            stage = stage,
        )
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

    fun stageFor(recipientId: String): DataReceiptEvent.Stage? {
        val canonicalRecipientId = DataReceiptIdentifierPolicy.canonicalOrNull(recipientId)
            ?: return null
        return recipientStages[canonicalRecipientId]
    }

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
        ): DataMessageReceiptProjection = DataMessageReceiptProjection(
            messageId = DataReceiptIdentifierPolicy.requireCanonical(messageId, "messageId"),
            conversationId = DataReceiptIdentifierPolicy.requireCanonical(
                conversationId,
                "conversationId",
            ),
            recipientStages = emptyMap(),
        )
    }
}
