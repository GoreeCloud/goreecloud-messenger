// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"path/filepath"
	"testing"
)

func TestNewTypingPrivacyRuntimeStoreMemory(t *testing.T) {
	store, err := NewTypingPrivacyRuntimeStore(TypingPrivacyPersistenceConfig{
		Mode:           TypingPrivacyPersistenceMemory,
		DefaultAllowed: false,
	})
	if err != nil {
		t.Fatalf("create memory store: %v", err)
	}
	preferences, err := store.GetTypingPreferences(context.Background(), "conversation-1", "user-1")
	if err != nil {
		t.Fatalf("get defaults: %v", err)
	}
	if preferences.PublishTyping || preferences.ObserveTyping {
		t.Fatal("expected configured deny-by-default memory preferences")
	}
}

func TestNewTypingPrivacyRuntimeStoreFilePersistsAcrossReopen(t *testing.T) {
	root := filepath.Join(t.TempDir(), "typing-privacy")
	config := TypingPrivacyPersistenceConfig{
		Mode:           TypingPrivacyPersistenceFile,
		RootDir:        root,
		DefaultAllowed: true,
	}
	store, err := NewTypingPrivacyRuntimeStore(config)
	if err != nil {
		t.Fatalf("create file store: %v", err)
	}
	want := TypingPrivacyPreferences{PublishTyping: false, ObserveTyping: true}
	if err := store.SetTypingPreferences(context.Background(), "conversation-1", "user-1", want); err != nil {
		t.Fatalf("set preferences: %v", err)
	}

	reopened, err := NewTypingPrivacyRuntimeStore(config)
	if err != nil {
		t.Fatalf("reopen file store: %v", err)
	}
	got, err := reopened.GetTypingPreferences(context.Background(), "conversation-1", "user-1")
	if err != nil {
		t.Fatalf("get reopened preferences: %v", err)
	}
	if got != want {
		t.Fatalf("reopened preferences = %+v, want %+v", got, want)
	}
}

func TestNewTypingPrivacyRuntimeStoreRejectsAmbiguousConfiguration(t *testing.T) {
	cases := []TypingPrivacyPersistenceConfig{
		{},
		{Mode: "remote"},
		{Mode: TypingPrivacyPersistenceMemory, RootDir: t.TempDir()},
		{Mode: TypingPrivacyPersistenceFile},
	}
	for _, config := range cases {
		if _, err := NewTypingPrivacyRuntimeStore(config); err == nil {
			t.Fatalf("expected configuration %+v to fail closed", config)
		}
	}
}
