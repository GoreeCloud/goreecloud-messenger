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
