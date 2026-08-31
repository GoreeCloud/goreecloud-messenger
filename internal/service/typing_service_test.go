// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func newTypingTestService(t *testing.T, now *time.Time) (*TypingService, *MemoryTypingPrivacyPolicy) {
	t.Helper()
	access := NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-typing",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-1", "user-2"},
	}); err != nil {
		t.Fatalf("set conversation: %v", err)
	}
	policy := NewMemoryTypingPrivacyPolicy(true)
	service, err := NewTypingService(NewMemoryTypingStore(), access, policy, func() time.Time { return *now })
	if err != nil {
		t.Fatalf("new typing service: %v", err)
	}
	return service, policy
}

func TestTypingServicePublishObserveAndIdle(t *testing.T) {
	now := time.Date(2026, 8, 30, 21, 0, 0, 0, time.UTC)
	service, _ := newTypingTestService(t, &now)

	if err := service.Publish(context.Background(), "user-1", domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       1,
		State:          domain.TypingStateTyping,
	}); err != nil {
		t.Fatalf("publish typing: %v", err)
	}

	active, err := service.List(context.Background(), "user-2", "conversation-typing")
	if err != nil {
		t.Fatalf("list typing: %v", err)
	}
	if len(active) != 1 || active[0].UserID != "user-1" || active[0].Sequence != 1 {
		t.Fatalf("unexpected active typing projection: %#v", active)
	}
	if want := now.Add(TypingIndicatorTTL); !active[0].ExpiresAt.Equal(want) {
		t.Fatalf("expires at = %v, want %v", active[0].ExpiresAt, want)
	}

	if err := service.Publish(context.Background(), "user-1", domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       2,
		State:          domain.TypingStateIdle,
	}); err != nil {
		t.Fatalf("publish idle: %v", err)
	}
	active, err = service.List(context.Background(), "user-2", "conversation-typing")
	if err != nil {
		t.Fatalf("list after idle: %v", err)
	}
	if len(active) != 0 {
		t.Fatalf("idle did not clear active typing: %#v", active)
	}
}

func TestTypingServiceExpiresAndRejectsStaleSignal(t *testing.T) {
	now := time.Date(2026, 8, 30, 21, 0, 0, 0, time.UTC)
	service, _ := newTypingTestService(t, &now)

	if err := service.Publish(context.Background(), "user-1", domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       3,
		State:          domain.TypingStateTyping,
	}); err != nil {
		t.Fatalf("publish typing: %v", err)
	}
	if err := service.Publish(context.Background(), "user-1", domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       3,
		State:          domain.TypingStateIdle,
	}); !errors.Is(err, ErrTypingStaleSignal) {
		t.Fatalf("stale signal error = %v, want %v", err, ErrTypingStaleSignal)
	}

	now = now.Add(TypingIndicatorTTL)
	active, err := service.List(context.Background(), "user-2", "conversation-typing")
	if err != nil {
		t.Fatalf("list expired typing: %v", err)
	}
	if len(active) != 0 {
		t.Fatalf("expired typing remained active: %#v", active)
	}
}

func TestTypingServiceEnforcesPrivacyPolicy(t *testing.T) {
	now := time.Date(2026, 8, 30, 21, 0, 0, 0, time.UTC)
	service, policy := newTypingTestService(t, &now)
	policy.SetPublish("conversation-typing", "user-1", false)

	err := service.Publish(context.Background(), "user-1", domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       1,
		State:          domain.TypingStateTyping,
	})
	if !errors.Is(err, ErrTypingPrivacyDenied) {
		t.Fatalf("publish privacy error = %v, want %v", err, ErrTypingPrivacyDenied)
	}

	policy.SetObserve("conversation-typing", "user-2", false)
	if _, err := service.List(context.Background(), "user-2", "conversation-typing"); !errors.Is(err, ErrTypingPrivacyDenied) {
		t.Fatalf("observe privacy error = %v, want %v", err, ErrTypingPrivacyDenied)
	}
}

func TestTypingServiceRejectsUserMismatchAndOutsider(t *testing.T) {
	now := time.Date(2026, 8, 30, 21, 0, 0, 0, time.UTC)
	service, _ := newTypingTestService(t, &now)

	err := service.Publish(context.Background(), "user-2", domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       1,
		State:          domain.TypingStateTyping,
	})
	if !errors.Is(err, ErrTypingUserMismatch) {
		t.Fatalf("user mismatch error = %v, want %v", err, ErrTypingUserMismatch)
	}
	if _, err := service.List(context.Background(), "outsider", "conversation-typing"); !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("outsider list error = %v, want %v", err, ErrConversationAccess)
	}
}
