// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileTypingPrivacyPolicyRejectsSymlinkedRecordRead(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), false)
	if err != nil {
		t.Fatal(err)
	}
	path, err := policy.recordPath("conversation-a", "user-a")
	if err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "outside.json")
	if err := os.WriteFile(
		outside,
		[]byte("{\"version\":1,\"publish_typing\":true,\"observe_typing\":true}\n"),
		0o600,
	); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, path); err != nil {
		t.Skipf("symlinks unavailable on this platform: %v", err)
	}

	if _, err := policy.GetTypingPreferences(context.Background(), "conversation-a", "user-a"); err == nil {
		t.Fatal("expected symlinked typing privacy record to fail closed")
	}
}

func TestFileTypingPrivacyPolicyRejectsBroadRecordPermissions(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), false)
	if err != nil {
		t.Fatal(err)
	}
	if err := policy.SetTypingPreferences(
		context.Background(),
		"conversation-a",
		"user-a",
		TypingPrivacyPreferences{PublishTyping: true, ObserveTyping: false},
	); err != nil {
		t.Fatal(err)
	}
	path, err := policy.recordPath("conversation-a", "user-a")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := policy.GetTypingPreferences(context.Background(), "conversation-a", "user-a"); err == nil {
		t.Fatal("expected broadly readable typing privacy record to fail closed")
	}
}

func TestFileTypingPrivacyPolicyRejectsPersistenceRootReplacement(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "typing")
	policy, err := NewFileTypingPrivacyPolicy(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(root); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, root); err != nil {
		t.Skipf("symlinks unavailable on this platform: %v", err)
	}

	if err := policy.SetTypingPreferences(
		context.Background(),
		"conversation-a",
		"user-a",
		TypingPrivacyPreferences{PublishTyping: true, ObserveTyping: true},
	); err == nil {
		t.Fatal("expected replaced typing privacy root to fail closed")
	}
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("replacement target was unexpectedly modified: %#v", entries)
	}
}
