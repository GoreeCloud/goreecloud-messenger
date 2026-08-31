// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
	"fmt"

	"github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

type TypingPrivacyDurabilityLevel string

const (
	TypingPrivacyDurabilityTransient         TypingPrivacyDurabilityLevel = "transient"
	TypingPrivacyDurabilitySingleNodeDurable TypingPrivacyDurabilityLevel = "single-node-durable"
)

// TypingPrivacyDurabilityStatus is a minimized runtime projection. It reports
// truthful persistence properties without exposing the configured file root or
// implying cross-device or production Privacy Shield acceptance.
type TypingPrivacyDurabilityStatus struct {
	PersistenceMode     service.TypingPrivacyPersistenceMode
	Level               TypingPrivacyDurabilityLevel
	DurableAfterRestart bool
	CrossDevice         bool
	ProductionReady     bool
}

func TypingPrivacyDurabilityStatusFor(
	mode service.TypingPrivacyPersistenceMode,
) (TypingPrivacyDurabilityStatus, error) {
	switch mode {
	case service.TypingPrivacyPersistenceMemory:
		return TypingPrivacyDurabilityStatus{
			PersistenceMode:     mode,
			Level:               TypingPrivacyDurabilityTransient,
			DurableAfterRestart: false,
			CrossDevice:         false,
			ProductionReady:     false,
		}, nil
	case service.TypingPrivacyPersistenceFile:
		return TypingPrivacyDurabilityStatus{
			PersistenceMode:     mode,
			Level:               TypingPrivacyDurabilitySingleNodeDurable,
			DurableAfterRestart: true,
			CrossDevice:         false,
			ProductionReady:     false,
		}, nil
	default:
		return TypingPrivacyDurabilityStatus{}, fmt.Errorf(
			"unsupported typing privacy persistence mode %q",
			mode,
		)
	}
}
