package com.goreecloud.messenger.client

import org.junit.Assert.assertArrayEquals
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class DataMessageSendCoordinatorTest {
    @Test
    fun blockedAuthorityReadinessNeverInvokesTransport() {
        var calls = 0
        val coordinator = coordinator(
            cryptography = DataMessagingReadiness.CryptographicState.NOT_ESTABLISHED,
            onTransportSubmit = {
                calls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )

        val result = coordinator.submit(message())

        assertEquals(0, calls)
        assertTrue(result is DataMessageSendCoordinator.Result.Blocked)
        val blocked = result as DataMessageSendCoordinator.Result.Blocked
        assertTrue(
            DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE in blocked.reasons,
        )
    }

    @Test
    fun activeE2eeWithoutAcceptedSecurityReviewNeverInvokesTransport() {
        var calls = 0
        val coordinator = coordinator(
            implementationReview = E2EEImplementationReviewState.NOT_ACCEPTED,
            onTransportSubmit = {
                calls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )

        val result = coordinator.submit(message())

        assertEquals(0, calls)
        assertEquals(
            DataMessageSendCoordinator.Result.Blocked(
                setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            ),
            result,
        )
    }

    @Test
    fun mismatchedAuthorizationAndE2eeScopesNeverInvokeTransport() {
        var calls = 0
        val coordinator = coordinator(
            authorizedConversationId = "conversation-1",
            e2eeConversationId = "conversation-2",
            onTransportSubmit = {
                calls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )

        val result = coordinator.submit(message())

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
        val coordinator = coordinator(
            authorizedConversationId = "conversation-2",
            e2eeConversationId = "conversation-2",
            onTransportSubmit = {
                calls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )

        val result = coordinator.submit(message())

        assertEquals(0, calls)
        // The stricter FR-005 cryptographic projection rejects conversation-2 as active E2EE
        // while resolving the prepared conversation-1 target, before the coordinator's later
        // verified-target comparison can be reached.
        assertEquals(
            DataMessageSendCoordinator.Result.Blocked(
                setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            ),
            result,
        )
    }

    @Test
    fun fullyVerifiedIndependentAuthoritiesInvokeOnlyInjectedDataTransport() {
        var calls = 0
        val coordinator = coordinator(
            onTransportSubmit = {
                calls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )

        val result = coordinator.submit(message())

        assertEquals(1, calls)
        assertTrue(result is DataMessageSendCoordinator.Result.Submitted)
        val submitted = result as DataMessageSendCoordinator.Result.Submitted
        assertEquals(CommunicationTransport.DATA, submitted.provenance.transport)
        assertEquals(CommunicationProtection.E2EE_ACTIVE, submitted.provenance.protection)
    }

    @Test
    fun failingAuthorityProviderNeverInvokesTransport() {
        var transportCalls = 0
        val resolver = DataMessagingAuthorityResolver(
            identityAuthority = GoreeCloudIdentitySessionAuthority {
                throw IllegalStateException("identity unavailable")
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority { conversationId ->
                ConversationAuthorizationEvidence(
                    state = DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
                    authorizedConversationId = conversationId,
                )
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority {
                DataMessagingReadiness.DataTransportState.AVAILABLE
            },
            e2eeSessionAuthority = E2EESessionAuthority { conversationId ->
                acceptedE2eeEvidence(conversationId)
            },
        )
        val coordinator = DataMessageSendCoordinator(
            authorityResolver = resolver,
            transport = EncryptedDataMessageTransport {
                transportCalls += 1
                EncryptedDataMessageTransport.Submission.Accepted
            },
        )

        val result = coordinator.submit(message())

        assertEquals(0, transportCalls)
        assertEquals(
            DataMessageSendCoordinator.Result.Blocked(
                setOf(DataMessagingReadiness.BlockReason.IDENTITY_NOT_AUTHENTICATED),
            ),
            result,
        )
    }

    @Test
    fun transportRejectionDoesNotInventFallbackSuccess() {
        val coordinator = coordinator(
            onTransportSubmit = {
                EncryptedDataMessageTransport.Submission.Rejected(
                    EncryptedDataMessageTransport.RejectionReason.TRANSPORT_UNAVAILABLE,
                )
            },
        )

        val result = coordinator.submit(message())

        assertEquals(
            DataMessageSendCoordinator.Result.TransportRejected(
                EncryptedDataMessageTransport.RejectionReason.TRANSPORT_UNAVAILABLE,
            ),
            result,
        )
    }

    @Test
    fun transportExceptionFailsClosedAsUnknownRejection() {
        val coordinator = coordinator(
            onTransportSubmit = {
                throw IllegalStateException("transport adapter failed")
            },
        )

        val result = coordinator.submit(message())

        assertEquals(
            DataMessageSendCoordinator.Result.TransportRejected(
                EncryptedDataMessageTransport.RejectionReason.UNKNOWN,
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

    private fun coordinator(
        identity: DataMessagingReadiness.IdentityState = DataMessagingReadiness.IdentityState.AUTHENTICATED,
        conversationAccess: DataMessagingReadiness.ConversationAccessState =
            DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
        authorizedConversationId: String? = "conversation-1",
        transport: DataMessagingReadiness.DataTransportState = DataMessagingReadiness.DataTransportState.AVAILABLE,
        cryptography: DataMessagingReadiness.CryptographicState =
            DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
        e2eeConversationId: String? = "conversation-1",
        implementationReview: E2EEImplementationReviewState = E2EEImplementationReviewState.ACCEPTED,
        deviceIdentity: E2EEDeviceIdentityState = E2EEDeviceIdentityState.ENROLLED,
        sessionEstablishment: E2EESessionEstablishmentState = E2EESessionEstablishmentState.ESTABLISHED,
        keyLifecycle: E2EEKeyLifecycleState = E2EEKeyLifecycleState.CURRENT,
        onTransportSubmit: (PreparedEncryptedDataMessage) -> EncryptedDataMessageTransport.Submission,
    ): DataMessageSendCoordinator {
        val resolver = DataMessagingAuthorityResolver(
            identityAuthority = GoreeCloudIdentitySessionAuthority { identity },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority {
                ConversationAuthorizationEvidence(
                    state = conversationAccess,
                    authorizedConversationId = authorizedConversationId,
                )
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority { transport },
            e2eeSessionAuthority = E2EESessionAuthority {
                E2EESessionEvidence(
                    state = cryptography,
                    e2eeConversationId = e2eeConversationId,
                    implementationReview = implementationReview,
                    deviceIdentity = deviceIdentity,
                    sessionEstablishment = sessionEstablishment,
                    keyLifecycle = keyLifecycle,
                )
            },
        )
        return DataMessageSendCoordinator(
            authorityResolver = resolver,
            transport = EncryptedDataMessageTransport { message -> onTransportSubmit(message) },
        )
    }

    private fun acceptedE2eeEvidence(conversationId: String): E2EESessionEvidence =
        E2EESessionEvidence(
            state = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
            e2eeConversationId = conversationId,
            implementationReview = E2EEImplementationReviewState.ACCEPTED,
            deviceIdentity = E2EEDeviceIdentityState.ENROLLED,
            sessionEstablishment = E2EESessionEstablishmentState.ESTABLISHED,
            keyLifecycle = E2EEKeyLifecycleState.CURRENT,
        )

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
