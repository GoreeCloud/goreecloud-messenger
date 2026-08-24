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
	ErrDuplicateAttachment  = errors.New("attachment already exists")
	ErrAttachmentNonceReuse = errors.New("attachment client nonce already used")
)

// AttachmentStore is the persistence boundary for encrypted GoreeCloud Data attachments.
type AttachmentStore interface {
	PutAttachment(context.Context, domain.DataAttachment) error
	GetAttachment(context.Context, string) (domain.DataAttachment, bool, error)
}

// AttachmentService validates encrypted Data attachments and conversation authorization.
type AttachmentService struct {
	store  AttachmentStore
	access ConversationAccess
}

func NewAttachmentService(store AttachmentStore, access ConversationAccess) (*AttachmentService, error) {
	if store == nil {
		return nil, errors.New("attachment store is required")
	}
	if access == nil {
		return nil, errors.New("conversation access verifier is required")
	}
	return &AttachmentService{store: store, access: access}, nil
}

func (s *AttachmentService) Submit(ctx context.Context, authenticatedUserID string, attachment domain.DataAttachment) error {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return errors.New("authenticated user id is required")
	}
	if err := attachment.Validate(); err != nil {
		return fmt.Errorf("validate Data attachment: %w", err)
	}
	if authenticatedUserID != attachment.SenderID {
		return ErrSenderMismatch
	}

	allowed, err := s.access.IsParticipant(ctx, attachment.ConversationID, authenticatedUserID)
	if err != nil {
		return fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return ErrConversationAccess
	}

	if err := s.store.PutAttachment(ctx, attachment); err != nil {
		return fmt.Errorf("persist Data attachment: %w", err)
	}
	return nil
}

func (s *AttachmentService) Get(ctx context.Context, authenticatedUserID, attachmentID string) (domain.DataAttachment, error) {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return domain.DataAttachment{}, errors.New("authenticated user id is required")
	}
	if strings.TrimSpace(attachmentID) == "" {
		return domain.DataAttachment{}, errors.New("attachment id is required")
	}

	attachment, ok, err := s.store.GetAttachment(ctx, attachmentID)
	if err != nil {
		return domain.DataAttachment{}, fmt.Errorf("load Data attachment: %w", err)
	}
	if !ok {
		return domain.DataAttachment{}, errors.New("attachment not found")
	}

	allowed, err := s.access.IsParticipant(ctx, attachment.ConversationID, authenticatedUserID)
	if err != nil {
		return domain.DataAttachment{}, fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return domain.DataAttachment{}, ErrConversationAccess
	}
	return cloneAttachment(attachment), nil
}

// MemoryAttachmentStore is a deterministic development store. It is not production persistence.
type MemoryAttachmentStore struct {
	mu          sync.RWMutex
	attachments map[string]domain.DataAttachment
	nonces      map[string]string
}

func NewMemoryAttachmentStore() *MemoryAttachmentStore {
	return &MemoryAttachmentStore{
		attachments: make(map[string]domain.DataAttachment),
		nonces:      make(map[string]string),
	}
}

func (s *MemoryAttachmentStore) PutAttachment(_ context.Context, attachment domain.DataAttachment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.attachments[attachment.AttachmentID]; exists {
		return ErrDuplicateAttachment
	}
	if existingAttachmentID, exists := s.nonces[attachment.ClientNonce]; exists {
		return fmt.Errorf("%w: nonce belongs to %s", ErrAttachmentNonceReuse, existingAttachmentID)
	}

	s.attachments[attachment.AttachmentID] = cloneAttachment(attachment)
	s.nonces[attachment.ClientNonce] = attachment.AttachmentID
	return nil
}

func (s *MemoryAttachmentStore) GetAttachment(_ context.Context, attachmentID string) (domain.DataAttachment, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	attachment, ok := s.attachments[attachmentID]
	if !ok {
		return domain.DataAttachment{}, false, nil
	}
	return cloneAttachment(attachment), true, nil
}

func cloneAttachment(attachment domain.DataAttachment) domain.DataAttachment {
	cloned := attachment
	cloned.Ciphertext = append([]byte(nil), attachment.Ciphertext...)
	return cloned
}
