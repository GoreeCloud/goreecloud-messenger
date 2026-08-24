// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

var (
	ErrDuplicateMessage = errors.New("message already exists")
	ErrNonceReuse        = errors.New("client nonce already used")
)

// DataStore is the persistence boundary for GoreeCloud Data envelopes.
type DataStore interface {
	Put(context.Context, domain.DataEnvelope) error
	ListConversation(context.Context, string) ([]domain.DataEnvelope, error)
}

// DataService validates encrypted GoreeCloud Data envelopes before persistence.
type DataService struct {
	store DataStore
}

func NewDataService(store DataStore) (*DataService, error) {
	if store == nil {
		return nil, errors.New("Data store is required")
	}
	return &DataService{store: store}, nil
}

func (s *DataService) Submit(ctx context.Context, envelope domain.DataEnvelope) error {
	if err := envelope.Validate(); err != nil {
		return fmt.Errorf("validate Data envelope: %w", err)
	}
	if err := s.store.Put(ctx, envelope); err != nil {
		return fmt.Errorf("persist Data envelope: %w", err)
	}
	return nil
}

func (s *DataService) ListConversation(ctx context.Context, conversationID string) ([]domain.DataEnvelope, error) {
	if conversationID == "" {
		return nil, errors.New("conversation id is required")
	}
	return s.store.ListConversation(ctx, conversationID)
}

// MemoryDataStore is a deterministic development store. It is not a production persistence implementation.
type MemoryDataStore struct {
	mu          sync.RWMutex
	messages    map[string]domain.DataEnvelope
	nonces      map[string]string
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

func cloneEnvelope(envelope domain.DataEnvelope) domain.DataEnvelope {
	cloned := envelope
	cloned.Ciphertext = append([]byte(nil), envelope.Ciphertext...)
	return cloned
}
