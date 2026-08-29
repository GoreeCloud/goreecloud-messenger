// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func TestDataServiceStoresAndListsThreadedReplies(t *testing.T) {
	service := newTestDataService(t)

	if err := service.Submit(context.Background(), "user-1", validEnvelope("msg-1", "nonce-1")); err != nil {
		t.Fatalf("submit root: %v", err)
	}
	firstReply := validEnvelope("msg-2", "nonce-2")
	firstReply.ReplyToMessageID = " msg-1 "
	firstReply.ThreadRootMessageID = " msg-1 "
	if err := service.Submit(context.Background(), "user-1", firstReply); err != nil {
		t.Fatalf("submit first thread reply: %v", err)
	}
	nestedReply := validEnvelope("msg-3", "nonce-3")
	nestedReply.ReplyToMessageID = "msg-2"
	nestedReply.ThreadRootMessageID = "msg-1"
	if err := service.Submit(context.Background(), "user-1", nestedReply); err != nil {
		t.Fatalf("submit nested thread reply: %v", err)
	}

	thread, err := service.ListThread(context.Background(), "user-2", "conv-1", " msg-1 ")
	if err != nil {
		t.Fatalf("list thread: %v", err)
	}
	if len(thread) != 3 {
		t.Fatalf("thread count = %d, want 3", len(thread))
	}
	if thread[0].MessageID != "msg-1" || thread[1].MessageID != "msg-2" || thread[2].MessageID != "msg-3" {
		t.Fatalf("unexpected thread order: %#v", thread)
	}
	if thread[1].ReplyToMessageID != "msg-1" || thread[1].ThreadRootMessageID != "msg-1" {
		t.Fatalf("unexpected first thread reply metadata: %#v", thread[1])
	}
	if thread[2].ReplyToMessageID != "msg-2" || thread[2].ThreadRootMessageID != "msg-1" {
		t.Fatalf("unexpected nested thread reply metadata: %#v", thread[2])
	}
}

func TestDataServiceRejectsThreadWithoutReplyTarget(t *testing.T) {
	service := newTestDataService(t)
	if err := service.Submit(context.Background(), "user-1", validEnvelope("msg-1", "nonce-1")); err != nil {
		t.Fatalf("submit root: %v", err)
	}

	threadReply := validEnvelope("msg-2", "nonce-2")
	threadReply.ThreadRootMessageID = "msg-1"
	err := service.Submit(context.Background(), "user-1", threadReply)
	if !errors.Is(err, ErrThreadTargetUnavailable) {
		t.Fatalf("expected thread target unavailable, got %v", err)
	}
}

func TestDataServiceRejectsReplyOutsideDeclaredThread(t *testing.T) {
	service := newTestDataService(t)
	if err := service.Submit(context.Background(), "user-1", validEnvelope("msg-1", "nonce-1")); err != nil {
		t.Fatalf("submit root: %v", err)
	}
	if err := service.Submit(context.Background(), "user-1", validEnvelope("msg-2", "nonce-2")); err != nil {
		t.Fatalf("submit unrelated message: %v", err)
	}

	threadReply := validEnvelope("msg-3", "nonce-3")
	threadReply.ThreadRootMessageID = "msg-1"
	threadReply.ReplyToMessageID = "msg-2"
	err := service.Submit(context.Background(), "user-1", threadReply)
	if !errors.Is(err, ErrThreadTargetUnavailable) {
		t.Fatalf("expected thread target unavailable, got %v", err)
	}
}

func TestListThreadUsesSameErrorForCrossConversationRoot(t *testing.T) {
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

	root := validEnvelope("msg-cross", "nonce-cross")
	root.ConversationID = "conv-2"
	if err := service.Submit(context.Background(), "user-1", root); err != nil {
		t.Fatalf("submit cross-conversation root: %v", err)
	}

	_, err = service.ListThread(context.Background(), "user-1", "conv-1", "msg-cross")
	if !errors.Is(err, ErrThreadTargetUnavailable) {
		t.Fatalf("expected privacy-bounded thread error, got %v", err)
	}
}
