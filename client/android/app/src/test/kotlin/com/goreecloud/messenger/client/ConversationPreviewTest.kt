package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class ConversationPreviewTest {
    @Test
    fun accessibilitySummaryIncludesProvenanceAndConversationState() {
        val preview = ConversationPreview(
            conversationId = "conversation-1",
            displayName = "Alex",
            secondaryIdentity = "@alex",
            previewText = "Hello",
            timestampLabel = "Now",
            provenance = CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.E2EE_ACTIVE,
            ),
            unreadCount = 3,
            isPinned = true,
            isMuted = true,
        )

        val summary = preview.accessibilitySummary()
        assertTrue(summary.contains("Alex"))
        assertTrue(summary.contains("@alex"))
        assertTrue(summary.contains("Data · E2EE"))
        assertTrue(summary.contains("3 unread"))
        assertTrue(summary.contains("Pinned"))
        assertTrue(summary.contains("Muted"))
    }

    @Test
    fun developmentCatalogIsDeterministicAndTransportTruthful() {
        val previews = DevelopmentConversationCatalog.previews()

        assertEquals(3, previews.size)
        assertEquals("development-alex", previews.first().conversationId)
        assertEquals("Data · E2EE", previews[0].provenance.displayLabel())
        assertEquals("Data", previews[1].provenance.displayLabel())
        assertEquals("SMS", previews[2].provenance.displayLabel())
    }

    @Test(expected = IllegalArgumentException::class)
    fun negativeUnreadCountIsRejected() {
        ConversationPreview(
            conversationId = "conversation-1",
            displayName = "Alex",
            secondaryIdentity = null,
            previewText = "Hello",
            timestampLabel = "Now",
            provenance = CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.UNKNOWN,
            ),
            unreadCount = -1,
            isPinned = false,
            isMuted = false,
        )
    }
}
