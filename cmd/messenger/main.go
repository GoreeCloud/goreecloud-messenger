// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"fmt"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	"github.com/GoreeCloud/goreecloud-messenger/internal/runtimeconfig"
)

func main() {
	receiptPersistence, err := runtimeconfig.ReceiptPersistenceFromEnvironment()
	if err != nil {
		panic(err)
	}
	receiptDiagnostic, err := runtimeconfig.ReceiptPersistenceDiagnosticFor(receiptPersistence)
	if err != nil {
		panic(err)
	}
	typingPrivacyPersistence, err := runtimeconfig.TypingPrivacyPersistenceFromEnvironment()
	if err != nil {
		panic(err)
	}
	typingPrivacyDurability, err := runtimeconfig.TypingPrivacyDurabilityStatusFor(typingPrivacyPersistence.Mode)
	if err != nil {
		panic(err)
	}

	message := domain.Message{
		ID:             "development-message",
		ConversationID: "development-conversation",
		SenderID:       "development-user",
		Body:           "GoreeCloud Messenger foundation",
		Transport:      domain.TransportData,
		Encryption:     domain.EncryptionE2EE,
		Delivery:       domain.DeliverySent,
		SentAt:         time.Now().UTC(),
	}

	if err := message.Validate(); err != nil {
		panic(err)
	}

	fmt.Printf(
		"Messenger development contract active: %s %s %s\n",
		message.ProvenanceLabel(),
		receiptDiagnostic.LogLine(),
		typingPrivacyDurability.LogLine(),
	)
}
