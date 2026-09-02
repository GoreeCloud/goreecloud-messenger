// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const typingPrivacyRecordVersion = 1
const typingPrivacyRecordMaxBytes int64 = 4096

type typingPrivacyRecord struct {
	Version       int  `json:"version"`
	PublishTyping bool `json:"publish_typing"`
	ObserveTyping bool `json:"observe_typing"`
}

func syncTypingPrivacyDirectory(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return err
	}
	defer directory.Close()
	return directory.Sync()
}

func validateTypingPrivacyDirectory(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return errors.New("typing privacy persistence root is not a protected directory")
	}
	if info.Mode().Perm()&0o077 != 0 {
		return errors.New("typing privacy persistence root permissions are broader than owner-only")
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	if filepath.Clean(resolved) != filepath.Clean(path) {
		return errors.New("typing privacy persistence root no longer resolves to its canonical path")
	}
	return nil
}

func validateTypingPrivacyRecordInfo(info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return errors.New("typing privacy preference is not a protected regular file")
	}
	if info.Mode().Perm()&0o077 != 0 {
		return errors.New("typing privacy preference permissions are broader than owner-only")
	}
	if info.Size() > typingPrivacyRecordMaxBytes {
		return errors.New("typing privacy preference exceeds the bounded record size")
	}
	return nil
}

func validateTypingPrivacyRecord(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	return validateTypingPrivacyRecordInfo(info)
}

func openValidatedTypingPrivacyRecord(path string) (*os.File, error) {
	expected, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if err := validateTypingPrivacyRecordInfo(expected); err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	opened, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if err := validateTypingPrivacyRecordInfo(opened); err != nil {
		_ = file.Close()
		return nil, err
	}
	if !os.SameFile(expected, opened) {
		_ = file.Close()
		return nil, errors.New("typing privacy preference changed between validation and open")
	}
	return file, nil
}

func readTypingPrivacyRecord(path string) (typingPrivacyRecord, error) {
	file, err := openValidatedTypingPrivacyRecord(path)
	if err != nil {
		return typingPrivacyRecord{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, typingPrivacyRecordMaxBytes+1))
	if err != nil {
		return typingPrivacyRecord{}, err
	}
	if int64(len(data)) > typingPrivacyRecordMaxBytes {
		return typingPrivacyRecord{}, errors.New("typing privacy preference exceeds the bounded record size")
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var record typingPrivacyRecord
	if err := decoder.Decode(&record); err != nil {
		return typingPrivacyRecord{}, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return typingPrivacyRecord{}, errors.New("typing privacy preference contains trailing JSON data")
		}
		return typingPrivacyRecord{}, err
	}
	return record, nil
}

func typingPrivacyContextError(ctx context.Context) error {
	if ctx == nil {
		return errors.New("typing privacy operation context is required")
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("typing privacy operation context ended: %w", err)
	}
	return nil
}

// FileTypingPrivacyPolicy is a single-node durable Development adapter for the
// minimized typing privacy preference contract. It stores only the two explicit
// booleans. Conversation/user identifiers are used only to derive a SHA-256 file
// key and are not written into the record body.
type FileTypingPrivacyPolicy struct {
	mu             sync.RWMutex
	rootDir        string
	defaultAllowed bool
}

func NewFileTypingPrivacyPolicy(rootDir string, defaultAllowed bool) (*FileTypingPrivacyPolicy, error) {
	rootDir = strings.TrimSpace(rootDir)
	if rootDir == "" {
		return nil, errors.New("typing privacy persistence root is required")
	}
	absolute, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("resolve typing privacy persistence root: %w", err)
	}
	if absolute == string(filepath.Separator) {
		return nil, errors.New("typing privacy persistence root cannot be filesystem root")
	}
	if err := os.MkdirAll(absolute, 0o700); err != nil {
		return nil, fmt.Errorf("create typing privacy persistence root: %w", err)
	}
	if err := os.Chmod(absolute, 0o700); err != nil {
		return nil, fmt.Errorf("protect typing privacy persistence root: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, fmt.Errorf("resolve typing privacy persistence root links: %w", err)
	}
	if err := validateTypingPrivacyDirectory(resolved); err != nil {
		return nil, fmt.Errorf("validate typing privacy persistence root: %w", err)
	}
	return &FileTypingPrivacyPolicy{rootDir: resolved, defaultAllowed: defaultAllowed}, nil
}

