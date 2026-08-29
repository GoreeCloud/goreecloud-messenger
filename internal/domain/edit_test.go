// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"testing"
	"time"
)

func validMessageEdit() MessageEdit {
	return MessageEdit{
		EditID:         "edit-1",
		MessageID:      "message-1",
		ConversationID: "conversation-1",
		EditorID:       "user-1",
		ClientNonce:    "edit-nonce-1",
		Ciphertext:     []byte("edited ciphertext"),
		Encryption:     EncryptionE2EE,
		EditedAt:       time.Date(2026, 8, 29, 13, 30, 0, 0, time.UTC),
	}
}

func TestMessageEditValidate(t *testing.T) {
	if err := validMessageEdit().Validate(); err != nil {
		t.Fatalf("valid edit rejected: %v", err)
	}
}

func TestMessageEditRejectsNonE2EE(t *testing.T) {
	edit := validMessageEdit()
	edit.Encryption = EncryptionNone
	if err := edit.Validate(); err == nil {
		t.Fatal("expected non-E2EE edit to be rejected")
	}
}

func TestMessageEditRejectsEmptyCiphertext(t *testing.T) {
	edit := validMessageEdit()
	edit.Ciphertext = nil
	if err := edit.Validate(); err == nil {
		t.Fatal("expected empty ciphertext to be rejected")
	}
}
