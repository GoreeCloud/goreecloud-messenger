package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Test

class E2EEAcceptanceExplanationTest {
    @Test
    fun acceptedEvidenceHasNoAcceptanceFailure() {
        assertNull(acceptedEvidence().acceptanceFailureFor(CONVERSATION_ID))
    }

    @Test
    fun reportsExactProtocolNeutralFailureWithoutChangingReadiness() {
        val cases = listOf(
            acceptedEvidence().copy(
                state = DataMessagingReadiness.CryptographicState.UNKNOWN,
            ) to E2EEAcceptanceFailure.CRYPTOGRAPHY_NOT_ACTIVE,
            acceptedEvidence().copy(
                e2eeConversationId = "other-conversation",
            ) to E2EEAcceptanceFailure.CONVERSATION_SCOPE_NOT_VERIFIED,
            acceptedEvidence().copy(
                implementationReview = E2EEImplementationReviewState.NOT_ACCEPTED,
            ) to E2EEAcceptanceFailure.IMPLEMENTATION_REVIEW_NOT_ACCEPTED,
            acceptedEvidence().copy(
                deviceIdentity = E2EEDeviceIdentityState.NOT_ENROLLED,
            ) to E2EEAcceptanceFailure.DEVICE_IDENTITY_NOT_ENROLLED,
            acceptedEvidence().copy(
                sessionEstablishment = E2EESessionEstablishmentState.NOT_ESTABLISHED,
            ) to E2EEAcceptanceFailure.SESSION_NOT_ESTABLISHED,
            acceptedEvidence().copy(
                keyLifecycle = E2EEKeyLifecycleState.NOT_CURRENT,
            ) to E2EEAcceptanceFailure.KEY_LIFECYCLE_NOT_CURRENT,
        )

        cases.forEach { (evidence, expectedFailure) ->
            assertEquals(expectedFailure, evidence.acceptanceFailureFor(CONVERSATION_ID))
        }
    }

    @Test
    fun presentationReasonsRemainProtocolNeutralAndBounded() {
        val expected = mapOf(
            E2EEAcceptanceFailure.CRYPTOGRAPHY_NOT_ACTIVE to
                "No verified active E2EE evidence is available.",
            E2EEAcceptanceFailure.CONVERSATION_SCOPE_NOT_VERIFIED to
                "E2EE evidence is not verified for this exact conversation.",
            E2EEAcceptanceFailure.IMPLEMENTATION_REVIEW_NOT_ACCEPTED to
                "The cryptographic implementation review is not accepted.",
            E2EEAcceptanceFailure.DEVICE_IDENTITY_NOT_ENROLLED to
                "The local cryptographic device identity is not enrolled.",
            E2EEAcceptanceFailure.SESSION_NOT_ESTABLISHED to
                "The conversation-scoped cryptographic session is not established.",
            E2EEAcceptanceFailure.KEY_LIFECYCLE_NOT_CURRENT to
                "The cryptographic key lifecycle is not current.",
        )

        assertEquals(E2EEAcceptanceFailure.entries.toSet(), expected.keys)
        expected.forEach { (failure, reason) ->
            assertEquals(reason, failure.presentationReason())
        }
    }

    @Test(expected = IllegalArgumentException::class)
    fun noncanonicalRequestedConversationIsRejectedBeforeExplanation() {
        acceptedEvidence().acceptanceFailureFor(" conversation-1")
    }

    private fun acceptedEvidence(): E2EESessionEvidence = E2EESessionEvidence(
        state = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
        e2eeConversationId = CONVERSATION_ID,
        implementationReview = E2EEImplementationReviewState.ACCEPTED,
        deviceIdentity = E2EEDeviceIdentityState.ENROLLED,
        sessionEstablishment = E2EESessionEstablishmentState.ESTABLISHED,
        keyLifecycle = E2EEKeyLifecycleState.CURRENT,
    )

    private companion object {
        const val CONVERSATION_ID = "conversation-1"
    }
}
