// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func TestFileReceiptStoreSurvivesReopenAndPreservesLatestState(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileReceiptStore(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	delivered := receipt("message-a", "conversation-a", "user-b", domain.ReceiptDelivered, 1)
	read := receipt("message-a", "conversation-a", "user-b", domain.ReceiptRead, 2)
	other := receipt("message-a", "conversation-a", "user-c", domain.ReceiptDelivered, 3)

	for _, value := range []domain.DeliveryReceipt{delivered, read, other} {
		if err := store.PutReceipt(ctx, value); err != nil {
			t.Fatal(err)
		}
	}

	reopened, err := NewFileReceiptStore(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := reopened.ListReceipts(ctx, "message-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("expected two latest receipt projections, got %d", len(got))
	}
	if got[0].UserID != "user-b" || got[0].State != domain.ReceiptRead {
		t.Fatalf("expected user-b read receipt after reopen, got %#v", got[0])
	}
	if got[1].UserID != "user-c" || got[1].State != domain.ReceiptDelivered {
		t.Fatalf("expected user-c delivered receipt after reopen, got %#v", got[1])
	}

	info, err := os.Stat(filepath.Join(root, "receipts.json"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected private receipt store permissions, got %o", info.Mode().Perm())
	}
}

func TestFileReceiptStoreRejectsRegressionWithoutChangingDurableState(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileReceiptStore(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := store.PutReceipt(ctx, receipt("message-a", "conversation-a", "user-b", domain.ReceiptRead, 2)); err != nil {
		t.Fatal(err)
	}
	if err := store.PutReceipt(ctx, receipt("message-a", "conversation-a", "user-b", domain.ReceiptDelivered, 3)); !errors.Is(err, ErrReceiptRegression) {
		t.Fatalf("expected receipt regression error, got %v", err)
	}

	reopened, err := NewFileReceiptStore(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := reopened.ListReceipts(ctx, "message-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].State != domain.ReceiptRead {
		t.Fatalf("regression changed durable state: %#v", got)
	}
}

func TestFileReceiptStoreFailsClosedOnCorruptOrUnsupportedState(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "receipts.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFileReceiptStore(root); err == nil {
		t.Fatal("expected corrupt store to fail closed")
	}

	unsupported := []byte("{\"version\":99,\"receipts\":[]}")
	if err := os.WriteFile(path, unsupported, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFileReceiptStore(root); err == nil {
		t.Fatal("expected unsupported store version to fail closed")
	}
}

func TestFileReceiptStoreHonorsCanceledContext(t *testing.T) {
	store, err := NewFileReceiptStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := store.PutReceipt(ctx, receipt("message-a", "conversation-a", "user-b", domain.ReceiptDelivered, 1)); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled write, got %v", err)
	}
	if _, err := store.ListReceipts(ctx, "message-a"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled read, got %v", err)
	}
}

func receipt(messageID, conversationID, userID string, state domain.ReceiptState, second int64) domain.DeliveryReceipt {
	return domain.DeliveryReceipt{
		MessageID:      messageID,
		ConversationID: conversationID,
		UserID:         userID,
		State:          state,
		ObservedAt:     time.Unix(second, 0).UTC(),
	}
}
