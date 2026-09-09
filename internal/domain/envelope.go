// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"time"
)

// DataEnvelope is the GoreeCloud-controlled wire-neutral unit accepted by the Data service.
// It intentionally carries no carrier transport fields because Data messages never silently
// become SMS, MMS, or RCS at this boundary.
type DataEnvelope struct {
	MessageID      string
	ConversationID string
	SenderID       string
	ClientNonce    string
	Ciphertext     []byte
	Encryption     EncryptionState
	CreatedAt      time.Time
}

func (e DataEnvelope) Validate() error {
	if err := ValidateOpaqueIdentifier(e.MessageID, "envelope message id"); err != nil {
		return err
	}
	if err := ValidateOpaqueIdentifier(e.ConversationID, "envelope conversation id"); err != nil {
		return err
	}
	if err := ValidateOpaqueIdentifier(e.SenderID, "envelope sender id"); err != nil {
		return err
	}
	if err := ValidateOpaqueIdentifier(e.ClientNonce, "envelope client nonce"); err != nil {
		return err
	}
	if len(e.Ciphertext) == 0 {
		return errors.New("envelope ciphertext is required")
	}
	if e.Encryption != EncryptionE2EE {
		return errors.New("GoreeCloud Data envelope requires active E2EE")
	}
	if e.CreatedAt.IsZero() {
		return errors.New("envelope creation timestamp is required")
	}
	return nil
}
