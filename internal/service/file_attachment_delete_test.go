// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileAttachmentStoreDeleteRemovesCiphertextAndKeepsMinimalTombstone(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileAttachmentStore(root)
	if err != nil {
		t.Fatalf("new file attachment store: %v", err)
	}
	attachment := testAttachment()
	if err := store.PutAttachment(context.Background(), attachment); err != nil {
		t.Fatalf("put attachment: %v", err)
	}

	deleted, err := store.DeleteAttachment(context.Background(), attachment.AttachmentID)
	if err != nil {
		t.Fatalf("delete attachment: %v", err)
	}
	if !deleted {
		t.Fatal("expected first delete to remove attachment")
	}

	key := attachmentStorageKey(attachment.AttachmentID)
	cipherPath := filepath.Join(root, "ciphertext", key+".bin")
	if _, err := os.Stat(cipherPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ciphertext still exists after delete: %v", err)
	}
	metadataPath := filepath.Join(root, "metadata", key+".json")
	contents, err := os.ReadFile(metadataPath)
	if err != nil {
		t.Fatalf("read deletion tombstone: %v", err)
	}
	var tombstone persistedAttachmentMetadata
	if err := json.Unmarshal(contents, &tombstone); err != nil {
		t.Fatalf("decode deletion tombstone: %v", err)
	}
	if !tombstone.Deleted || tombstone.AttachmentID != attachment.AttachmentID || tombstone.ClientNonce != attachment.ClientNonce {
		t.Fatalf("unexpected deletion tombstone: %#v", tombstone)
	}
	if tombstone.ConversationID != "" || tombstone.SenderID != "" || tombstone.Filename != "" || tombstone.MIMEType != "" || tombstone.CiphertextBytes != 0 {
		t.Fatalf("deletion tombstone retained user-facing metadata: %#v", tombstone)
	}

	if _, ok, err := store.GetAttachment(context.Background(), attachment.AttachmentID); err != nil || ok {
		t.Fatalf("deleted attachment should be unavailable, ok=%v err=%v", ok, err)
	}
	listed, err := store.ListAttachmentMetadata(context.Background(), attachment.ConversationID, 10)
	if err != nil {
		t.Fatalf("list attachment metadata: %v", err)
	}
	if len(listed) != 0 {
		t.Fatalf("deleted attachment remained discoverable: %#v", listed)
	}

	restarted, err := NewFileAttachmentStore(root)
	if err != nil {
		t.Fatalf("reopen file attachment store: %v", err)
	}
	reusedID := attachment
	reusedID.ClientNonce = "nonce-2"
	if err := restarted.PutAttachment(context.Background(), reusedID); !errors.Is(err, ErrDuplicateAttachment) {
		t.Fatalf("expected tombstoned id to remain reserved after restart, got %v", err)
	}
	reusedNonce := attachment
	reusedNonce.AttachmentID = "attachment-2"
	if err := restarted.PutAttachment(context.Background(), reusedNonce); !errors.Is(err, ErrAttachmentNonceReuse) {
		t.Fatalf("expected tombstoned nonce to remain reserved after restart, got %v", err)
	}
}
