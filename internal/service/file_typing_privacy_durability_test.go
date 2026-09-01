// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileTypingPrivacyPolicyResolvesPersistenceRootSymlink(t *testing.T) {
	base := t.TempDir()
	realRoot := filepath.Join(base, "real")
	if err := os.Mkdir(realRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	linkedRoot := filepath.Join(base, "linked")
	if err := os.Symlink(realRoot, linkedRoot); err != nil {
		t.Skipf("symlinks unavailable on this platform: %v", err)
	}

	policy, err := NewFileTypingPrivacyPolicy(linkedRoot, false)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(realRoot)
	if err != nil {
		t.Fatal(err)
	}
	if policy.rootDir != resolved {
		t.Fatalf("rootDir = %q, want canonical %q", policy.rootDir, resolved)
	}
	if err := policy.SetTypingPreferences(
		context.Background(),
		"conversation-a",
		"user-a",
		TypingPrivacyPreferences{PublishTyping: true, ObserveTyping: false},
	); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir(realRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || filepath.Ext(entries[0].Name()) != ".json" {
		t.Fatalf("expected one committed record and no temporary artifacts, got %#v", entries)
	}
}
