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

// LogLine returns only bounded categorical durability truth. It deliberately
// excludes configured storage paths, preference contents, identities, and any
// other user-derived value.
func (status TypingPrivacyDurabilityStatus) LogLine() string {
	return fmt.Sprintf(
		"typing_privacy_persistence=%s typing_privacy_durability=%s typing_privacy_restart_durable=%t typing_privacy_cross_device=%t typing_privacy_production_ready=%t",
		status.PersistenceMode,
		status.Level,
		status.DurableAfterRestart,
		status.CrossDevice,
		status.ProductionReady,
	)
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
