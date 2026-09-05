package com.goreecloud.messenger.client

/**
 * Client-side policy boundary for a future GoreeCloud Data message-send path.
 *
 * This type does not authenticate a user, authorize a conversation, establish network transport,
 * derive cryptographic state, hold keys, encrypt content, persist messages, or send anything. It
 * only combines independently verified states supplied by their responsible runtime authorities.
 *
 * Conversation authorization is deliberately scoped to one exact conversation identifier. A bare
 * VERIFIED_PARTICIPANT enum without that scope is insufficient readiness evidence and fails closed.
 *
 * The Development Android client currently supplies none of those production authorities, so this
 * contract must not be interpreted as a working message composer or E2EE implementation.
 */
object DataMessagingReadiness {
    enum class IdentityState {
        AUTHENTICATED,
        UNAUTHENTICATED,
        UNKNOWN,
    }

    enum class ConversationAccessState {
        VERIFIED_PARTICIPANT,
        NOT_PARTICIPANT,
        UNKNOWN,
    }

    enum class DataTransportState {
        AVAILABLE,
        UNAVAILABLE,
        UNKNOWN,
    }

    enum class CryptographicState {
        E2EE_ACTIVE,
        E2EE_UNAVAILABLE,
        NOT_ESTABLISHED,
        UNKNOWN,
    }

    data class Evidence(
        val identity: IdentityState,
        val conversationAccess: ConversationAccessState,
        val transport: DataTransportState,
        val cryptography: CryptographicState,
        val verifiedConversationId: String? = null,
    )

    enum class BlockReason {
        IDENTITY_NOT_AUTHENTICATED,
        CONVERSATION_ACCESS_NOT_VERIFIED,
        DATA_TRANSPORT_NOT_AVAILABLE,
        E2EE_NOT_VERIFIED_ACTIVE,
    }

    sealed interface Result {
        data class Ready(
            val verifiedConversationId: String,
            val provenance: CommunicationProvenance = CommunicationProvenance(
                transport = CommunicationTransport.DATA,
                protection = CommunicationProtection.E2EE_ACTIVE,
            ),
        ) : Result

        data class Blocked(val reasons: Set<BlockReason>) : Result
    }

    /**
     * Require every independently owned prerequisite before a future client may expose an eligible
     * encrypted GoreeCloud Data send operation.
     *
     * All missing, negative, unavailable, unknown, or unscoped states fail closed. There is
     * deliberately no downgrade to carrier messaging and no conversion from transport availability
     * to E2EE state.
     */
    fun evaluate(evidence: Evidence): Result {
        val verifiedConversationId = evidence.verifiedConversationId
            ?.trim()
            ?.takeIf { it.isNotEmpty() }

        val reasons = buildSet {
            if (evidence.identity != IdentityState.AUTHENTICATED) {
                add(BlockReason.IDENTITY_NOT_AUTHENTICATED)
            }
            if (
                evidence.conversationAccess != ConversationAccessState.VERIFIED_PARTICIPANT ||
                verifiedConversationId == null
            ) {
                add(BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED)
            }
            if (evidence.transport != DataTransportState.AVAILABLE) {
                add(BlockReason.DATA_TRANSPORT_NOT_AVAILABLE)
            }
            if (evidence.cryptography != CryptographicState.E2EE_ACTIVE) {
                add(BlockReason.E2EE_NOT_VERIFIED_ACTIVE)
            }
        }

        return if (reasons.isEmpty()) {
            Result.Ready(verifiedConversationId = checkNotNull(verifiedConversationId))
        } else {
            Result.Blocked(reasons)
        }
    }
}
