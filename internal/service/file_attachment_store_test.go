// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileAttachmentStorePersistsOpaqueAttachmentAcrossRestart(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileAttachmentStore(root)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	attachment := testAttachment()
	attachment.AttachmentID = "../../attachment-one"
	attachment.ClientNonce = "restart-nonce"
	attachment.Ciphertext = []byte("opaque-ciphertext")

	if err := store.PutAttachment(context.Background(), attachment); err != nil {
		t.Fatalf("put attachment: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "..", "attachment-one")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("attachment identifier escaped hashed storage boundary: %v", err)
	}

	restarted, err := NewFileAttachmentStore(root)
	if err != nil {
		t.Fatalf("restart store: %v", err)
	}
	loaded, ok, err := restarted.GetAttachment(context.Background(), attachment.AttachmentID)
	if err != nil || !ok {
		t.Fatalf("get after restart: ok=%v err=%v", ok, err)
	}
	if string(loaded.Ciphertext) != string(attachment.Ciphertext) {
		t.Fatalf("ciphertext changed across restart: %q", loaded.Ciphertext)
	}
	if loaded.Filename != attachment.Filename || loaded.MIMEType != attachment.MIMEType {
		t.Fatalf("metadata changed across restart: %#v", loaded)
	}

	listed, err := restarted.ListAttachmentMetadata(context.Background(), attachment.ConversationID, 10)
	if err != nil {
		t.Fatalf("list metadata: %v", err)
	}
	if len(listed) != 1 || listed[0].AttachmentID != attachment.AttachmentID {
		t.Fatalf("unexpected listing: %#v", listed)
	}
	if listed[0].CiphertextBytes != len(attachment.Ciphertext) {
		t.Fatalf("ciphertext byte count = %d", listed[0].CiphertextBytes)
	}
}

func TestFileAttachmentStoreEnforcesNonceUniquenessAcrossRestart(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileAttachmentStore(root)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	first := testAttachment()
	first.AttachmentID = "attachment-one"
	first.ClientNonce = "shared-nonce"
	if err := store.PutAttachment(context.Background(), first); err != nil {
		t.Fatalf("put first: %v", err)
	}

	restarted, err := NewFileAttachmentStore(root)
	if err != nil {
		t.Fatalf("restart store: %v", err)
	}
	second := testAttachment()
	second.AttachmentID = "attachment-two"
	second.ClientNonce = first.ClientNonce
	if err := restarted.PutAttachment(context.Background(), second); !errors.Is(err, ErrAttachmentNonceReuse) {
		t.Fatalf("expected nonce reuse rejection, got %v", err)
	}
}

func TestFileAttachmentStoreFailsClosedWhenCommittedCiphertextLengthDrifts(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileAttachmentStore(root)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	attachment := testAttachment()
	attachment.AttachmentID = "attachment-corrupt"
	attachment.ClientNonce = "corrupt-nonce"
	attachment.Ciphertext = []byte("committed-ciphertext")
	if err := store.PutAttachment(context.Background(), attachment); err != nil {
		t.Fatalf("put attachment: %v", err)
	}

	cipherPath := filepath.Join(root, "ciphertext", attachmentStorageKey(attachment.AttachmentID)+".bin")
	if err := os.WriteFile(cipherPath, []byte("short"), attachmentFileMode); err != nil {
		t.Fatalf("corrupt ciphertext fixture: %v", err)
	}
	if _, ok, err := store.GetAttachment(context.Background(), attachment.AttachmentID); err == nil || ok {
		t.Fatalf("expected fail-closed length mismatch, ok=%v err=%v", ok, err)
	}
}

func TestFileAttachmentStoreUsesPrivateDirectoryAndFileModes(t *testing.T) {
	root := filepath.Join(t.TempDir(), "store")
	store, err := NewFileAttachmentStore(root)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	attachment := testAttachment()
	attachment.AttachmentID = "attachment-mode"
	attachment.ClientNonce = "mode-nonce"
	if err := store.PutAttachment(context.Background(), attachment); err != nil {
		t.Fatalf("put attachment: %v", err)
	}

	for _, path := range []string{root, filepath.Join(root, "metadata"), filepath.Join(root, "ciphertext")} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat directory %s: %v", path, err)
		}
		if info.Mode().Perm() != attachmentDirectoryMode {
			t.Fatalf("directory mode %s = %o", path, info.Mode().Perm())
		}
	}
	key := attachmentStorageKey(attachment.AttachmentID)
	for _, path := range []string{
		filepath.Join(root, "metadata", key+".json"),
		filepath.Join(root, "ciphertext", key+".bin"),
	} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat file %s: %v", path, err)
		}
		if info.Mode().Perm() != attachmentFileMode {
			t.Fatalf("file mode %s = %o", path, info.Mode().Perm())
		}
	}
}
