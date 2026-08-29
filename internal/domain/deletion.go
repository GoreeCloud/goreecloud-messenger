// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"strings"
	"time"
)

// MessageDeletionScope identifies the semantic scope of a Messenger deletion control event.
type MessageDeletionScope string

const (
	MessageDeletionEveryone MessageDeletionScope = "everyone"
)

// MessageDeletion is an immutable delete-for-everyone tombstone for an existing GoreeCloud Data message.
// It carries no plaintext or ciphertext. A tombstone instructs authorized clients to suppress the
// message; it is not proof that every remote device, backup, or retained encrypted object was erased.
type MessageDeletion struct {
	DeletionID     string
	MessageID      string
	ConversationID string
	DeleterID      string
	ClientNonce    string
	Scope          MessageDeletionScope
	DeletedAt      time.Time
}

func (d MessageDeletion) Validate() error {
	if strings.TrimSpace(d.DeletionID) == "" {
		return errors.New("deletion id is required")
	}
	if strings.TrimSpace(d.MessageID) == "" {
		return errors.New("deletion message id is required")
	}
	if strings.TrimSpace(d.ConversationID) == "" {
		return errors.New("deletion conversation id is required")
	}
	if strings.TrimSpace(d.DeleterID) == "" {
		return errors.New("deletion deleter id is required")
	}
	if strings.TrimSpace(d.ClientNonce) == "" {
		return errors.New("deletion client nonce is required")
	}
	if d.Scope != MessageDeletionEveryone {
		return errors.New("unsupported message deletion scope")
	}
	if d.DeletedAt.IsZero() {
		return errors.New("deletion timestamp is required")
	}
	return nil
}
