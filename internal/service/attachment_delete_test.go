// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"
)

func TestAttachmentDeleteIsIdempotentAndPreservesReplayReservations(t *testing.T) {
	service, _ := testAttachmentService(t)
	attachment := testAttachment()
	if err := service.Submit(context.Background(), "user-a", attachment); err != nil {
		t.Fatalf("submit attachment: %v", err)
	}

	if err := service.Delete(context.Background(), "user-b", attachment.AttachmentID); err != nil {
		t.Fatalf("delete attachment as participant: %v", err)
	}
	if _, err := service.Get(context.Background(), "user-a", attachment.AttachmentID); !errors.Is(err, ErrAttachmentNotFound) {
		t.Fatalf("expected deleted attachment to be unavailable, got %v", err)
	}
	listed, err := service.List(context.Background(), "user-a", attachment.ConversationID, 10)
	if err != nil {
		t.Fatalf("list attachments: %v", err)
	}
	if len(listed) != 0 {
		t.Fatalf("deleted attachment remained discoverable: %#v", listed)
	}
	if err := service.Delete(context.Background(), "user-a", attachment.AttachmentID); err != nil {
		t.Fatalf("repeat delete should be idempotent: %v", err)
	}

	reusedID := attachment
	reusedID.ClientNonce = "nonce-2"
	if err := service.Submit(context.Background(), "user-a", reusedID); !errors.Is(err, ErrDuplicateAttachment) {
		t.Fatalf("expected deleted attachment id to stay reserved, got %v", err)
	}

	reusedNonce := attachment
	reusedNonce.AttachmentID = "attachment-2"
	if err := service.Submit(context.Background(), "user-a", reusedNonce); !errors.Is(err, ErrAttachmentNonceReuse) {
		t.Fatalf("expected deleted attachment nonce to stay reserved, got %v", err)
	}
}

func TestAttachmentDeleteRejectsConversationOutsider(t *testing.T) {
	service, _ := testAttachmentService(t)
	attachment := testAttachment()
	if err := service.Submit(context.Background(), "user-a", attachment); err != nil {
		t.Fatalf("submit attachment: %v", err)
	}
	if err := service.Delete(context.Background(), "user-c", attachment.AttachmentID); !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("expected conversation access denial, got %v", err)
	}
	if _, err := service.Get(context.Background(), "user-a", attachment.AttachmentID); err != nil {
		t.Fatalf("unauthorized delete changed attachment state: %v", err)
	}
}
