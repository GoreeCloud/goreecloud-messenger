// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"strings"
	"time"
)

// MessageEdit is an immutable encrypted revision of an existing GoreeCloud Data message.
// The server transports ciphertext and revision metadata only; it does not interpret plaintext.
type MessageEdit struct {
	EditID         string
	MessageID      string
	ConversationID string
	EditorID       string
	ClientNonce    string
	Ciphertext     []byte
	Encryption     EncryptionState
	EditedAt       time.Time
}

func (e MessageEdit) Validate() error {
	if strings.TrimSpace(e.EditID) == "" {
		return errors.New("edit id is required")
	}
	if strings.TrimSpace(e.MessageID) == "" {
		return errors.New("edit message id is required")
	}
	if strings.TrimSpace(e.ConversationID) == "" {
		return errors.New("edit conversation id is required")
	}
	if strings.TrimSpace(e.EditorID) == "" {
		return errors.New("edit editor id is required")
	}
	if strings.TrimSpace(e.ClientNonce) == "" {
		return errors.New("edit client nonce is required")
	}
	if len(e.Ciphertext) == 0 {
		return errors.New("edit ciphertext is required")
	}
	if e.Encryption != EncryptionE2EE {
		return errors.New("GoreeCloud Data message edit requires active E2EE")
	}
	if e.EditedAt.IsZero() {
		return errors.New("edit timestamp is required")
	}
	return nil
}
