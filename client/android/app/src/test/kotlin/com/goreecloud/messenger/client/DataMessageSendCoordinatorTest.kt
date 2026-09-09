package com.goreecloud.messenger.client

import org.junit.Assert.assertArrayEquals
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class DataMessageSendCoordinatorTest {
    private val readyEvidence = DataMessagingReadiness.Evidence(
        identity = DataMessagingReadiness.IdentityState.AUTHENTICATED,
        conversationAccess = DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
        transport = DataMessagingReadiness.DataTransportState.AVAILABLE,
        cryptography = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
        authorizedConversationId = "conversation-1",
        e2eeConversationId = "conversation-1",
    )

    @Test
    fun blockedReadinessNeverInvokesTransport() {
        var calls = 0
        val coordinator = DataMessageSendCoordinator(
            transport = EncryptedDataMessageTransport {
                calls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )
        val blockedEvidence = readyEvidence.copy(
            cryptography = DataMessagingReadiness.CryptographicState.NOT_ESTABLISHED,
        )

        val result = coordinator.submit(blockedEvidence, message())

        assertEquals(0, calls)
        assertTrue(result is DataMessageSendCoordinator.Result.Blocked)
        val blocked = result as DataMessageSendCoordinator.Result.Blocked
        assertTrue(
            DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE in blocked.reasons,
        )
    }

    @Test
    fun mismatchedAuthorizationAndE2eeScopesNeverInvokeTransport() {
        var calls = 0
        val coordinator = DataMessageSendCoordinator(
            transport = EncryptedDataMessageTransport {
                calls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )
        val mismatchedEvidence = readyEvidence.copy(
            e2eeConversationId = "conversation-2",
        )

        val result = coordinator.submit(mismatchedEvidence, message())

        assertEquals(0, calls)
        assertEquals(
            DataMessageSendCoordinator.Result.Blocked(
                setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            ),
            result,
        )
    }

    @Test
    fun verifiedDifferentConversationNeverInvokesTransportForPreparedTarget() {
        var calls = 0
        val coordinator = DataMessageSendCoordinator(
            transport = EncryptedDataMessageTransport {
                calls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )
        val differentConversationEvidence = readyEvidence.copy(
            authorizedConversationId = "conversation-2",
            e2eeConversationId = "conversation-2",
        )

        val result = coordinator.submit(differentConversationEvidence, message())

        assertEquals(0, calls)
        assertEquals(
            DataMessageSendCoordinator.Result.Blocked(
                setOf(DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED),
            ),
            result,
        )
    }

    @Test
    fun fullyVerifiedReadinessInvokesOnlyInjectedDataTransport() {
        var calls = 0
        val coordinator = DataMessageSendCoordinator(
            transport = EncryptedDataMessageTransport {
                calls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )

        val result = coordinator.submit(readyEvidence, message())

        assertEquals(1, calls)
        assertTrue(result is DataMessageSendCoordinator.Result.Submitted)
        val submitted = result as DataMessageSendCoordinator.Result.Submitted
        assertEquals(CommunicationTransport.DATA, submitted.provenance.transport)
        assertEquals(CommunicationProtection.E2EE_ACTIVE, submitted.provenance.protection)
    }

    @Test
    fun transportRejectionDoesNotInventFallbackSuccess() {
        val coordinator = DataMessageSendCoordinator(
            transport = EncryptedDataMessageTransport {
                EncryptedDataMessageTransport.Submission.Rejected(
                    EncryptedDataMessageTransport.RejectionReason.TRANSPORT_UNAVAILABLE,
                )
            },
        )

        val result = coordinator.submit(readyEvidence, message())

        assertEquals(
            DataMessageSendCoordinator.Result.TransportRejected(
                EncryptedDataMessageTransport.RejectionReason.TRANSPORT_UNAVAILABLE,
            ),
            result,
        )
    }

    @Test
    fun preparedCiphertextIsDefensivelyCopied() {
        val source = byteArrayOf(1, 2, 3, 4)
        val prepared = message(ciphertext = source)
        source[0] = 99

        val firstRead = prepared.ciphertextCopy()
        firstRead[1] = 88

        assertArrayEquals(byteArrayOf(1, 2, 3, 4), prepared.ciphertextCopy())
    }

    @Test
    fun preparedMessageRequiresExactBoundedOpaqueIdentifiers() {
        val oversized = "x".repeat(DataReceiptIdentifierPolicy.MAX_IDENTIFIER_LENGTH + 1)
        val invalid = listOf(" conversation-1", "conversation-1 ", "conversation\n1", oversized)

        invalid.forEach { conversationId ->
            try {
                PreparedEncryptedDataMessage.create(
                    messageId = "message-1",
                    conversationId = conversationId,
                    clientNonce = "nonce-1",
                    ciphertext = byteArrayOf(1),
                )
                throw AssertionError("invalid conversation identifier was accepted: $conversationId")
            } catch (_: IllegalArgumentException) {
                // Expected fail-closed construction.
            }
        }

        val prepared = PreparedEncryptedDataMessage.create(
            messageId = "message / opaque",
            conversationId = "conversation:one/two",
            clientNonce = "nonce / opaque",
            ciphertext = byteArrayOf(1),
        )
        assertEquals("message / opaque", prepared.messageId)
        assertEquals("conversation:one/two", prepared.conversationId)
        assertEquals("nonce / opaque", prepared.clientNonce)
    }

    @Test(expected = IllegalArgumentException::class)
    fun preparedMessageRejectsEmptyCiphertext() {
        message(ciphertext = byteArrayOf())
    }

    private fun message(
        ciphertext: ByteArray = byteArrayOf(10, 20, 30),
    ): PreparedEncryptedDataMessage =
        PreparedEncryptedDataMessage.create(
            messageId = "message-1",
            conversationId = "conversation-1",
            clientNonce = "nonce-1",
            ciphertext = ciphertext,
        )
}
