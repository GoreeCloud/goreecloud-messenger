// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"testing"
	"time"
)

func TestMessageValidationAllowsE2EEData(t *testing.T) {
	m := Message{
		ID:             "msg-1",
		ConversationID: "conv-1",
		SenderID:       "user-1",
		Body:           "hello",
		Transport:      TransportData,
		Encryption:     EncryptionE2EE,
		Delivery:       DeliverySent,
		SentAt:         time.Now(),
	}

	if err := m.Validate(); err != nil {
		t.Fatalf("expected valid message, got %v", err)
	}
	if got := m.ProvenanceLabel(); got != "E2EE · Data" {
		t.Fatalf("unexpected provenance label %q", got)
	}
}

func TestMessageValidationRejectsFalseE2EEOnSMS(t *testing.T) {
	m := Message{
		ID:             "msg-2",
		ConversationID: "conv-1",
		SenderID:       "user-1",
		Body:           "fallback",
		Transport:      TransportSMS,
		Encryption:     EncryptionE2EE,
		Delivery:       DeliveryPending,
		SentAt:         time.Now(),
	}

	if err := m.Validate(); err == nil {
		t.Fatal("expected E2EE assertion on SMS to be rejected")
	}
}

func TestMessageValidationRejectsUnknownDeliveryState(t *testing.T) {
	m := Message{
		ID:             "msg-3",
		ConversationID: "conv-1",
		SenderID:       "user-1",
		Body:           "hello",
		Transport:      TransportData,
		Encryption:     EncryptionNone,
		Delivery:       DeliveryState("lost"),
		SentAt:         time.Now(),
	}

	if err := m.Validate(); err == nil {
		t.Fatal("expected unsupported delivery state to fail")
	}
}

func TestIdentityRequiresUsernamePrefix(t *testing.T) {
	identity := Identity{UserID: "user-1", Username: "alex"}
	if err := identity.Validate(); err == nil {
		t.Fatal("expected username without @ prefix to fail")
	}
}

func TestIdentityRejectsBareAtSign(t *testing.T) {
	identity := Identity{UserID: "user-1", Username: "@"}
	if err := identity.Validate(); err == nil {
		t.Fatal("expected bare @ username to fail")
	}
}
