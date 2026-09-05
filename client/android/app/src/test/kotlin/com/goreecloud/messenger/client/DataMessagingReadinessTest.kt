package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class DataMessagingReadinessTest {
    private val fullyReady = DataMessagingReadiness.Evidence(
        identity = DataMessagingReadiness.IdentityState.AUTHENTICATED,
        conversationAccess = DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
        transport = DataMessagingReadiness.DataTransportState.AVAILABLE,
        cryptography = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
        authorizedConversationId = "conversation-1",
        e2eeConversationId = "conversation-1",
    )

    @Test
    fun everyVerifiedPrerequisiteProducesOnlyDataE2eeReadiness() {
        val result = DataMessagingReadiness.evaluate(fullyReady)

        assertTrue(result is DataMessagingReadiness.Result.Ready)
        val ready = result as DataMessagingReadiness.Result.Ready
        assertEquals("conversation-1", ready.verifiedConversationId)
        assertEquals(CommunicationTransport.DATA, ready.provenance.transport)
        assertEquals(CommunicationProtection.E2EE_ACTIVE, ready.provenance.protection)
        assertEquals("Data · E2EE", ready.provenance.displayLabel())
    }

    @Test
    fun independentConversationScopesAreNormalizedBeforeReadiness() {
        val result = DataMessagingReadiness.evaluate(
            fullyReady.copy(
                authorizedConversationId = "  conversation-1  ",
                e2eeConversationId = "\tconversation-1\n",
            ),
        )

        assertEquals(
            "conversation-1",
            (result as DataMessagingReadiness.Result.Ready).verifiedConversationId,
        )
    }

    @Test
    fun participantStateWithoutAuthorizationScopeFailsClosed() {
        val result = DataMessagingReadiness.evaluate(
            fullyReady.copy(authorizedConversationId = "   "),
        )

        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    @Test
    fun activeE2eeWithoutCryptographicConversationScopeFailsClosed() {
        val result = DataMessagingReadiness.evaluate(
            fullyReady.copy(e2eeConversationId = "   "),
        )

        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    @Test
    fun authorizationAndE2eeScopesMustIdentifySameConversation() {
        val result = DataMessagingReadiness.evaluate(
            fullyReady.copy(e2eeConversationId = "conversation-2"),
        )

        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    @Test
    fun unknownEvidenceFailsClosedAcrossEveryAuthority() {
        val result = DataMessagingReadiness.evaluate(
            DataMessagingReadiness.Evidence(
                identity = DataMessagingReadiness.IdentityState.UNKNOWN,
                conversationAccess = DataMessagingReadiness.ConversationAccessState.UNKNOWN,
                transport = DataMessagingReadiness.DataTransportState.UNKNOWN,
                cryptography = DataMessagingReadiness.CryptographicState.UNKNOWN,
            )
        )

        assertTrue(result is DataMessagingReadiness.Result.Blocked)
        assertEquals(
            DataMessagingReadiness.BlockReason.entries.toSet(),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    @Test
    fun transportAvailabilityCannotUpgradeUnverifiedCryptography() {
        val result = DataMessagingReadiness.evaluate(
            fullyReady.copy(cryptography = DataMessagingReadiness.CryptographicState.NOT_ESTABLISHED)
        )

        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    @Test
    fun activeCryptographyCannotBypassIdentityOrConversationAuthorization() {
        val result = DataMessagingReadiness.evaluate(
            fullyReady.copy(
                identity = DataMessagingReadiness.IdentityState.UNAUTHENTICATED,
                conversationAccess = DataMessagingReadiness.ConversationAccessState.NOT_PARTICIPANT,
            )
        )

        assertEquals(
            setOf(
                DataMessagingReadiness.BlockReason.IDENTITY_NOT_AUTHENTICATED,
                DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED,
            ),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    @Test
    fun unavailableDataTransportDoesNotFallBackToCarrierMessaging() {
        val result = DataMessagingReadiness.evaluate(
            fullyReady.copy(transport = DataMessagingReadiness.DataTransportState.UNAVAILABLE)
        )

        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    @Test
    fun e2eeUnavailableIsBlockedRatherThanPresentedAsEncryptedData() {
        val result = DataMessagingReadiness.evaluate(
            fullyReady.copy(cryptography = DataMessagingReadiness.CryptographicState.E2EE_UNAVAILABLE)
        )

        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }
}
