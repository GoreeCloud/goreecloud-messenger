// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
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

func TestTypingPrivacyDurabilityStatusRejectsUnknownMode(t *testing.T) {
	if _, err := TypingPrivacyDurabilityStatusFor(service.TypingPrivacyPersistenceMode("cluster")); err == nil {
		t.Fatal("expected unsupported persistence mode to fail")
	}
}
