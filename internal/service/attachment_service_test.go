// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func testAttachment() domain.DataAttachment {
	return domain.DataAttachment{
		AttachmentID:   "attachment-1",
		ConversationID: "conversation-1",
		SenderID:       "user-a",
		ClientNonce:    "nonce-1",
		Filename:       "photo.jpg",
		MIMEType:       "image/jpeg",
		Ciphertext:     []byte{1, 2, 3, 4},
	}
}

func testAttachmentService(t *testing.T) (*AttachmentService, *MemoryAttachmentStore) {
	t.Helper()
	store := NewMemoryAttachmentStore()
	access := NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-1",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-a", "user-b"},
	}); err != nil {
		t.Fatalf("set conversation: %v", err)
	}
	service, err := NewAttachmentService(store, access)
	if err != nil {
		t.Fatalf("new attachment service: %v", err)
	}
	return service, store
}

func TestAttachmentServiceSubmitAndGet(t *testing.T) {
	service, _ := testAttachmentService(t)
	attachment := testAttachment()
	if err := service.Submit(context.Background(), "user-a", attachment); err != nil {
		t.Fatalf("submit attachment: %v", err)
	}

	got, err := service.Get(context.Background(), "user-b", attachment.AttachmentID)
	if err != nil {
		t.Fatalf("get attachment: %v", err)
	}
	if got.Filename != attachment.Filename || got.MIMEType != attachment.MIMEType {
		t.Fatalf("unexpected attachment metadata: %#v", got)
	}
	got.Ciphertext[0] = 99
	again, err := service.Get(context.Background(), "user-a", attachment.AttachmentID)
	if err != nil {
		t.Fatalf("get attachment again: %v", err)
	}
	if again.Ciphertext[0] != 1 {
		t.Fatal("ciphertext mutation leaked into stored attachment")
	}
}

func TestAttachmentServiceRejectsSenderMismatch(t *testing.T) {
	service, _ := testAttachmentService(t)
	err := service.Submit(context.Background(), "user-b", testAttachment())
	if !errors.Is(err, ErrSenderMismatch) {
		t.Fatalf("expected sender mismatch, got %v", err)
	}
}

func TestAttachmentServiceRejectsConversationOutsider(t *testing.T) {
	service, _ := testAttachmentService(t)
	attachment := testAttachment()
	attachment.SenderID = "user-c"
	err := service.Submit(context.Background(), "user-c", attachment)
	if !errors.Is(err, ErrConversationAccess) {
		t.Fatalf("expected conversation access error, got %v", err)
	}
}

func TestAttachmentServiceRejectsDuplicateAndNonceReuse(t *testing.T) {
	service, _ := testAttachmentService(t)
	attachment := testAttachment()
	if err := service.Submit(context.Background(), "user-a", attachment); err != nil {
		t.Fatalf("submit attachment: %v", err)
	}
	if err := service.Submit(context.Background(), "user-a", attachment); !errors.Is(err, ErrDuplicateAttachment) {
		t.Fatalf("expected duplicate attachment, got %v", err)
	}

	second := attachment
	second.AttachmentID = "attachment-2"
	if err := service.Submit(context.Background(), "user-a", second); !errors.Is(err, ErrAttachmentNonceReuse) {
		t.Fatalf("expected nonce reuse, got %v", err)
	}
}

func TestAttachmentValidationRejectsEmptyCiphertext(t *testing.T) {
	attachment := testAttachment()
	attachment.Ciphertext = nil
	if err := attachment.Validate(); err == nil {
		t.Fatal("expected empty ciphertext validation error")
	}
}
