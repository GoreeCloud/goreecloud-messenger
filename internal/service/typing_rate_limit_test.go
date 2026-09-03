// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func TestTypingServiceRateLimitsRepeatedTypingButNeverDelaysIdle(t *testing.T) {
	now := time.Date(2026, 9, 2, 22, 0, 0, 0, time.UTC)
	service, _ := newTypingTestService(t, &now)

	publish := func(sequence uint64, state domain.TypingState) error {
		return service.Publish(context.Background(), "user-1", domain.TypingSignal{
			ConversationID: "conversation-typing",
			UserID:         "user-1",
			Sequence:       sequence,
			State:          state,
		})
	}

	if err := publish(1, domain.TypingStateTyping); err != nil {
		t.Fatalf("initial typing publish: %v", err)
	}

	now = now.Add(100 * time.Millisecond)
	if err := publish(2, domain.TypingStateTyping); !errors.Is(err, ErrTypingRateLimited) {
		t.Fatalf("rapid typing error = %v, want %v", err, ErrTypingRateLimited)
	}

	// A stop signal must never be delayed by the typing=true flood guard.
	if err := publish(3, domain.TypingStateIdle); err != nil {
		t.Fatalf("idle publish after rate limit: %v", err)
	}

	// Successful idle clears the per-participant reservation, so a genuine restart is immediate.
	if err := publish(4, domain.TypingStateTyping); err != nil {
		t.Fatalf("typing restart after idle: %v", err)
	}
}

func TestTypingServiceAllowsTypingAtMinimumInterval(t *testing.T) {
	now := time.Date(2026, 9, 2, 22, 0, 0, 0, time.UTC)
	service, _ := newTypingTestService(t, &now)

	first := domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       1,
		State:          domain.TypingStateTyping,
	}
	if err := service.Publish(context.Background(), "user-1", first); err != nil {
		t.Fatalf("initial typing publish: %v", err)
	}

	now = now.Add(TypingPublishMinInterval)
	first.Sequence = 2
	if err := service.Publish(context.Background(), "user-1", first); err != nil {
		t.Fatalf("typing publish at minimum interval: %v", err)
	}
}

func TestTypingServiceDoesNotRetainFailedTypingReservation(t *testing.T) {
	now := time.Date(2026, 9, 2, 22, 0, 0, 0, time.UTC)
	service, _ := newTypingTestService(t, &now)

	if err := service.Publish(context.Background(), "user-1", domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       2,
		State:          domain.TypingStateIdle,
	}); err != nil {
		t.Fatalf("seed sequence: %v", err)
	}

	// This typing signal reserves a rate slot but fails the store's monotonic sequence check.
	if err := service.Publish(context.Background(), "user-1", domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       1,
		State:          domain.TypingStateTyping,
	}); !errors.Is(err, ErrTypingStaleSignal) {
		t.Fatalf("stale typing error = %v, want %v", err, ErrTypingStaleSignal)
	}

	// A valid signal at the same instant must still be accepted because the failed reservation rolled back.
	if err := service.Publish(context.Background(), "user-1", domain.TypingSignal{
		ConversationID: "conversation-typing",
		UserID:         "user-1",
		Sequence:       3,
		State:          domain.TypingStateTyping,
	}); err != nil {
		t.Fatalf("valid typing after failed reservation: %v", err)
	}
}
