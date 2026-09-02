// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"os"
	"path/filepath"
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
	payload := []byte(`{"version":1,"publish_typing":true,"observe_typing":false,"unexpected":true}`)
	if err := os.WriteFile(path, payload, 0o600); err != nil {
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
	payload := []byte(`{"version":1,"publish_typing":true,"observe_typing":false} {"version":1}`)
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("write trailing record: %v", err)
	}

	_, err = policy.GetTypingPreferences(context.Background(), "conversation-3", "user-3")
	if err == nil || !strings.Contains(err.Error(), "trailing JSON data") {
		t.Fatalf("expected trailing-data rejection, got %v", err)
	}
}

func TestReadTypingPrivacyRecordRejectsSymlinkPathEvenWhenTargetIsValid(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target.json")
	payload := []byte(`{"version":1,"publish_typing":true,"observe_typing":false}`)
	if err := os.WriteFile(target, payload, 0o600); err != nil {
		t.Fatalf("write target record: %v", err)
	}
	link := filepath.Join(root, "record.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("create record symlink: %v", err)
	}

	_, err := readTypingPrivacyRecord(link)
	if err == nil || !strings.Contains(err.Error(), "not a protected regular file") {
		t.Fatalf("expected symlink rejection at read boundary, got %v", err)
	}
}

func TestReadTypingPrivacyRecordRejectsBroadOpenedPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "record.json")
	payload := []byte(`{"version":1,"publish_typing":true,"observe_typing":false}`)
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write broad-permission record: %v", err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatalf("set broad record permissions: %v", err)
	}

	_, err := readTypingPrivacyRecord(path)
	if err == nil || !strings.Contains(err.Error(), "permissions are broader than owner-only") {
		t.Fatalf("expected opened-record permission rejection, got %v", err)
	}
}
