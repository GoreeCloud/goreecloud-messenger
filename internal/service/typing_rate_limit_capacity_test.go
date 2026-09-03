// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"fmt"
	"testing"
	"time"
)

func TestTypingPublishLimiterRejectsNewReservationAtCapacity(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(fmt.Sprintf("conversation-%d", index), "user")] = now
	}

	if service.reserveTypingPublish("conversation-over-capacity", "user", now) {
		t.Fatal("new typing reservation unexpectedly succeeded at capacity")
	}
	if got := len(service.lastTypingPublish); got != TypingPublishReservationLimit {
		t.Fatalf("reservation count = %d, want %d", got, TypingPublishReservationLimit)
	}
}

func TestTypingPublishLimiterAllowsExistingKeyRefreshAtCapacity(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(fmt.Sprintf("conversation-%d", index), "user")] = now
	}

	refreshAt := now.Add(TypingPublishMinInterval)
	if !service.reserveTypingPublish("conversation-0", "user", refreshAt) {
		t.Fatal("existing reservation refresh was rejected solely because the limiter was at capacity")
	}
	if got := len(service.lastTypingPublish); got != TypingPublishReservationLimit {
		t.Fatalf("reservation count = %d, want %d", got, TypingPublishReservationLimit)
	}
}

func TestTypingPublishLimiterPrunesExpiredStateBeforeCapacityCheck(t *testing.T) {
	now := time.Date(2026, 9, 3, 7, 0, 0, 0, time.UTC)
	service := &TypingService{lastTypingPublish: make(map[string]time.Time)}
	for index := 0; index < TypingPublishReservationLimit; index++ {
		service.lastTypingPublish[typingStateKey(fmt.Sprintf("conversation-%d", index), "user")] = now
	}

	later := now.Add(TypingIndicatorTTL)
	if !service.reserveTypingPublish("conversation-after-expiry", "user", later) {
		t.Fatal("reservation after TTL cleanup was rejected")
	}
	if got := len(service.lastTypingPublish); got != 1 {
		t.Fatalf("reservation count after cleanup = %d, want 1", got)
	}
}
