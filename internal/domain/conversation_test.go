// SPDX-License-Identifier: AGPL-3.0-only

package domain

import "testing"

func TestConversationValidationAllowsGroup(t *testing.T) {
	conversation := Conversation{
		ID:             "conv-1",
		Kind:           ConversationGroup,
		ParticipantIDs: []string{"user-1", "user-2", "user-3"},
	}

	if err := conversation.Validate(); err != nil {
		t.Fatalf("expected valid group conversation, got %v", err)
	}
}

func TestConversationValidationRejectsDuplicateParticipants(t *testing.T) {
	conversation := Conversation{
		ID:             "conv-2",
		Kind:           ConversationDirect,
		ParticipantIDs: []string{"user-1", "user-1"},
	}

	if err := conversation.Validate(); err == nil {
		t.Fatal("expected duplicate participant ids to fail")
	}
}

func TestTransportTransitionRequiresConfirmationWhenLeavingData(t *testing.T) {
	transition := TransportTransition{
		From: TransportData,
		To:   TransportSMS,
	}

	if err := transition.Validate(); err == nil {
		t.Fatal("expected unconfirmed Data-to-SMS transition to fail")
	}

	transition.UserConfirmed = true
	if err := transition.Validate(); err != nil {
		t.Fatalf("expected confirmed transition to succeed, got %v", err)
	}
}
