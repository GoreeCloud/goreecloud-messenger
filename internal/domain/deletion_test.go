// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"testing"
	"time"
)

func validMessageDeletion() MessageDeletion {
	return MessageDeletion{
		DeletionID:     "deletion-1",
		MessageID:      "message-1",
		ConversationID: "conversation-1",
		DeleterID:      "user-1",
		ClientNonce:    "delete-nonce-1",
		Scope:          MessageDeletionEveryone,
		DeletedAt:      time.Date(2026, 8, 29, 14, 0, 0, 0, time.UTC),
	}
}

func TestMessageDeletionValidate(t *testing.T) {
	if err := validMessageDeletion().Validate(); err != nil {
		t.Fatalf("valid deletion rejected: %v", err)
	}
}

func TestMessageDeletionRejectsUnsupportedScope(t *testing.T) {
	deletion := validMessageDeletion()
	deletion.Scope = "local"
	if err := deletion.Validate(); err == nil {
		t.Fatal("expected unsupported deletion scope to be rejected")
	}
}

func TestMessageDeletionRequiresReplayNonce(t *testing.T) {
	deletion := validMessageDeletion()
	deletion.ClientNonce = ""
	if err := deletion.Validate(); err == nil {
		t.Fatal("expected missing client nonce to be rejected")
	}
}
