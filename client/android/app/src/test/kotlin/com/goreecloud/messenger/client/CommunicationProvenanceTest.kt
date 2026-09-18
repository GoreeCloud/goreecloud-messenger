package com.goreecloud.messenger.client

import org.junit.Assert.assertEquals
import org.junit.Assert.assertThrows
import org.junit.Test

class CommunicationProvenanceTest {
    @Test
    fun verifiedDataE2eeHasExplicitCombinedLabel() {
        assertEquals(
            "Data · E2EE",
            CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.E2EE_ACTIVE,
            ).displayLabel(),
        )
    }

    @Test
    fun unverifiedDataNeverLooksEncrypted() {
        assertEquals(
            "Data · Protection not verified",
            CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.UNKNOWN,
            ).displayLabel(),
        )
        assertEquals(
            "Data · E2EE unavailable",
            CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.UNAVAILABLE,
            ).displayLabel(),
        )
        assertEquals(
            "Data · Not end-to-end encrypted",
            CommunicationProvenance(
                CommunicationTransport.DATA,
                CommunicationProtection.NOT_E2EE,
            ).displayLabel(),
        )
    }

    @Test
    fun carrierTransportsRetainTheirActualTransportLabels() {
        assertEquals(
            "SMS",
            CommunicationProvenance(
                CommunicationTransport.SMS,
                CommunicationProtection.UNKNOWN,
            ).displayLabel(),
        )
        assertEquals(
            "MMS",
            CommunicationProvenance(
                CommunicationTransport.MMS,
                CommunicationProtection.NOT_E2EE,
            ).displayLabel(),
        )
        assertEquals(
            "RCS",
            CommunicationProvenance(
                CommunicationTransport.RCS,
                CommunicationProtection.UNKNOWN,
            ).displayLabel(),
        )
    }

    @Test
    fun carrierTransportCannotClaimGoreecloudE2ee() {
        listOf(
            CommunicationTransport.SMS,
            CommunicationTransport.MMS,
            CommunicationTransport.RCS,
        ).forEach { transport ->
            assertThrows(IllegalArgumentException::class.java) {
                CommunicationProvenance(transport, CommunicationProtection.E2EE_ACTIVE)
            }
        }
    }
}
