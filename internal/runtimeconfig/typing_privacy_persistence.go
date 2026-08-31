// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

const (
	typingPrivacyPersistenceEnv = "GOREECLOUD_MESSENGER_TYPING_PRIVACY_PERSISTENCE"
	typingPrivacyRootEnv        = "GOREECLOUD_MESSENGER_TYPING_PRIVACY_ROOT"
)

// TypingPrivacyPersistenceFromEnvironment derives the Development typing privacy
// persistence selection from explicit process configuration. There is no implicit
// memory fallback and the privacy default remains deny.
func TypingPrivacyPersistenceFromEnvironment() (service.TypingPrivacyPersistenceConfig, error) {
	return typingPrivacyPersistenceFromLookup(os.LookupEnv)
}

func typingPrivacyPersistenceFromLookup(lookup lookupEnvironment) (service.TypingPrivacyPersistenceConfig, error) {
	if lookup == nil {
		return service.TypingPrivacyPersistenceConfig{}, errors.New("environment lookup is required")
	}

	modeValue, ok := lookup(typingPrivacyPersistenceEnv)
	if !ok || strings.TrimSpace(modeValue) == "" {
		return service.TypingPrivacyPersistenceConfig{}, fmt.Errorf("%s is required", typingPrivacyPersistenceEnv)
	}
	mode := service.TypingPrivacyPersistenceMode(strings.ToLower(strings.TrimSpace(modeValue)))

	rootValue, rootSet := lookup(typingPrivacyRootEnv)
	root := strings.TrimSpace(rootValue)

	switch mode {
	case service.TypingPrivacyPersistenceMemory:
		if rootSet {
			return service.TypingPrivacyPersistenceConfig{}, fmt.Errorf("%s must be unset in memory mode", typingPrivacyRootEnv)
		}
		return service.TypingPrivacyPersistenceConfig{
			Mode:           mode,
			DefaultAllowed: false,
		}, nil

	case service.TypingPrivacyPersistenceFile:
		if !rootSet || root == "" {
			return service.TypingPrivacyPersistenceConfig{}, fmt.Errorf("%s is required in file mode", typingPrivacyRootEnv)
		}
		if !filepath.IsAbs(root) {
			return service.TypingPrivacyPersistenceConfig{}, fmt.Errorf("%s must be an absolute path", typingPrivacyRootEnv)
		}
		cleaned := filepath.Clean(root)
		if cleaned == string(filepath.Separator) {
			return service.TypingPrivacyPersistenceConfig{}, fmt.Errorf("%s must not be the filesystem root", typingPrivacyRootEnv)
		}
		return service.TypingPrivacyPersistenceConfig{
			Mode:           mode,
			RootDir:        cleaned,
			DefaultAllowed: false,
		}, nil

	default:
		return service.TypingPrivacyPersistenceConfig{}, fmt.Errorf("unsupported %s value %q", typingPrivacyPersistenceEnv, modeValue)
	}
}
