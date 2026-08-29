// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

const (
	attachmentDirectoryMode os.FileMode = 0o700
	attachmentFileMode      os.FileMode = 0o600
)

type persistedAttachmentMetadata struct {
	AttachmentID    string `json:"attachment_id"`
	ConversationID  string `json:"conversation_id,omitempty"`
	SenderID        string `json:"sender_id,omitempty"`
	ClientNonce     string `json:"client_nonce"`
	Filename        string `json:"filename,omitempty"`
	MIMEType        string `json:"mime_type,omitempty"`
	CiphertextBytes int    `json:"ciphertext_bytes,omitempty"`
	Deleted         bool   `json:"deleted,omitempty"`
}

// FileAttachmentStore is a durable local store for already-encrypted GoreeCloud Data
// attachment bytes. Attachment identifiers are hashed before becoming filenames, raw
// ciphertext is never interpreted, and metadata is the commit marker for a completed write.
// Deletion removes ciphertext first and then replaces metadata with a minimal tombstone so
// attachment identifiers and client nonces remain permanently reserved against replay.
type FileAttachmentStore struct {
	mu          sync.RWMutex
	root        string
	metadataDir string
	cipherDir   string
}

func NewFileAttachmentStore(root string) (*FileAttachmentStore, error) {
	if root == "" {
		return nil, errors.New("attachment store root is required")
	}
	store := &FileAttachmentStore{
		root:        root,
		metadataDir: filepath.Join(root, "metadata"),
		cipherDir:   filepath.Join(root, "ciphertext"),
	}
	for _, dir := range []string{store.root, store.metadataDir, store.cipherDir} {
		if err := os.MkdirAll(dir, attachmentDirectoryMode); err != nil {
			return nil, fmt.Errorf("create attachment store directory: %w", err)
		}
		if err := os.Chmod(dir, attachmentDirectoryMode); err != nil {
			return nil, fmt.Errorf("restrict attachment store directory: %w", err)
		}
	}
	return store, nil
}

func (s *FileAttachmentStore) PutAttachment(ctx context.Context, attachment domain.DataAttachment) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := attachment.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok, err := s.readMetadataLocked(attachment.AttachmentID); err != nil {
		return err
	} else if ok {
		return ErrDuplicateAttachment
	}

	metadataFiles, err := os.ReadDir(s.metadataDir)
	if err != nil {
		return fmt.Errorf("read attachment metadata directory: %w", err)
	}
	for _, entry := range metadataFiles {
		if entry.IsDir() {
			continue
		}
		metadata, err := s.readMetadataPathLocked(filepath.Join(s.metadataDir, entry.Name()))
		if err != nil {
			return err
		}
		if metadata.ClientNonce == attachment.ClientNonce {
			return fmt.Errorf("%w: nonce belongs to %s", ErrAttachmentNonceReuse, metadata.AttachmentID)
		}
	}

	key := attachmentStorageKey(attachment.AttachmentID)
	cipherPath := filepath.Join(s.cipherDir, key+".bin")
	metadataPath := filepath.Join(s.metadataDir, key+".json")
	cipherTemp := cipherPath + ".tmp"
	metadataTemp := metadataPath + ".tmp"
	cleanup := func() {
		_ = os.Remove(cipherTemp)
		_ = os.Remove(metadataTemp)
	}
	defer cleanup()

	if err := writePrivateFile(cipherTemp, attachment.Ciphertext); err != nil {
		return fmt.Errorf("write attachment ciphertext: %w", err)
	}
	metadataBytes, err := json.Marshal(persistedAttachmentMetadata{
		AttachmentID:    attachment.AttachmentID,
		ConversationID:  attachment.ConversationID,
		SenderID:        attachment.SenderID,
		ClientNonce:     attachment.ClientNonce,
		Filename:        attachment.Filename,
		MIMEType:        attachment.MIMEType,
		CiphertextBytes: len(attachment.Ciphertext),
	})
	if err != nil {
		return fmt.Errorf("encode attachment metadata: %w", err)
	}
	if err := writePrivateFile(metadataTemp, metadataBytes); err != nil {
		return fmt.Errorf("write attachment metadata: %w", err)
	}

	if err := os.Rename(cipherTemp, cipherPath); err != nil {
		return fmt.Errorf("commit attachment ciphertext: %w", err)
	}
	if err := os.Chmod(cipherPath, attachmentFileMode); err != nil {
		_ = os.Remove(cipherPath)
		return fmt.Errorf("restrict attachment ciphertext: %w", err)
	}
	if err := os.Rename(metadataTemp, metadataPath); err != nil {
		_ = os.Remove(cipherPath)
		return fmt.Errorf("commit attachment metadata: %w", err)
	}
	if err := os.Chmod(metadataPath, attachmentFileMode); err != nil {
		_ = os.Remove(metadataPath)
		_ = os.Remove(cipherPath)
		return fmt.Errorf("restrict attachment metadata: %w", err)
	}
	return nil
}

