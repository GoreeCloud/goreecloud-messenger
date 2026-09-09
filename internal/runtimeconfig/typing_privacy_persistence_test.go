// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
	"path/filepath"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestTypingPrivacyPersistenceFromLookup(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    service.TypingPrivacyPersistenceConfig
		wantErr bool
	}{
		{name: "missing mode", env: map[string]string{}, wantErr: true},
		{name: "unsupported mode", env: map[string]string{typingPrivacyPersistenceEnv: "database"}, wantErr: true},
		{name: "memory", env: map[string]string{typingPrivacyPersistenceEnv: " MEMORY "}, want: service.TypingPrivacyPersistenceConfig{Mode: service.TypingPrivacyPersistenceMemory, DefaultAllowed: false}},
		{name: "memory rejects root", env: map[string]string{typingPrivacyPersistenceEnv: "memory", typingPrivacyRootEnv: "/private/typing"}, wantErr: true},
		{name: "memory rejects explicitly blank root", env: map[string]string{typingPrivacyPersistenceEnv: "memory", typingPrivacyRootEnv: ""}, wantErr: true},
		{name: "file requires root", env: map[string]string{typingPrivacyPersistenceEnv: "file"}, wantErr: true},
		{name: "file rejects relative root", env: map[string]string{typingPrivacyPersistenceEnv: "file", typingPrivacyRootEnv: "private/typing"}, wantErr: true},
		{name: "file rejects filesystem root", env: map[string]string{typingPrivacyPersistenceEnv: "file", typingPrivacyRootEnv: string(filepath.Separator)}, wantErr: true},
		{name: "file", env: map[string]string{typingPrivacyPersistenceEnv: " FILE ", typingPrivacyRootEnv: filepath.Join(string(filepath.Separator), "var", "lib", "goreecloud", "messenger", "typing")}, want: service.TypingPrivacyPersistenceConfig{Mode: service.TypingPrivacyPersistenceFile, RootDir: filepath.Join(string(filepath.Separator), "var", "lib", "goreecloud", "messenger", "typing"), DefaultAllowed: false}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lookup := func(key string) (string, bool) {
				value, ok := tc.env[key]
				return value, ok
			}
			got, err := typingPrivacyPersistenceFromLookup(lookup)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got config %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("config mismatch: got %+v want %+v", got, tc.want)
			}
		})
	}
}

func TestTypingPrivacyPersistenceFromLookupRequiresLookup(t *testing.T) {
	if _, err := typingPrivacyPersistenceFromLookup(nil); err == nil {
		t.Fatal("expected nil lookup to fail")
	}
}
