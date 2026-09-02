// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFileTypingPrivacyPolicyHeldWriteLockFailsClosedWithoutReplacingPreference(t *testing.T) {
	root := t.TempDir()
	first, err := NewFileTypingPrivacyPolicy(root, false)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewFileTypingPrivacyPolicy(root, false)
	if err != nil {
		t.Fatal(err)
	}
	initial := TypingPrivacyPreferences{PublishTyping: false, ObserveTyping: true}
	if err := first.SetTypingPreferences(context.Background(), "conversation-a", "user-a", initial); err != nil {
		t.Fatal(err)
	}
	path, err := first.recordPath("conversation-a", "user-a")
	if err != nil {
		t.Fatal(err)
	}

	release, err := acquireTypingPrivacyWriteLock(context.Background(), path)
	if err != nil {
		t.Fatalf("acquire external write lock: %v", err)
	}
	defer func() { _ = release() }()

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	err = second.SetTypingPreferences(
		ctx,
		"conversation-a",
		"user-a",
		TypingPrivacyPreferences{PublishTyping: true, ObserveTyping: false},
	)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("SetTypingPreferences error = %v, want context deadline", err)
	}

	got, err := first.GetTypingPreferences(context.Background(), "conversation-a", "user-a")
	if err != nil {
		t.Fatal(err)
	}
	if got != initial {
		t.Fatalf("contended write replaced committed preference: got %#v want %#v", got, initial)
	}
}

func TestTypingPrivacyWriteLockIsOwnerOnlyAndIdentifierMinimized(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), false)
	if err != nil {
		t.Fatal(err)
	}
	path, err := policy.recordPath("conversation-private", "user-private")
	if err != nil {
		t.Fatal(err)
	}
	release, err := acquireTypingPrivacyWriteLock(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	lockPath := path + ".lock"
	info, err := os.Lstat(lockPath)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() || info.Mode().Perm() != 0o700 {
		t.Fatalf("write lock mode = %v, want owner-only directory 0700", info.Mode())
	}
	for _, forbidden := range []string{"conversation-private", "user-private"} {
		if strings.Contains(lockPath, forbidden) {
			t.Fatalf("write lock path exposed identifier %q", forbidden)
		}
	}
	if err := release(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(lockPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("released write lock still exists: %v", err)
	}
}

func TestTypingPrivacyWriteLockFailsClosedWhenStaleLockExists(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), false)
	if err != nil {
		t.Fatal(err)
	}
	path, err := policy.recordPath("conversation-a", "user-a")
	if err != nil {
		t.Fatal(err)
	}
	lockPath := path + ".lock"
	if err := os.Mkdir(lockPath, 0o700); err != nil {
		t.Fatal(err)
	}

	started := time.Now()
	_, err = acquireTypingPrivacyWriteLock(context.Background(), path)
	if err == nil || !strings.Contains(err.Error(), "write lock is already held") {
		t.Fatalf("expected bounded stale-lock rejection, got %v", err)
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("stale write lock rejection took too long: %v", elapsed)
	}
}
