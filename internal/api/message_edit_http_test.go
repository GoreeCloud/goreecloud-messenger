// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func newTestMessageEditHTTPHandler(t *testing.T, userID string) http.Handler {
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
		MessageID:      "message-1",
		ConversationID: "conversation-1",
		SenderID:       "user-1",
		ClientNonce:    "nonce-1",
		Ciphertext:     []byte("original ciphertext"),
		Encryption:     domain.EncryptionE2EE,
		CreatedAt:      time.Date(2026, 8, 29, 13, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("seed message: %v", err)
	}

	editService, err := messagingservice.NewMessageEditService(store, messagingservice.NewMemoryMessageEditStore(), access)
	if err != nil {
		t.Fatalf("new message edit service: %v", err)
	}
	handler, err := NewMessageEditHTTPHandler(editService, testAuthenticator{userID: userID})
	if err != nil {
		t.Fatalf("new message edit HTTP handler: %v", err)
	}
	return handler.Routes()
}

func TestMessageEditHTTPRoundTripPreservesCiphertext(t *testing.T) {
	handler := newTestMessageEditHTTPHandler(t, "user-1")
	ciphertext := base64.StdEncoding.EncodeToString([]byte("opaque edited ciphertext"))
	body := `{"edit_id":"edit-1","conversation_id":"conversation-1","editor_id":"user-1","client_nonce":"edit-nonce-1","ciphertext":"` + ciphertext + `","encryption":"e2ee","edited_at":"2026-08-29T13:30:00Z"}`

	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/edits", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("edit status = %d, body = %s", recorder.Code, recorder.Body.String())
	}

	handler = newTestMessageEditHTTPHandlerWithSeededEdit(t, "user-2")
	list := httptest.NewRequest(http.MethodGet, "/v1/data/messages/message-1/edits", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, list)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", listRecorder.Code, listRecorder.Body.String())
	}
	response := listRecorder.Body.String()
	if !strings.Contains(response, `"ciphertext":"`+ciphertext+`"`) {
		t.Fatalf("response did not preserve edit ciphertext: %s", response)
	}
	if strings.Contains(response, "opaque edited ciphertext") {
		t.Fatalf("response exposed decoded edit content: %s", response)
	}
}

func newTestMessageEditHTTPHandlerWithSeededEdit(t *testing.T, userID string) http.Handler {
	t.Helper()

	store := messagingservice.NewMemoryDataStore()
	access := messagingservice.NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{ID: "conversation-1", Kind: domain.ConversationDirect, ParticipantIDs: []string{"user-1", "user-2"}}); err != nil {
		t.Fatalf("set conversation: %v", err)
	}
	dataService, err := messagingservice.NewDataService(store, access)
	if err != nil {
		t.Fatalf("new Data service: %v", err)
	}
	if err := dataService.Submit(context.Background(), "user-1", domain.DataEnvelope{
		MessageID: "message-1", ConversationID: "conversation-1", SenderID: "user-1", ClientNonce: "nonce-1",
		Ciphertext: []byte("original ciphertext"), Encryption: domain.EncryptionE2EE,
		CreatedAt: time.Date(2026, 8, 29, 13, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("seed message: %v", err)
	}
	editStore := messagingservice.NewMemoryMessageEditStore()
	editService, err := messagingservice.NewMessageEditService(store, editStore, access)
	if err != nil {
		t.Fatalf("new message edit service: %v", err)
	}
	if err := editService.Record(context.Background(), "user-1", domain.MessageEdit{
		EditID: "edit-1", MessageID: "message-1", ConversationID: "conversation-1", EditorID: "user-1", ClientNonce: "edit-nonce-1",
		Ciphertext: []byte("opaque edited ciphertext"), Encryption: domain.EncryptionE2EE,
		EditedAt: time.Date(2026, 8, 29, 13, 30, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("seed edit: %v", err)
	}
	handler, err := NewMessageEditHTTPHandler(editService, testAuthenticator{userID: userID})
	if err != nil {
		t.Fatalf("new message edit HTTP handler: %v", err)
	}
	return handler.Routes()
}

func TestMessageEditHTTPRejectsNonSender(t *testing.T) {
	handler := newTestMessageEditHTTPHandler(t, "user-2")
	body := `{"edit_id":"edit-1","conversation_id":"conversation-1","editor_id":"user-2","client_nonce":"edit-nonce-1","ciphertext":"Y2lwaGVydGV4dA==","encryption":"e2ee","edited_at":"2026-08-29T13:30:00Z"}`

	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/edits", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestMessageEditHTTPRejectsUnknownFields(t *testing.T) {
	handler := newTestMessageEditHTTPHandler(t, "user-1")
	body := `{"edit_id":"edit-1","conversation_id":"conversation-1","editor_id":"user-1","client_nonce":"edit-nonce-1","ciphertext":"Y2lwaGVydGV4dA==","encryption":"e2ee","edited_at":"2026-08-29T13:30:00Z","plaintext":"must-not-be-accepted"}`

	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/edits", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestMessageEditHTTPRejectsMissingAuthentication(t *testing.T) {
	handler := newTestMessageEditHTTPHandler(t, "")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/messages/message-1/edits", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
