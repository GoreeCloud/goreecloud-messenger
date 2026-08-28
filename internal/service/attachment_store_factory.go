// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"errors"
	"fmt"
	"strings"
)

type AttachmentStoreMode string

const (
	AttachmentStoreMemory AttachmentStoreMode = "memory"
	AttachmentStoreFile   AttachmentStoreMode = "file"
)

var ErrAttachmentStoreMode = errors.New("attachment store mode is invalid")

// NewAttachmentStore creates the explicitly selected encrypted-attachment persistence boundary.
// Memory remains the development default; durable file mode requires an explicit private root.
func NewAttachmentStore(mode AttachmentStoreMode, root string) (AttachmentStore, error) {
	switch AttachmentStoreMode(strings.TrimSpace(string(mode))) {
	case "", AttachmentStoreMemory:
		return NewMemoryAttachmentStore(), nil
	case AttachmentStoreFile:
		if strings.TrimSpace(root) == "" {
			return nil, errors.New("durable attachment store root is required")
		}
		store, err := NewFileAttachmentStore(root)
		if err != nil {
			return nil, fmt.Errorf("create durable attachment store: %w", err)
		}
		return store, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrAttachmentStoreMode, mode)
	}
}
