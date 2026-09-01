// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestFileTypingPrivacyPolicyRejectsCanceledReadBeforeDefaultFallback(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), true)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := policy.GetTypingPreferences(ctx, "conversation-a", "user-a"); !errors.Is(err, context.Canceled) {
		t.Fatalf("GetTypingPreferences error = %v, want context.Canceled", err)
	}
}

func TestFileTypingPrivacyPolicyCanceledWriteDoesNotReplaceCommittedPreference(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), false)
	if err != nil {
		t.Fatal(err)
	}
	if err := policy.SetTypingPreferences(
		context.Background(),
		"conversation-a",
		"user-a",
		TypingPrivacyPreferences{PublishTyping: false, ObserveTyping: false},
	); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := policy.SetTypingPreferences(
		ctx,
		"conversation-a",
		"user-a",
		TypingPrivacyPreferences{PublishTyping: true, ObserveTyping: true},
	); !errors.Is(err, context.Canceled) {
		t.Fatalf("SetTypingPreferences error = %v, want context.Canceled", err)
	}

	preferences, err := policy.GetTypingPreferences(context.Background(), "conversation-a", "user-a")
	if err != nil {
		t.Fatal(err)
	}
	if preferences.PublishTyping || preferences.ObserveTyping {
		t.Fatalf("canceled write replaced committed privacy preference: %#v", preferences)
	}
	entries, err := os.ReadDir(policy.rootDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("canceled write left unexpected persistence artifacts: %#v", entries)
	}
}

func TestFileTypingPrivacyPolicyRejectsNilOperationContext(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := policy.GetTypingPreferences(nil, "conversation-a", "user-a"); err == nil {
		t.Fatal("expected nil read context to fail closed")
	}
	if err := policy.SetTypingPreferences(
		nil,
		"conversation-a",
		"user-a",
		TypingPrivacyPreferences{},
	); err == nil {
		t.Fatal("expected nil write context to fail closed")
	}
}
