// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestRuntimeReceiptStoreSelectionFailsClosedWithoutExplicitMode(t *testing.T) {
	for _, config := range []ReceiptPersistenceConfig{
		{},
		{Mode: ReceiptPersistenceMode("unknown")},
		{Mode: ReceiptPersistenceFile},
	} {
		if _, err := newRuntimeReceiptStore(config); err == nil {
			t.Fatalf("expected config %+v to fail closed", config)
		}
	}
}

func TestRuntimeReceiptStoreSelectionSupportsExplicitMemoryMode(t *testing.T) {
	store, err := newRuntimeReceiptStore(ReceiptPersistenceConfig{Mode: ReceiptPersistenceMemory})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := store.(*messagingservice.MemoryReceiptStore); !ok {
		t.Fatalf("expected MemoryReceiptStore, got %T", store)
	}
}

func TestRuntimeReceiptStoreSelectionPersistsFileModeAcrossReopen(t *testing.T) {
	root := t.TempDir()
	config := ReceiptPersistenceConfig{Mode: ReceiptPersistenceFile, Root: root}
	first, err := newRuntimeReceiptStore(config)
	if err != nil {
		t.Fatal(err)
	}

	receipt := domain.DeliveryReceipt{
		MessageID:      "message-a",
		ConversationID: "conversation-a",
		UserID:         "user-b",
		State:          domain.ReceiptRead,
		ObservedAt:     time.Date(2026, 8, 29, 20, 0, 0, 0, time.UTC),
	}
	if err := first.PutReceipt(context.Background(), receipt); err != nil {
		t.Fatal(err)
	}

	second, err := newRuntimeReceiptStore(config)
	if err != nil {
		t.Fatal(err)
	}
	persisted, err := second.ListReceipts(context.Background(), receipt.MessageID)
	if err != nil {
		t.Fatal(err)
	}
	if len(persisted) != 1 || persisted[0] != receipt {
		t.Fatalf("expected reopened store to preserve receipt, got %+v", persisted)
	}
}

func TestConfiguredDataRuntimeRequiresAuthoritativeDependencies(t *testing.T) {
	if _, err := NewConfiguredDataRuntimeHandler(DataRuntimeDependencies{}, ReceiptPersistenceConfig{Mode: ReceiptPersistenceMemory}); err == nil {
		t.Fatal("expected missing runtime dependencies to fail closed")
	}
}
