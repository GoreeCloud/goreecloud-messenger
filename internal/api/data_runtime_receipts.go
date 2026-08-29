// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"errors"
	"fmt"

	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

// ReceiptPersistenceMode identifies the receipt persistence implementation selected at
// the Data runtime composition boundary. Selection is explicit so a process never silently
// changes receipt durability because of an implicit environment fallback.
type ReceiptPersistenceMode string

const (
	ReceiptPersistenceMemory ReceiptPersistenceMode = "memory"
	ReceiptPersistenceFile   ReceiptPersistenceMode = "file"
)

// ReceiptPersistenceConfig configures only the receipt-store boundary. File mode is the
// current durable single-node development option and does not imply replicated production
// persistence, backup guarantees, or multi-writer safety.
type ReceiptPersistenceConfig struct {
	Mode ReceiptPersistenceMode
	Root string
}

// DataRuntimeDependencies are the already-authoritative service boundaries required to
// compose the current Data HTTP runtime while selecting receipt persistence here.
type DataRuntimeDependencies struct {
	Data               *messagingservice.DataService
	MessageLookup      messagingservice.MessageLookup
	ConversationAccess messagingservice.ConversationAccess
	Attachments        *messagingservice.AttachmentService
	Auth               Authenticator
}

// NewConfiguredDataRuntimeHandler composes the existing runtime with an explicitly selected
// receipt store. It closes the prior gap where FileReceiptStore existed but the application
// composition layer could only receive a pre-built ReceiptService from an external caller.
func NewConfiguredDataRuntimeHandler(
	dependencies DataRuntimeDependencies,
	config ReceiptPersistenceConfig,
) (*DataRuntimeHandler, error) {
	if dependencies.Data == nil || dependencies.MessageLookup == nil || dependencies.ConversationAccess == nil || dependencies.Attachments == nil || dependencies.Auth == nil {
		return nil, errors.New("Data, message lookup, conversation access, attachment, and authentication boundaries are required")
	}

	store, err := newRuntimeReceiptStore(config)
	if err != nil {
		return nil, err
	}
	receipts, err := messagingservice.NewReceiptService(
		dependencies.MessageLookup,
		store,
		dependencies.ConversationAccess,
	)
	if err != nil {
		return nil, fmt.Errorf("compose receipt service: %w", err)
	}
	return NewDataRuntimeHandler(
		dependencies.Data,
		receipts,
		dependencies.Attachments,
		dependencies.Auth,
	)
}

func newRuntimeReceiptStore(config ReceiptPersistenceConfig) (messagingservice.ReceiptStore, error) {
	switch config.Mode {
	case ReceiptPersistenceMemory:
		return messagingservice.NewMemoryReceiptStore(), nil
	case ReceiptPersistenceFile:
		store, err := messagingservice.NewFileReceiptStore(config.Root)
		if err != nil {
			return nil, fmt.Errorf("open durable receipt store: %w", err)
		}
		return store, nil
	case "":
		return nil, errors.New("receipt persistence mode is required")
	default:
		return nil, fmt.Errorf("unsupported receipt persistence mode %q", config.Mode)
	}
}
