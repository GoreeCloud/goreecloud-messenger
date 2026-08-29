// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"testing"
	"time"
)

func validReaction() MessageReaction {
	return MessageReaction{
		ReactionID:     "reaction-1",
		MessageID:      "message-1",
		ConversationID: "conversation-1",
		ReactorID:      "user-2",
		ClientNonce:    "nonce-1",
		Operation:      ReactionSet,
		Ciphertext:     []byte("opaque encrypted reaction"),
		Encryption:     EncryptionE2EE,
		ReactedAt:      time.Now().UTC(),
	}
}

func TestMessageReactionValidatesEncryptedSet(t *testing.T) {
	if err := validReaction().Validate(); err != nil {
		t.Fatalf("validate reaction: %v", err)
	}
}

func TestMessageReactionRejectsPlaintextlessSet(t *testing.T) {
	reaction := validReaction()
	reaction.Ciphertext = nil
	if err := reaction.Validate(); err == nil {
		t.Fatal("expected set reaction without ciphertext to fail")
	}
}

func TestMessageReactionClearRejectsCiphertext(t *testing.T) {
	reaction := validReaction()
	reaction.Operation = ReactionClear
	if err := reaction.Validate(); err == nil {
		t.Fatal("expected clear reaction carrying ciphertext to fail")
	}
}

func TestMessageReactionClearAllowsNoCiphertext(t *testing.T) {
	reaction := validReaction()
	reaction.Operation = ReactionClear
	reaction.Ciphertext = nil
	if err := reaction.Validate(); err != nil {
		t.Fatalf("validate clear reaction: %v", err)
	}
}
