// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"testing"
	"time"
)

func TestDataEnvelopeValidationAllowsEncryptedPayload(t *testing.T) {
	envelope := DataEnvelope{
		MessageID:      "msg-1",
		ConversationID: "conv-1",
		SenderID:       "user-1",
		ClientNonce:    "nonce-1",
		Ciphertext:     []byte("ciphertext"),
		Encryption:     EncryptionE2EE,
		CreatedAt:      time.Now().UTC(),
	}

	if err := envelope.Validate(); err != nil {
		t.Fatalf("expected valid Data envelope, got %v", err)
	}
}

func TestDataEnvelopeValidationRejectsUnencryptedPayload(t *testing.T) {
	envelope := DataEnvelope{
		MessageID:      "msg-2",
		ConversationID: "conv-1",
		SenderID:       "user-1",
		ClientNonce:    "nonce-2",
		Ciphertext:     []byte("plaintext-like-data"),
		Encryption:     EncryptionNone,
		CreatedAt:      time.Now().UTC(),
	}

	if err := envelope.Validate(); err == nil {
		t.Fatal("expected Data envelope without E2EE to fail")
	}
}