func (p *FileTypingPrivacyPolicy) GetTypingPreferences(
	ctx context.Context,
	conversationID,
	userID string,
) (TypingPrivacyPreferences, error) {
	if err := typingPrivacyContextError(ctx); err != nil {
		return TypingPrivacyPreferences{}, err
	}
	path, err := p.recordPath(conversationID, userID)
	if err != nil {
		return TypingPrivacyPreferences{}, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()
	if err := typingPrivacyContextError(ctx); err != nil {
		return TypingPrivacyPreferences{}, err
	}
	if err := validateTypingPrivacyDirectory(p.rootDir); err != nil {
		return TypingPrivacyPreferences{}, fmt.Errorf("validate typing privacy persistence root: %w", err)
	}
	if err := validateTypingPrivacyRecord(path); errors.Is(err, os.ErrNotExist) {
		if err := typingPrivacyContextError(ctx); err != nil {
			return TypingPrivacyPreferences{}, err
		}
		return TypingPrivacyPreferences{
			PublishTyping: p.defaultAllowed,
			ObserveTyping: p.defaultAllowed,
		}, nil
	} else if err != nil {
		return TypingPrivacyPreferences{}, fmt.Errorf("validate typing privacy preference: %w", err)
	}
	if err := typingPrivacyContextError(ctx); err != nil {
		return TypingPrivacyPreferences{}, err
	}
	record, err := readTypingPrivacyRecord(path)
	if err != nil {
		return TypingPrivacyPreferences{}, fmt.Errorf("decode typing privacy preference: %w", err)
	}
	if err := typingPrivacyContextError(ctx); err != nil {
		return TypingPrivacyPreferences{}, err
	}
	if record.Version != typingPrivacyRecordVersion {
		return TypingPrivacyPreferences{}, fmt.Errorf("unsupported typing privacy preference version %d", record.Version)
	}
	return TypingPrivacyPreferences{
		PublishTyping: record.PublishTyping,
		ObserveTyping: record.ObserveTyping,
	}, nil
}

func (p *FileTypingPrivacyPolicy) SetTypingPreferences(
	ctx context.Context,
	conversationID,
	userID string,
	preferences TypingPrivacyPreferences,
) error {
	if err := typingPrivacyContextError(ctx); err != nil {
		return err
	}
	path, err := p.recordPath(conversationID, userID)
	if err != nil {
		return err
	}
	record := typingPrivacyRecord{
		Version:       typingPrivacyRecordVersion,
		PublishTyping: preferences.PublishTyping,
		ObserveTyping: preferences.ObserveTyping,
	}
	bytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encode typing privacy preference: %w", err)
	}
	bytes = append(bytes, '\n')

	p.mu.Lock()
	defer p.mu.Unlock()
	if err := typingPrivacyContextError(ctx); err != nil {
		return err
	}
	if err := validateTypingPrivacyDirectory(p.rootDir); err != nil {
		return fmt.Errorf("validate typing privacy persistence root: %w", err)
	}
	temp, err := os.CreateTemp(p.rootDir, ".typing-privacy-*")
	if err != nil {
		return fmt.Errorf("create typing privacy temporary record: %w", err)
	}
	tempPath := temp.Name()
	removeTemp := true
	defer func() {
		_ = temp.Close()
		if removeTemp {
			_ = os.Remove(tempPath)
		}
	}()
	if err := temp.Chmod(0o600); err != nil {
		return fmt.Errorf("protect typing privacy temporary record: %w", err)
	}
	if _, err := temp.Write(bytes); err != nil {
		return fmt.Errorf("write typing privacy temporary record: %w", err)
	}
	if err := temp.Sync(); err != nil {
		return fmt.Errorf("sync typing privacy temporary record: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close typing privacy temporary record: %w", err)
	}
	if err := typingPrivacyContextError(ctx); err != nil {
		return err
	}
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("commit typing privacy preference: %w", err)
	}
	removeTemp = false
	if err := syncTypingPrivacyDirectory(p.rootDir); err != nil {
		return fmt.Errorf("sync typing privacy persistence root: %w", err)
	}
	return nil
}

func (p *FileTypingPrivacyPolicy) CanPublishTyping(ctx context.Context, conversationID, userID string) (bool, error) {
	preferences, err := p.GetTypingPreferences(ctx, conversationID, userID)
	if err != nil {
		return false, err
	}
	return preferences.PublishTyping, nil
}

func (p *FileTypingPrivacyPolicy) CanObserveTyping(ctx context.Context, conversationID, userID string) (bool, error) {
	preferences, err := p.GetTypingPreferences(ctx, conversationID, userID)
	if err != nil {
		return false, err
	}
	return preferences.ObserveTyping, nil
}

func (p *FileTypingPrivacyPolicy) recordPath(conversationID, userID string) (string, error) {
	conversationID = strings.TrimSpace(conversationID)
	userID = strings.TrimSpace(userID)
	if conversationID == "" || userID == "" {
		return "", errors.New("conversation id and user id are required for typing privacy persistence")
	}
	digest := sha256.Sum256([]byte(typingStateKey(conversationID, userID)))
	return filepath.Join(p.rootDir, hex.EncodeToString(digest[:])+".json"), nil
}
