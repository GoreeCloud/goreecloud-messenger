// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

const (
	TypingIndicatorTTL                       = 10 * time.Second
	TypingPublishMinInterval                 = 250 * time.Millisecond
	TypingPublishReservationLimit            = 4096
	TypingPublishParticipantReservationLimit = 64
)

var (
	ErrTypingUserMismatch  = errors.New("authenticated user does not match typing user")
	ErrTypingPrivacyDenied = errors.New("typing indicator privacy policy denied operation")
	ErrTypingStaleSignal   = errors.New("typing signal does not advance current state")
	ErrTypingRateLimited   = errors.New("typing signal publish rate exceeded")
)

// TypingPrivacyPolicy controls whether a participant may publish or observe typing state.
type TypingPrivacyPolicy interface {
	CanPublishTyping(context.Context, string, string) (bool, error)
	CanObserveTyping(context.Context, string, string) (bool, error)
}

// TypingStore keeps only the latest sequence and short-lived active projection.
type TypingStore interface {
	ApplyTyping(context.Context, domain.TypingSignal, time.Time) error
	ListActiveTyping(context.Context, string, time.Time) ([]domain.ActiveTyping, error)
}

// TypingService authorizes content-free, privacy-controlled ephemeral typing state.
type TypingService struct {
	store  TypingStore
	access ConversationAccess
	policy TypingPrivacyPolicy
	now    func() time.Time

	publishMu         sync.Mutex
	lastTypingPublish map[string]time.Time
}

func NewTypingService(store TypingStore, access ConversationAccess, policy TypingPrivacyPolicy, now func() time.Time) (*TypingService, error) {
	if store == nil || access == nil || policy == nil || now == nil {
		return nil, errors.New("typing store, conversation access, privacy policy, and clock are required")
	}
	return &TypingService{
		store:             store,
		access:            access,
		policy:            policy,
		now:               now,
		lastTypingPublish: make(map[string]time.Time),
	}, nil
}

func (s *TypingService) Publish(ctx context.Context, authenticatedUserID string, signal domain.TypingSignal) error {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return errors.New("authenticated user id is required")
	}
	if err := signal.Validate(); err != nil {
		return fmt.Errorf("validate typing signal: %w", err)
	}
	if authenticatedUserID != signal.UserID {
		return ErrTypingUserMismatch
	}

	allowed, err := s.access.IsParticipant(ctx, signal.ConversationID, authenticatedUserID)
	if err != nil {
		return fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return ErrConversationAccess
	}
	allowed, err = s.policy.CanPublishTyping(ctx, signal.ConversationID, authenticatedUserID)
	if err != nil {
		return fmt.Errorf("evaluate typing publish privacy: %w", err)
	}
	if !allowed {
		return ErrTypingPrivacyDenied
	}

	now := s.now().UTC()
	reservedTypingPublish := false
	if signal.State == domain.TypingStateTyping {
		if !s.reserveTypingPublish(signal.ConversationID, authenticatedUserID, now) {
			return ErrTypingRateLimited
		}
		reservedTypingPublish = true
	}

	expiresAt := time.Time{}
	if signal.State == domain.TypingStateTyping {
		expiresAt = now.Add(TypingIndicatorTTL)
	}
	if err := s.store.ApplyTyping(ctx, signal, expiresAt); err != nil {
		if reservedTypingPublish {
			s.releaseTypingPublish(signal.ConversationID, authenticatedUserID, now)
		}
		return fmt.Errorf("apply typing signal: %w", err)
	}
	if signal.State == domain.TypingStateIdle {
		s.clearTypingPublish(signal.ConversationID, authenticatedUserID)
	}
	return nil
}

func (s *TypingService) List(ctx context.Context, authenticatedUserID, conversationID string) ([]domain.ActiveTyping, error) {
	if strings.TrimSpace(authenticatedUserID) == "" {
		return nil, errors.New("authenticated user id is required")
	}
	conversationID = strings.TrimSpace(conversationID)
	if conversationID == "" {
		return nil, errors.New("conversation id is required")
	}

	allowed, err := s.access.IsParticipant(ctx, conversationID, authenticatedUserID)
	if err != nil {
		return nil, fmt.Errorf("verify conversation access: %w", err)
	}
	if !allowed {
		return nil, ErrConversationAccess
	}
	allowed, err = s.policy.CanObserveTyping(ctx, conversationID, authenticatedUserID)
	if err != nil {
		return nil, fmt.Errorf("evaluate typing observe privacy: %w", err)
	}
	if !allowed {
		return nil, ErrTypingPrivacyDenied
	}

	active, err := s.store.ListActiveTyping(ctx, conversationID, s.now().UTC())
	if err != nil {
		return nil, fmt.Errorf("list active typing: %w", err)
	}
	filtered := make([]domain.ActiveTyping, 0, len(active))
	for _, indicator := range active {
		if indicator.UserID == authenticatedUserID {
			continue
		}
		filtered = append(filtered, indicator)
	}
	return filtered, nil
}

