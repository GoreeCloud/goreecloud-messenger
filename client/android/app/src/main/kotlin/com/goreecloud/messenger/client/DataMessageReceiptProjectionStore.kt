package com.goreecloud.messenger.client

/**
 * Bounded client-local owner of receipt projections for messages the client has explicitly
 * registered. Receipt events cannot create projection authority on their own.
 *
 * This store is intentionally process-local Development state. It is not transport,
 * authentication, persistence, delivery evidence, or server authority.
 */
class DataMessageReceiptProjectionStore(
    private val capacity: Int = DEFAULT_CAPACITY,
    private val recipientCapacity: Int = DEFAULT_RECIPIENT_CAPACITY,
) {
    data class Scope(
        val messageId: String,
        val conversationId: String,
    ) {
        companion object {
            fun create(messageId: String, conversationId: String): Scope = Scope(
                messageId = DataReceiptIdentifierPolicy.requireCanonical(messageId, "messageId"),
                conversationId = DataReceiptIdentifierPolicy.requireCanonical(
                    conversationId,
                    "conversationId",
                ),
            )

            fun lookup(messageId: String, conversationId: String): Scope? {
                val canonicalMessageId = DataReceiptIdentifierPolicy.canonicalOrNull(messageId)
                    ?: return null
                val canonicalConversationId = DataReceiptIdentifierPolicy.canonicalOrNull(conversationId)
                    ?: return null
                return Scope(canonicalMessageId, canonicalConversationId)
            }
        }
    }

    sealed interface RegisterResult {
        data class Registered(val projection: DataMessageReceiptProjection) : RegisterResult
        data class Existing(val projection: DataMessageReceiptProjection) : RegisterResult
        data object CapacityReached : RegisterResult
    }

    sealed interface ApplyResult {
        data class Applied(val projection: DataMessageReceiptProjection) : ApplyResult
        data class Ignored(
            val reason: DataMessageReceiptProjection.IgnoreReason,
            val projection: DataMessageReceiptProjection,
        ) : ApplyResult
        data class Rejected(val reason: DataMessageReceiptProjection.RejectReason) : ApplyResult
        data object UnknownScope : ApplyResult
        data object RecipientCapacityReached : ApplyResult
    }

    private val projections = LinkedHashMap<Scope, DataMessageReceiptProjection>()

    init {
        require(capacity in 1..MAX_CAPACITY) {
            "capacity must be between 1 and $MAX_CAPACITY"
        }
        require(recipientCapacity in 1..MAX_RECIPIENT_CAPACITY) {
            "recipientCapacity must be between 1 and $MAX_RECIPIENT_CAPACITY"
        }
    }

    @Synchronized
    fun register(messageId: String, conversationId: String): RegisterResult {
        val scope = Scope.create(messageId, conversationId)
        projections[scope]?.let { return RegisterResult.Existing(it) }
        if (projections.size >= capacity) return RegisterResult.CapacityReached

        val projection = DataMessageReceiptProjection.empty(scope.messageId, scope.conversationId)
        projections[scope] = projection
        return RegisterResult.Registered(projection)
    }

    @Synchronized
    fun apply(event: DataReceiptEvent): ApplyResult {
        // DataReceiptEvent construction already proves these identifiers are canonical. Avoid any
        // normalization at the store boundary so altered spellings cannot collapse onto a scope.
        val scope = Scope(event.messageId, event.conversationId)
        val current = projections[scope] ?: return ApplyResult.UnknownScope
        val isNewRecipient = current.stageFor(event.recipientId) == null
        if (isNewRecipient && current.snapshot().size >= recipientCapacity) {
            return ApplyResult.RecipientCapacityReached
        }

        return when (val result = current.apply(event)) {
            is DataMessageReceiptProjection.ApplyResult.Applied -> {
                projections[scope] = result.projection
                ApplyResult.Applied(result.projection)
            }
            is DataMessageReceiptProjection.ApplyResult.Ignored -> ApplyResult.Ignored(
                reason = result.reason,
                projection = result.projection,
            )
            is DataMessageReceiptProjection.ApplyResult.Rejected -> ApplyResult.Rejected(result.reason)
        }
    }

    @Synchronized
    fun projectionFor(messageId: String, conversationId: String): DataMessageReceiptProjection? {
        val scope = Scope.lookup(messageId, conversationId) ?: return null
        return projections[scope]
    }

    @Synchronized
    fun remove(messageId: String, conversationId: String): Boolean {
        val scope = Scope.lookup(messageId, conversationId) ?: return false
        return projections.remove(scope) != null
    }

    @Synchronized
    fun registeredCount(): Int = projections.size

    companion object {
        const val DEFAULT_CAPACITY = 256
        const val MAX_CAPACITY = 4096
        const val DEFAULT_RECIPIENT_CAPACITY = 256
        const val MAX_RECIPIENT_CAPACITY = 4096
    }
}
