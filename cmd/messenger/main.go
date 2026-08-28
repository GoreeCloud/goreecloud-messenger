// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	"github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

const (
	attachmentStoreModeEnv = "GOREECLOUD_MESSENGER_ATTACHMENT_STORE"
	attachmentStoreRootEnv = "GOREECLOUD_MESSENGER_ATTACHMENT_ROOT"
)

func main() {
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

	mode := service.AttachmentStoreMode(os.Getenv(attachmentStoreModeEnv))
	if _, err := service.NewAttachmentStore(mode, os.Getenv(attachmentStoreRootEnv)); err != nil {
		panic(err)
	}
	modeLabel := string(mode)
	if modeLabel == "" {
		modeLabel = string(service.AttachmentStoreMemory)
	}

	fmt.Printf(
		"Messenger development contract active: %s; attachment-store=%s\n",
		message.ProvenanceLabel(),
		modeLabel,
	)
}