func (s *TypingService) reserveTypingPublish(conversationID, userID string, now time.Time) bool {
	s.publishMu.Lock()
	defer s.publishMu.Unlock()

	for key, lastPublish := range s.lastTypingPublish {
		if !lastPublish.Add(TypingIndicatorTTL).After(now) {
			delete(s.lastTypingPublish, key)
		}
	}

	key := typingStateKey(conversationID, userID)
	lastPublish, exists := s.lastTypingPublish[key]
	if exists && now.Before(lastPublish.Add(TypingPublishMinInterval)) {
		return false
	}
	if !exists {
		if len(s.lastTypingPublish) >= TypingPublishReservationLimit {
			return false
		}
		if participantTypingPublishReservationCount(s.lastTypingPublish, userID) >= TypingPublishParticipantReservationLimit {
			return false
		}
	}
	s.lastTypingPublish[key] = now
	return true
}

func participantTypingPublishReservationCount(reservations map[string]time.Time, userID string) int {
	participantSuffix := "." + typingStateComponent(userID)
	count := 0
	for key := range reservations {
		if strings.HasSuffix(key, participantSuffix) {
			count++
		}
	}
	return count
}

func (s *TypingService) releaseTypingPublish(conversationID, userID string, reservedAt time.Time) {
	s.publishMu.Lock()
	defer s.publishMu.Unlock()
	key := typingStateKey(conversationID, userID)
	if current, ok := s.lastTypingPublish[key]; ok && current.Equal(reservedAt) {
		delete(s.lastTypingPublish, key)
	}
}

func (s *TypingService) clearTypingPublish(conversationID, userID string) {
	s.publishMu.Lock()
	defer s.publishMu.Unlock()
	delete(s.lastTypingPublish, typingStateKey(conversationID, userID))
}

// MemoryTypingStore is an ephemeral deterministic Development implementation.
type MemoryTypingStore struct {
	mu       sync.Mutex
	sequence map[string]uint64
	active   map[string]domain.ActiveTyping
}

func NewMemoryTypingStore() *MemoryTypingStore {
	return &MemoryTypingStore{
		sequence: make(map[string]uint64),
		active:   make(map[string]domain.ActiveTyping),
	}
}

func (s *MemoryTypingStore) ApplyTyping(_ context.Context, signal domain.TypingSignal, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := typingStateKey(signal.ConversationID, signal.UserID)
	if signal.Sequence <= s.sequence[key] {
		return ErrTypingStaleSignal
	}
	s.sequence[key] = signal.Sequence
	if signal.State == domain.TypingStateIdle {
		delete(s.active, key)
		return nil
	}
	s.active[key] = domain.ActiveTyping{
		ConversationID: signal.ConversationID,
		UserID:         signal.UserID,
		Sequence:       signal.Sequence,
		ExpiresAt:      expiresAt,
	}
	return nil
}

func (s *MemoryTypingStore) ListActiveTyping(_ context.Context, conversationID string, now time.Time) ([]domain.ActiveTyping, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]domain.ActiveTyping, 0)
	prefix := typingStateComponent(conversationID) + "."
	for key, indicator := range s.active {
		if !indicator.ExpiresAt.After(now) {
			delete(s.active, key)
			continue
		}
		if strings.HasPrefix(key, prefix) {
			result = append(result, indicator)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].UserID < result[j].UserID
	})
	return result, nil
}

// MemoryTypingPrivacyPolicy provides explicit Development controls for publishing and observing typing state.
type MemoryTypingPrivacyPolicy struct {
	mu             sync.RWMutex
	defaultAllowed bool
	publish        map[string]bool
	observe        map[string]bool
}

func NewMemoryTypingPrivacyPolicy(defaultAllowed bool) *MemoryTypingPrivacyPolicy {
	return &MemoryTypingPrivacyPolicy{
		defaultAllowed: defaultAllowed,
		publish:        make(map[string]bool),
		observe:        make(map[string]bool),
	}
}

func (p *MemoryTypingPrivacyPolicy) SetPublish(conversationID, userID string, allowed bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.publish[typingStateKey(conversationID, userID)] = allowed
}

func (p *MemoryTypingPrivacyPolicy) SetObserve(conversationID, userID string, allowed bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.observe[typingStateKey(conversationID, userID)] = allowed
}

func (p *MemoryTypingPrivacyPolicy) CanPublishTyping(_ context.Context, conversationID, userID string) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	allowed, ok := p.publish[typingStateKey(conversationID, userID)]
	if !ok {
		return p.defaultAllowed, nil
	}
	return allowed, nil
}

func (p *MemoryTypingPrivacyPolicy) CanObserveTyping(_ context.Context, conversationID, userID string) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	allowed, ok := p.observe[typingStateKey(conversationID, userID)]
	if !ok {
		return p.defaultAllowed, nil
	}
	return allowed, nil
}

func typingStateComponent(value string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(value))
}

func typingStateKey(conversationID, userID string) string {
	return typingStateComponent(conversationID) + "." + typingStateComponent(userID)
}
