// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"fmt"
	"testing"
	"time"
)

func TestTypingPublishLimiterRejectsNewReservationAtGlobalCapacity(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(
			fmt.Sprintf("conversation-%d", index),
			fmt.Sprintf("user-%d", index),
		)] = now
	}

	if service.reserveTypingPublish("conversation-over-capacity", "new-user", now) {
		t.Fatal("new typing reservation unexpectedly succeeded at global capacity")
	}
	if got := len(service.lastTypingPublish); got != TypingPublishReservationLimit {
		t.Fatalf("reservation count = %d, want %d", got, TypingPublishReservationLimit)
	}
}

func TestTypingPublishLimiterAllowsExistingKeyRefreshAtGlobalCapacity(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(
			fmt.Sprintf("conversation-%d", index),
			fmt.Sprintf("user-%d", index),
		)] = now
	}

	refreshAt := now.Add(TypingPublishMinInterval)
	if !service.reserveTypingPublish("conversation-0", "user-0", refreshAt) {
		t.Fatal("existing reservation refresh was rejected solely because the limiter was at global capacity")
	}
	if got := len(service.lastTypingPublish); got != TypingPublishReservationLimit {
		t.Fatalf("reservation count = %d, want %d", got, TypingPublishReservationLimit)
	}
}

func TestTypingPublishLimiterPrunesExpiredStateBeforeGlobalCapacityCheck(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(
			fmt.Sprintf("conversation-%d", index),
			fmt.Sprintf("user-%d", index),
		)] = now
	}

	later := now.Add(TypingIndicatorTTL)
	if !service.reserveTypingPublish("conversation-after-expiry", "user-after-expiry", later) {
		t.Fatal("reservation after TTL cleanup was rejected")
	}
	if got := len(service.lastTypingPublish); got != 1 {
		t.Fatalf("reservation count after cleanup = %d, want 1", got)
	}
}

func TestTypingPublishLimiterRejectsNewConversationAtParticipantCapacity(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishParticipantReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(fmt.Sprintf("conversation-%d", index), "user-a")] = now
	}

	if service.reserveTypingPublish("conversation-over-participant-capacity", "user-a", now) {
		t.Fatal("participant unexpectedly acquired a reservation above its active capacity")
	}
	if got := participantTypingPublishReservationCount(service.lastTypingPublish, "user-a"); got != TypingPublishParticipantReservationLimit {
		t.Fatalf("participant reservation count = %d, want %d", got, TypingPublishParticipantReservationLimit)
	}
}

func TestTypingPublishLimiterParticipantCapacityDoesNotBlockAnotherParticipant(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishParticipantReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(fmt.Sprintf("conversation-%d", index), "user-a")] = now
	}

	if !service.reserveTypingPublish("conversation-b", "user-b", now) {
		t.Fatal("participant capacity for user-a incorrectly blocked user-b")
	}
	if got := participantTypingPublishReservationCount(service.lastTypingPublish, "user-b"); got != 1 {
		t.Fatalf("user-b reservation count = %d, want 1", got)
	}
}

func TestTypingPublishLimiterAllowsExistingParticipantKeyRefreshAtParticipantCapacity(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishParticipantReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(fmt.Sprintf("conversation-%d", index), "user-a")] = now
	}

	refreshAt := now.Add(TypingPublishMinInterval)
	if !service.reserveTypingPublish("conversation-0", "user-a", refreshAt) {
		t.Fatal("existing participant reservation refresh was rejected solely because the participant was at capacity")
	}
}

func TestTypingPublishLimiterParticipantCapacityIsReleasedByIdleClear(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishParticipantReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(fmt.Sprintf("conversation-%d", index), "user-a")] = now
	}

	service.clearTypingPublish("conversation-0", "user-a")
	if !service.reserveTypingPublish("conversation-new", "user-a", now) {
		t.Fatal("participant reservation was not available after idle clear")
	}
}
