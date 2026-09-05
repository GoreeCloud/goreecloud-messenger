package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Test

class DataMessageReceiptProjectionTest {
    @Test
    fun deliveryThenReadAdvancesMonotonicallyForRecipient() {
        val initial = DataMessageReceiptProjection.empty("message-1", "conversation-1")
        val delivered = DataReceiptEvent.create(
            messageId = "message-1",
            conversationId = "conversation-1",
            recipientId = "recipient-1",
            stage = DataReceiptEvent.Stage.DELIVERED,
        )
        val read = DataReceiptEvent.create(
            messageId = "message-1",
            conversationId = "conversation-1",
            recipientId = "recipient-1",
            stage = DataReceiptEvent.Stage.READ,
        )

        val afterDelivery = (initial.apply(delivered) as DataMessageReceiptProjection.ApplyResult.Applied).projection
        val afterRead = (afterDelivery.apply(read) as DataMessageReceiptProjection.ApplyResult.Applied).projection

        assertEquals(DataReceiptEvent.Stage.DELIVERED, afterDelivery.stageFor("recipient-1"))
        assertEquals(DataReceiptEvent.Stage.READ, afterRead.stageFor("recipient-1"))
        assertNull(initial.stageFor("recipient-1"))
    }

    @Test
    fun duplicateAndStaleReceiptCannotRegressProjection() {
        val initial = DataMessageReceiptProjection.empty("message-1", "conversation-1")
        val read = DataReceiptEvent.create(
            messageId = "message-1",
            conversationId = "conversation-1",
            recipientId = "recipient-1",
            stage = DataReceiptEvent.Stage.READ,
        )
        val afterRead = (initial.apply(read) as DataMessageReceiptProjection.ApplyResult.Applied).projection

        val duplicate = afterRead.apply(read)
        assertTrue(duplicate is DataMessageReceiptProjection.ApplyResult.Ignored)
        assertEquals(
            DataMessageReceiptProjection.IgnoreReason.DUPLICATE,
            (duplicate as DataMessageReceiptProjection.ApplyResult.Ignored).reason,
        )

        val stale = afterRead.apply(
            DataReceiptEvent.create(
                messageId = "message-1",
                conversationId = "conversation-1",
                recipientId = "recipient-1",
                stage = DataReceiptEvent.Stage.DELIVERED,
            ),
        )
        assertTrue(stale is DataMessageReceiptProjection.ApplyResult.Ignored)
        assertEquals(
            DataMessageReceiptProjection.IgnoreReason.STALE,
            (stale as DataMessageReceiptProjection.ApplyResult.Ignored).reason,
        )
        assertEquals(DataReceiptEvent.Stage.READ, stale.projection.stageFor("recipient-1"))
    }

    @Test
    fun crossMessageAndCrossConversationReceiptsAreRejected() {
        val projection = DataMessageReceiptProjection.empty("message-1", "conversation-1")

        val wrongMessage = projection.apply(
            DataReceiptEvent.create(
                messageId = "message-2",
                conversationId = "conversation-1",
                recipientId = "recipient-1",
                stage = DataReceiptEvent.Stage.DELIVERED,
            ),
        )
        assertTrue(wrongMessage is DataMessageReceiptProjection.ApplyResult.Rejected)
        assertEquals(
            DataMessageReceiptProjection.RejectReason.MESSAGE_SCOPE_MISMATCH,
            (wrongMessage as DataMessageReceiptProjection.ApplyResult.Rejected).reason,
        )

        val wrongConversation = projection.apply(
            DataReceiptEvent.create(
                messageId = "message-1",
                conversationId = "conversation-2",
                recipientId = "recipient-1",
                stage = DataReceiptEvent.Stage.DELIVERED,
            ),
        )
        assertTrue(wrongConversation is DataMessageReceiptProjection.ApplyResult.Rejected)
        assertEquals(
            DataMessageReceiptProjection.RejectReason.CONVERSATION_SCOPE_MISMATCH,
            (wrongConversation as DataMessageReceiptProjection.ApplyResult.Rejected).reason,
        )
        assertTrue(projection.snapshot().isEmpty())
    }

    @Test
    fun recipientsAdvanceIndependentlyAndSnapshotsDoNotExposeMutableAuthority() {
        var projection = DataMessageReceiptProjection.empty(" message-1 ", " conversation-1 ")
        val first = DataReceiptEvent.create(
            messageId = "message-1",
            conversationId = "conversation-1",
            recipientId = " recipient-1 ",
            stage = DataReceiptEvent.Stage.READ,
        )
        val second = DataReceiptEvent.create(
            messageId = "message-1",
            conversationId = "conversation-1",
            recipientId = "recipient-2",
            stage = DataReceiptEvent.Stage.DELIVERED,
        )

        projection = (projection.apply(first) as DataMessageReceiptProjection.ApplyResult.Applied).projection
        projection = (projection.apply(second) as DataMessageReceiptProjection.ApplyResult.Applied).projection

        assertEquals(DataReceiptEvent.Stage.READ, projection.stageFor(" recipient-1 "))
        assertEquals(DataReceiptEvent.Stage.DELIVERED, projection.stageFor("recipient-2"))
        assertEquals(2, projection.snapshot().size)
    }
}
