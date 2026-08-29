// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
	"path/filepath"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/api"
)

func TestReceiptPersistenceEnvironmentRequiresExplicitMode(t *testing.T) {
	if _, err := receiptPersistenceFromLookup(environment(nil)); err == nil {
		t.Fatal("expected missing receipt persistence mode to fail closed")
	}
}

func TestReceiptPersistenceEnvironmentAcceptsExplicitMemoryWithoutRoot(t *testing.T) {
	config, err := receiptPersistenceFromLookup(environment(map[string]string{
		receiptPersistenceEnv: " memory ",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if config.Mode != api.ReceiptPersistenceMemory || config.Root != "" {
		t.Fatalf("unexpected memory config: %#v", config)
	}
}

func TestReceiptPersistenceEnvironmentRejectsIgnoredRootInMemoryMode(t *testing.T) {
	_, err := receiptPersistenceFromLookup(environment(map[string]string{
		receiptPersistenceEnv: "memory",
		receiptRootEnv:        filepath.Join(string(filepath.Separator), "var", "lib", "goreecloud", "messenger"),
	}))
	if err == nil {
		t.Fatal("expected memory mode with a durable root to fail closed")
	}
}

func TestReceiptPersistenceEnvironmentAcceptsAbsoluteFileRoot(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "var", "lib", "goreecloud", "messenger", "receipts")
	config, err := receiptPersistenceFromLookup(environment(map[string]string{
		receiptPersistenceEnv: "FILE",
		receiptRootEnv:        root + string(filepath.Separator) + ".",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if config.Mode != api.ReceiptPersistenceFile || config.Root != root {
		t.Fatalf("unexpected file config: %#v", config)
	}
}

func TestReceiptPersistenceEnvironmentRejectsUnsafeOrAmbiguousFileConfiguration(t *testing.T) {
	for name, values := range map[string]map[string]string{
		"missing root": {
			receiptPersistenceEnv: "file",
		},
		"relative root": {
			receiptPersistenceEnv: "file",
			receiptRootEnv:        "./receipts",
		},
		"filesystem root": {
			receiptPersistenceEnv: "file",
			receiptRootEnv:        string(filepath.Separator),
		},
		"unknown mode": {
			receiptPersistenceEnv: "database",
		},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := receiptPersistenceFromLookup(environment(values)); err == nil {
				t.Fatal("expected configuration to fail closed")
			}
		})
	}
}

func environment(values map[string]string) lookupEnvironment {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
