// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"testing"
	"time"
)

func TestDeliveryReceiptValidation(t *testing.T) {
	receipt := DeliveryReceipt{
		MessageID:      "message-1",
		ConversationID: "conversation-1",
		UserID:         "user-2",
		State:          ReceiptDelivered,
		ObservedAt:     time.Now().UTC(),
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("valid receipt rejected: %v", err)
	}
}

func TestDeliveryReceiptRejectsInvalidState(t *testing.T) {
	receipt := DeliveryReceipt{
		MessageID:      "message-1",
		ConversationID: "conversation-1",
		UserID:         "user-2",
		State:          ReceiptState("seen-ish"),
		ObservedAt:     time.Now().UTC(),
	}
	if err := receipt.Validate(); err == nil {
		t.Fatal("expected invalid receipt state to be rejected")
	}
}

func TestReceiptStateRankIsMonotonic(t *testing.T) {
	if ReceiptDelivered.Rank() >= ReceiptRead.Rank() {
		t.Fatal("read must rank after delivered")
	}
}
