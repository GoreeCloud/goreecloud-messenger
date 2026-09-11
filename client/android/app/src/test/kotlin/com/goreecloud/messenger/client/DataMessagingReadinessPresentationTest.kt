package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class DataMessagingReadinessPresentationTest {
    @Test
    fun allMissingPrerequisitesRemainExplicitAndOrdered() {
        val result = DataMessagingReadiness.Result.Blocked(
            DataMessagingReadiness.BlockReason.entries.toSet(),
        )

        val presentation = DataMessagingReadinessPresentationPolicy.present(result)

        assertEquals(
            listOf(
                "Identity authentication",
                "Conversation authorization",
                "GoreeCloud Data transport",
                "Verified active E2EE",
            ),
            presentation.checklist.map { it.label },
        )
        assertTrue(presentation.checklist.none { it.verified })
        assertTrue(presentation.body().contains("No Send action is exposed."))
    }

    @Test
    fun partialEvidenceDoesNotTurnWholeResultReady() {
        val result = DataMessagingReadiness.Result.Blocked(
            setOf(
                DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE,
                DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE,
            ),
        )

        val presentation = DataMessagingReadinessPresentationPolicy.present(result)

        assertEquals(listOf(true, true, false, false), presentation.checklist.map { it.verified })
        assertTrue(presentation.body().contains("Verified — Identity authentication"))
        assertTrue(presentation.body().contains("Not verified — Verified active E2EE"))
        assertFalse(presentation.body().contains("prerequisites are independently verified ("))
    }

    @Test
    fun blockedIndependentAuthoritiesCanExplainUnknownEvidence() {
        val result = DataMessagingReadiness.Result.Blocked(
            DataMessagingReadiness.BlockReason.entries.toSet(),
        )
        val evidence = DataMessagingReadiness.Evidence(
            identity = DataMessagingReadiness.IdentityState.UNKNOWN,
            conversationAccess = DataMessagingReadiness.ConversationAccessState.UNKNOWN,
            transport = DataMessagingReadiness.DataTransportState.UNKNOWN,
            cryptography = DataMessagingReadiness.CryptographicState.UNKNOWN,
        )

        val presentation = DataMessagingReadinessPresentationPolicy.present(
            result = result,
            evidence = evidence,
            e2eeFailure = E2EEAcceptanceFailure.CRYPTOGRAPHY_NOT_ACTIVE,
        )

        assertEquals(
            listOf(
                "No verified GoreeCloud Identity session evidence is available.",
                "No verified conversation-participant authorization evidence is available.",
                "No verified GoreeCloud Data transport availability evidence is available.",
                "No verified active E2EE evidence is available.",
            ),
            presentation.checklist.map { it.detail },
        )
        assertEquals(4, presentation.body().split("Why — ").size - 1)
    }

    @Test
    fun explicitNegativeAuthorityStatesUseBoundedReasons() {
        val result = DataMessagingReadiness.Result.Blocked(
            setOf(
                DataMessagingReadiness.BlockReason.IDENTITY_NOT_AUTHENTICATED,
                DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED,
                DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE,
            ),
        )
        val evidence = DataMessagingReadiness.Evidence(
            identity = DataMessagingReadiness.IdentityState.UNAUTHENTICATED,
            conversationAccess = DataMessagingReadiness.ConversationAccessState.NOT_PARTICIPANT,
            transport = DataMessagingReadiness.DataTransportState.UNAVAILABLE,
            cryptography = DataMessagingReadiness.CryptographicState.UNKNOWN,
        )

        val presentation = DataMessagingReadinessPresentationPolicy.present(
            result = result,
            evidence = evidence,
        )

        assertEquals(
            "GoreeCloud Identity reports that this client is not authenticated.",
            presentation.checklist[0].detail,
        )
        assertEquals(
            "This identity is not verified as a participant in the exact conversation.",
            presentation.checklist[1].detail,
        )
        assertEquals(
            "GoreeCloud Data transport is reported unavailable.",
            presentation.checklist[2].detail,
        )
        assertEquals(null, presentation.checklist[3].detail)
    }

    @Test
    fun missingConversationScopeGetsExactScopeExplanation() {
        val result = DataMessagingReadiness.Result.Blocked(
            setOf(DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED),
        )
        val evidence = DataMessagingReadiness.Evidence(
            identity = DataMessagingReadiness.IdentityState.AUTHENTICATED,
            conversationAccess = DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
            transport = DataMessagingReadiness.DataTransportState.AVAILABLE,
            cryptography = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
            authorizedConversationId = null,
            e2eeConversationId = "conversation-1",
        )

        val presentation = DataMessagingReadinessPresentationPolicy.present(
            result = result,
            evidence = evidence,
        )

        assertEquals(
            "Conversation authorization is not verified for an exact canonical conversation.",
            presentation.checklist[1].detail,
        )
    }

    @Test
    fun contradictoryAuxiliaryFailureDetailsAreSuppressedForVerifiedItems() {
        val result = DataMessagingReadiness.Result.Blocked(
            setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
        )
        val evidence = DataMessagingReadiness.Evidence(
            identity = DataMessagingReadiness.IdentityState.UNAUTHENTICATED,
            conversationAccess = DataMessagingReadiness.ConversationAccessState.NOT_PARTICIPANT,
            transport = DataMessagingReadiness.DataTransportState.UNAVAILABLE,
            cryptography = DataMessagingReadiness.CryptographicState.UNKNOWN,
        )

        val presentation = DataMessagingReadinessPresentationPolicy.present(
            result = result,
            evidence = evidence,
            e2eeFailure = E2EEAcceptanceFailure.SESSION_NOT_ESTABLISHED,
        )

        assertTrue(presentation.checklist.take(3).all { it.verified })
        assertTrue(presentation.checklist.take(3).all { it.detail == null })
        assertEquals(
            "The conversation-scoped cryptographic session is not established.",
            presentation.checklist.last().detail,
        )
        assertEquals(1, presentation.body().split("Why — ").size - 1)
    }

    @Test
    fun blockedE2eeItemCanExplainWhyWithoutChangingReadiness() {
        val result = DataMessagingReadiness.Result.Blocked(
            setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
        )

        val presentation = DataMessagingReadinessPresentationPolicy.present(
            result = result,
            e2eeFailure = E2EEAcceptanceFailure.SESSION_NOT_ESTABLISHED,
        )

        val e2ee = presentation.checklist.last()
        assertFalse(e2ee.verified)
        assertEquals(
            "The conversation-scoped cryptographic session is not established.",
            e2ee.detail,
        )
        assertTrue(
            presentation.body().contains(
                "Why — The conversation-scoped cryptographic session is not established.",
            ),
        )
        assertTrue(presentation.body().contains("Send remains unavailable"))
    }

    @Test
    fun e2eeFailureDetailIsSuppressedWhenE2eeIsVerified() {
        val result = DataMessagingReadiness.Result.Blocked(
            setOf(DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE),
        )

        val presentation = DataMessagingReadinessPresentationPolicy.present(
            result = result,
            e2eeFailure = E2EEAcceptanceFailure.KEY_LIFECYCLE_NOT_CURRENT,
        )

        val e2ee = presentation.checklist.last()
        assertTrue(e2ee.verified)
        assertEquals(null, e2ee.detail)
        assertFalse(presentation.body().contains("Why —"))
    }

    @Test
    fun readyResultShowsEveryPrerequisiteVerifiedWithProvenance() {
        val result = DataMessagingReadiness.Result.Ready(verifiedConversationId = "conversation-1")

        val presentation = DataMessagingReadinessPresentationPolicy.present(result)

        assertTrue(presentation.checklist.all { it.verified })
        assertTrue(presentation.body().contains("Data · E2EE"))
        assertFalse(presentation.body().contains("Not verified —"))
    }
}
