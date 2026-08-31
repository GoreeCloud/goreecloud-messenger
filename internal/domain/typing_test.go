// SPDX-License-Identifier: AGPL-3.0-only

package domain

import "testing"

func TestTypingSignalValidate(t *testing.T) {
	valid := TypingSignal{ConversationID: "conversation-1", UserID: "user-1", Sequence: 1, State: TypingStateTyping}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid typing signal rejected: %v", err)
	}

	invalidSequence := valid
	invalidSequence.Sequence = 0
	if err := invalidSequence.Validate(); err == nil {
		t.Fatal("zero typing sequence accepted")
	}

	invalidState := valid
	invalidState.State = TypingState("drafting")
	if err := invalidState.Validate(); err == nil {
		t.Fatal("unsupported typing state accepted")
	}
}
