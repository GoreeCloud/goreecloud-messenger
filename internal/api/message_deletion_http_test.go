// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func newTestMessageDeletionHTTPHandler(t *testing.T, userID string) http.Handler {
	t.Helper()

	store := messagingservice.NewMemoryDataStore()
	access := messagingservice.NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-1",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-1", "user-2"},
	}); err != nil {
		t.Fatalf("set conversation: %v", err)
	}
	dataService, err := messagingservice.NewDataService(store, access)
	if err != nil {
		t.Fatalf("new Data service: %v", err)
	}
	if err := dataService.Submit(context.Background(), "user-1", domain.DataEnvelope{
		MessageID: "message-1", ConversationID: "conversation-1", SenderID: "user-1", ClientNonce: "nonce-1",
		Ciphertext: []byte("opaque ciphertext"), Encryption: domain.EncryptionE2EE,
		CreatedAt: time.Date(2026, 8, 29, 13, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("seed message: %v", err)
	}
	deletionService, err := messagingservice.NewMessageDeletionService(store, messagingservice.NewMemoryMessageDeletionStore(), access)
	if err != nil {
		t.Fatalf("new deletion service: %v", err)
	}
	handler, err := NewMessageDeletionHTTPHandler(deletionService, testAuthenticator{userID: userID})
	if err != nil {
		t.Fatalf("new deletion HTTP handler: %v", err)
	}
	return handler.Routes()
}

func TestMessageDeletionHTTPRoundTripMinimizesMetadata(t *testing.T) {
	handler := newTestMessageDeletionHTTPHandler(t, "user-1")
	body := `{"deletion_id":"deletion-1","conversation_id":"conversation-1","deleter_id":"user-1","client_nonce":"delete-nonce-1","scope":"everyone","deleted_at":"2026-08-29T14:00:00Z"}`

	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/deletions", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("deletion status = %d, body = %s", recorder.Code, recorder.Body.String())
	}

	list := httptest.NewRequest(http.MethodGet, "/v1/data/messages/message-1/deletions", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, list)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", listRecorder.Code, listRecorder.Body.String())
	}
	response := listRecorder.Body.String()
	if !strings.Contains(response, `"deletion_id":"deletion-1"`) || !strings.Contains(response, `"scope":"everyone"`) {
		t.Fatalf("response missing deletion tombstone: %s", response)
	}
	if strings.Contains(response, "client_nonce") || strings.Contains(response, "delete-nonce-1") {
		t.Fatalf("response leaked replay nonce: %s", response)
	}
	if strings.Contains(response, "ciphertext") || strings.Contains(response, "opaque ciphertext") {
		t.Fatalf("response exposed message content state: %s", response)
	}
}

func TestMessageDeletionHTTPRejectsNonSender(t *testing.T) {
	handler := newTestMessageDeletionHTTPHandler(t, "user-2")
	body := `{"deletion_id":"deletion-1","conversation_id":"conversation-1","deleter_id":"user-2","client_nonce":"delete-nonce-1","scope":"everyone","deleted_at":"2026-08-29T14:00:00Z"}`
	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/deletions", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestMessageDeletionHTTPRejectsUnknownFields(t *testing.T) {
	handler := newTestMessageDeletionHTTPHandler(t, "user-1")
	body := `{"deletion_id":"deletion-1","conversation_id":"conversation-1","deleter_id":"user-1","client_nonce":"delete-nonce-1","scope":"everyone","deleted_at":"2026-08-29T14:00:00Z","plaintext":"must-not-be-accepted"}`
	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/deletions", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestMessageDeletionHTTPRejectsMissingAuthentication(t *testing.T) {
	handler := newTestMessageDeletionHTTPHandler(t, "")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/messages/message-1/deletions", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
