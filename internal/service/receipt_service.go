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
	ErrReceiptUserMismatch  = errors.New("authenticated user does not match receipt user")
	ErrReceiptRegression    = errors.New("receipt state cannot move backwards")
	ErrReceiptPrivacyDenied = errors.New("read receipt privacy policy denied")
	ErrSelfReceipt          = errors.New("sender cannot acknowledge own message")
	ErrMessageNotFound      = errors.New("message not found")
)

type MessageLookup interface {
	Get(context.Context, string) (domain.DataEnvelope, bool, error)
}

type ReceiptStore interface {
	PutReceipt(context.Context, domain.DeliveryReceipt) error
	ListReceipts(context.Context, string) ([]domain.DeliveryReceipt, error)
}

// ReceiptPrivacyPolicy controls publication and observation of read state.
// Delivered state remains operational delivery metadata and is not gated by this policy.
type ReceiptPrivacyPolicy interface {
	CanPublishRead(context.Context, string, string) (bool, error)
	CanObserveRead(context.Context, string, string) (bool, error)
}

type allowAllReceiptPrivacyPolicy struct{}

func (allowAllReceiptPrivacyPolicy) CanPublishRead(context.Context, string, string) (bool, error) {
	return true, nil
}

func (allowAllReceiptPrivacyPolicy) CanObserveRead(context.Context, string, string) (bool, error) {
	return true, nil
}

type ReceiptService struct {
	messages MessageLookup
	receipts ReceiptStore
	access   ConversationAccess
	privacy  ReceiptPrivacyPolicy
}

func NewReceiptService(messages MessageLookup, receipts ReceiptStore, access ConversationAccess) (*ReceiptService, error) {
	return NewReceiptServiceWithPrivacy(messages, receipts, access, allowAllReceiptPrivacyPolicy{})
}

func NewReceiptServiceWithPrivacy(messages MessageLookup, receipts ReceiptStore, access ConversationAccess, privacy ReceiptPrivacyPolicy) (*ReceiptService, error) {
	if messages == nil || receipts == nil || access == nil || privacy == nil {
		return nil, errors.New("message lookup, receipt store, conversation access, and receipt privacy policy are required")
	}
	return &ReceiptService{messages: messages, receipts: receipts, access: access, privacy: privacy}, nil
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
	if receipt.State == domain.ReceiptRead {
		allowed, err := s.privacy.CanPublishRead(ctx, receipt.ConversationID, authenticatedUserID)
		if err != nil {
			return fmt.Errorf("verify read receipt publication privacy: %w", err)
		}
		if !allowed {
			return ErrReceiptPrivacyDenied
		}
	}
	return s.receipts.PutReceipt(ctx, receipt)
}

func (s *ReceiptService) List(ctx context.Context, authenticatedUserID, messageID string) ([]domain.DeliveryReceipt, error) {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return nil, errors.New("authenticated user id is required")
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

	receipts, err := s.receipts.ListReceipts(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("list receipts: %w", err)
	}

	hasRead := false
	for _, receipt := range receipts {
		if receipt.State == domain.ReceiptRead {
			hasRead = true
			break
		}
	}
	if !hasRead {
		return receipts, nil
	}

	canObserveRead, err := s.privacy.CanObserveRead(ctx, message.ConversationID, authenticatedUserID)
	if err != nil {
		return nil, fmt.Errorf("verify read receipt observation privacy: %w", err)
	}

	visible := make([]domain.DeliveryReceipt, 0, len(receipts))
	for _, receipt := range receipts {
		if receipt.State != domain.ReceiptRead {
			visible = append(visible, receipt)
			continue
		}
		if !canObserveRead {
			continue
		}
		canPublishRead, err := s.privacy.CanPublishRead(ctx, message.ConversationID, receipt.UserID)
		if err != nil {
			return nil, fmt.Errorf("verify stored read receipt publication privacy: %w", err)
		}
		if canPublishRead {
			visible = append(visible, receipt)
		}
	}
	return visible, nil
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

// MemoryReceiptPrivacyPolicy is a deterministic Development policy implementation.
type MemoryReceiptPrivacyPolicy struct {
	mu             sync.RWMutex
	defaultAllowed bool
	publish        map[string]bool
	observe        map[string]bool
}

func NewMemoryReceiptPrivacyPolicy(defaultAllowed bool) *MemoryReceiptPrivacyPolicy {
	return &MemoryReceiptPrivacyPolicy{
		defaultAllowed: defaultAllowed,
		publish:        make(map[string]bool),
		observe:        make(map[string]bool),
	}
}

func (p *MemoryReceiptPrivacyPolicy) SetPublish(conversationID, userID string, allowed bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.publish[receiptPrivacyKey(conversationID, userID)] = allowed
}

func (p *MemoryReceiptPrivacyPolicy) SetObserve(conversationID, userID string, allowed bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.observe[receiptPrivacyKey(conversationID, userID)] = allowed
}

func (p *MemoryReceiptPrivacyPolicy) CanPublishRead(_ context.Context, conversationID, userID string) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	allowed, ok := p.publish[receiptPrivacyKey(conversationID, userID)]
	if !ok {
		return p.defaultAllowed, nil
	}
	return allowed, nil
}

func (p *MemoryReceiptPrivacyPolicy) CanObserveRead(_ context.Context, conversationID, userID string) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	allowed, ok := p.observe[receiptPrivacyKey(conversationID, userID)]
	if !ok {
		return p.defaultAllowed, nil
	}
	return allowed, nil
}

func receiptPrivacyKey(conversationID, userID string) string {
	return strings.TrimSpace(conversationID) + "\x00" + strings.TrimSpace(userID)
}
