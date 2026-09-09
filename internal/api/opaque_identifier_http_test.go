// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestDataHTTPRejectsBoundaryWhitespaceConversationPathInsteadOfNormalizing(t *testing.T) {
	handler := newTestHandler(t, "user-1")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/%20conversation-1/messages", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestDataHTTPRejectsBoundaryWhitespaceReceiptMessagePathInsteadOfNormalizing(t *testing.T) {
	handler := newTestHandlerWithSeededMessage(t, "user-2")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/messages/%20message-1/receipts", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestDataHTTPRejectsNoncanonicalAuthenticatedIdentity(t *testing.T) {
	handler := newTestHandler(t, "user-1 ")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/conversation-1/messages", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestAttachmentHTTPRejectsBoundaryWhitespaceConversationPathInsteadOfNormalizing(t *testing.T) {
	handler := newTestAttachmentHandler(t, "user-1")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/%20conversation-1/attachments", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestAttachmentHTTPRejectsBoundaryWhitespaceAttachmentPathInsteadOfNormalizing(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-1",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-1", "user-2"},
	}); err != nil {
		t.Fatalf("set conversation: %v", err)
	}
	service, err := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)
	if err != nil {
		t.Fatalf("new attachment service: %v", err)
	}
	if err := service.Submit(t.Context(), "user-1", domain.DataAttachment{
		AttachmentID:   "attachment-1",
		ConversationID: "conversation-1",
		SenderID:       "user-1",
		ClientNonce:    "nonce-1",
		Filename:       "photo.jpg",
		MIMEType:       "image/jpeg",
		Ciphertext:     []byte("ciphertext"),
	}); err != nil {
		t.Fatalf("seed attachment: %v", err)
	}
	handler, err := NewAttachmentHTTPHandler(service, testAuthenticator{userID: "user-1"})
	if err != nil {
		t.Fatalf("new attachment handler: %v", err)
	}
	request := httptest.NewRequest(http.MethodGet, "/v1/data/attachments/%20attachment-1", nil)
	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
