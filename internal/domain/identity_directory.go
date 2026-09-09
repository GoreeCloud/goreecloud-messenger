// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	MaxIdentityHandleLength      = 32
	MaxIdentitySubjectLength     = 255
	MaxIdentityDisplayNameLength = 160
)

var identityHandlePattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9._-]{0,30}[a-z0-9])?$`)

// IdentityDirectoryProjection is the minimized GoreeCloud Identity projection
// Messenger may consume after an authorized exact-handle resolution succeeds.
// It intentionally excludes email addresses, phone numbers, discovery policy,
// credentials, sessions, roles, and all application-specific authorization.
type IdentityDirectoryProjection struct {
	Subject     string
	Handle      string
	DisplayName string
}

// ValidateIdentityDirectoryProjection validates only the projection returned by
// the Identity-owned directory authority. Messenger never canonicalizes an
// Identity handle into a different value; the resolver must return Identity's
// exact canonical lowercase handle.
func ValidateIdentityDirectoryProjection(value IdentityDirectoryProjection) error {
	if err := ValidateOpaqueIdentifier(value.Subject, "identity subject"); err != nil {
		return err
	}
	if len([]rune(value.Subject)) > MaxIdentitySubjectLength {
		return fmt.Errorf("identity subject must not exceed %d characters", MaxIdentitySubjectLength)
	}
	if !identityHandlePattern.MatchString(value.Handle) || len(value.Handle) > MaxIdentityHandleLength {
		return fmt.Errorf("identity handle is not canonical")
	}
	if !utf8.ValidString(value.DisplayName) {
		return fmt.Errorf("identity display name must contain valid Unicode")
	}
	if strings.TrimSpace(value.DisplayName) != value.DisplayName {
		return fmt.Errorf("identity display name must already be canonical")
	}
	if len([]rune(value.DisplayName)) > MaxIdentityDisplayNameLength {
		return fmt.Errorf("identity display name must not exceed %d characters", MaxIdentityDisplayNameLength)
	}
	for _, character := range value.DisplayName {
		if character >= 0x00 && character <= 0x1f || character == 0x7f {
			return fmt.Errorf("identity display name must not contain control characters")
		}
	}
	return nil
}
