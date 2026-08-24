// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"fmt"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func main() {
	message := domain.Message{
		ID:             "development-message",
		ConversationID: "development-conversation",
		SenderID:       "development-user",
		Body:           "GoreeCloud Messenger foundation",
		Transport:      domain.TransportData,
		Encryption:     domain.EncryptionE2EE,
		SentAt:         time.Now().UTC(),
	}

	if err := message.Validate(); err != nil {
		panic(err)
	}

	fmt.Printf("Messenger development contract active: %s\n", message.ProvenanceLabel())
}
