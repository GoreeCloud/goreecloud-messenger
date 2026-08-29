// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func newReactionTestServices(t *testing.T) (*DataService, *MessageReactionService) {
	t.Helper()

	store := NewMemoryDataStore()
	access := NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conv-1",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-1", "user-2"},
	}); err != nil {
		t.Fatalf("set conversation: %v", err)
	}
	if err := access.SetConversation(domain.Conversation{
		ID:             "conv-2",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-1", "user-3"},
	}); err != nil {
		t.Fatalf("set second conversation: %v", err)
	}

	data, err := NewDataService(store, access)
	if err != nil {
		t.Fatalf("new data service: %v", err)
	}
	reactions, err := NewMessageReactionService(store, NewMemoryMessageReactionStore(), access)
	if err != nil {
		t.Fatalf("new reaction service: %v", err)
	}
	return data, reactions
}

func reactionEnvelope(id, conversationID, senderID, nonce string, createdAt time.Time) domain.DataEnvelope {
	return domain.DataEnvelope{
		MessageID:      id,
		ConversationID: conversationID,
		SenderID:       senderID,
		ClientNonce:    nonce,
		Ciphertext:     []byte("opaque message ciphertext"),
		Encryption:     domain.EncryptionE2EE,
		CreatedAt:      createdAt,
	}
}

func reactionEvent(id, nonce, operation string, reactedAt time.Time) domain.MessageReaction {
	reaction := domain.MessageReaction{
		ReactionID:     id,
		MessageID:      "msg-1",
		ConversationID: "conv-1",
		ReactorID:      "user-2",
		ClientNonce:    nonce,
		Operation:      domain.ReactionOperation(operation),
		Encryption:     domain.EncryptionE2EE,
		ReactedAt:      reactedAt,
	}
	if reaction.Operation == domain.ReactionSet {
		reaction.Ciphertext = []byte("encrypted reaction value")
	}
	return reaction
}

func TestMessageReactionSetReplaceAndClear(t *testing.T) {
	data, reactions := newReactionTestServices(t)
	createdAt := time.Date(2026, 8, 29, 14, 0, 0, 0, time.UTC)
	if err := data.Submit(context.Background(), "user-1", reactionEnvelope("msg-1", "conv-1", "user-1", "message-nonce-1", createdAt)); err != nil {
		t.Fatalf("submit message: %v", err)
	}

	first := reactionEvent("reaction-1", "reaction-nonce-1", string(domain.ReactionSet), createdAt.Add(time.Minute))
	if err := reactions.Record(context.Background(), "user-2", first); err != nil {
		t.Fatalf("record first reaction: %v", err)
	}
	first.Ciphertext[0] = 'X'

	second := reactionEvent("reaction-2", "reaction-nonce-2", string(domain.ReactionSet), createdAt.Add(2*time.Minute))
	second.Ciphertext = []byte("encrypted replacement")
	if err := reactions.Record(context.Background(), "user-2", second); err != nil {
		t.Fatalf("replace reaction: %v", err)
	}

	current, err := reactions.List(context.Background(), "user-1", "msg-1")
	if err != nil {
		t.Fatalf("list reactions: %v", err)
	}
	if len(current) != 1 || current[0].ReactionID != "reaction-2" || string(current[0].Ciphertext) != "encrypted replacement" {
		t.Fatalf("unexpected current reaction: %#v", current)
	}

	clear := reactionEvent("reaction-3", "reaction-nonce-3", string(domain.ReactionClear), createdAt.Add(3*time.Minute))
	if err := reactions.Record(context.Background(), "user-2", clear); err != nil {
		t.Fatalf("clear reaction: %v", err)
	}
	current, err = reactions.List(context.Background(), "user-1", "msg-1")
	if err != nil {
		t.Fatalf("list cleared reactions: %v", err)
	}
	if len(current) != 0 {
		t.Fatalf("cleared reaction remained active: %#v", current)
	}
}

func TestMessageReactionRejectsStaleAndNonceReuse(t *testing.T) {
	data, reactions := newReactionTestServices(t)
	createdAt := time.Date(2026, 8, 29, 14, 0, 0, 0, time.UTC)
	if err := data.Submit(context.Background(), "user-1", reactionEnvelope("msg-1", "conv-1", "user-1", "message-nonce-1", createdAt)); err != nil {
		t.Fatalf("submit message: %v", err)
	}

	first := reactionEvent("reaction-1", "reaction-nonce-1", string(domain.ReactionSet), createdAt.Add(2*time.Minute))
	if err := reactions.Record(context.Background(), "user-2", first); err != nil {
		t.Fatalf("record reaction: %v", err)
	}

	stale := reactionEvent("reaction-2", "reaction-nonce-2", string(domain.ReactionSet), createdAt.Add(time.Minute))
	if err := reactions.Record(context.Background(), "user-2", stale); !errors.Is(err, ErrReactionStale) {
		t.Fatalf("expected stale reaction error, got %v", err)
	}

	reused := reactionEvent("reaction-3", "reaction-nonce-1", string(domain.ReactionSet), createdAt.Add(3*time.Minute))
	if err := reactions.Record(context.Background(), "user-2", reused); !errors.Is(err, ErrReactionNonceReuse) {
		t.Fatalf("expected nonce reuse error, got %v", err)
	}
}

func TestMessageReactionUsesBoundedTargetErrorAcrossConversations(t *testing.T) {
	data, reactions := newReactionTestServices(t)
	createdAt := time.Date(2026, 8, 29, 14, 0, 0, 0, time.UTC)
	if err := data.Submit(context.Background(), "user-1", reactionEnvelope("msg-cross", "conv-2", "user-1", "message-nonce-cross", createdAt)); err != nil {
		t.Fatalf("submit cross-conversation message: %v", err)
	}

	reaction := reactionEvent("reaction-1", "reaction-nonce-1", string(domain.ReactionSet), createdAt.Add(time.Minute))
	reaction.MessageID = "msg-cross"
	if err := reactions.Record(context.Background(), "user-2", reaction); !errors.Is(err, ErrReactionTargetUnavailable) {
		t.Fatalf("expected bounded target error, got %v", err)
	}

	reaction.MessageID = "missing"
	reaction.ReactionID = "reaction-2"
	reaction.ClientNonce = "reaction-nonce-2"
	if err := reactions.Record(context.Background(), "user-2", reaction); !errors.Is(err, ErrReactionTargetUnavailable) {
		t.Fatalf("expected same missing-target error, got %v", err)
	}
}

func TestMessageReactionRejectsNonparticipant(t *testing.T) {
	data, reactions := newReactionTestServices(t)
	createdAt := time.Date(2026, 8, 29, 14, 0, 0, 0, time.UTC)
	if err := data.Submit(context.Background(), "user-1", reactionEnvelope("msg-1", "conv-1", "user-1", "message-nonce-1", createdAt)); err != nil {
		t.Fatalf("submit message: %v", err)
	}

	reaction := reactionEvent("reaction-1", "reaction-nonce-1", string(domain.ReactionSet), createdAt.Add(time.Minute))
	reaction.ReactorID = "user-3"
	if err := reactions.Record(context.Background(), "user-3", reaction); !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("expected conversation access error, got %v", err)
	}
}
