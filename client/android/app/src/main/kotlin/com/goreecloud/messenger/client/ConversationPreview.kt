package com.goreecloud.messenger.client

data class ConversationPreview(
    val conversationId: String,
    val displayName: String,
    val secondaryIdentity: String?,
    val previewText: String,
    val timestampLabel: String,
    val provenance: CommunicationProvenance,
    val unreadCount: Int,
    val isPinned: Boolean,
    val isMuted: Boolean,
) {
    init {
        require(conversationId.isNotBlank()) { "conversationId must not be blank" }
        require(displayName.isNotBlank()) { "displayName must not be blank" }
        require(previewText.isNotBlank()) { "previewText must not be blank" }
        require(timestampLabel.isNotBlank()) { "timestampLabel must not be blank" }
        require(unreadCount >= 0) { "unreadCount must not be negative" }
        require(secondaryIdentity?.isNotBlank() != false) { "secondaryIdentity must be null or nonblank" }
    }

    fun accessibilitySummary(): String = buildString {
        append(displayName)
        secondaryIdentity?.let { append(", ").append(it) }
        append(". ").append(provenance.displayLabel())
        append(". ").append(previewText)
        append(". ").append(timestampLabel)
        if (unreadCount > 0) append(". ").append(unreadCount).append(" unread")
        if (isPinned) append(". Pinned")
        if (isMuted) append(". Muted")
    }
}

/**
 * Static Development-only presentation data. These records are not accounts, persisted messages,
 * server state, or evidence that a transport is available. They exist solely to exercise the
 * native conversation-list information architecture while live client authorities remain absent.
 */
object DevelopmentConversationCatalog {
    fun previews(): List<ConversationPreview> = listOf(
        ConversationPreview(
            conversationId = "development-alex",
            displayName = "Alex",
            secondaryIdentity = "@alex",
            previewText = "I'll be there around six.",
            timestampLabel = "6:02 PM",
            provenance = CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.E2EE_ACTIVE,
            ),
            unreadCount = 2,
            isPinned = true,
            isMuted = false,
        ),
        ConversationPreview(
            conversationId = "development-team",
            displayName = "GoreeCloud Team",
            secondaryIdentity = "Development group",
            previewText = "Android client structure is taking shape.",
            timestampLabel = "5:41 PM",
            provenance = CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.UNKNOWN,
            ),
            unreadCount = 0,
            isPinned = false,
            isMuted = false,
        ),
        ConversationPreview(
            conversationId = "development-taylor",
            displayName = "Taylor",
            secondaryIdentity = null,
            previewText = "Carrier transport example only.",
            timestampLabel = "Yesterday",
            provenance = CommunicationProvenance(
                CommunicationTransport.SMS,
                CommunicationProtection.UNKNOWN,
            ),
            unreadCount = 0,
            isPinned = false,
            isMuted = true,
        ),
    )
}
