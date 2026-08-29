// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

var (
	ErrReactionUserMismatch      = errors.New("authenticated user does not match reaction reactor")
	ErrReactionTargetUnavailable = errors.New("reaction target unavailable")
	ErrDuplicateReaction         = errors.New("message reaction event already exists")
	ErrReactionNonceReuse        = errors.New("message reaction client nonce already used")
	ErrReactionBeforeMessage     = errors.New("reaction timestamp precedes original message")
	ErrReactionStale             = errors.New("reaction event does not advance current state")
)

// MessageReactionStore persists immutable reaction events and returns the current active per-user projection.
type MessageReactionStore interface {
	PutReaction(context.Context, domain.MessageReaction) error
	ListCurrentReactions(context.Context, string) ([]domain.MessageReaction, error)
}

// MessageReactionService authorizes opaque encrypted reactions without decrypting reaction values or message content.
type MessageReactionService struct {
	messages  MessageLookup
	reactions MessageReactionStore
	access    ConversationAccess
}

func NewMessageReactionService(messages MessageLookup, reactions MessageReactionStore, access ConversationAccess) (*MessageReactionService, error) {
	if messages == nil || reactions == nil || access == nil {
		return nil, errors.New("message lookup, reaction store, and conversation access are required")
	}
	return &MessageReactionService{messages: messages, reactions: reactions, access: access}, nil
}

func (s *MessageReactionService) Record(ctx context.Context, authenticatedUserID string, reaction domain.MessageReaction) error {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return errors.New("authenticated user id is required")
	}
	if err := reaction.Validate(); err != nil {
		return fmt.Errorf("validate message reaction: %w", err)
	}
	if authenticatedUserID != reaction.ReactorID {
		return ErrReactionUserMismatch
	}

	allowed, err := s.access.IsParticipant(ctx, reaction.ConversationID, authenticatedUserID)
	if err != nil {
		return fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return ErrConversationAccess
	}

	message, found, err := s.messages.Get(ctx, reaction.MessageID)
	if err != nil {
		return fmt.Errorf("lookup reaction target: %w", err)
	}
	if !found || message.ConversationID != reaction.ConversationID {
		return ErrReactionTargetUnavailable
	}
	if reaction.ReactedAt.Before(message.CreatedAt) {
		return ErrReactionBeforeMessage
	}

	if err := s.reactions.PutReaction(ctx, reaction); err != nil {
		return fmt.Errorf("persist message reaction: %w", err)
	}
	return nil
}

func (s *MessageReactionService) List(ctx context.Context, authenticatedUserID, messageID string) ([]domain.MessageReaction, error) {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return nil, errors.New("authenticated user id is required")
	}
	messageID = strings.TrimSpace(messageID)
	if messageID == "" {
		return nil, errors.New("message id is required")
	}

	message, found, err := s.messages.Get(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("lookup reaction target: %w", err)
	}
	if !found {
		return nil, ErrReactionTargetUnavailable
	}
	allowed, err := s.access.IsParticipant(ctx, message.ConversationID, authenticatedUserID)
	if err != nil {
		return nil, fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return nil, ErrConversationAccess
	}
	return s.reactions.ListCurrentReactions(ctx, messageID)
}

// MemoryMessageReactionStore is a deterministic Development store for encrypted reaction events.
type MemoryMessageReactionStore struct {
	mu     sync.RWMutex
	events map[string]domain.MessageReaction
	nonces map[string]string
	latest map[string]domain.MessageReaction
}

func NewMemoryMessageReactionStore() *MemoryMessageReactionStore {
	return &MemoryMessageReactionStore{
		events: make(map[string]domain.MessageReaction),
		nonces: make(map[string]string),
		latest: make(map[string]domain.MessageReaction),
	}
}

func (s *MemoryMessageReactionStore) PutReaction(_ context.Context, reaction domain.MessageReaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[reaction.ReactionID]; exists {
		return ErrDuplicateReaction
	}
	nonceKey := reaction.ReactorID + "\x00" + reaction.ClientNonce
	if existingReactionID, exists := s.nonces[nonceKey]; exists {
		return fmt.Errorf("%w: nonce belongs to %s", ErrReactionNonceReuse, existingReactionID)
	}
	stateKey := reaction.MessageID + "\x00" + reaction.ReactorID
	if existing, exists := s.latest[stateKey]; exists && !reaction.ReactedAt.After(existing.ReactedAt) {
		return ErrReactionStale
	}

	stored := cloneMessageReaction(reaction)
	s.events[reaction.ReactionID] = stored
	s.nonces[nonceKey] = reaction.ReactionID
	s.latest[stateKey] = stored
	return nil
}

func (s *MemoryMessageReactionStore) ListCurrentReactions(_ context.Context, messageID string) ([]domain.MessageReaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.MessageReaction, 0)
	for key, reaction := range s.latest {
		if !strings.HasPrefix(key, messageID+"\x00") || reaction.Operation != domain.ReactionSet {
			continue
		}
		result = append(result, cloneMessageReaction(reaction))
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].ReactedAt.Equal(result[j].ReactedAt) {
			return result[i].ReactorID < result[j].ReactorID
		}
		return result[i].ReactedAt.Before(result[j].ReactedAt)
	})
	return result, nil
}

func cloneMessageReaction(reaction domain.MessageReaction) domain.MessageReaction {
	cloned := reaction
	cloned.Ciphertext = append([]byte(nil), reaction.Ciphertext...)
	return cloned
}
