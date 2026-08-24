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

func TestDataServiceSubmitAndList(t *testing.T) {
	store := NewMemoryDataStore()
	service, err := NewDataService(store)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}

	envelope := validEnvelope("msg-1", "nonce-1")
	if err := service.Submit(context.Background(), envelope); err != nil {
		t.Fatalf("submit envelope: %v", err)
	}

	messages, err := service.ListConversation(context.Background(), "conv-1")
	if err != nil {
		t.Fatalf("list conversation: %v", err)
	}
	if len(messages) != 1 || messages[0].MessageID != "msg-1" {
		t.Fatalf("unexpected stored messages: %#v", messages)
	}
}

func TestDataServiceRejectsNonceReuse(t *testing.T) {
	store := NewMemoryDataStore()
	service, err := NewDataService(store)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}

	if err := service.Submit(context.Background(), validEnvelope("msg-1", "nonce-1")); err != nil {
		t.Fatalf("first submit: %v", err)
	}

	err = service.Submit(context.Background(), validEnvelope("msg-2", "nonce-1"))
	if !errors.Is(err, ErrNonceReuse) {
		t.Fatalf("expected nonce reuse error, got %v", err)
	}
}

func TestDataServiceDoesNotExposeMutableCiphertext(t *testing.T) {
	store := NewMemoryDataStore()
	service, err := NewDataService(store)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}

	envelope := validEnvelope("msg-1", "nonce-1")
	if err := service.Submit(context.Background(), envelope); err != nil {
		t.Fatalf("submit envelope: %v", err)
	}
	envelope.Ciphertext[0] = 'X'

	messages, err := service.ListConversation(context.Background(), "conv-1")
	if err != nil {
		t.Fatalf("list conversation: %v", err)
	}
	if string(messages[0].Ciphertext) != "encrypted payload" {
		t.Fatalf("stored ciphertext was mutated: %q", messages[0].Ciphertext)
	}
}
