package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Test

class DataMessageReceiptProjectionStoreTest {
    @Test
    fun receiptCannotCreateProjectionAuthority() {
        val store = DataMessageReceiptProjectionStore(capacity = 2)
        val event = DataReceiptEvent.create(
            messageId = "message-1",
            conversationId = "conversation-1",
            recipientId = "recipient-1",
            stage = DataReceiptEvent.Stage.DELIVERED,
        )

        assertTrue(store.apply(event) is DataMessageReceiptProjectionStore.ApplyResult.UnknownScope)
        assertEquals(0, store.registeredCount())
    }

    @Test
    fun registeredProjectionAdvancesAndCanBeExplicitlyRemoved() {
        val store = DataMessageReceiptProjectionStore(capacity = 2)
        assertTrue(
            store.register(" message-1 ", " conversation-1 ") is
                DataMessageReceiptProjectionStore.RegisterResult.Registered,
        )

        val result = store.apply(
            DataReceiptEvent.create(
                messageId = "message-1",
                conversationId = "conversation-1",
                recipientId = "recipient-1",
                stage = DataReceiptEvent.Stage.READ,
            ),
        )
        assertTrue(result is DataMessageReceiptProjectionStore.ApplyResult.Applied)
        assertEquals(
            DataReceiptEvent.Stage.READ,
            store.projectionFor("message-1", "conversation-1")?.stageFor("recipient-1"),
        )

        assertTrue(store.remove("message-1", "conversation-1"))
        assertNull(store.projectionFor("message-1", "conversation-1"))
    }

    @Test
    fun capacityFailsClosedWithoutEvictingExistingProjection() {
        val store = DataMessageReceiptProjectionStore(capacity = 1)
        store.register("message-1", "conversation-1")

        assertTrue(
            store.register("message-2", "conversation-1") is
                DataMessageReceiptProjectionStore.RegisterResult.CapacityReached,
        )
        assertEquals(1, store.registeredCount())
        assertTrue(store.projectionFor("message-1", "conversation-1") != null)
        assertNull(store.projectionFor("message-2", "conversation-1"))
    }

    @Test
    fun recipientCapacityFailsClosedWithoutDroppingExistingRecipientState() {
        val store = DataMessageReceiptProjectionStore(capacity = 1, recipientCapacity = 1)
        store.register("message-1", "conversation-1")
        store.apply(
            DataReceiptEvent.create(
                messageId = "message-1",
                conversationId = "conversation-1",
                recipientId = "recipient-1",
                stage = DataReceiptEvent.Stage.DELIVERED,
            ),
        )

        val overflow = store.apply(
            DataReceiptEvent.create(
                messageId = "message-1",
                conversationId = "conversation-1",
                recipientId = "recipient-2",
                stage = DataReceiptEvent.Stage.READ,
            ),
        )
        assertTrue(
            overflow is DataMessageReceiptProjectionStore.ApplyResult.RecipientCapacityReached,
        )
        assertEquals(
            DataReceiptEvent.Stage.DELIVERED,
            store.projectionFor("message-1", "conversation-1")?.stageFor("recipient-1"),
        )
        assertNull(store.projectionFor("message-1", "conversation-1")?.stageFor("recipient-2"))

        val existingRecipientAdvance = store.apply(
            DataReceiptEvent.create(
                messageId = "message-1",
                conversationId = "conversation-1",
                recipientId = "recipient-1",
                stage = DataReceiptEvent.Stage.READ,
            ),
        )
        assertTrue(existingRecipientAdvance is DataMessageReceiptProjectionStore.ApplyResult.Applied)
    }

    @Test
    fun sameMessageIdInAnotherConversationDoesNotShareState() {
        val store = DataMessageReceiptProjectionStore(capacity = 2)
        store.register("message-1", "conversation-1")
        store.register("message-1", "conversation-2")

        store.apply(
            DataReceiptEvent.create(
                messageId = "message-1",
                conversationId = "conversation-1",
                recipientId = "recipient-1",
                stage = DataReceiptEvent.Stage.READ,
            ),
        )

        assertEquals(
            DataReceiptEvent.Stage.READ,
            store.projectionFor("message-1", "conversation-1")?.stageFor("recipient-1"),
        )
        assertNull(store.projectionFor("message-1", "conversation-2")?.stageFor("recipient-1"))
    }
}
