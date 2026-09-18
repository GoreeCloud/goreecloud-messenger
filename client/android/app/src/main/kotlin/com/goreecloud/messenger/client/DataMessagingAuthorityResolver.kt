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
 * Protocol-neutral acceptance state for one future GoreeCloud Data transport prerequisite.
 *
 * These values are minimized projections only. ACCEPTED does not carry endpoints, credentials,
 * tokens, certificates, retry schedules, or other implementation details across this boundary.
 */
enum class DataTransportAcceptanceState {
    ACCEPTED,
    NOT_ACCEPTED,
    UNKNOWN,
}

/**
 * Minimized evidence supplied by a future GoreeCloud Data transport authority.
 *
 * A provider may project [DataMessagingReadiness.DataTransportState.AVAILABLE] into readiness only
 * after its responsible runtime has independently accepted required configuration, authentication
 * binding, protected-channel operation, and bounded timeout/retry/failure behavior. Test fixtures
 * can exercise these states but do not create production transport acceptance.
 */
data class DataTransportEvidence(
    val state: DataMessagingReadiness.DataTransportState,
    val configuration: DataTransportAcceptanceState = DataTransportAcceptanceState.UNKNOWN,
    val authenticationBinding: DataTransportAcceptanceState = DataTransportAcceptanceState.UNKNOWN,
    val channelProtection: DataTransportAcceptanceState = DataTransportAcceptanceState.UNKNOWN,
    val failurePolicy: DataTransportAcceptanceState = DataTransportAcceptanceState.UNKNOWN,
) {
    /**
     * Return only the transport state allowed to participate in messaging readiness.
     *
     * A bare AVAILABLE claim fails closed to UNKNOWN. Negative/unavailable states are preserved and
     * are never upgraded by positive auxiliary facts.
     */
    fun readinessProjection(): DataTransportEvidence {
        if (state != DataMessagingReadiness.DataTransportState.AVAILABLE) {
            return this
        }

        val accepted =
            configuration == DataTransportAcceptanceState.ACCEPTED &&
                authenticationBinding == DataTransportAcceptanceState.ACCEPTED &&
                channelProtection == DataTransportAcceptanceState.ACCEPTED &&
                failurePolicy == DataTransportAcceptanceState.ACCEPTED

        return if (accepted) this else copy(state = DataMessagingReadiness.DataTransportState.UNKNOWN)
    }
}

/**
 * Narrow acceptance projection for a future GoreeCloud Data transport implementation.
 * Availability alone is insufficient and does not imply authentication, conversation access,
 * E2EE readiness, endpoint approval, or production acceptance.
 */
fun interface GoreeCloudDataTransportAuthority {
    fun evidence(): DataTransportEvidence
}

/**
 * Whether the responsible cryptographic authority has accepted the implementation under the
 * applicable security review. Test fixtures may exercise this state, but they do not create a real
 * security review or production acceptance.
 */
enum class E2EEImplementationReviewState {
    ACCEPTED,
    NOT_ACCEPTED,
    UNKNOWN,
}

/**
 * Whether the cryptographic authority has an enrolled cryptographic identity for the local device.
 * This is not a user-facing manual verification claim and carries no device key material.
 */
enum class E2EEDeviceIdentityState {
    ENROLLED,
    NOT_ENROLLED,
    UNKNOWN,
}

/** Conversation-scoped cryptographic session establishment state. */
enum class E2EESessionEstablishmentState {
    ESTABLISHED,
    NOT_ESTABLISHED,
    UNKNOWN,
}

/**
 * Whether the responsible authority accepts the current session/device key lifecycle state.
 * CURRENT includes any rotation or key-change processing that the eventual reviewed implementation
 * requires before it may call the conversation actively protected.
 */
enum class E2EEKeyLifecycleState {
    CURRENT,
    NOT_CURRENT,
    UNKNOWN,
}

/**
 * Conversation-scoped cryptographic evidence supplied by a future reviewed E2EE authority.
 *
 * A provider may project [DataMessagingReadiness.CryptographicState.E2EE_ACTIVE] only together
 * with explicit acceptance evidence that the implementation review, local cryptographic device
 * identity, session establishment, and key lifecycle are all current for the exact conversation.
 * No key material, session secret, algorithm identifier, ciphertext, or device secret crosses this
 * boundary. Group and multi-device implementations remain responsible for resolving their complete
 * participant/device state before returning these minimized acceptance projections.
 */
data class E2EESessionEvidence(
    val state: DataMessagingReadiness.CryptographicState,
    val e2eeConversationId: String? = null,
    val implementationReview: E2EEImplementationReviewState = E2EEImplementationReviewState.UNKNOWN,
    val deviceIdentity: E2EEDeviceIdentityState = E2EEDeviceIdentityState.UNKNOWN,
    val sessionEstablishment: E2EESessionEstablishmentState = E2EESessionEstablishmentState.UNKNOWN,
    val keyLifecycle: E2EEKeyLifecycleState = E2EEKeyLifecycleState.UNKNOWN,
) {
    /**
     * Return the minimized readiness projection allowed to leave the cryptographic authority seam.
     *
     * A contradictory bare E2EE_ACTIVE claim fails closed to UNKNOWN unless all protocol-neutral
     * acceptance facts are positive and the exact canonical conversation scope matches the request.
     * Negative/non-established states remain negative; this function never upgrades them.
     */
    fun readinessProjectionFor(expectedConversationId: String): E2EESessionEvidence {
        if (state != DataMessagingReadiness.CryptographicState.E2EE_ACTIVE) {
            return this
        }

        val canonicalConversationId = e2eeConversationId
            ?.let(DataReceiptIdentifierPolicy::canonicalOrNull)
        val accepted =
            canonicalConversationId == expectedConversationId &&
                implementationReview == E2EEImplementationReviewState.ACCEPTED &&
                deviceIdentity == E2EEDeviceIdentityState.ENROLLED &&
                sessionEstablishment == E2EESessionEstablishmentState.ESTABLISHED &&
                keyLifecycle == E2EEKeyLifecycleState.CURRENT

        return if (accepted) {
            copy(e2eeConversationId = canonicalConversationId)
        } else {
            copy(
                state = DataMessagingReadiness.CryptographicState.UNKNOWN,
                e2eeConversationId = canonicalConversationId,
            )
        }
    }
}

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
            dataTransportAuthority.evidence().readinessProjection().state
        } catch (_: Exception) {
            DataMessagingReadiness.DataTransportState.UNKNOWN
        }

        val cryptography = try {
            e2eeSessionAuthority
                .stateFor(targetConversationId)
                .readinessProjectionFor(targetConversationId)
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
