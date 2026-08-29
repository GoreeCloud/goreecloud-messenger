// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"strings"
	"time"
)

type ReactionOperation string

const (
	ReactionSet   ReactionOperation = "set"
	ReactionClear ReactionOperation = "clear"
)

func (o ReactionOperation) Valid() bool {
	return o == ReactionSet || o == ReactionClear
}

// MessageReaction is an immutable encrypted reaction-state event for one GoreeCloud Data message.
// A set event carries opaque reaction ciphertext. A clear event intentionally carries no reaction value.
type MessageReaction struct {
	ReactionID     string
	MessageID      string
	ConversationID string
	ReactorID      string
	ClientNonce    string
	Operation      ReactionOperation
	Ciphertext     []byte
	Encryption     EncryptionState
	ReactedAt      time.Time
}

func (r MessageReaction) Validate() error {
	if strings.TrimSpace(r.ReactionID) == "" {
		return errors.New("reaction id is required")
	}
	if strings.TrimSpace(r.MessageID) == "" {
		return errors.New("reaction message id is required")
	}
	if strings.TrimSpace(r.ConversationID) == "" {
		return errors.New("reaction conversation id is required")
	}
	if strings.TrimSpace(r.ReactorID) == "" {
		return errors.New("reaction reactor id is required")
	}
	if strings.TrimSpace(r.ClientNonce) == "" {
		return errors.New("reaction client nonce is required")
	}
	if !r.Operation.Valid() {
		return errors.New("reaction operation is invalid")
	}
	if r.Operation == ReactionSet && len(r.Ciphertext) == 0 {
		return errors.New("set reaction ciphertext is required")
	}
	if r.Operation == ReactionClear && len(r.Ciphertext) != 0 {
		return errors.New("clear reaction must not carry ciphertext")
	}
	if r.Encryption != EncryptionE2EE {
		return errors.New("GoreeCloud Data reaction requires active E2EE")
	}
	if r.ReactedAt.IsZero() {
		return errors.New("reaction timestamp is required")
	}
	return nil
}
