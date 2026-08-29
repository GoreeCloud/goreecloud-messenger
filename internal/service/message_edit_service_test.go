// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func newTestMessageEditService(t *testing.T) (*MessageEditService, *MemoryDataStore) {
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

	dataService, err := NewDataService(store, access)
	if err != nil {
		t.Fatalf("new Data service: %v", err)
	}
	message := domain.DataEnvelope{
		MessageID:      "msg-1",
		ConversationID: "conv-1",
		SenderID:       "user-1",
		ClientNonce:    "nonce-1",
		Ciphertext:     []byte("original ciphertext"),
		Encryption:     domain.EncryptionE2EE,
		CreatedAt:      time.Date(2026, 8, 29, 13, 0, 0, 0, time.UTC),
	}
	if err := dataService.Submit(context.Background(), "user-1", message); err != nil {
		t.Fatalf("seed message: %v", err)
	}

	editService, err := NewMessageEditService(store, NewMemoryMessageEditStore(), access)
	if err != nil {
		t.Fatalf("new edit service: %v", err)
	}
	return editService, store
}

func validMessageEditServiceEdit(id, nonce string) domain.MessageEdit {
	return domain.MessageEdit{
		EditID:         id,
		MessageID:      "msg-1",
		ConversationID: "conv-1",
		EditorID:       "user-1",
		ClientNonce:    nonce,
		Ciphertext:     []byte("edited ciphertext"),
		Encryption:     domain.EncryptionE2EE,
		EditedAt:       time.Date(2026, 8, 29, 13, 30, 0, 0, time.UTC),
	}
}

func TestMessageEditServiceRecordsAndListsOpaqueEdit(t *testing.T) {
	service, _ := newTestMessageEditService(t)
	edit := validMessageEditServiceEdit("edit-1", "edit-nonce-1")

	if err := service.Record(context.Background(), "user-1", edit); err != nil {
		t.Fatalf("record edit: %v", err)
	}
	edit.Ciphertext[0] = 'X'

	edits, err := service.List(context.Background(), "user-2", "msg-1")
	if err != nil {
		t.Fatalf("list edits: %v", err)
	}
	if len(edits) != 1 || edits[0].EditID != "edit-1" {
		t.Fatalf("unexpected edits: %#v", edits)
	}
	if string(edits[0].Ciphertext) != "edited ciphertext" {
		t.Fatalf("stored ciphertext was mutated: %q", edits[0].Ciphertext)
	}
}

func TestMessageEditServiceRejectsNonSender(t *testing.T) {
	service, _ := newTestMessageEditService(t)
	edit := validMessageEditServiceEdit("edit-1", "edit-nonce-1")
	edit.EditorID = "user-2"

	err := service.Record(context.Background(), "user-2", edit)
	if !errors.Is(err, ErrEditNotSender) {
		t.Fatalf("expected original-sender enforcement, got %v", err)
	}
}

func TestMessageEditServiceRejectsNonceReuse(t *testing.T) {
	service, _ := newTestMessageEditService(t)
	if err := service.Record(context.Background(), "user-1", validMessageEditServiceEdit("edit-1", "edit-nonce-1")); err != nil {
		t.Fatalf("first edit: %v", err)
	}

	err := service.Record(context.Background(), "user-1", validMessageEditServiceEdit("edit-2", "edit-nonce-1"))
	if !errors.Is(err, ErrEditNonceReuse) {
		t.Fatalf("expected edit nonce reuse error, got %v", err)
	}
}

func TestMessageEditServiceRejectsEditBeforeOriginalMessage(t *testing.T) {
	service, _ := newTestMessageEditService(t)
	edit := validMessageEditServiceEdit("edit-1", "edit-nonce-1")
	edit.EditedAt = time.Date(2026, 8, 29, 12, 59, 59, 0, time.UTC)

	err := service.Record(context.Background(), "user-1", edit)
	if !errors.Is(err, ErrEditBeforeMessage) {
		t.Fatalf("expected edit-before-message error, got %v", err)
	}
}

func TestMessageEditServiceRejectsConversationMismatch(t *testing.T) {
	service, _ := newTestMessageEditService(t)
	edit := validMessageEditServiceEdit("edit-1", "edit-nonce-1")
	edit.ConversationID = "other-conversation"

	err := service.Record(context.Background(), "user-1", edit)
	if !errors.Is(err, ErrMessageNotFound) {
		t.Fatalf("expected message-not-found boundary, got %v", err)
	}
}
