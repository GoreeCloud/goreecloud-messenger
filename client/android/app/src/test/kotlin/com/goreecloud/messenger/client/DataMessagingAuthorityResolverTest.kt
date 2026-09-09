package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class DataMessagingAuthorityResolverTest {
    @Test
    fun resolvesEachIndependentAuthorityForExactRequestedConversation() {
        val requestedScopes = mutableListOf<String>()
        val resolver = DataMessagingAuthorityResolver(
            identityAuthority = GoreeCloudIdentitySessionAuthority {
                DataMessagingReadiness.IdentityState.AUTHENTICATED
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority { conversationId ->
                requestedScopes += "authorization:$conversationId"
                ConversationAuthorizationEvidence(
                    state = DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
                    authorizedConversationId = conversationId,
                )
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority {
                DataMessagingReadiness.DataTransportState.AVAILABLE
            },
            e2eeSessionAuthority = E2EESessionAuthority { conversationId ->
                requestedScopes += "e2ee:$conversationId"
                E2EESessionEvidence(
                    state = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
                    e2eeConversationId = conversationId,
                )
            },
        )

        val evidence = resolver.evidenceFor("conversation 1")
        val readiness = DataMessagingReadiness.evaluate(evidence)

        assertEquals(
            listOf("authorization:conversation 1", "e2ee:conversation 1"),
            requestedScopes,
        )
        assertEquals(
            "conversation 1",
            (readiness as DataMessagingReadiness.Result.Ready).verifiedConversationId,
        )
    }

    @Test
    fun oneProviderFailureFailsClosedWithoutUpgradingFromOtherAuthorities() {
        val resolver = DataMessagingAuthorityResolver(
            identityAuthority = GoreeCloudIdentitySessionAuthority {
                DataMessagingReadiness.IdentityState.AUTHENTICATED
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority {
                throw IllegalStateException("authorization unavailable")
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority {
                DataMessagingReadiness.DataTransportState.AVAILABLE
            },
            e2eeSessionAuthority = E2EESessionAuthority { conversationId ->
                E2EESessionEvidence(
                    state = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
                    e2eeConversationId = conversationId,
                )
            },
        )

        val result = DataMessagingReadiness.evaluate(resolver.evidenceFor("conversation-1"))

        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    @Test
    fun transportFailureCannotBorrowReadinessFromIdentityAuthorizationOrE2ee() {
        val resolver = DataMessagingAuthorityResolver(
            identityAuthority = GoreeCloudIdentitySessionAuthority {
                DataMessagingReadiness.IdentityState.AUTHENTICATED
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority { conversationId ->
                ConversationAuthorizationEvidence(
                    state = DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
                    authorizedConversationId = conversationId,
                )
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority {
                throw IllegalStateException("transport unavailable")
            },
            e2eeSessionAuthority = E2EESessionAuthority { conversationId ->
                E2EESessionEvidence(
                    state = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
                    e2eeConversationId = conversationId,
                )
            },
        )

        val result = DataMessagingReadiness.evaluate(resolver.evidenceFor("conversation-1"))

        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    @Test
    fun invalidRequestedScopeFailsBeforeAnyAuthorityIsQueried() {
        var calls = 0
        val resolver = DataMessagingAuthorityResolver(
            identityAuthority = GoreeCloudIdentitySessionAuthority {
                calls += 1
                DataMessagingReadiness.IdentityState.AUTHENTICATED
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority {
                calls += 1
                ConversationAuthorizationEvidence(DataMessagingReadiness.ConversationAccessState.UNKNOWN)
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority {
                calls += 1
                DataMessagingReadiness.DataTransportState.UNKNOWN
            },
            e2eeSessionAuthority = E2EESessionAuthority {
                calls += 1
                E2EESessionEvidence(DataMessagingReadiness.CryptographicState.UNKNOWN)
            },
        )

        try {
            resolver.evidenceFor(" conversation-1")
            throw AssertionError("noncanonical scope was accepted")
        } catch (_: IllegalArgumentException) {
            // Expected before any authority is queried.
        }

        assertEquals(0, calls)
    }

    @Test
    fun mismatchedPositiveAuthorityScopesRemainBlocked() {
        val resolver = DataMessagingAuthorityResolver(
            identityAuthority = GoreeCloudIdentitySessionAuthority {
                DataMessagingReadiness.IdentityState.AUTHENTICATED
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority {
                ConversationAuthorizationEvidence(
                    state = DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
                    authorizedConversationId = "conversation-1",
                )
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority {
                DataMessagingReadiness.DataTransportState.AVAILABLE
            },
            e2eeSessionAuthority = E2EESessionAuthority {
                E2EESessionEvidence(
                    state = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
                    e2eeConversationId = "conversation-2",
                )
            },
        )

        val result = DataMessagingReadiness.evaluate(resolver.evidenceFor("conversation-1"))

        assertTrue(result is DataMessagingReadiness.Result.Blocked)
        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }
}
