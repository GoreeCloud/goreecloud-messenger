package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class ConversationMessagePreviewTest {
    @Test
    fun developmentConversationMessagesRemainUnverifiedAndDeterministic() {
        val messages = DevelopmentConversationMessageCatalog.messages()
        assertEquals(2, messages.size)
        assertTrue(messages.all { it.provenance.protection != CommunicationProtection.E2EE_ACTIVE })
        assertEquals("Data · Protection not verified · Development preview", messages[0].statusLabel())
        assertEquals("Data · Protection not verified · Delivered · Development preview", messages[1].statusLabel())
    }

    @Test(expected = IllegalArgumentException::class)
    fun incomingMessageCannotClaimSenderReceipt() {
        ConversationMessagePreview(
            messageId = "message",
            senderLabel = "Alex",
            body = "Hello",
            timestampLabel = "Now",
            direction = MessageDirection.INCOMING,
            provenance = CommunicationProvenance(CommunicationTransport.DATA, CommunicationProtection.UNKNOWN),
            receiptState = MessageReceiptState.READ,
        )
    }
}
