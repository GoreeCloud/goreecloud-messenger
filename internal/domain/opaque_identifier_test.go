// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"strings"
	"testing"
)

func TestValidateOpaqueIdentifierPreservesExactInternalContent(t *testing.T) {
	for _, value := range []string{
		"conversation-1",
		"conversation:one/two",
		"conversation one two",
		"message@device#1",
		"scope-😀",
	} {
		if err := ValidateOpaqueIdentifier(value, "identifier"); err != nil {
			t.Fatalf("ValidateOpaqueIdentifier(%q) = %v", value, err)
		}
	}
}

func TestValidateOpaqueIdentifierRejectsNoncanonicalOrUnsafeValues(t *testing.T) {
	for _, value := range []string{
		"",
		" ",
		" conversation-1",
		"conversation-1 ",
		"conversation\n1",
		"conversation\x7f1",
		string([]byte{0xff}),
		strings.Repeat("a", MaxOpaqueIdentifierUTF16Units+1),
	} {
		if err := ValidateOpaqueIdentifier(value, "identifier"); err == nil {
			t.Fatalf("ValidateOpaqueIdentifier(%q) unexpectedly succeeded", value)
		}
	}
}

func TestValidateOpaqueIdentifierCountsUTF16CodeUnits(t *testing.T) {
	// U+1F600 occupies two UTF-16 code units, matching Kotlin String.length.
	withinLimit := strings.Repeat("😀", MaxOpaqueIdentifierUTF16Units/2)
	if err := ValidateOpaqueIdentifier(withinLimit, "identifier"); err != nil {
		t.Fatalf("within-limit supplementary identifier rejected: %v", err)
	}
	tooLong := withinLimit + "a"
	if err := ValidateOpaqueIdentifier(tooLong, "identifier"); err == nil {
		t.Fatal("identifier exceeding the UTF-16 limit unexpectedly succeeded")
	}
}
