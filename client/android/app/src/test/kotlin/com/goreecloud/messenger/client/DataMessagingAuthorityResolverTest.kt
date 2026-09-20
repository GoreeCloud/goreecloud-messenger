package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class DataMessagingAuthorityResolverTest {
    @Test
    fun bareAuthenticatedIdentityClaimFailsClosed() {
        val evidence = resolver(
            identityEvidence = IdentitySessionEvidence(
                state = DataMessagingReadiness.IdentityState.AUTHENTICATED,
            ),
        ).evidenceFor("conversation-1")

        assertIdentityBlocked(evidence)
    }

    @Test
    fun missingIdentitySessionBindingFailsClosed() {
        val evidence = resolver(
            identityEvidence = acceptedIdentityEvidence().copy(
                sessionBinding = IdentityBindingState.NOT_BOUND,
            ),
        ).evidenceFor("conversation-1")

        assertIdentityBlocked(evidence)
    }

    @Test
    fun missingIdentityDeviceBindingFailsClosed() {
        val evidence = resolver(
            identityEvidence = acceptedIdentityEvidence().copy(
                deviceBinding = IdentityBindingState.NOT_BOUND,
            ),
        ).evidenceFor("conversation-1")

        assertIdentityBlocked(evidence)
    }

    @Test
    fun unauthenticatedIdentityIsNotUpgradedByPositiveBindings() {
        val evidence = resolver(
            identityEvidence = acceptedIdentityEvidence().copy(
                state = DataMessagingReadiness.IdentityState.UNAUTHENTICATED,
            ),
        ).evidenceFor("conversation-1")

        assertEquals(DataMessagingReadiness.IdentityState.UNAUTHENTICATED, evidence.identity)
        assertIdentityBlocked(evidence)
    }

    @Test
    fun bareVerifiedParticipantClaimFailsClosed() {
        val evidence = resolver(
            authorizationProvider = { conversationId ->
                ConversationAuthorizationEvidence(
                    state = DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
                    authorizedConversationId = conversationId,
                )
            },
        ).evidenceFor("conversation-1")

        assertConversationBlocked(evidence)
    }

    @Test
    fun missingConversationIdentityBindingFailsClosed() {
        val evidence = resolver(
            authorizationProvider = { conversationId ->
                acceptedConversationAuthorizationEvidence(conversationId).copy(
                    identityBinding = ConversationAuthorizationAcceptanceState.NOT_ACCEPTED,
                )
            },
        ).evidenceFor("conversation-1")

        assertConversationBlocked(evidence)
    }

    @Test
    fun staleConversationAuthorizationDecisionFailsClosed() {
        val evidence = resolver(
            authorizationProvider = { conversationId ->
                acceptedConversationAuthorizationEvidence(conversationId).copy(
                    decisionFreshness = ConversationAuthorizationAcceptanceState.NOT_ACCEPTED,
                )
            },
        ).evidenceFor("conversation-1")

        assertConversationBlocked(evidence)
    }

    @Test
    fun mismatchedConversationAuthorizationScopeFailsClosed() {
        val evidence = resolver(
            authorizationProvider = {
                acceptedConversationAuthorizationEvidence("conversation-2")
            },
        ).evidenceFor("conversation-1")

        assertConversationBlocked(evidence)
    }

    @Test
    fun resolvesEachIndependentAuthorityForExactRequestedConversation() {
        val requestedScopes = mutableListOf<String>()
        val resolver = resolver(
            e2eeProvider = { conversationId ->
                requestedScopes += "e2ee:$conversationId"
                acceptedE2eeEvidence(conversationId)
            },
            authorizationObserver = { conversationId ->
                requestedScopes += "authorization:$conversationId"
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
    fun bareAvailableTransportClaimFailsClosed() {
        val resolver = resolver(
            transportEvidence = DataTransportEvidence(
                state = DataMessagingReadiness.DataTransportState.AVAILABLE,
            ),
        )

        assertTransportBlocked(resolver.evidenceFor("conversation-1"))
    }

    @Test
    fun everyMissingTransportAcceptanceFactFailsClosed() {
        val rejected = DataTransportAcceptanceState.NOT_ACCEPTED
        val candidates = listOf(
            acceptedTransportEvidence().copy(configuration = rejected),
            acceptedTransportEvidence().copy(authenticationBinding = rejected),
            acceptedTransportEvidence().copy(channelProtection = rejected),
            acceptedTransportEvidence().copy(failurePolicy = rejected),
        )

        candidates.forEach { evidence ->
            assertTransportBlocked(
                resolver(transportEvidence = evidence).evidenceFor("conversation-1"),
            )
        }
    }

    @Test
    fun unavailableTransportIsNotUpgradedByPositiveAcceptanceFacts() {
        val resolver = resolver(
            transportEvidence = acceptedTransportEvidence().copy(
                state = DataMessagingReadiness.DataTransportState.UNAVAILABLE,
            ),
        )

        val evidence = resolver.evidenceFor("conversation-1")

        assertEquals(DataMessagingReadiness.DataTransportState.UNAVAILABLE, evidence.transport)
        assertTransportBlocked(evidence)
    }

    @Test
    fun bareActiveClaimWithoutAcceptanceEvidenceFailsClosed() {
        val resolver = resolver(
            e2eeProvider = { conversationId ->
                E2EESessionEvidence(
                    state = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
                    e2eeConversationId = conversationId,
                )
            },
        )

        assertE2eeBlocked(resolver.evidenceFor("conversation-1"))
    }

    @Test
    fun rejectedImplementationReviewFailsClosed() {
        val resolver = resolver(
            e2eeProvider = { conversationId ->
                acceptedE2eeEvidence(conversationId).copy(
                    implementationReview = E2EEImplementationReviewState.NOT_ACCEPTED,
                )
            },
        )

        assertE2eeBlocked(resolver.evidenceFor("conversation-1"))
    }

    @Test
    fun missingDeviceIdentityEnrollmentFailsClosed() {
        val resolver = resolver(
            e2eeProvider = { conversationId ->
                acceptedE2eeEvidence(conversationId).copy(
                    deviceIdentity = E2EEDeviceIdentityState.NOT_ENROLLED,
                )
            },
        )

        assertE2eeBlocked(resolver.evidenceFor("conversation-1"))
    }

    @Test
    fun unestablishedSessionFailsClosed() {
        val resolver = resolver(
            e2eeProvider = { conversationId ->
                acceptedE2eeEvidence(conversationId).copy(
                    sessionEstablishment = E2EESessionEstablishmentState.NOT_ESTABLISHED,
                )
            },
        )

        assertE2eeBlocked(resolver.evidenceFor("conversation-1"))
    }

    @Test
    fun noncurrentKeyLifecycleFailsClosed() {
        val resolver = resolver(
            e2eeProvider = { conversationId ->
                acceptedE2eeEvidence(conversationId).copy(
                    keyLifecycle = E2EEKeyLifecycleState.NOT_CURRENT,
                )
            },
        )

        assertE2eeBlocked(resolver.evidenceFor("conversation-1"))
    }

    @Test
    fun noncanonicalE2eeScopeFailsClosedInsteadOfAliasing() {
        val resolver = resolver(
            e2eeProvider = {
                acceptedE2eeEvidence(" conversation-1")
            },
        )

        assertE2eeBlocked(resolver.evidenceFor("conversation-1"))
    }

    @Test
    fun negativeCryptographicStateIsNotUpgradedByPositiveAcceptanceFacts() {
        val resolver = resolver(
            e2eeProvider = { conversationId ->
                acceptedE2eeEvidence(conversationId).copy(
                    state = DataMessagingReadiness.CryptographicState.NOT_ESTABLISHED,
                )
            },
        )

        val evidence = resolver.evidenceFor("conversation-1")

        assertEquals(
            DataMessagingReadiness.CryptographicState.NOT_ESTABLISHED,
            evidence.cryptography,
        )
        assertE2eeBlocked(evidence)
    }

    @Test
    fun oneProviderFailureFailsClosedWithoutUpgradingFromOtherAuthorities() {
        val resolver = DataMessagingAuthorityResolver(
            identityAuthority = GoreeCloudIdentitySessionAuthority {
                acceptedIdentityEvidence()
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority {
                throw IllegalStateException("authorization unavailable")
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority {
                acceptedTransportEvidence()
            },
            e2eeSessionAuthority = E2EESessionAuthority { conversationId ->
                acceptedE2eeEvidence(conversationId)
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
                acceptedIdentityEvidence()
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority { conversationId ->
                acceptedConversationAuthorizationEvidence(conversationId)
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority {
                throw IllegalStateException("transport unavailable")
            },
            e2eeSessionAuthority = E2EESessionAuthority { conversationId ->
                acceptedE2eeEvidence(conversationId)
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
                acceptedIdentityEvidence()
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority {
                calls += 1
                ConversationAuthorizationEvidence(DataMessagingReadiness.ConversationAccessState.UNKNOWN)
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority {
                calls += 1
                DataTransportEvidence(DataMessagingReadiness.DataTransportState.UNKNOWN)
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
        val resolver = resolver(
            e2eeProvider = {
                acceptedE2eeEvidence("conversation-2")
            },
        )

        val result = DataMessagingReadiness.evaluate(resolver.evidenceFor("conversation-1"))

        assertTrue(result is DataMessagingReadiness.Result.Blocked)
        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    private fun resolver(
        identityEvidence: IdentitySessionEvidence = acceptedIdentityEvidence(),
        e2eeProvider: (String) -> E2EESessionEvidence = { conversationId ->
            acceptedE2eeEvidence(conversationId)
        },
        transportEvidence: DataTransportEvidence = acceptedTransportEvidence(),
        authorizationProvider: (String) -> ConversationAuthorizationEvidence = { conversationId ->
            acceptedConversationAuthorizationEvidence(conversationId)
        },
        authorizationObserver: (String) -> Unit = {},
    ): DataMessagingAuthorityResolver =
        DataMessagingAuthorityResolver(
            identityAuthority = GoreeCloudIdentitySessionAuthority {
                identityEvidence
            },
            conversationAuthorizationAuthority = ConversationAuthorizationAuthority { conversationId ->
                authorizationObserver(conversationId)
                authorizationProvider(conversationId)
            },
            dataTransportAuthority = GoreeCloudDataTransportAuthority { transportEvidence },
            e2eeSessionAuthority = E2EESessionAuthority(e2eeProvider),
        )

    private fun acceptedIdentityEvidence(): IdentitySessionEvidence =
        IdentitySessionEvidence(
            state = DataMessagingReadiness.IdentityState.AUTHENTICATED,
            sessionBinding = IdentityBindingState.BOUND,
            deviceBinding = IdentityBindingState.BOUND,
        )

    private fun acceptedConversationAuthorizationEvidence(
        conversationId: String,
    ): ConversationAuthorizationEvidence =
        ConversationAuthorizationEvidence(
            state = DataMessagingReadiness.ConversationAccessState.VERIFIED_PARTICIPANT,
            authorizedConversationId = conversationId,
            identityBinding = ConversationAuthorizationAcceptanceState.ACCEPTED,
            decisionFreshness = ConversationAuthorizationAcceptanceState.ACCEPTED,
        )

    private fun acceptedTransportEvidence(): DataTransportEvidence =
        DataTransportEvidence(
            state = DataMessagingReadiness.DataTransportState.AVAILABLE,
            configuration = DataTransportAcceptanceState.ACCEPTED,
            authenticationBinding = DataTransportAcceptanceState.ACCEPTED,
            channelProtection = DataTransportAcceptanceState.ACCEPTED,
            failurePolicy = DataTransportAcceptanceState.ACCEPTED,
        )

    private fun acceptedE2eeEvidence(conversationId: String): E2EESessionEvidence =
        E2EESessionEvidence(
            state = DataMessagingReadiness.CryptographicState.E2EE_ACTIVE,
            e2eeConversationId = conversationId,
            implementationReview = E2EEImplementationReviewState.ACCEPTED,
            deviceIdentity = E2EEDeviceIdentityState.ENROLLED,
            sessionEstablishment = E2EESessionEstablishmentState.ESTABLISHED,
            keyLifecycle = E2EEKeyLifecycleState.CURRENT,
        )

    private fun assertIdentityBlocked(evidence: DataMessagingReadiness.Evidence) {
        val result = DataMessagingReadiness.evaluate(evidence)
        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.IDENTITY_NOT_AUTHENTICATED),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    private fun assertConversationBlocked(evidence: DataMessagingReadiness.Evidence) {
        val result = DataMessagingReadiness.evaluate(evidence)
        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.CONVERSATION_ACCESS_NOT_VERIFIED),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    private fun assertTransportBlocked(evidence: DataMessagingReadiness.Evidence) {
        val result = DataMessagingReadiness.evaluate(evidence)
        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.DATA_TRANSPORT_NOT_AVAILABLE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }

    private fun assertE2eeBlocked(evidence: DataMessagingReadiness.Evidence) {
        val result = DataMessagingReadiness.evaluate(evidence)
        assertEquals(
            setOf(DataMessagingReadiness.BlockReason.E2EE_NOT_VERIFIED_ACTIVE),
            (result as DataMessagingReadiness.Result.Blocked).reasons,
        )
    }
}
