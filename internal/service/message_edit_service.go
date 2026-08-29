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
	ErrEditUserMismatch  = errors.New("authenticated user does not match edit editor")
	ErrEditNotSender     = errors.New("only the original sender may edit a message")
	ErrDuplicateEdit     = errors.New("message edit already exists")
	ErrEditNonceReuse    = errors.New("message edit client nonce already used")
	ErrEditBeforeMessage = errors.New("message edit timestamp precedes original message")
)

// MessageEditStore persists immutable encrypted edit revisions separately from original envelopes.
type MessageEditStore interface {
	PutEdit(context.Context, domain.MessageEdit) error
	ListEdits(context.Context, string) ([]domain.MessageEdit, error)
}

// MessageEditService authorizes encrypted sent-message revisions without decrypting message content.
type MessageEditService struct {
	messages MessageLookup
	edits    MessageEditStore
	access   ConversationAccess
}

func NewMessageEditService(messages MessageLookup, edits MessageEditStore, access ConversationAccess) (*MessageEditService, error) {
	if messages == nil || edits == nil || access == nil {
		return nil, errors.New("message lookup, edit store, and conversation access are required")
	}
	return &MessageEditService{messages: messages, edits: edits, access: access}, nil
}

func (s *MessageEditService) Record(ctx context.Context, authenticatedUserID string, edit domain.MessageEdit) error {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return errors.New("authenticated user id is required")
	}
	if err := edit.Validate(); err != nil {
		return fmt.Errorf("validate message edit: %w", err)
	}
	if authenticatedUserID != edit.EditorID {
		return ErrEditUserMismatch
	}

	message, found, err := s.messages.Get(ctx, edit.MessageID)
	if err != nil {
		return fmt.Errorf("lookup message: %w", err)
	}
	if !found || message.ConversationID != edit.ConversationID {
		return ErrMessageNotFound
	}

	allowed, err := s.access.IsParticipant(ctx, edit.ConversationID, authenticatedUserID)
	if err != nil {
		return fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return ErrConversationAccess
	}
	if message.SenderID != authenticatedUserID {
		return ErrEditNotSender
	}
	if edit.EditedAt.Before(message.CreatedAt) {
		return ErrEditBeforeMessage
	}

	if err := s.edits.PutEdit(ctx, edit); err != nil {
		return fmt.Errorf("persist message edit: %w", err)
	}
	return nil
}

func (s *MessageEditService) List(ctx context.Context, authenticatedUserID, messageID string) ([]domain.MessageEdit, error) {
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
	return s.edits.ListEdits(ctx, messageID)
}

// MemoryMessageEditStore is a deterministic Development store for encrypted edit revisions.
type MemoryMessageEditStore struct {
	mu        sync.RWMutex
	edits     map[string]domain.MessageEdit
	nonces    map[string]string
	byMessage map[string][]string
}

func NewMemoryMessageEditStore() *MemoryMessageEditStore {
	return &MemoryMessageEditStore{
		edits:     make(map[string]domain.MessageEdit),
		nonces:    make(map[string]string),
		byMessage: make(map[string][]string),
	}
}

func (s *MemoryMessageEditStore) PutEdit(_ context.Context, edit domain.MessageEdit) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.edits[edit.EditID]; exists {
		return ErrDuplicateEdit
	}
	nonceKey := edit.EditorID + "\x00" + edit.ClientNonce
	if existingEditID, exists := s.nonces[nonceKey]; exists {
		return fmt.Errorf("%w: nonce belongs to %s", ErrEditNonceReuse, existingEditID)
	}

	s.edits[edit.EditID] = cloneMessageEdit(edit)
	s.nonces[nonceKey] = edit.EditID
	s.byMessage[edit.MessageID] = append(s.byMessage[edit.MessageID], edit.EditID)
	return nil
}

func (s *MemoryMessageEditStore) ListEdits(_ context.Context, messageID string) ([]domain.MessageEdit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.byMessage[messageID]
	result := make([]domain.MessageEdit, 0, len(ids))
	for _, id := range ids {
		result = append(result, cloneMessageEdit(s.edits[id]))
	}
	return result, nil
}

func cloneMessageEdit(edit domain.MessageEdit) domain.MessageEdit {
	cloned := edit
	cloned.Ciphertext = append([]byte(nil), edit.Ciphertext...)
	return cloned
}
