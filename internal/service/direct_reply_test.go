// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func TestDataServiceStoresDirectReplyReference(t *testing.T) {
	service := newTestDataService(t)

	if err := service.Submit(context.Background(), "user-1", validEnvelope("msg-1", "nonce-1")); err != nil {
		t.Fatalf("submit original: %v", err)
	}
	reply := validEnvelope("msg-2", "nonce-2")
	reply.ReplyToMessageID = "  msg-1  "
	if err := service.Submit(context.Background(), "user-1", reply); err != nil {
		t.Fatalf("submit reply: %v", err)
	}

	messages, err := service.ListConversation(context.Background(), "user-2", "conv-1")
	if err != nil {
		t.Fatalf("list conversation: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("message count = %d, want 2", len(messages))
	}
	if messages[1].ReplyToMessageID != "msg-1" {
		t.Fatalf("reply target = %q, want msg-1", messages[1].ReplyToMessageID)
	}
}

func TestDataServiceRejectsUnavailableReplyTarget(t *testing.T) {
	service := newTestDataService(t)
	reply := validEnvelope("msg-2", "nonce-2")
	reply.ReplyToMessageID = "missing-message"

	err := service.Submit(context.Background(), "user-1", reply)
	if !errors.Is(err, ErrReplyTargetUnavailable) {
		t.Fatalf("expected reply target unavailable, got %v", err)
	}
}

func TestDataServiceUsesSameErrorForCrossConversationReplyTarget(t *testing.T) {
	store := NewMemoryDataStore()
	access := NewMemoryConversationAccess()
	for _, conversation := range []domain.Conversation{
		{ID: "conv-1", Kind: domain.ConversationDirect, ParticipantIDs: []string{"user-1", "user-2"}},
		{ID: "conv-2", Kind: domain.ConversationDirect, ParticipantIDs: []string{"user-1", "user-3"}},
	} {
		if err := access.SetConversation(conversation); err != nil {
			t.Fatalf("set conversation %s: %v", conversation.ID, err)
		}
	}
	service, err := NewDataService(store, access)
	if err != nil {
		t.Fatalf("new Data service: %v", err)
	}

	target := validEnvelope("msg-cross", "nonce-cross")
	target.ConversationID = "conv-2"
	if err := service.Submit(context.Background(), "user-1", target); err != nil {
		t.Fatalf("submit cross-conversation target: %v", err)
	}

	reply := validEnvelope("msg-2", "nonce-2")
	reply.ReplyToMessageID = "msg-cross"
	err = service.Submit(context.Background(), "user-1", reply)
	if !errors.Is(err, ErrReplyTargetUnavailable) {
		t.Fatalf("expected privacy-bounded reply error, got %v", err)
	}
}
