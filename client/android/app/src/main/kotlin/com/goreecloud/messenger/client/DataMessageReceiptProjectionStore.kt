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
) {
    data class Scope(
        val messageId: String,
        val conversationId: String,
    ) {
        companion object {
            fun create(messageId: String, conversationId: String): Scope {
                val normalizedMessageId = messageId.trim()
                val normalizedConversationId = conversationId.trim()
                require(normalizedMessageId.isNotEmpty()) { "messageId must not be blank" }
                require(normalizedConversationId.isNotEmpty()) { "conversationId must not be blank" }
                return Scope(normalizedMessageId, normalizedConversationId)
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
    }

    private val projections = LinkedHashMap<Scope, DataMessageReceiptProjection>()

    init {
        require(capacity in 1..MAX_CAPACITY) {
            "capacity must be between 1 and $MAX_CAPACITY"
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
        val scope = Scope.create(event.messageId, event.conversationId)
        val current = projections[scope] ?: return ApplyResult.UnknownScope
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
    fun projectionFor(messageId: String, conversationId: String): DataMessageReceiptProjection? =
        projections[Scope.create(messageId, conversationId)]

    @Synchronized
    fun remove(messageId: String, conversationId: String): Boolean =
        projections.remove(Scope.create(messageId, conversationId)) != null

    @Synchronized
    fun registeredCount(): Int = projections.size

    companion object {
        const val DEFAULT_CAPACITY = 256
        const val MAX_CAPACITY = 4096
    }
}
