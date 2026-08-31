// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	"github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

type allowTypingConversationAccess struct{}

func (allowTypingConversationAccess) IsParticipant(context.Context, string, string) (bool, error) {
	return true, nil
}

func TestTypingPresenceRuntimeSharesPrivacyStoreAcrossServices(t *testing.T) {
	now := time.Date(2026, time.August, 31, 17, 45, 0, 0, time.UTC)
	runtime, err := typingPresenceRuntimeFromConfig(
		service.TypingPrivacyPersistenceConfig{
			Mode:           service.TypingPrivacyPersistenceMemory,
			DefaultAllowed: false,
		},
		service.NewMemoryTypingStore(),
		allowTypingConversationAccess{},
		func() time.Time { return now },
	)
	if err != nil {
		t.Fatalf("compose typing presence runtime: %v", err)
	}
	if runtime.PersistenceMode != service.TypingPrivacyPersistenceMemory {
		t.Fatalf("persistence mode mismatch: %q", runtime.PersistenceMode)
	}

	ctx := context.Background()
	aliceSignal := domain.TypingSignal{
		ConversationID: "conversation-1",
		UserID:         "alice",
		Sequence:       1,
		State:          domain.TypingStateTyping,
	}
	if err := runtime.Typing.Publish(ctx, "alice", aliceSignal); !errors.Is(err, service.ErrTypingPrivacyDenied) {
		t.Fatalf("default-deny typing publish should fail with privacy denial, got %v", err)
	}

	if _, err := runtime.Preferences.Update(ctx, "alice", "conversation-1", service.TypingPrivacyPreferences{
		PublishTyping: true,
		ObserveTyping: true,
	}); err != nil {
		t.Fatalf("update alice typing privacy preferences: %v", err)
	}
	if err := runtime.Typing.Publish(ctx, "alice", aliceSignal); err != nil {
		t.Fatalf("typing publish should observe preference update through shared store: %v", err)
	}

	if _, err := runtime.Preferences.Update(ctx, "bob", "conversation-1", service.TypingPrivacyPreferences{
		PublishTyping: true,
		ObserveTyping: true,
	}); err != nil {
		t.Fatalf("update bob typing privacy preferences: %v", err)
	}
	if err := runtime.Typing.Publish(ctx, "bob", domain.TypingSignal{
		ConversationID: "conversation-1",
		UserID:         "bob",
		Sequence:       1,
		State:          domain.TypingStateTyping,
	}); err != nil {
		t.Fatalf("publish bob typing state: %v", err)
	}

	active, err := runtime.Typing.List(ctx, "alice", "conversation-1")
	if err != nil {
		t.Fatalf("list active typing state: %v", err)
	}
	if len(active) != 1 || active[0].UserID != "bob" {
		t.Fatalf("expected alice to observe bob through shared privacy store, got %+v", active)
	}
}

func TestTypingPresenceRuntimeRejectsInvalidComposition(t *testing.T) {
	_, err := typingPresenceRuntimeFromConfig(
		service.TypingPrivacyPersistenceConfig{Mode: service.TypingPrivacyPersistenceMemory},
		nil,
		allowTypingConversationAccess{},
		time.Now,
	)
	if err == nil {
		t.Fatal("expected nil typing store to fail")
	}
}
