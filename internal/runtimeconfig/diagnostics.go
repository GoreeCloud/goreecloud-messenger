// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
	"errors"
	"fmt"

	"github.com/GoreeCloud/goreecloud-messenger/internal/api"
)

const runtimeConfigurationSource = "environment"

type ReceiptPersistenceDiagnostic struct {
	Mode                string
	Durability          string
	ConfigurationSource string
}

// ReceiptPersistenceDiagnosticFor returns a deliberately minimized operational
// projection. It never contains the configured durable filesystem root.
func ReceiptPersistenceDiagnosticFor(config api.ReceiptPersistenceConfig) (ReceiptPersistenceDiagnostic, error) {
	switch config.Mode {
	case api.ReceiptPersistenceMemory:
		if config.Root != "" {
			return ReceiptPersistenceDiagnostic{}, errors.New("memory receipt persistence must not carry a durable root")
		}
		return ReceiptPersistenceDiagnostic{
			Mode:                string(api.ReceiptPersistenceMemory),
			Durability:          "process-local",
			ConfigurationSource: runtimeConfigurationSource,
		}, nil
	case api.ReceiptPersistenceFile:
		if config.Root == "" {
			return ReceiptPersistenceDiagnostic{}, errors.New("file receipt persistence requires configured durable storage")
		}
		return ReceiptPersistenceDiagnostic{
			Mode:                string(api.ReceiptPersistenceFile),
			Durability:          "single-node-durable",
			ConfigurationSource: runtimeConfigurationSource,
		}, nil
	default:
		return ReceiptPersistenceDiagnostic{}, fmt.Errorf("unsupported receipt persistence mode")
	}
}

func (diagnostic ReceiptPersistenceDiagnostic) LogLine() string {
	return fmt.Sprintf(
		"receipt_persistence=%s receipt_durability=%s receipt_config_source=%s",
		diagnostic.Mode,
		diagnostic.Durability,
		diagnostic.ConfigurationSource,
	)
}
