// SPDX-License-Identifier: AGPL-3.0-only

package domain

import "testing"

func TestCallValidationAllowsEncryptedVideoData(t *testing.T) {
	call := CallSession{
		ID:             "call-1",
		ConversationID: "conv-1",
		Kind:           CallVideo,
		Transport:      TransportData,
		Encryption:     EncryptionE2EE,
	}

	if err := call.Validate(); err != nil {
		t.Fatalf("expected valid encrypted video call, got %v", err)
	}
	if got := call.ProvenanceLabel(); got != "Video · E2EE · Data" {
		t.Fatalf("unexpected call provenance label %q", got)
	}
}

func TestCallValidationRejectsCarrierTransport(t *testing.T) {
	call := CallSession{
		ID:             "call-2",
		ConversationID: "conv-1",
		Kind:           CallVoice,
		Transport:      TransportRCS,
		Encryption:     EncryptionUnknown,
	}

	if err := call.Validate(); err == nil {
		t.Fatal("expected GoreeCloud call on RCS transport to fail")
	}
}
