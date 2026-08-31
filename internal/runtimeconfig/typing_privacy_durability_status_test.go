// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
	"strings"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestTypingPrivacyDurabilityStatusMemoryIsTruthfullyTransient(t *testing.T) {
	status, err := TypingPrivacyDurabilityStatusFor(service.TypingPrivacyPersistenceMemory)
	if err != nil {
		t.Fatalf("memory durability status: %v", err)
	}
	if status.Level != TypingPrivacyDurabilityTransient {
		t.Fatalf("unexpected durability level: %q", status.Level)
	}
	if status.DurableAfterRestart || status.CrossDevice || status.ProductionReady {
		t.Fatalf("memory status overclaimed durability: %+v", status)
	}
}

func TestTypingPrivacyDurabilityStatusFileIsOnlySingleNodeDurable(t *testing.T) {
	status, err := TypingPrivacyDurabilityStatusFor(service.TypingPrivacyPersistenceFile)
	if err != nil {
		t.Fatalf("file durability status: %v", err)
	}
	if status.Level != TypingPrivacyDurabilitySingleNodeDurable {
		t.Fatalf("unexpected durability level: %q", status.Level)
	}
	if !status.DurableAfterRestart {
		t.Fatal("file persistence should report restart durability")
	}
	if status.CrossDevice || status.ProductionReady {
		t.Fatalf("file status overclaimed distributed or production readiness: %+v", status)
	}
}

func TestTypingPrivacyDurabilityStatusLogLineIsBounded(t *testing.T) {
	status, err := TypingPrivacyDurabilityStatusFor(service.TypingPrivacyPersistenceFile)
	if err != nil {
		t.Fatalf("file durability status: %v", err)
	}

	line := status.LogLine()
	want := "typing_privacy_persistence=file typing_privacy_durability=single-node-durable typing_privacy_restart_durable=true typing_privacy_cross_device=false typing_privacy_production_ready=false"
	if line != want {
		t.Fatalf("unexpected log line:\n got: %q\nwant: %q", line, want)
	}
	for _, forbidden := range []string{"/", "\\", "root", "preference", "user"} {
		if strings.Contains(strings.ToLower(line), forbidden) {
			t.Fatalf("diagnostic leaked forbidden detail %q: %q", forbidden, line)
		}
	}
}

func TestTypingPrivacyDurabilityStatusRejectsUnknownMode(t *testing.T) {
	if _, err := TypingPrivacyDurabilityStatusFor(service.TypingPrivacyPersistenceMode("cluster")); err == nil {
		t.Fatal("expected unsupported persistence mode to fail")
	}
}
