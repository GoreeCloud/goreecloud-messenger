// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

var (
	ErrDuplicateMessage       = errors.New("message already exists")
	ErrNonceReuse             = errors.New("client nonce already used")
	ErrSenderMismatch         = errors.New("authenticated user does not match envelope sender")
	ErrConversationAccess     = errors.New("user is not a conversation participant")
	ErrReplyTargetUnavailable = errors.New("reply target unavailable")
)

// DataStore is the persistence boundary for GoreeCloud Data envelopes.
type DataStore interface {
	Put(context.Context, domain.DataEnvelope) error
	Get(context.Context, string) (domain.DataEnvelope, bool, error)
	ListConversation(context.Context, string) ([]domain.DataEnvelope, error)
}

// ConversationAccess verifies membership without trusting client-supplied authorization state.
type ConversationAccess interface {
	IsParticipant(context.Context, string, string) (bool, error)
}

// DataService validates encrypted GoreeCloud Data envelopes and authorization before persistence.
type DataService struct {
	store  DataStore
	access ConversationAccess
}

func NewDataService(store DataStore, access ConversationAccess) (*DataService, error) {
	if store == nil {
		return nil, errors.New("Data store is required")
	}
	if access == nil {
		return nil, errors.New("conversation access verifier is required")
	}
	return &DataService{store: store, access: access}, nil
}

func (s *DataService) Submit(ctx context.Context, authenticatedUserID string, envelope domain.DataEnvelope) error {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return errors.New("authenticated user id is required")
	}
	if err := envelope.Validate(); err != nil {
		return fmt.Errorf("validate Data envelope: %w", err)
	}
	if authenticatedUserID != envelope.SenderID {
		return ErrSenderMismatch
	}

	allowed, err := s.access.IsParticipant(ctx, envelope.ConversationID, authenticatedUserID)
	if err != nil {
		return fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return ErrConversationAccess
	}

	envelope.ReplyToMessageID = strings.TrimSpace(envelope.ReplyToMessageID)
	if envelope.ReplyToMessageID != "" {
		target, found, err := s.store.Get(ctx, envelope.ReplyToMessageID)
		if err != nil {
			return fmt.Errorf("load reply target: %w", err)
		}
		if !found || target.ConversationID != envelope.ConversationID || target.MessageID == envelope.MessageID {
			return ErrReplyTargetUnavailable
		}
	}

	if err := s.store.Put(ctx, envelope); err != nil {
		return fmt.Errorf("persist Data envelope: %w", err)
	}
	return nil
}

func (s *DataService) ListConversation(ctx context.Context, authenticatedUserID, conversationID string) ([]domain.DataEnvelope, error) {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return nil, errors.New("authenticated user id is required")
	}
	if strings.TrimSpace(conversationID) == "" {
		return nil, errors.New("conversation id is required")
	}

	allowed, err := s.access.IsParticipant(ctx, conversationID, authenticatedUserID)
	if err != nil {
		return nil, fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return nil, ErrConversationAccess
	}
	return s.store.ListConversation(ctx, conversationID)
}

// MemoryDataStore is a deterministic development store. It is not a production persistence implementation.
type MemoryDataStore struct {
	mu           sync.RWMutex
	messages     map[string]domain.DataEnvelope
	nonces       map[string]string
	conversation map[string][]string
}

func NewMemoryDataStore() *MemoryDataStore {
	return &MemoryDataStore{
		messages:     make(map[string]domain.DataEnvelope),
		nonces:       make(map[string]string),
		conversation: make(map[string][]string),
	}
}

func (s *MemoryDataStore) Put(_ context.Context, envelope domain.DataEnvelope) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.messages[envelope.MessageID]; exists {
		return ErrDuplicateMessage
	}
	if existingMessageID, exists := s.nonces[envelope.ClientNonce]; exists {
		return fmt.Errorf("%w: nonce belongs to %s", ErrNonceReuse, existingMessageID)
	}

	s.messages[envelope.MessageID] = cloneEnvelope(envelope)
	s.nonces[envelope.ClientNonce] = envelope.MessageID
	s.conversation[envelope.ConversationID] = append(s.conversation[envelope.ConversationID], envelope.MessageID)
	return nil
}

func (s *MemoryDataStore) Get(_ context.Context, messageID string) (domain.DataEnvelope, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	envelope, ok := s.messages[messageID]
	if !ok {
		return domain.DataEnvelope{}, false, nil
	}
	return cloneEnvelope(envelope), true, nil
}

func (s *MemoryDataStore) ListConversation(_ context.Context, conversationID string) ([]domain.DataEnvelope, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.conversation[conversationID]
	result := make([]domain.DataEnvelope, 0, len(ids))
	for _, id := range ids {
		result = append(result, cloneEnvelope(s.messages[id]))
	}
	return result, nil
}

// MemoryConversationAccess is a deterministic development membership verifier.
type MemoryConversationAccess struct {
	mu      sync.RWMutex
	members map[string]map[string]struct{}
}

func NewMemoryConversationAccess() *MemoryConversationAccess {
	return &MemoryConversationAccess{members: make(map[string]map[string]struct{})}
}

func (a *MemoryConversationAccess) SetConversation(conversation domain.Conversation) error {
	if err := conversation.Validate(); err != nil {
		return fmt.Errorf("validate conversation: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	members := make(map[string]struct{}, len(conversation.ParticipantIDs))
	for _, participantID := range conversation.ParticipantIDs {
		members[participantID] = struct{}{}
	}
	a.members[conversation.ID] = members
	return nil
}

func (a *MemoryConversationAccess) IsParticipant(_ context.Context, conversationID, userID string) (bool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	members, exists := a.members[conversationID]
	if !exists {
		return false, nil
	}
	_, allowed := members[userID]
	return allowed, nil
}

func cloneEnvelope(envelope domain.DataEnvelope) domain.DataEnvelope {
	cloned := envelope
	cloned.Ciphertext = append([]byte(nil), envelope.Ciphertext...)
	return cloned
}
