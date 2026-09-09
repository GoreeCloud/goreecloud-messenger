package com.goreecloud.messenger.client

/**
 * Narrow client projection of the future GoreeCloud Identity session authority.
 *
 * Implementations own the evidence used to determine whether Messenger currently has an
 * authenticated Identity session. This interface does not carry credentials, tokens, user secrets,
 * or permission to synthesize authentication from local UI state.
 */
fun interface GoreeCloudIdentitySessionAuthority {
    fun authenticationState(): DataMessagingReadiness.IdentityState
}

/**
 * Conversation-scoped authorization evidence supplied by the future authorization authority.
 * A positive state without its exact authority-owned conversation scope remains insufficient.
 */
data class ConversationAuthorizationEvidence(
    val state: DataMessagingReadiness.ConversationAccessState,
    val authorizedConversationId: String? = null,
)

fun interface ConversationAuthorizationAuthority {
    fun accessFor(conversationId: String): ConversationAuthorizationEvidence
}

/**
 * Narrow availability projection for a future GoreeCloud Data transport implementation.
 * Availability alone does not imply authentication, conversation access, or E2EE readiness.
 */
fun interface GoreeCloudDataTransportAuthority {
    fun availability(): DataMessagingReadiness.DataTransportState
}

/**
 * Conversation-scoped cryptographic evidence supplied by a future reviewed E2EE authority.
 * This projection contains no key material, session secret, algorithm claim, or ciphertext.
 */
data class E2EESessionEvidence(
    val state: DataMessagingReadiness.CryptographicState,
    val e2eeConversationId: String? = null,
)

fun interface E2EESessionAuthority {
    fun stateFor(conversationId: String): E2EESessionEvidence
}

/**
 * Resolves independently owned messaging authority projections for one exact conversation.
 *
 * The resolver does not create authority. Each provider remains responsible for its own evidence.
 * Provider failures fail closed to UNKNOWN only for that provider; one successful authority cannot
 * upgrade another missing or failed authority. Exact scope agreement is still enforced by
 * [DataMessagingReadiness].
 */
class DataMessagingAuthorityResolver(
    private val identityAuthority: GoreeCloudIdentitySessionAuthority,
    private val conversationAuthorizationAuthority: ConversationAuthorizationAuthority,
    private val dataTransportAuthority: GoreeCloudDataTransportAuthority,
    private val e2eeSessionAuthority: E2EESessionAuthority,
) {
    fun evidenceFor(conversationId: String): DataMessagingReadiness.Evidence {
        val targetConversationId = DataReceiptIdentifierPolicy.requireCanonical(
            conversationId,
            "conversationId",
        )

        val identity = try {
            identityAuthority.authenticationState()
        } catch (_: Exception) {
            DataMessagingReadiness.IdentityState.UNKNOWN
        }

        val authorization = try {
            conversationAuthorizationAuthority.accessFor(targetConversationId)
        } catch (_: Exception) {
            ConversationAuthorizationEvidence(
                state = DataMessagingReadiness.ConversationAccessState.UNKNOWN,
            )
        }

        val transport = try {
            dataTransportAuthority.availability()
        } catch (_: Exception) {
            DataMessagingReadiness.DataTransportState.UNKNOWN
        }

        val cryptography = try {
            e2eeSessionAuthority.stateFor(targetConversationId)
        } catch (_: Exception) {
            E2EESessionEvidence(
                state = DataMessagingReadiness.CryptographicState.UNKNOWN,
            )
        }

        return DataMessagingReadiness.Evidence(
            identity = identity,
            conversationAccess = authorization.state,
            transport = transport,
            cryptography = cryptography.state,
            authorizedConversationId = authorization.authorizedConversationId,
            e2eeConversationId = cryptography.e2eeConversationId,
        )
    }
}
