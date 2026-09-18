package com.goreecloud.messenger.client

/**
 * Protocol-neutral explanation for why an E2EE authority projection cannot currently be accepted
 * for one exact conversation. These values contain no key material, credentials, secrets, or
 * protocol identifiers and do not grant authority; they only explain already-supplied evidence.
 */
enum class E2EEAcceptanceFailure {
    CRYPTOGRAPHY_NOT_ACTIVE,
    CONVERSATION_SCOPE_NOT_VERIFIED,
    IMPLEMENTATION_REVIEW_NOT_ACCEPTED,
    DEVICE_IDENTITY_NOT_ENROLLED,
    SESSION_NOT_ESTABLISHED,
    KEY_LIFECYCLE_NOT_CURRENT,
}

/**
 * Bounded user-facing explanation copy for the disconnected Development shell.
 *
 * These strings deliberately describe missing acceptance evidence rather than claiming a concrete
 * cryptographic protocol, key state, threat outcome, or security guarantee.
 */
internal fun E2EEAcceptanceFailure.presentationReason(): String = when (this) {
    E2EEAcceptanceFailure.CRYPTOGRAPHY_NOT_ACTIVE ->
        "No verified active E2EE evidence is available."
    E2EEAcceptanceFailure.CONVERSATION_SCOPE_NOT_VERIFIED ->
        "E2EE evidence is not verified for this exact conversation."
    E2EEAcceptanceFailure.IMPLEMENTATION_REVIEW_NOT_ACCEPTED ->
        "The cryptographic implementation review is not accepted."
    E2EEAcceptanceFailure.DEVICE_IDENTITY_NOT_ENROLLED ->
        "The local cryptographic device identity is not enrolled."
    E2EEAcceptanceFailure.SESSION_NOT_ESTABLISHED ->
        "The conversation-scoped cryptographic session is not established."
    E2EEAcceptanceFailure.KEY_LIFECYCLE_NOT_CURRENT ->
        "The cryptographic key lifecycle is not current."
}

/**
 * Return the first protocol-neutral acceptance failure in the same fail-closed order used by the
 * current E2EE readiness contract. A null result means the supplied evidence satisfies the current
 * acceptance projection for the exact canonical conversation; it does not independently prove
 * production security or create authorization.
 */
fun E2EESessionEvidence.acceptanceFailureFor(
    expectedConversationId: String,
): E2EEAcceptanceFailure? {
    val targetConversationId = DataReceiptIdentifierPolicy.requireCanonical(
        expectedConversationId,
        "expectedConversationId",
    )

    if (state != DataMessagingReadiness.CryptographicState.E2EE_ACTIVE) {
        return E2EEAcceptanceFailure.CRYPTOGRAPHY_NOT_ACTIVE
    }

    val canonicalConversationId = e2eeConversationId
        ?.let(DataReceiptIdentifierPolicy::canonicalOrNull)
    if (canonicalConversationId != targetConversationId) {
        return E2EEAcceptanceFailure.CONVERSATION_SCOPE_NOT_VERIFIED
    }
    if (implementationReview != E2EEImplementationReviewState.ACCEPTED) {
        return E2EEAcceptanceFailure.IMPLEMENTATION_REVIEW_NOT_ACCEPTED
    }
    if (deviceIdentity != E2EEDeviceIdentityState.ENROLLED) {
        return E2EEAcceptanceFailure.DEVICE_IDENTITY_NOT_ENROLLED
    }
    if (sessionEstablishment != E2EESessionEstablishmentState.ESTABLISHED) {
        return E2EEAcceptanceFailure.SESSION_NOT_ESTABLISHED
    }
    if (keyLifecycle != E2EEKeyLifecycleState.CURRENT) {
        return E2EEAcceptanceFailure.KEY_LIFECYCLE_NOT_CURRENT
    }

    return null
}
