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
	ConversationID  string `json:"conversation_id"`
	SenderID        string `json:"sender_id"`
	ClientNonce     string `json:"client_nonce"`
	Filename        string `json:"filename"`
	MIMEType        string `json:"mime_type"`
	CiphertextBytes int    `json:"ciphertext_bytes"`
}

// FileAttachmentStore is a durable local store for already-encrypted GoreeCloud Data
// attachment bytes. Attachment identifiers are hashed before becoming filenames, raw
// ciphertext is never interpreted, and metadata is the commit marker for a completed write.
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
		if metadata.ConversationID != conversationID {
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
	if metadata.AttachmentID == "" || metadata.ConversationID == "" || metadata.ClientNonce == "" || metadata.CiphertextBytes < 1 {
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
