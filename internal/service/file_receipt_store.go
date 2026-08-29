// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

const receiptStoreVersion = 1

type receiptStoreDocument struct {
	Version  int                      `json:"version"`
	Receipts []domain.DeliveryReceipt `json:"receipts"`
}

// FileReceiptStore is a durable single-node Development ReceiptStore.
//
// It persists the latest receipt state per message/user to one private JSON document
// using write-temp-and-rename replacement. It is intentionally not a distributed
// database, replication protocol, backup-erasure guarantee, or multi-writer store.
type FileReceiptStore struct {
	mu       sync.RWMutex
	path     string
	receipts map[string]domain.DeliveryReceipt
}

func NewFileReceiptStore(root string) (*FileReceiptStore, error) {
	if root == "" {
		return nil, errors.New("receipt store root is required")
	}
	if err := os.MkdirAll(root, 0o700); err != nil {
		return nil, fmt.Errorf("create receipt store root: %w", err)
	}
	store := &FileReceiptStore{
		path:     filepath.Join(root, "receipts.json"),
		receipts: make(map[string]domain.DeliveryReceipt),
	}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *FileReceiptStore) PutReceipt(ctx context.Context, receipt domain.DeliveryReceipt) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := receipt.Validate(); err != nil {
		return fmt.Errorf("validate receipt: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := receiptKey(receipt.MessageID, receipt.UserID)
	if existing, ok := s.receipts[key]; ok && receipt.State.Rank() < existing.State.Rank() {
		return ErrReceiptRegression
	}

	next := cloneReceiptMap(s.receipts)
	next[key] = receipt
	if err := s.persist(next); err != nil {
		return err
	}
	s.receipts = next
	return nil
}

func (s *FileReceiptStore) ListReceipts(ctx context.Context, messageID string) ([]domain.DeliveryReceipt, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.DeliveryReceipt, 0)
	for _, receipt := range s.receipts {
		if receipt.MessageID == messageID {
			result = append(result, receipt)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].UserID != result[j].UserID {
			return result[i].UserID < result[j].UserID
		}
		if !result[i].ObservedAt.Equal(result[j].ObservedAt) {
			return result[i].ObservedAt.Before(result[j].ObservedAt)
		}
		return result[i].State < result[j].State
	})
	return result, nil
}

func (s *FileReceiptStore) load() error {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read receipt store: %w", err)
	}
	var document receiptStoreDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return fmt.Errorf("decode receipt store: %w", err)
	}
	if document.Version != receiptStoreVersion {
		return fmt.Errorf("unsupported receipt store version %d", document.Version)
	}
	loaded := make(map[string]domain.DeliveryReceipt, len(document.Receipts))
	for _, receipt := range document.Receipts {
		if err := receipt.Validate(); err != nil {
			return fmt.Errorf("invalid persisted receipt: %w", err)
		}
		key := receiptKey(receipt.MessageID, receipt.UserID)
		if prior, ok := loaded[key]; ok && receipt.State.Rank() < prior.State.Rank() {
			return errors.New("persisted receipt state regresses")
		}
		loaded[key] = receipt
	}
	s.receipts = loaded
	return nil
}

func (s *FileReceiptStore) persist(receipts map[string]domain.DeliveryReceipt) error {
	values := make([]domain.DeliveryReceipt, 0, len(receipts))
	for _, receipt := range receipts {
		values = append(values, receipt)
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].MessageID != values[j].MessageID {
			return values[i].MessageID < values[j].MessageID
		}
		return values[i].UserID < values[j].UserID
	})
	data, err := json.MarshalIndent(receiptStoreDocument{Version: receiptStoreVersion, Receipts: values}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode receipt store: %w", err)
	}
	data = append(data, '\n')

	temporary, err := os.CreateTemp(filepath.Dir(s.path), ".receipts-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary receipt store: %w", err)
	}
	temporaryName := temporary.Name()
	cleanup := func() {
		_ = temporary.Close()
		_ = os.Remove(temporaryName)
	}
	if err := temporary.Chmod(0o600); err != nil {
		cleanup()
		return fmt.Errorf("protect temporary receipt store: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		cleanup()
		return fmt.Errorf("write temporary receipt store: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		cleanup()
		return fmt.Errorf("sync temporary receipt store: %w", err)
	}
	if err := temporary.Close(); err != nil {
		_ = os.Remove(temporaryName)
		return fmt.Errorf("close temporary receipt store: %w", err)
	}
	if err := os.Rename(temporaryName, s.path); err != nil {
		_ = os.Remove(temporaryName)
		return fmt.Errorf("replace receipt store: %w", err)
	}
	if err := os.Chmod(s.path, 0o600); err != nil {
		return fmt.Errorf("protect receipt store: %w", err)
	}
	return nil
}

func receiptKey(messageID, userID string) string {
	return messageID + "\x00" + userID
}

func cloneReceiptMap(source map[string]domain.DeliveryReceipt) map[string]domain.DeliveryReceipt {
	result := make(map[string]domain.DeliveryReceipt, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}