func (s *FileAttachmentStore) GetAttachment(ctx context.Context, attachmentID string) (domain.DataAttachment, bool, error) {
	if err := ctx.Err(); err != nil {
		return domain.DataAttachment{}, false, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	metadata, ok, err := s.readMetadataLocked(attachmentID)
	if err != nil || !ok {
		return domain.DataAttachment{}, ok, err
	}
	if metadata.Deleted {
		return domain.DataAttachment{}, false, nil
	}
	ciphertext, err := os.ReadFile(filepath.Join(s.cipherDir, attachmentStorageKey(attachmentID)+".bin"))
	if err != nil {
		return domain.DataAttachment{}, false, fmt.Errorf("read attachment ciphertext: %w", err)
	}
	if len(ciphertext) != metadata.CiphertextBytes {
		return domain.DataAttachment{}, false, errors.New("attachment ciphertext length does not match committed metadata")
	}
	return domain.DataAttachment{
		AttachmentID:   metadata.AttachmentID,
		ConversationID: metadata.ConversationID,
		SenderID:       metadata.SenderID,
		ClientNonce:    metadata.ClientNonce,
		Filename:       metadata.Filename,
		MIMEType:       metadata.MIMEType,
		Ciphertext:     append([]byte(nil), ciphertext...),
	}, true, nil
}

func (s *FileAttachmentStore) DeleteAttachment(ctx context.Context, attachmentID string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	metadata, ok, err := s.readMetadataLocked(attachmentID)
	if err != nil || !ok {
		return false, err
	}
	if metadata.Deleted {
		return false, nil
	}

	key := attachmentStorageKey(attachmentID)
	cipherPath := filepath.Join(s.cipherDir, key+".bin")
	if err := os.Remove(cipherPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("remove attachment ciphertext: %w", err)
	}

	// Preserve only replay-prevention state. User-facing metadata is removed from the tombstone.
	metadata.ConversationID = ""
	metadata.SenderID = ""
	metadata.Filename = ""
	metadata.MIMEType = ""
	metadata.CiphertextBytes = 0
	metadata.Deleted = true
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return false, fmt.Errorf("encode attachment deletion tombstone: %w", err)
	}
	metadataPath := filepath.Join(s.metadataDir, key+".json")
	metadataTemp := metadataPath + ".delete.tmp"
	defer func() { _ = os.Remove(metadataTemp) }()
	if err := writePrivateFile(metadataTemp, metadataBytes); err != nil {
		return false, fmt.Errorf("write attachment deletion tombstone: %w", err)
	}
	if err := os.Rename(metadataTemp, metadataPath); err != nil {
		return false, fmt.Errorf("commit attachment deletion tombstone: %w", err)
	}
	if err := os.Chmod(metadataPath, attachmentFileMode); err != nil {
		return false, fmt.Errorf("restrict attachment deletion tombstone: %w", err)
	}
	return true, nil
}

func (s *FileAttachmentStore) ListAttachmentMetadata(
	ctx context.Context,
	conversationID string,
	limit int,
) ([]domain.DataAttachmentMetadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.metadataDir)
	if err != nil {
		return nil, fmt.Errorf("read attachment metadata directory: %w", err)
	}
	items := make([]domain.DataAttachmentMetadata, 0, limit)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		metadata, err := s.readMetadataPathLocked(filepath.Join(s.metadataDir, entry.Name()))
		if err != nil {
			return nil, err
		}
		if metadata.Deleted || metadata.ConversationID != conversationID {
			continue
		}
		items = append(items, domain.DataAttachmentMetadata{
			AttachmentID:    metadata.AttachmentID,
			ConversationID:  metadata.ConversationID,
			SenderID:        metadata.SenderID,
			Filename:        metadata.Filename,
			MIMEType:        metadata.MIMEType,
			CiphertextBytes: metadata.CiphertextBytes,
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].AttachmentID < items[j].AttachmentID })
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}

func (s *FileAttachmentStore) readMetadataLocked(attachmentID string) (persistedAttachmentMetadata, bool, error) {
	path := filepath.Join(s.metadataDir, attachmentStorageKey(attachmentID)+".json")
	metadata, err := s.readMetadataPathLocked(path)
	if errors.Is(err, os.ErrNotExist) {
		return persistedAttachmentMetadata{}, false, nil
	}
	if err != nil {
		return persistedAttachmentMetadata{}, false, err
	}
	if metadata.AttachmentID != attachmentID {
		return persistedAttachmentMetadata{}, false, errors.New("attachment metadata identity mismatch")
	}
	return metadata, true, nil
}

func (s *FileAttachmentStore) readMetadataPathLocked(path string) (persistedAttachmentMetadata, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return persistedAttachmentMetadata{}, err
	}
	var metadata persistedAttachmentMetadata
	if err := json.Unmarshal(contents, &metadata); err != nil {
		return persistedAttachmentMetadata{}, fmt.Errorf("decode attachment metadata: %w", err)
	}
	if metadata.AttachmentID == "" || metadata.ClientNonce == "" {
		return persistedAttachmentMetadata{}, errors.New("committed attachment metadata is incomplete")
	}
	if metadata.Deleted {
		if metadata.CiphertextBytes != 0 {
			return persistedAttachmentMetadata{}, errors.New("deleted attachment tombstone retains ciphertext length")
		}
		return metadata, nil
	}
	if metadata.ConversationID == "" || metadata.SenderID == "" || metadata.Filename == "" || metadata.MIMEType == "" || metadata.CiphertextBytes < 1 {
		return persistedAttachmentMetadata{}, errors.New("committed attachment metadata is incomplete")
	}
	return metadata, nil
}

func attachmentStorageKey(attachmentID string) string {
	digest := sha256.Sum256([]byte(attachmentID))
	return hex.EncodeToString(digest[:])
}

func writePrivateFile(path string, contents []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, attachmentFileMode)
	if err != nil {
		return err
	}
	if _, err := file.Write(contents); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
