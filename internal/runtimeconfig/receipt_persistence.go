// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoreeCloud/goreecloud-messenger/internal/api"
)

const (
	receiptPersistenceEnv = "GOREECLOUD_MESSENGER_RECEIPT_PERSISTENCE"
	receiptRootEnv        = "GOREECLOUD_MESSENGER_RECEIPT_ROOT"
)

type lookupEnvironment func(string) (string, bool)

// ReceiptPersistenceFromEnvironment derives the receipt persistence selection from the
// executable process environment. Configuration is explicit: there is no implicit memory mode.
func ReceiptPersistenceFromEnvironment() (api.ReceiptPersistenceConfig, error) {
	return receiptPersistenceFromLookup(os.LookupEnv)
}

func receiptPersistenceFromLookup(lookup lookupEnvironment) (api.ReceiptPersistenceConfig, error) {
	if lookup == nil {
		return api.ReceiptPersistenceConfig{}, errors.New("environment lookup is required")
	}

	modeValue, ok := lookup(receiptPersistenceEnv)
	if !ok || strings.TrimSpace(modeValue) == "" {
		return api.ReceiptPersistenceConfig{}, fmt.Errorf("%s is required", receiptPersistenceEnv)
	}
	mode := api.ReceiptPersistenceMode(strings.ToLower(strings.TrimSpace(modeValue)))

	rootValue, rootSet := lookup(receiptRootEnv)
	root := strings.TrimSpace(rootValue)

	switch mode {
	case api.ReceiptPersistenceMemory:
		if rootSet && root != "" {
			return api.ReceiptPersistenceConfig{}, fmt.Errorf("%s must be unset in memory mode", receiptRootEnv)
		}
		return api.ReceiptPersistenceConfig{Mode: mode}, nil

	case api.ReceiptPersistenceFile:
		if !rootSet || root == "" {
			return api.ReceiptPersistenceConfig{}, fmt.Errorf("%s is required in file mode", receiptRootEnv)
		}
		if !filepath.IsAbs(root) {
			return api.ReceiptPersistenceConfig{}, fmt.Errorf("%s must be an absolute path", receiptRootEnv)
		}
		cleaned := filepath.Clean(root)
		if cleaned == string(filepath.Separator) {
			return api.ReceiptPersistenceConfig{}, fmt.Errorf("%s must not be the filesystem root", receiptRootEnv)
		}
		return api.ReceiptPersistenceConfig{Mode: mode, Root: cleaned}, nil

	default:
		return api.ReceiptPersistenceConfig{}, fmt.Errorf("unsupported %s value %q", receiptPersistenceEnv, modeValue)
	}
}
