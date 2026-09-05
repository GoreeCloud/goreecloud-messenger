package com.goreecloud.messenger.client

/**
 * Opaque client-prepared GoreeCloud Data message ready for a future transport adapter.
 *
 * This type contains ciphertext and protocol identifiers only. It deliberately does not carry a
 * caller-authored sender identity, credential, key, plaintext body, carrier fallback, or transport
 * endpoint. Authenticated identity remains owned by the runtime authentication authority.
 */
class PreparedEncryptedDataMessage private constructor(
    val messageId: String,
    val conversationId: String,
    val clientNonce: String,
    ciphertext: ByteArray,
) {
    private val encryptedBytes = ciphertext.copyOf()

    val ciphertextSizeBytes: Int
        get() = encryptedBytes.size

    fun ciphertextCopy(): ByteArray = encryptedBytes.copyOf()

    companion object {
        fun create(
            messageId: String,
            conversationId: String,
            clientNonce: String,
            ciphertext: ByteArray,
        ): PreparedEncryptedDataMessage {
            require(messageId.isNotBlank()) { "messageId must not be blank" }
            require(conversationId.isNotBlank()) { "conversationId must not be blank" }
            require(clientNonce.isNotBlank()) { "clientNonce must not be blank" }
            require(ciphertext.isNotEmpty()) { "ciphertext must not be empty" }

            return PreparedEncryptedDataMessage(
                messageId = messageId.trim(),
                conversationId = conversationId.trim(),
                clientNonce = clientNonce.trim(),
                ciphertext = ciphertext,
            )
        }
    }
}

/**
 * Future transport seam for already-prepared encrypted Data messages.
 *
 * The current Development client does not provide a production implementation. A later adapter
 * may implement this interface only after its own Identity/session, authorization, network,
 * protocol, and deployment requirements are satisfied.
 */
fun interface EncryptedDataMessageTransport {
    fun submit(message: PreparedEncryptedDataMessage): Submission

    sealed interface Submission {
        data object Accepted : Submission

        data class Rejected(val reason: RejectionReason) : Submission
    }

    enum class RejectionReason {
        TRANSPORT_UNAVAILABLE,
        AUTHORIZATION_REJECTED,
        PROTOCOL_REJECTED,
        UNKNOWN,
    }
}

/**
 * Enforces the fail-closed readiness policy at the final client seam before any injected Data
 * transport can be invoked.
 *
 * Conversation authorization must be verified for the exact conversation carried by the prepared
 * message. A valid participant decision for another conversation cannot be reused at this seam.
 *
 * This coordinator does not authenticate, authorize, encrypt, persist, retry, queue, synchronize,
 * or send on its own. It simply refuses to call the supplied transport unless all four independent
 * readiness authorities are positively verified and the authorization scope matches the attempted
 * operation.
 */
class DataMessageSendCoordinator(
    private val transport: EncryptedDataMessageTransport,
) {
    sealed interface Result {
        data class Blocked(val reasons: Set<DataMessagingReadiness.BlockReason>) : Result

        data class Submitted(
            val provenance: CommunicationProvenance,
        ) : Result

        data class TransportRejected(
            val reason: EncryptedDataMessageTransport.RejectionReason,
        ) : Result
    }

    fun submit(
        evidence: DataMessagingReadiness.Evidence,
        message: PreparedEncryptedDataMessage,
    ): Result =
        when (val readiness = DataMessagingReadiness.evaluate(evidence)) {
            is DataMessagingReadiness.Result.Blocked -> Result.Blocked(readiness.reasons)
            is DataMessagingReadiness.Result.Ready -> {
                if (readiness.verifiedConversationId != message.conversationId) {
                    Result.Blocked(
                        setOf(DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED),
                    )
                } else {
                    when (val submission = transport.submit(message)) {
                        EncryptedDataMessageTransport.Submission.Accepted ->
                            Result.Submitted(readiness.provenance)

                        is EncryptedDataMessageTransport.Submission.Rejected ->
                            Result.TransportRejected(submission.reason)
                    }
                }
            }
        }
}
