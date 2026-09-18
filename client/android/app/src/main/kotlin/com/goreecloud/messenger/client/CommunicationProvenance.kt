package com.goreecloud.messenger.client

enum class CommunicationTransport(val label: String) {
    DATA("Data"),
    SMS("SMS"),
    MMS("MMS"),
    RCS("RCS"),
}

enum class CommunicationProtection {
    E2EE_ACTIVE,
    NOT_E2EE,
    UNAVAILABLE,
    UNKNOWN,
}

data class CommunicationProvenance(
    val transport: CommunicationTransport,
    val protection: CommunicationProtection,
) {
    init {
        require(protection != CommunicationProtection.E2EE_ACTIVE || transport == CommunicationTransport.DATA) {
            "the current Messenger client may claim GoreeCloud E2EE only for verified Data transport"
        }
    }

    fun displayLabel(): String = when (protection) {
        CommunicationProtection.E2EE_ACTIVE -> "${transport.label} · E2EE"
        CommunicationProtection.NOT_E2EE -> when (transport) {
            CommunicationTransport.DATA -> "Data · Not end-to-end encrypted"
            else -> transport.label
        }
        CommunicationProtection.UNAVAILABLE -> when (transport) {
            CommunicationTransport.DATA -> "Data · E2EE unavailable"
            else -> transport.label
        }
        CommunicationProtection.UNKNOWN -> when (transport) {
            CommunicationTransport.DATA -> "Data · Protection not verified"
            else -> transport.label
        }
    }
}
