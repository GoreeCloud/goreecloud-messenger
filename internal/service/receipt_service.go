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
	ErrReceiptUserMismatch = errors.New("authenticated user does not match receipt user")
	ErrReceiptRegression   = errors.New("receipt state cannot move backwards")
	ErrSelfReceipt         = errors.New("sender cannot acknowledge own message")
	ErrMessageNotFound     = errors.New("message not found")
)

type MessageLookup interface {
	Get(context.Context, string) (domain.DataEnvelope, bool, error)
}

type ReceiptStore interface {
	PutReceipt(context.Context, domain.DeliveryReceipt) error
	ListReceipts(context.Context, string) ([]domain.DeliveryReceipt, error)
}

type ReceiptService struct {
	messages MessageLookup
	receipts ReceiptStore
	access   ConversationAccess
}

func NewReceiptService(messages MessageLookup, receipts ReceiptStore, access ConversationAccess) (*ReceiptService, error) {
	if messages == nil || receipts == nil || access == nil {
		return nil, errors.New("message lookup, receipt store, and conversation access are required")
	}
	return &ReceiptService{messages: messages, receipts: receipts, access: access}, nil
}

func (s *ReceiptService) Record(ctx context.Context, authenticatedUserID string, receipt domain.DeliveryReceipt) error {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return errors.New("authenticated user id is required")
	}
	if err := receipt.Validate(); err != nil {
		return fmt.Errorf("validate receipt: %w", err)
	}
	if authenticatedUserID != receipt.UserID {
		return ErrReceiptUserMismatch
	}
	allowed, err := s.access.IsParticipant(ctx, receipt.ConversationID, authenticatedUserID)
	if err != nil {
		return fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return ErrConversationAccess
	}
	message, found, err := s.messages.Get(ctx, receipt.MessageID)
	if err != nil {
		return fmt.Errorf("lookup message: %w", err)
	}
	if !found || message.ConversationID != receipt.ConversationID {
		return ErrMessageNotFound
	}
	if message.SenderID == authenticatedUserID {
		return ErrSelfReceipt
	}
	return s.receipts.PutReceipt(ctx, receipt)
}

func (s *ReceiptService) List(ctx context.Context, authenticatedUserID, messageID string) ([]domain.DeliveryReceipt, error) {
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
	return s.receipts.ListReceipts(ctx, messageID)
}

type MemoryReceiptStore struct {
	mu       sync.RWMutex
	receipts map[string]map[string]domain.DeliveryReceipt
}

func NewMemoryReceiptStore() *MemoryReceiptStore {
	return &MemoryReceiptStore{receipts: make(map[string]map[string]domain.DeliveryReceipt)}
}

func (s *MemoryReceiptStore) PutReceipt(_ context.Context, receipt domain.DeliveryReceipt) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	byUser := s.receipts[receipt.MessageID]
	if byUser == nil {
		byUser = make(map[string]domain.DeliveryReceipt)
		s.receipts[receipt.MessageID] = byUser
	}
	if existing, ok := byUser[receipt.UserID]; ok && receipt.State.Rank() < existing.State.Rank() {
		return ErrReceiptRegression
	}
	byUser[receipt.UserID] = receipt
	return nil
}

func (s *MemoryReceiptStore) ListReceipts(_ context.Context, messageID string) ([]domain.DeliveryReceipt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	byUser := s.receipts[messageID]
	result := make([]domain.DeliveryReceipt, 0, len(byUser))
	for _, receipt := range byUser {
		result = append(result, receipt)
	}
	return result, nil
}
