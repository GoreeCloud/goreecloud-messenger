// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func TestAttachmentServiceListReturnsBoundedMetadataOnly(t *testing.T) {
	service, _ := testAttachmentService(t)
	first := testAttachment()
	first.AttachmentID = "attachment-b"
	first.ClientNonce = "nonce-b"
	first.Ciphertext = []byte{1, 2, 3}
	second := testAttachment()
	second.AttachmentID = "attachment-a"
	second.ClientNonce = "nonce-a"
	second.Ciphertext = []byte{4, 5}

	for _, attachment := range []domain.DataAttachment{first, second} {
		if err := service.Submit(context.Background(), "user-a", attachment); err != nil {
			t.Fatalf("submit attachment %s: %v", attachment.AttachmentID, err)
		}
	}

	items, err := service.List(context.Background(), "user-b", "conversation-1", 1)
	if err != nil {
		t.Fatalf("list attachments: %v", err)
	}
	if len(items) != 1 || items[0].AttachmentID != "attachment-a" {
		t.Fatalf("unexpected bounded deterministic listing: %#v", items)
	}
	if items[0].CiphertextBytes != 2 {
		t.Fatalf("ciphertext byte count = %d, want 2", items[0].CiphertextBytes)
	}
}

func TestAttachmentServiceListRejectsOutsiderBeforeStoreRead(t *testing.T) {
	service, _ := testAttachmentService(t)
	_, err := service.List(context.Background(), "user-c", "conversation-1", 10)
	if !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("expected conversation access error, got %v", err)
	}
}

func TestAttachmentServiceListRejectsUnboundedLimit(t *testing.T) {
	service, _ := testAttachmentService(t)
	_, err := service.List(context.Background(), "user-a", "conversation-1", MaxAttachmentListResults+1)
	if !errors.Is(err, ErrAttachmentListLimit) {
		t.Fatalf("expected list limit error, got %v", err)
	}
}
