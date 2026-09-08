// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"errors"
	"fmt"
	"strings"
)

type TypingPrivacyPersistenceMode string

const (
	TypingPrivacyPersistenceMemory TypingPrivacyPersistenceMode = "memory"
	TypingPrivacyPersistenceFile   TypingPrivacyPersistenceMode = "file"
)

// TypingPrivacyRuntimeStore is the shared policy/preference boundary required by
// the typing service and the authenticated typing-preference service.
type TypingPrivacyRuntimeStore interface {
	TypingPrivacyPolicy
	TypingPrivacyPreferenceStore
}

// TypingPrivacyPersistenceConfig selects one explicit Development persistence
// implementation. It does not imply production Privacy Shield acceptance.
type TypingPrivacyPersistenceConfig struct {
	Mode           TypingPrivacyPersistenceMode
	RootDir        string
	DefaultAllowed bool
}

func NewTypingPrivacyRuntimeStore(config TypingPrivacyPersistenceConfig) (TypingPrivacyRuntimeStore, error) {
	mode := TypingPrivacyPersistenceMode(strings.ToLower(strings.TrimSpace(string(config.Mode))))
	root := strings.TrimSpace(config.RootDir)

	switch mode {
	case TypingPrivacyPersistenceMemory:
		if root != "" {
			return nil, errors.New("typing privacy memory persistence must not configure a file root")
		}
		return NewMemoryTypingPrivacyPolicy(config.DefaultAllowed), nil
	case TypingPrivacyPersistenceFile:
		if root == "" {
			return nil, errors.New("typing privacy file persistence requires a root")
		}
		store, err := NewFileTypingPrivacyPolicy(root, config.DefaultAllowed)
		if err != nil {
			return nil, fmt.Errorf("initialize typing privacy file persistence: %w", err)
		}
		return store, nil
	case "":
		return nil, errors.New("typing privacy persistence mode is required")
	default:
		return nil, fmt.Errorf("unsupported typing privacy persistence mode %q", mode)
	}
}
