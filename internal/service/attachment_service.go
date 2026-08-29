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

const MaxAttachmentListResults = 100

var (
	ErrDuplicateAttachment  = errors.New("attachment already exists")
	ErrAttachmentNonceReuse = errors.New("attachment client nonce already used")
	ErrAttachmentNotFound   = errors.New("attachment not found")
	ErrAttachmentListLimit  = errors.New("attachment list limit is invalid")
)

// AttachmentStore is the persistence boundary for encrypted GoreeCloud Data attachments.
// Deletion removes retrievable ciphertext while retaining a minimal replay-prevention tombstone.
type AttachmentStore interface {
	PutAttachment(context.Context, domain.DataAttachment) error
	GetAttachment(context.Context, string) (domain.DataAttachment, bool, error)
	ListAttachmentMetadata(context.Context, string, int) ([]domain.DataAttachmentMetadata, error)
	DeleteAttachment(context.Context, string) (bool, error)
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
		return domain.DataAttachment{}, ErrAttachmentNotFound
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

// Delete removes retrievable ciphertext for an attachment after authenticating conversation access.
// Missing attachments are treated as already deleted so callers can safely retry without learning
// whether a prior request completed. Stores keep only the minimum tombstone state required to prevent
// attachment-id and client-nonce replay.
func (s *AttachmentService) Delete(ctx context.Context, authenticatedUserID, attachmentID string) error {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return errors.New("authenticated user id is required")
	}
	if strings.TrimSpace(attachmentID) == "" {
		return errors.New("attachment id is required")
	}

	if _, err := s.Get(ctx, authenticatedUserID, attachmentID); err != nil {
		if errors.Is(err, ErrAttachmentNotFound) {
			return nil
		}
		return err
	}
	if _, err := s.store.DeleteAttachment(ctx, attachmentID); err != nil {
		return fmt.Errorf("delete Data attachment: %w", err)
	}
	return nil
}

// List returns a bounded metadata-only projection for one conversation. Authorization is checked
// before the store is queried so callers cannot use the listing boundary to probe conversation state.
func (s *AttachmentService) List(
	ctx context.Context,
	authenticatedUserID string,
	conversationID string,
	limit int,
) ([]domain.DataAttachmentMetadata, error) {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return nil, errors.New("authenticated user id is required")
	}
	if strings.TrimSpace(conversationID) == "" {
		return nil, errors.New("conversation id is required")
	}
	if limit < 1 || limit > MaxAttachmentListResults {
		return nil, ErrAttachmentListLimit
	}

	allowed, err := s.access.IsParticipant(ctx, conversationID, authenticatedUserID)
	if err != nil {
		return nil, fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return nil, ErrConversationAccess
	}

	metadata, err := s.store.ListAttachmentMetadata(ctx, conversationID, limit)
	if err != nil {
		return nil, fmt.Errorf("list Data attachment metadata: %w", err)
	}
	return append([]domain.DataAttachmentMetadata(nil), metadata...), nil
}

// MemoryAttachmentStore is a deterministic development store. It is not production persistence.
type MemoryAttachmentStore struct {
	mu          sync.RWMutex
	attachments map[string]domain.DataAttachment
	nonces      map[string]string
	tombstones  map[string]struct{}
}

func NewMemoryAttachmentStore() *MemoryAttachmentStore {
	return &MemoryAttachmentStore{
		attachments: make(map[string]domain.DataAttachment),
		nonces:      make(map[string]string),
		tombstones:  make(map[string]struct{}),
	}
}

func (s *MemoryAttachmentStore) PutAttachment(_ context.Context, attachment domain.DataAttachment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.attachments[attachment.AttachmentID]; exists {
		return ErrDuplicateAttachment
	}
	if _, exists := s.tombstones[attachment.AttachmentID]; exists {
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

func (s *MemoryAttachmentStore) DeleteAttachment(_ context.Context, attachmentID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, deleted := s.tombstones[attachmentID]; deleted {
		return false, nil
	}
	if _, exists := s.attachments[attachmentID]; !exists {
		return false, nil
	}
	delete(s.attachments, attachmentID)
	s.tombstones[attachmentID] = struct{}{}
	// Intentionally retain the nonce reservation. Deletion must not reopen replay state.
	return true, nil
}

func (s *MemoryAttachmentStore) ListAttachmentMetadata(
	_ context.Context,
	conversationID string,
	limit int,
) ([]domain.DataAttachmentMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metadata := make([]domain.DataAttachmentMetadata, 0, limit)
	for _, attachment := range s.attachments {
		if attachment.ConversationID == conversationID {
			metadata = append(metadata, attachment.Metadata())
		}
	}
	sort.Slice(metadata, func(i, j int) bool {
		return metadata[i].AttachmentID < metadata[j].AttachmentID
	})
	if len(metadata) > limit {
		metadata = metadata[:limit]
	}
	return metadata, nil
}

func cloneAttachment(attachment domain.DataAttachment) domain.DataAttachment {
	cloned := attachment
	cloned.Ciphertext = append([]byte(nil), attachment.Ciphertext...)
	return cloned
}
