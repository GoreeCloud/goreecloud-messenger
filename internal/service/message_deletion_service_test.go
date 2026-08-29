// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func newTestMessageDeletionService(t *testing.T) (*MessageDeletionService, *DataService) {
	t.Helper()

	store := NewMemoryDataStore()
	access := NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-1",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-1", "user-2"},
	}); err != nil {
		t.Fatalf("set conversation: %v", err)
	}
	dataService, err := NewDataService(store, access)
	if err != nil {
		t.Fatalf("new Data service: %v", err)
	}
	if err := dataService.Submit(context.Background(), "user-1", domain.DataEnvelope{
		MessageID: "message-1", ConversationID: "conversation-1", SenderID: "user-1", ClientNonce: "nonce-1",
		Ciphertext: []byte("opaque ciphertext"), Encryption: domain.EncryptionE2EE,
		CreatedAt: time.Date(2026, 8, 29, 13, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("seed message: %v", err)
	}
	service, err := NewMessageDeletionService(store, NewMemoryMessageDeletionStore(), access)
	if err != nil {
		t.Fatalf("new message deletion service: %v", err)
	}
	return service, dataService
}

func validDeletionRecord() domain.MessageDeletion {
	return domain.MessageDeletion{
		DeletionID: "deletion-1", MessageID: "message-1", ConversationID: "conversation-1",
		DeleterID: "user-1", ClientNonce: "delete-nonce-1", Scope: domain.MessageDeletionEveryone,
		DeletedAt: time.Date(2026, 8, 29, 14, 0, 0, 0, time.UTC),
	}
}

func TestMessageDeletionServiceRecordAndList(t *testing.T) {
	service, _ := newTestMessageDeletionService(t)
	deletion := validDeletionRecord()
	if err := service.Record(context.Background(), "user-1", deletion); err != nil {
		t.Fatalf("record deletion: %v", err)
	}

	deletions, err := service.List(context.Background(), "user-2", "message-1")
	if err != nil {
		t.Fatalf("list deletion: %v", err)
	}
	if len(deletions) != 1 || deletions[0].DeletionID != "deletion-1" {
		t.Fatalf("unexpected deletion state: %#v", deletions)
	}
}

func TestMessageDeletionServiceRejectsNonSender(t *testing.T) {
	service, _ := newTestMessageDeletionService(t)
	deletion := validDeletionRecord()
	deletion.DeleterID = "user-2"
	if err := service.Record(context.Background(), "user-2", deletion); !errors.Is(err, ErrDeletionNotSender) {
		t.Fatalf("expected non-sender rejection, got %v", err)
	}
}

func TestMessageDeletionServiceRejectsDeletionBeforeMessage(t *testing.T) {
	service, _ := newTestMessageDeletionService(t)
	deletion := validDeletionRecord()
	deletion.DeletedAt = time.Date(2026, 8, 29, 12, 59, 59, 0, time.UTC)
	if err := service.Record(context.Background(), "user-1", deletion); !errors.Is(err, ErrDeletionBeforeMessage) {
		t.Fatalf("expected timestamp rejection, got %v", err)
	}
}

func TestMessageDeletionServiceRejectsSecondTombstone(t *testing.T) {
	service, _ := newTestMessageDeletionService(t)
	if err := service.Record(context.Background(), "user-1", validDeletionRecord()); err != nil {
		t.Fatalf("record first deletion: %v", err)
	}
	second := validDeletionRecord()
	second.DeletionID = "deletion-2"
	second.ClientNonce = "delete-nonce-2"
	second.DeletedAt = second.DeletedAt.Add(time.Minute)
	if err := service.Record(context.Background(), "user-1", second); !errors.Is(err, ErrMessageAlreadyDeleted) {
		t.Fatalf("expected second tombstone rejection, got %v", err)
	}
}

func TestMemoryMessageDeletionStoreRejectsNonceReuse(t *testing.T) {
	store := NewMemoryMessageDeletionStore()
	first := validDeletionRecord()
	if err := store.PutDeletion(context.Background(), first); err != nil {
		t.Fatalf("put first deletion: %v", err)
	}
	second := first
	second.DeletionID = "deletion-2"
	second.MessageID = "message-2"
	if err := store.PutDeletion(context.Background(), second); !errors.Is(err, ErrDeletionNonceReuse) {
		t.Fatalf("expected nonce reuse rejection, got %v", err)
	}
}

func TestMessageDeletionServiceRejectsConversationOutsider(t *testing.T) {
	service, _ := newTestMessageDeletionService(t)
	if _, err := service.List(context.Background(), "user-3", "message-1"); !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("expected conversation access rejection, got %v", err)
	}
}
