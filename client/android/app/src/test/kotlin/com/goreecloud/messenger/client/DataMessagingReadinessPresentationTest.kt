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
    fun readyResultShowsEveryPrerequisiteVerifiedWithProvenance() {
        val result = DataMessagingReadiness.Result.Ready(verifiedConversationId = "conversation-1")

        val presentation = DataMessagingReadinessPresentationPolicy.present(result)

        assertTrue(presentation.checklist.all { it.verified })
        assertTrue(presentation.body().contains("Data · E2EE"))
        assertFalse(presentation.body().contains("Not verified —"))
    }
}
