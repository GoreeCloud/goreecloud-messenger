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

func newReactionHTTPHandler(t *testing.T, userID string) http.Handler {
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
	data, err := messagingservice.NewDataService(store, access)
	if err != nil {
		t.Fatalf("new data service: %v", err)
	}
	if err := data.Submit(context.Background(), "user-1", domain.DataEnvelope{
		MessageID:      "message-1",
		ConversationID: "conversation-1",
		SenderID:       "user-1",
		ClientNonce:    "message-nonce-1",
		Ciphertext:     []byte("opaque message ciphertext"),
		Encryption:     domain.EncryptionE2EE,
		CreatedAt:      time.Date(2026, 8, 29, 14, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("seed message: %v", err)
	}
	reactions, err := messagingservice.NewMessageReactionService(store, messagingservice.NewMemoryMessageReactionStore(), access)
	if err != nil {
		t.Fatalf("new reaction service: %v", err)
	}
	handler, err := NewMessageReactionHTTPHandler(reactions, testAuthenticator{userID: userID})
	if err != nil {
		t.Fatalf("new reaction HTTP handler: %v", err)
	}
	return handler.Routes()
}

func TestMessageReactionRoundTripKeepsValueOpaqueAndMinimizesMetadata(t *testing.T) {
	handler := newReactionHTTPHandler(t, "user-2")
	plainReaction := "thumbs-up"
	ciphertext := base64.StdEncoding.EncodeToString([]byte(plainReaction))
	body := `{"reaction_id":"reaction-1","conversation_id":"conversation-1","reactor_id":"user-2","client_nonce":"reaction-nonce-1","operation":"set","ciphertext":"` + ciphertext + `","encryption":"e2ee","reacted_at":"2026-08-29T14:01:00Z"}`

	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/reactions", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("record status = %d, body = %s", recorder.Code, recorder.Body.String())
	}

	list := httptest.NewRequest(http.MethodGet, "/v1/data/messages/message-1/reactions", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, list)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", listRecorder.Code, listRecorder.Body.String())
	}
	response := listRecorder.Body.String()
	if !strings.Contains(response, `"ciphertext":"`+ciphertext+`"`) {
		t.Fatalf("response missing encrypted reaction: %s", response)
	}
	if strings.Contains(response, plainReaction) {
		t.Fatalf("response exposed decoded reaction value: %s", response)
	}
	if strings.Contains(response, "reaction-nonce-1") || strings.Contains(response, `"operation"`) {
		t.Fatalf("response exposed replay/control metadata not needed by active projection: %s", response)
	}
	if listRecorder.Header().Get("Cache-Control") != "no-store" || listRecorder.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("missing response protections: headers=%v", listRecorder.Header())
	}
}

func TestMessageReactionClearRemovesActiveProjection(t *testing.T) {
	handler := newReactionHTTPHandler(t, "user-2")
	setBody := `{"reaction_id":"reaction-1","conversation_id":"conversation-1","reactor_id":"user-2","client_nonce":"reaction-nonce-1","operation":"set","ciphertext":"ZW5jcnlwdGVk","encryption":"e2ee","reacted_at":"2026-08-29T14:01:00Z"}`
	setRequest := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/reactions", strings.NewReader(setBody))
	setRecorder := httptest.NewRecorder()
	handler.ServeHTTP(setRecorder, setRequest)
	if setRecorder.Code != http.StatusAccepted {
		t.Fatalf("set status = %d, body = %s", setRecorder.Code, setRecorder.Body.String())
	}

	clearBody := `{"reaction_id":"reaction-2","conversation_id":"conversation-1","reactor_id":"user-2","client_nonce":"reaction-nonce-2","operation":"clear","encryption":"e2ee","reacted_at":"2026-08-29T14:02:00Z"}`
	clearRequest := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/reactions", strings.NewReader(clearBody))
	clearRecorder := httptest.NewRecorder()
	handler.ServeHTTP(clearRecorder, clearRequest)
	if clearRecorder.Code != http.StatusAccepted {
		t.Fatalf("clear status = %d, body = %s", clearRecorder.Code, clearRecorder.Body.String())
	}

	list := httptest.NewRequest(http.MethodGet, "/v1/data/messages/message-1/reactions", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, list)
	if listRecorder.Code != http.StatusOK || strings.TrimSpace(listRecorder.Body.String()) != "[]" {
		t.Fatalf("cleared list status = %d, body = %s", listRecorder.Code, listRecorder.Body.String())
	}
}

func TestMessageReactionRejectsUnknownFields(t *testing.T) {
	handler := newReactionHTTPHandler(t, "user-2")
	body := `{"reaction_id":"reaction-1","conversation_id":"conversation-1","reactor_id":"user-2","client_nonce":"reaction-nonce-1","operation":"set","ciphertext":"ZW5jcnlwdGVk","encryption":"e2ee","reacted_at":"2026-08-29T14:01:00Z","plaintext_reaction":"👍"}`
	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/reactions", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
