// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestAttachmentCiphertextDownloadReturnsExactOpaqueBytes(t *testing.T) {
	handler := newTestAttachmentHandler(t, "user-1")
	ciphertext := []byte{0x00, 0x01, 0x7f, 0x80, 0xff, 'G', 'C'}
	body := `{"attachment_id":"attachment-raw-1","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"attachment-raw-nonce-1","filename":"private-photo.jpg","mime_type":"image/jpeg","ciphertext":"` + base64.StdEncoding.EncodeToString(ciphertext) + `"}`

	submit := httptest.NewRequest(http.MethodPost, "/v1/data/attachments", strings.NewReader(body))
	submitRecorder := httptest.NewRecorder()
	handler.ServeHTTP(submitRecorder, submit)
	if submitRecorder.Code != http.StatusAccepted {
		t.Fatalf("submit status = %d, body = %s", submitRecorder.Code, submitRecorder.Body.String())
	}

	download := httptest.NewRequest(http.MethodGet, "/v1/data/attachments/attachment-raw-1/ciphertext", nil)
	downloadRecorder := httptest.NewRecorder()
	handler.ServeHTTP(downloadRecorder, download)
	if downloadRecorder.Code != http.StatusOK {
		t.Fatalf("download status = %d, body = %s", downloadRecorder.Code, downloadRecorder.Body.String())
	}
	if got := downloadRecorder.Body.Bytes(); string(got) != string(ciphertext) {
		t.Fatalf("ciphertext bytes changed: got %v want %v", got, ciphertext)
	}
	if got := downloadRecorder.Header().Get("Content-Type"); got != "application/octet-stream" {
		t.Fatalf("content type = %q", got)
	}
	if got := downloadRecorder.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("cache control = %q", got)
	}
	if got := downloadRecorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("nosniff = %q", got)
	}
	if got := downloadRecorder.Header().Get("Content-Length"); got != "7" {
		t.Fatalf("content length = %q", got)
	}
	for name, values := range downloadRecorder.Header() {
		for _, value := range values {
			if strings.Contains(value, "private-photo.jpg") || strings.Contains(value, "image/jpeg") {
				t.Fatalf("raw ciphertext header %q leaked plaintext attachment metadata: %q", name, value)
			}
		}
	}
}

func TestAttachmentCiphertextDownloadRequiresAuthentication(t *testing.T) {
	handler := newTestAttachmentHandler(t, "")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/attachments/attachment-1/ciphertext", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestAttachmentCiphertextDownloadRejectsConversationOutsider(t *testing.T) {
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
		AttachmentID:   "attachment-raw-1",
		ConversationID: "conversation-1",
		SenderID:       "user-1",
		ClientNonce:    "attachment-raw-nonce-1",
		Filename:       "private-photo.jpg",
		MIMEType:       "image/jpeg",
		Ciphertext:     []byte("opaque-ciphertext"),
	}); err != nil {
		t.Fatalf("seed attachment: %v", err)
	}
	handler, err := NewAttachmentHTTPHandler(service, testAuthenticator{userID: "user-3"})
	if err != nil {
		t.Fatalf("new attachment handler: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/v1/data/attachments/attachment-raw-1/ciphertext", nil)
	recorder := httptest.NewRecorder()
	handler.Routes().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
