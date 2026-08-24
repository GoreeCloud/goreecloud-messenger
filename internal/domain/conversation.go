// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"fmt"
	"strings"
)

// ConversationKind distinguishes direct and group messaging.
type ConversationKind string

const (
	ConversationDirect ConversationKind = "direct"
	ConversationGroup  ConversationKind = "group"
)

func (k ConversationKind) Valid() bool {
	switch k {
	case ConversationDirect, ConversationGroup:
		return true
	default:
		return false
	}
}

// Conversation is transport-neutral. Individual messages retain their own transport provenance.
type Conversation struct {
	ID             string
	Kind           ConversationKind
	ParticipantIDs []string
}

func (c Conversation) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return errors.New("conversation id is required")
	}
	if !c.Kind.Valid() {
		return fmt.Errorf("unsupported conversation kind %q", c.Kind)
	}

	minimumParticipants := 2
	if len(c.ParticipantIDs) < minimumParticipants {
		return errors.New("conversation requires at least two participants")
	}

	seen := make(map[string]struct{}, len(c.ParticipantIDs))
	for _, participantID := range c.ParticipantIDs {
		participantID = strings.TrimSpace(participantID)
		if participantID == "" {
			return errors.New("conversation participant id is required")
		}
		if _, exists := seen[participantID]; exists {
			return errors.New("conversation participant ids must be unique")
		}
		seen[participantID] = struct{}{}
	}

	return nil
}

// TransportTransition records a user-visible move from one transport to another.
type TransportTransition struct {
	From          Transport
	To            Transport
	UserConfirmed bool
}

func (t TransportTransition) Validate() error {
	if !t.From.Valid() || !t.To.Valid() {
		return errors.New("transport transition requires valid source and destination transports")
	}
	if t.From == t.To {
		return errors.New("transport transition requires different source and destination transports")
	}
	if t.From == TransportData && t.To != TransportData && !t.UserConfirmed {
		return errors.New("leaving GoreeCloud Data transport requires explicit user confirmation")
	}
	return nil
}
