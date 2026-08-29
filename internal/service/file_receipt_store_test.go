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
	rootInfo, err := os.Stat(root)
	if err != nil {
		t.Fatal(err)
	}
	if rootInfo.Mode().Perm() != 0o700 {
		t.Fatalf("expected private receipt root permissions, got %o", rootInfo.Mode().Perm())
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

func TestFileReceiptStoreFailsClosedOnCorruptUnsupportedOrUnsafeState(t *testing.T) {
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

	valid := []byte("{\"version\":1,\"receipts\":[]}")
	if err := os.WriteFile(path, valid, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFileReceiptStore(root); err == nil {
		t.Fatal("expected permissive persisted-state permissions to fail closed")
	}
}

func TestFileReceiptStoreProtectsRootAndRejectsSymlinkBoundaries(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "receipts")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFileReceiptStore(root); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o700 {
		t.Fatalf("expected root to be tightened to 0700, got %o", info.Mode().Perm())
	}

	target := filepath.Join(base, "target")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	rootLink := filepath.Join(base, "root-link")
	if err := os.Symlink(target, rootLink); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFileReceiptStore(rootLink); err == nil {
		t.Fatal("expected symlink receipt root to fail closed")
	}

	stateRoot := filepath.Join(base, "state-root")
	if err := os.Mkdir(stateRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	stateTarget := filepath.Join(base, "state-target.json")
	if err := os.WriteFile(stateTarget, []byte("{\"version\":1,\"receipts\":[]}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(stateTarget, filepath.Join(stateRoot, "receipts.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFileReceiptStore(stateRoot); err == nil {
		t.Fatal("expected symlink persisted state to fail closed")
	}
}

func TestFileReceiptStorePoisonsProcessStateWhenPostRenameDurabilityIsUnknown(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileReceiptStore(root)
	if err != nil {
		t.Fatal(err)
	}

	originalSync := syncReceiptStoreDirectory
	syncReceiptStoreDirectory = func(string) error { return errors.New("simulated directory sync failure") }
	ctx := context.Background()
	writeErr := store.PutReceipt(ctx, receipt("message-a", "conversation-a", "user-b", domain.ReceiptDelivered, 1))
	if !errors.Is(writeErr, ErrReceiptDurabilityUnknown) {
		t.Fatalf("expected durability-unknown error, got %v", writeErr)
	}
	if _, err := store.ListReceipts(ctx, "message-a"); !errors.Is(err, ErrReceiptDurabilityUnknown) {
		t.Fatalf("expected poisoned store to refuse reads, got %v", err)
	}
	if err := store.PutReceipt(ctx, receipt("message-b", "conversation-b", "user-b", domain.ReceiptDelivered, 2)); !errors.Is(err, ErrReceiptDurabilityUnknown) {
		t.Fatalf("expected poisoned store to refuse writes, got %v", err)
	}

	// Reopening is the explicit reconciliation boundary after an ambiguous post-rename failure.
	syncReceiptStoreDirectory = originalSync
	reopened, err := NewFileReceiptStore(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := reopened.ListReceipts(ctx, "message-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].State != domain.ReceiptDelivered {
		t.Fatalf("expected reopen to reconcile renamed durable state, got %#v", got)
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
