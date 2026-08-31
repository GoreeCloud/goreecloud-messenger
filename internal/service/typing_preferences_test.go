// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func typingPreferenceServiceForTest(t *testing.T, defaultAllowed bool) (*TypingPrivacyPreferenceService, *MemoryTypingPrivacyPolicy) {
	t.Helper()
	access := NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-a",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-a", "user-b"},
	}); err != nil {
		t.Fatal(err)
	}
	policy := NewMemoryTypingPrivacyPolicy(defaultAllowed)
	service, err := NewTypingPrivacyPreferenceService(policy, access)
	if err != nil {
		t.Fatal(err)
	}
	return service, policy
}

func TestTypingPrivacyPreferencesUseEffectiveDefaultsAndExplicitOverrides(t *testing.T) {
	service, policy := typingPreferenceServiceForTest(t, true)

	initial, err := service.Get(context.Background(), "user-a", "conversation-a")
	if err != nil {
		t.Fatal(err)
	}
	if !initial.PublishTyping || !initial.ObserveTyping {
		t.Fatalf("expected enabled defaults, got %+v", initial)
	}

	updated, err := service.Update(context.Background(), "user-a", "conversation-a", TypingPrivacyPreferences{
		PublishTyping: false,
		ObserveTyping: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.PublishTyping || !updated.ObserveTyping {
		t.Fatalf("unexpected update projection: %+v", updated)
	}

	publish, err := policy.CanPublishTyping(context.Background(), "conversation-a", "user-a")
	if err != nil || publish {
		t.Fatalf("expected publish disabled, allowed=%v err=%v", publish, err)
	}
	observe, err := policy.CanObserveTyping(context.Background(), "conversation-a", "user-a")
	if err != nil || !observe {
		t.Fatalf("expected observe enabled, allowed=%v err=%v", observe, err)
	}
}

func TestTypingPrivacyPreferencesResetToPolicyOwnedDefault(t *testing.T) {
	service, policy := typingPreferenceServiceForTest(t, false)
	if _, err := service.Update(context.Background(), "user-a", "conversation-a", TypingPrivacyPreferences{
		PublishTyping: true,
		ObserveTyping: true,
	}); err != nil {
		t.Fatal(err)
	}

	reset, err := service.Reset(context.Background(), "user-a", "conversation-a")
	if err != nil {
		t.Fatal(err)
	}
	if reset.PublishTyping || reset.ObserveTyping {
		t.Fatalf("expected policy-owned disabled defaults after reset, got %+v", reset)
	}
	publish, err := policy.CanPublishTyping(context.Background(), "conversation-a", "user-a")
	if err != nil || publish {
		t.Fatalf("expected reset publish default false, allowed=%v err=%v", publish, err)
	}
	observe, err := policy.CanObserveTyping(context.Background(), "conversation-a", "user-a")
	if err != nil || observe {
		t.Fatalf("expected reset observe default false, allowed=%v err=%v", observe, err)
	}
}

func TestTypingPrivacyPreferencesRejectConversationOutsider(t *testing.T) {
	service, _ := typingPreferenceServiceForTest(t, true)

	_, err := service.Get(context.Background(), "outsider", "conversation-a")
	if !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("expected conversation access error, got %v", err)
	}
	_, err = service.Update(context.Background(), "outsider", "conversation-a", TypingPrivacyPreferences{})
	if !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("expected conversation access error, got %v", err)
	}
	_, err = service.Reset(context.Background(), "outsider", "conversation-a")
	if !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("expected conversation access error for reset, got %v", err)
	}
}
