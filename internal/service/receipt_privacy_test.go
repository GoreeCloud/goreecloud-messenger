// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func newReceiptPrivacyTestService(t *testing.T) (*ReceiptService, *MemoryReceiptPrivacyPolicy) {
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
	if err := dataService.Submit(context.Background(), "user-1", validEnvelope("msg-privacy", "nonce-privacy")); err != nil {
		t.Fatalf("seed message: %v", err)
	}

	privacy := NewMemoryReceiptPrivacyPolicy(true)
	service, err := NewReceiptServiceWithPrivacy(store, NewMemoryReceiptStore(), access, privacy)
	if err != nil {
		t.Fatalf("create receipt service: %v", err)
	}
	return service, privacy
}

func privacyReceipt(state domain.ReceiptState) domain.DeliveryReceipt {
	receipt := validReceipt(state)
	receipt.MessageID = "msg-privacy"
	return receipt
}

func TestReceiptServiceRejectsReadWhenPublishPrivacyDisabled(t *testing.T) {
	service, privacy := newReceiptPrivacyTestService(t)
	privacy.SetPublish("conv-1", "user-2", false)

	if err := service.Record(context.Background(), "user-2", privacyReceipt(domain.ReceiptDelivered)); err != nil {
		t.Fatalf("record delivered: %v", err)
	}
	err := service.Record(context.Background(), "user-2", privacyReceipt(domain.ReceiptRead))
	if !errors.Is(err, ErrReceiptPrivacyDenied) {
		t.Fatalf("expected privacy denial, got %v", err)
	}

	receipts, err := service.List(context.Background(), "user-1", "msg-privacy")
	if err != nil {
		t.Fatalf("list receipts: %v", err)
	}
	if len(receipts) != 1 || receipts[0].State != domain.ReceiptDelivered {
		t.Fatalf("delivery state should remain visible: %#v", receipts)
	}
}

func TestReceiptServiceFiltersReadWhenObserverPrivacyDisabled(t *testing.T) {
	service, privacy := newReceiptPrivacyTestService(t)
	if err := service.Record(context.Background(), "user-2", privacyReceipt(domain.ReceiptRead)); err != nil {
		t.Fatalf("record read: %v", err)
	}
	privacy.SetObserve("conv-1", "user-1", false)

	receipts, err := service.List(context.Background(), "user-1", "msg-privacy")
	if err != nil {
		t.Fatalf("list receipts: %v", err)
	}
	if len(receipts) != 0 {
		t.Fatalf("observer privacy should hide read projection: %#v", receipts)
	}
}

func TestReceiptServiceHidesReadAfterPublisherDisablesSharing(t *testing.T) {
	service, privacy := newReceiptPrivacyTestService(t)
	if err := service.Record(context.Background(), "user-2", privacyReceipt(domain.ReceiptRead)); err != nil {
		t.Fatalf("record read: %v", err)
	}
	privacy.SetPublish("conv-1", "user-2", false)

	receipts, err := service.List(context.Background(), "user-1", "msg-privacy")
	if err != nil {
		t.Fatalf("list receipts: %v", err)
	}
	if len(receipts) != 0 {
		t.Fatalf("current publisher privacy should hide stored read projection: %#v", receipts)
	}
}

func TestReceiptServiceDeliveryUnaffectedByReadObservationPrivacy(t *testing.T) {
	service, privacy := newReceiptPrivacyTestService(t)
	privacy.SetPublish("conv-1", "user-2", false)
	privacy.SetObserve("conv-1", "user-1", false)

	if err := service.Record(context.Background(), "user-2", privacyReceipt(domain.ReceiptDelivered)); err != nil {
		t.Fatalf("record delivered: %v", err)
	}
	receipts, err := service.List(context.Background(), "user-1", "msg-privacy")
	if err != nil {
		t.Fatalf("list receipts: %v", err)
	}
	if len(receipts) != 1 || receipts[0].State != domain.ReceiptDelivered {
		t.Fatalf("delivery should not be gated by read privacy: %#v", receipts)
	}
}
