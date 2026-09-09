// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"fmt"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// MaxOpaqueIdentifierUTF16Units matches the bounded opaque identifier limit used by
// the native Android Messenger client. Authority-owned identifiers are compared
// exactly after validation; this package never trims or otherwise rewrites them.
const MaxOpaqueIdentifierUTF16Units = 512

// ValidateOpaqueIdentifier enforces the shared Development identifier boundary for
// Data message, conversation, receipt, attachment, participant, typing, and replay scopes.
// Internal whitespace and punctuation remain valid opaque content. Invalid UTF-8,
// boundary whitespace, C0/DEL controls, blank values, and oversized identifiers fail closed.
func ValidateOpaqueIdentifier(value, label string) error {
	if value == "" {
		return fmt.Errorf("%s is required", label)
	}
	if !utf8.ValidString(value) {
		return fmt.Errorf("%s must contain valid Unicode", label)
	}
	if strings.TrimSpace(value) != value {
		return fmt.Errorf("%s must already be canonical", label)
	}
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", label)
	}
	if len(utf16.Encode([]rune(value))) > MaxOpaqueIdentifierUTF16Units {
		return fmt.Errorf("%s must not exceed %d UTF-16 code units", label, MaxOpaqueIdentifierUTF16Units)
	}
	for _, character := range value {
		if character >= 0x00 && character <= 0x1f || character == 0x7f {
			return fmt.Errorf("%s must not contain control characters", label)
		}
	}
	return nil
}
