// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func newReceiptTestServices(t *testing.T) (*DataService, *ReceiptService) {
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
		t.Fatalf("create Data service: %v", err)
	}
	receiptService, err := NewReceiptService(store, NewMemoryReceiptStore(), access)
	if err != nil {
		t.Fatalf("create receipt service: %v", err)
	}
	if err := dataService.Submit(context.Background(), "user-1", validEnvelope("msg-1", "nonce-1")); err != nil {
		t.Fatalf("submit message: %v", err)
	}
	return dataService, receiptService
}

func validReceipt(state domain.ReceiptState) domain.DeliveryReceipt {
	return domain.DeliveryReceipt{
		MessageID:      "msg-1",
		ConversationID: "conv-1",
		UserID:         "user-2",
		State:          state,
		ObservedAt:     time.Now().UTC(),
	}
}

func TestReceiptServiceRecordsDeliveredThenRead(t *testing.T) {
	_, service := newReceiptTestServices(t)
	if err := service.Record(context.Background(), "user-2", validReceipt(domain.ReceiptDelivered)); err != nil {
		t.Fatalf("record delivered: %v", err)
	}
	if err := service.Record(context.Background(), "user-2", validReceipt(domain.ReceiptRead)); err != nil {
		t.Fatalf("record read: %v", err)
	}
	receipts, err := service.List(context.Background(), "user-1", "msg-1")
	if err != nil {
		t.Fatalf("list receipts: %v", err)
	}
	if len(receipts) != 1 || receipts[0].State != domain.ReceiptRead {
		t.Fatalf("unexpected receipts: %#v", receipts)
	}
}

func TestReceiptServiceRejectsRegression(t *testing.T) {
	_, service := newReceiptTestServices(t)
	if err := service.Record(context.Background(), "user-2", validReceipt(domain.ReceiptRead)); err != nil {
		t.Fatalf("record read: %v", err)
	}
	err := service.Record(context.Background(), "user-2", validReceipt(domain.ReceiptDelivered))
	if !errors.Is(err, ErrReceiptRegression) {
		t.Fatalf("expected regression error, got %v", err)
	}
}

func TestReceiptServiceRejectsSelfReceipt(t *testing.T) {
	_, service := newReceiptTestServices(t)
	receipt := validReceipt(domain.ReceiptDelivered)
	receipt.UserID = "user-1"
	err := service.Record(context.Background(), "user-1", receipt)
	if !errors.Is(err, ErrSelfReceipt) {
		t.Fatalf("expected self receipt error, got %v", err)
	}
}

func TestReceiptServiceRejectsUserMismatch(t *testing.T) {
	_, service := newReceiptTestServices(t)
	err := service.Record(context.Background(), "user-1", validReceipt(domain.ReceiptDelivered))
	if !errors.Is(err, ErrReceiptUserMismatch) {
		t.Fatalf("expected receipt user mismatch, got %v", err)
	}
}
