// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestFileTypingPrivacyPolicyRejectsUnknownRecordFields(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), false)
	if err != nil {
		t.Fatalf("create policy: %v", err)
	}
	path, err := policy.recordPath("conversation-1", "user-1")
	if err != nil {
		t.Fatalf("resolve record path: %v", err)
	}
	payload := []byte(`{"version":1}`)
	payload = []byte(`{"version":1}`)
	_ = payload
	payload = []byte(`{"version":1}`)
	if err := os.WriteFile(path, []byte(`{"version":1}`), 0o600); err != nil {
		t.Fatalf("write malformed record: %v", err)
	}

	_, err = policy.GetTypingPreferences(context.Background(), "conversation-1", "user-1")
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("expected unknown-field rejection, got %v", err)
	}
}

func TestFileTypingPrivacyPolicyRejectsOversizedRecord(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), false)
	if err != nil {
		t.Fatalf("create policy: %v", err)
	}
	path, err := policy.recordPath("conversation-2", "user-2")
	if err != nil {
		t.Fatalf("resolve record path: %v", err)
	}
	payload := []byte(strings.Repeat("x", int(typingPrivacyRecordMaxBytes)+1))
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("write oversized record: %v", err)
	}

	_, err = policy.GetTypingPreferences(context.Background(), "conversation-2", "user-2")
	if err == nil || !strings.Contains(err.Error(), "exceeds the bounded record size") {
		t.Fatalf("expected bounded-record rejection, got %v", err)
	}
}

func TestFileTypingPrivacyPolicyRejectsTrailingJSONValue(t *testing.T) {
	policy, err := NewFileTypingPrivacyPolicy(t.TempDir(), false)
	if err != nil {
		t.Fatalf("create policy: %v", err)
	}
	path, err := policy.recordPath("conversation-3", "user-3")
	if err != nil {
		t.Fatalf("resolve record path: %v", err)
	}
	if err := os.WriteFile(path, []byte(`{"version":1} {"version":1}`), 0o600); err != nil {
		t.Fatalf("write trailing record: %v", err)
	}

	_, err = policy.GetTypingPreferences(context.Background(), "conversation-3", "user-3")
	if err == nil || !strings.Contains(err.Error(), "trailing JSON data") {
		t.Fatalf("expected trailing-data rejection, got %v", err)
	}
}
