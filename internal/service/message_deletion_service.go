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
	ErrDeletionUserMismatch  = errors.New("authenticated user does not match deletion deleter")
	ErrDeletionNotSender     = errors.New("only the original sender may delete a message for everyone")
	ErrDuplicateDeletion     = errors.New("message deletion already exists")
	ErrDeletionNonceReuse    = errors.New("message deletion client nonce already used")
	ErrDeletionBeforeMessage = errors.New("message deletion timestamp precedes original message")
	ErrMessageAlreadyDeleted = errors.New("message already has a delete-for-everyone tombstone")
)

// MessageDeletionStore persists one immutable delete-for-everyone tombstone per message.
type MessageDeletionStore interface {
	PutDeletion(context.Context, domain.MessageDeletion) error
	GetDeletion(context.Context, string) (domain.MessageDeletion, bool, error)
}

// MessageDeletionService authorizes delete-for-everyone tombstones without reading message plaintext.
type MessageDeletionService struct {
	messages  MessageLookup
	deletions MessageDeletionStore
	access    ConversationAccess
}

func NewMessageDeletionService(messages MessageLookup, deletions MessageDeletionStore, access ConversationAccess) (*MessageDeletionService, error) {
	if messages == nil || deletions == nil || access == nil {
		return nil, errors.New("message lookup, deletion store, and conversation access are required")
	}
	return &MessageDeletionService{messages: messages, deletions: deletions, access: access}, nil
}

func (s *MessageDeletionService) Record(ctx context.Context, authenticatedUserID string, deletion domain.MessageDeletion) error {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return errors.New("authenticated user id is required")
	}
	if err := deletion.Validate(); err != nil {
		return fmt.Errorf("validate message deletion: %w", err)
	}
	if authenticatedUserID != deletion.DeleterID {
		return ErrDeletionUserMismatch
	}

	message, found, err := s.messages.Get(ctx, deletion.MessageID)
	if err != nil {
		return fmt.Errorf("lookup message: %w", err)
	}
	if !found || message.ConversationID != deletion.ConversationID {
		return ErrMessageNotFound
	}

	allowed, err := s.access.IsParticipant(ctx, deletion.ConversationID, authenticatedUserID)
	if err != nil {
		return fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return ErrConversationAccess
	}
	if message.SenderID != authenticatedUserID {
		return ErrDeletionNotSender
	}
	if deletion.DeletedAt.Before(message.CreatedAt) {
		return ErrDeletionBeforeMessage
	}

	if err := s.deletions.PutDeletion(ctx, deletion); err != nil {
		return fmt.Errorf("persist message deletion: %w", err)
	}
	return nil
}

func (s *MessageDeletionService) List(ctx context.Context, authenticatedUserID, messageID string) ([]domain.MessageDeletion, error) {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return nil, errors.New("authenticated user id is required")
	}
	if strings.TrimSpace(messageID) == "" {
		return nil, errors.New("message id is required")
	}

	message, found, err := s.messages.Get(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("lookup message: %w", err)
	}
	if !found {
		return nil, ErrMessageNotFound
	}
	allowed, err := s.access.IsParticipant(ctx, message.ConversationID, authenticatedUserID)
	if err != nil {
		return nil, fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return nil, ErrConversationAccess
	}

	deletion, found, err := s.deletions.GetDeletion(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("load message deletion: %w", err)
	}
	if !found {
		return []domain.MessageDeletion{}, nil
	}
	return []domain.MessageDeletion{deletion}, nil
}

// MemoryMessageDeletionStore is a deterministic Development store for delete-for-everyone tombstones.
type MemoryMessageDeletionStore struct {
	mu          sync.RWMutex
	deletions   map[string]domain.MessageDeletion
	nonces      map[string]string
	byMessageID map[string]string
}

func NewMemoryMessageDeletionStore() *MemoryMessageDeletionStore {
	return &MemoryMessageDeletionStore{
		deletions:   make(map[string]domain.MessageDeletion),
		nonces:      make(map[string]string),
		byMessageID: make(map[string]string),
	}
}

func (s *MemoryMessageDeletionStore) PutDeletion(_ context.Context, deletion domain.MessageDeletion) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.deletions[deletion.DeletionID]; exists {
		return ErrDuplicateDeletion
	}
	nonceKey := deletion.DeleterID + "\x00" + deletion.ClientNonce
	if existingDeletionID, exists := s.nonces[nonceKey]; exists {
		return fmt.Errorf("%w: nonce belongs to %s", ErrDeletionNonceReuse, existingDeletionID)
	}
	if existingDeletionID, exists := s.byMessageID[deletion.MessageID]; exists {
		return fmt.Errorf("%w: tombstone is %s", ErrMessageAlreadyDeleted, existingDeletionID)
	}

	s.deletions[deletion.DeletionID] = deletion
	s.nonces[nonceKey] = deletion.DeletionID
	s.byMessageID[deletion.MessageID] = deletion.DeletionID
	return nil
}

func (s *MemoryMessageDeletionStore) GetDeletion(_ context.Context, messageID string) (domain.MessageDeletion, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	deletionID, ok := s.byMessageID[messageID]
	if !ok {
		return domain.MessageDeletion{}, false, nil
	}
	return s.deletions[deletionID], true, nil
}
