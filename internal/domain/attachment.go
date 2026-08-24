// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"fmt"
	"strings"
)

const (
	MaxAttachmentCiphertextBytes = 25 * 1024 * 1024
	MaxAttachmentFilenameBytes   = 255
	MaxAttachmentMIMETypeBytes   = 127
)

// DataAttachment describes encrypted attachment bytes carried by GoreeCloud Data.
// Ciphertext is opaque to the service and must never be treated as plaintext content.
type DataAttachment struct {
	AttachmentID  string
	ConversationID string
	SenderID      string
	ClientNonce   string
	Filename      string
	MIMEType      string
	Ciphertext    []byte
}

func (a DataAttachment) Validate() error {
	if strings.TrimSpace(a.AttachmentID) == "" {
		return errors.New("attachment id is required")
	}
	if strings.TrimSpace(a.ConversationID) == "" {
		return errors.New("conversation id is required")
	}
	if strings.TrimSpace(a.SenderID) == "" {
		return errors.New("sender id is required")
	}
	if strings.TrimSpace(a.ClientNonce) == "" {
		return errors.New("client nonce is required")
	}
	if len(a.Filename) == 0 || len(a.Filename) > MaxAttachmentFilenameBytes {
		return fmt.Errorf("filename length must be between 1 and %d bytes", MaxAttachmentFilenameBytes)
	}
	if len(a.MIMEType) == 0 || len(a.MIMEType) > MaxAttachmentMIMETypeBytes {
		return fmt.Errorf("MIME type length must be between 1 and %d bytes", MaxAttachmentMIMETypeBytes)
	}
	if len(a.Ciphertext) == 0 {
		return errors.New("attachment ciphertext is required")
	}
	if len(a.Ciphertext) > MaxAttachmentCiphertextBytes {
		return fmt.Errorf("attachment ciphertext exceeds %d bytes", MaxAttachmentCiphertextBytes)
	}
	return nil
}
