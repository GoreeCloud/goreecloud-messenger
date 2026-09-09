// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"time"
)

// TypingState is ephemeral conversation-presence metadata and never message content.
type TypingState string

const (
	TypingStateTyping TypingState = "typing"
	TypingStateIdle   TypingState = "idle"
)

// TypingSignal is an authenticated, content-free typing-state update.
type TypingSignal struct {
	ConversationID string
	UserID         string
	Sequence       uint64
	State          TypingState
}

func (s TypingSignal) Validate() error {
	if err := ValidateOpaqueIdentifier(s.ConversationID, "typing conversation id"); err != nil {
		return err
	}
	if err := ValidateOpaqueIdentifier(s.UserID, "typing user id"); err != nil {
		return err
	}
	if s.Sequence == 0 {
		return errors.New("typing sequence must be greater than zero")
	}
	switch s.State {
	case TypingStateTyping, TypingStateIdle:
		return nil
	default:
		return errors.New("unsupported typing state")
	}
}

// ActiveTyping is the short-lived server projection visible to authorized observers.
type ActiveTyping struct {
	ConversationID string
	UserID         string
	Sequence       uint64
	ExpiresAt      time.Time
}
