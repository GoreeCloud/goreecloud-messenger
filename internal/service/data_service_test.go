// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func validEnvelope(id, nonce string) domain.DataEnvelope {
	return domain.DataEnvelope{
		MessageID:      id,
		ConversationID: "conv-1",
		SenderID:       "user-1",
		ClientNonce:    nonce,
		Ciphertext:     []byte("encrypted payload"),
		Encryption:     domain.EncryptionE2EE,
		CreatedAt:      time.Now().UTC(),
	}
}

func newTestDataService(t *testing.T) *DataService {
	t.Helper()

	access := NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conv-1",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-1", "user-2"},
	}); err != nil {
		t.Fatalf("set conversation: %v", err)
	}

	service, err := NewDataService(NewMemoryDataStore(), access)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}
	return service
}

func TestDataServiceSubmitAndList(t *testing.T) {
	service := newTestDataService(t)
	envelope := validEnvelope("msg-1", "nonce-1")

	if err := service.Submit(context.Background(), "user-1", envelope); err != nil {
		t.Fatalf("submit envelope: %v", err)
	}

	messages, err := service.ListConversation(context.Background(), "user-2", "conv-1")
	if err != nil {
		t.Fatalf("list conversation: %v", err)
	}
	if len(messages) != 1 || messages[0].MessageID != "msg-1" {
		t.Fatalf("unexpected stored messages: %#v", messages)
	}
}

func TestDataServiceRejectsNonceReuse(t *testing.T) {
	service := newTestDataService(t)

	if err := service.Submit(context.Background(), "user-1", validEnvelope("msg-1", "nonce-1")); err != nil {
		t.Fatalf("first submit: %v", err)
	}

	err := service.Submit(context.Background(), "user-1", validEnvelope("msg-2", "nonce-1"))
	if !errors.Is(err, ErrNonceReuse) {
		t.Fatalf("expected nonce reuse error, got %v", err)
	}
}

func TestDataServiceRejectsSenderMismatch(t *testing.T) {
	service := newTestDataService(t)

	err := service.Submit(context.Background(), "user-2", validEnvelope("msg-1", "nonce-1"))
	if !errors.Is(err, ErrSenderMismatch) {
		t.Fatalf("expected sender mismatch, got %v", err)
	}
}

func TestDataServiceRejectsCrossConversationAccess(t *testing.T) {
	service := newTestDataService(t)

	_, err := service.ListConversation(context.Background(), "user-3", "conv-1")
	if !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("expected conversation access error, got %v", err)
	}
}

func TestDataServiceDoesNotExposeMutableCiphertext(t *testing.T) {
	service := newTestDataService(t)
	envelope := validEnvelope("msg-1", "nonce-1")

	if err := service.Submit(context.Background(), "user-1", envelope); err != nil {
		t.Fatalf("submit envelope: %v", err)
	}
	envelope.Ciphertext[0] = 'X'

	messages, err := service.ListConversation(context.Background(), "user-1", "conv-1")
	if err != nil {
		t.Fatalf("list conversation: %v", err)
	}
	if string(messages[0].Ciphertext) != "encrypted payload" {
		t.Fatalf("stored ciphertext was mutated: %q", messages[0].Ciphertext)
	}
}
