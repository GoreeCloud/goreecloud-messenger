package com.goreecloud.messenger.client

enum class MessageDirection { INCOMING, OUTGOING }
enum class MessageReceiptState { NONE, SENT, DELIVERED, READ }

data class ConversationMessagePreview(
    val messageId: String,
    val senderLabel: String,
    val body: String,
    val timestampLabel: String,
    val direction: MessageDirection,
    val provenance: CommunicationProvenance,
    val receiptState: MessageReceiptState,
    val replyToLabel: String? = null,
) {
    init {
        require(messageId.isNotBlank()) { "messageId must not be blank" }
        require(senderLabel.isNotBlank()) { "senderLabel must not be blank" }
        require(body.isNotBlank()) { "body must not be blank" }
        require(timestampLabel.isNotBlank()) { "timestampLabel must not be blank" }
        require(replyToLabel?.isNotBlank() != false) { "replyToLabel must be null or nonblank" }
        require(direction == MessageDirection.OUTGOING || receiptState == MessageReceiptState.NONE) {
            "incoming Development previews must not claim sender-side receipt state"
        }
    }

    fun statusLabel(): String = buildList {
        add(provenance.displayLabel())
        when (receiptState) {
            MessageReceiptState.NONE -> Unit
            MessageReceiptState.SENT -> add("Sent")
            MessageReceiptState.DELIVERED -> add("Delivered")
            MessageReceiptState.READ -> add("Read")
        }
        add("Development preview")
    }.joinToString(" · ")
}

object DevelopmentConversationMessageCatalog {
    fun messages(): List<ConversationMessagePreview> = listOf(
        ConversationMessagePreview(
            messageId = "dev-message-1",
            senderLabel = "Alex",
            body = "Are we still on for six?",
            timestampLabel = "5:58 PM",
            direction = MessageDirection.INCOMING,
            provenance = CommunicationProvenance(CommunicationTransport.DATA, CommunicationProtection.UNKNOWN),
            receiptState = MessageReceiptState.NONE,
        ),
        ConversationMessagePreview(
            messageId = "dev-message-2",
            senderLabel = "You",
            body = "Yes — I’ll be there around six.",
            timestampLabel = "6:02 PM",
            direction = MessageDirection.OUTGOING,
            provenance = CommunicationProvenance(CommunicationTransport.DATA, CommunicationProtection.UNKNOWN),
            receiptState = MessageReceiptState.DELIVERED,
            replyToLabel = "Replying to Alex",
        ),
    )
}
