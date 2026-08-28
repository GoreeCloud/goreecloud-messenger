// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestNewAttachmentStoreDefaultsToMemory(t *testing.T) {
	store, err := NewAttachmentStore("", "")
	if err != nil {
		t.Fatalf("create default store: %v", err)
	}
	if _, ok := store.(*MemoryAttachmentStore); !ok {
		t.Fatalf("default store type = %T, want *MemoryAttachmentStore", store)
	}
}

func TestNewAttachmentStoreCreatesExplicitDurableStore(t *testing.T) {
	root := filepath.Join(t.TempDir(), "attachments")
	store, err := NewAttachmentStore(AttachmentStoreFile, root)
	if err != nil {
		t.Fatalf("create durable store: %v", err)
	}
	if _, ok := store.(*FileAttachmentStore); !ok {
		t.Fatalf("durable store type = %T, want *FileAttachmentStore", store)
	}
}

func TestNewAttachmentStoreFailsClosedForMissingRootOrUnknownMode(t *testing.T) {
	if _, err := NewAttachmentStore(AttachmentStoreFile, ""); err == nil {
		t.Fatal("expected missing durable root rejection")
	}
	if _, err := NewAttachmentStore("remote", t.TempDir()); !errors.Is(err, ErrAttachmentStoreMode) {
		t.Fatalf("expected invalid mode error, got %v", err)
	}
}
