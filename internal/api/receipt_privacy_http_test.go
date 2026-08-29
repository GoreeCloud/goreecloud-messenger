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

func newReceiptPrivacyHTTPFixture(t *testing.T) (*messagingservice.DataService, *messagingservice.ReceiptService, *messagingservice.MemoryReceiptPrivacyPolicy) {
	t.Helper()

	store := messagingservice.NewMemoryDataStore()
	access := messagingservice.NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-receipt-privacy",
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
		MessageID:      "message-receipt-privacy",
		ConversationID: "conversation-receipt-privacy",
		SenderID:       "user-1",
		ClientNonce:    "nonce-receipt-privacy",
		Ciphertext:     []byte("ciphertext"),
		Encryption:     domain.EncryptionE2EE,
		CreatedAt:      time.Date(2026, 8, 29, 14, 30, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("seed message: %v", err)
	}

	privacy := messagingservice.NewMemoryReceiptPrivacyPolicy(true)
	receiptService, err := messagingservice.NewReceiptServiceWithPrivacy(store, messagingservice.NewMemoryReceiptStore(), access, privacy)
	if err != nil {
		t.Fatalf("new receipt service: %v", err)
	}
	return dataService, receiptService, privacy
}

func newReceiptPrivacyHTTPHandler(t *testing.T, dataService *messagingservice.DataService, receiptService *messagingservice.ReceiptService, userID string) http.Handler {
	t.Helper()
	handler, err := NewHandler(dataService, receiptService, testAuthenticator{userID: userID})
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}
	return handler.Routes()
}

func TestReadReceiptHTTPRejectsPrivacyDeniedPublicationButAllowsDelivery(t *testing.T) {
	dataService, receiptService, privacy := newReceiptPrivacyHTTPFixture(t)
	privacy.SetPublish("conversation-receipt-privacy", "user-2", false)
	handler := newReceiptPrivacyHTTPHandler(t, dataService, receiptService, "user-2")

	readBody := `{"conversation_id":"conversation-receipt-privacy","user_id":"user-2","state":"read","observed_at":"2026-08-29T14:31:00Z"}`
	readBody = strings.ReplaceAll(readBody, `\"`, `"`)
	readRequest := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-receipt-privacy/receipts", strings.NewReader(readBody))
	readRecorder := httptest.NewRecorder()
	handler.ServeHTTP(readRecorder, readRequest)
	if readRecorder.Code != http.StatusForbidden {
		t.Fatalf("read privacy status = %d, body = %s", readRecorder.Code, readRecorder.Body.String())
	}

	deliveredBody := `{"conversation_id":"conversation-receipt-privacy","user_id":"user-2","state":"delivered","observed_at":"2026-08-29T14:31:05Z"}`
	deliveredBody = strings.ReplaceAll(deliveredBody, `\"`, `"`)
	deliveredRequest := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-receipt-privacy/receipts", strings.NewReader(deliveredBody))
	deliveredRecorder := httptest.NewRecorder()
	handler.ServeHTTP(deliveredRecorder, deliveredRequest)
	if deliveredRecorder.Code != http.StatusAccepted {
		t.Fatalf("delivery status = %d, body = %s", deliveredRecorder.Code, deliveredRecorder.Body.String())
	}
}

func TestReadReceiptHTTPFiltersReadWhenObserverPrivacyDisabled(t *testing.T) {
	dataService, receiptService, privacy := newReceiptPrivacyHTTPFixture(t)
	publisher := newReceiptPrivacyHTTPHandler(t, dataService, receiptService, "user-2")
	observer := newReceiptPrivacyHTTPHandler(t, dataService, receiptService, "user-1")

	readBody := `{"conversation_id":"conversation-receipt-privacy","user_id":"user-2","state":"read","observed_at":"2026-08-29T14:31:00Z"}`
	readBody = strings.ReplaceAll(readBody, `\"`, `"`)
	readRequest := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-receipt-privacy/receipts", strings.NewReader(readBody))
	readRecorder := httptest.NewRecorder()
	publisher.ServeHTTP(readRecorder, readRequest)
	if readRecorder.Code != http.StatusAccepted {
		t.Fatalf("record read status = %d, body = %s", readRecorder.Code, readRecorder.Body.String())
	}

	privacy.SetObserve("conversation-receipt-privacy", "user-1", false)
	listRequest := httptest.NewRequest(http.MethodGet, "/v1/data/messages/message-receipt-privacy/receipts", nil)
	listRecorder := httptest.NewRecorder()
	observer.ServeHTTP(listRecorder, listRequest)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", listRecorder.Code, listRecorder.Body.String())
	}
	if strings.Contains(listRecorder.Body.String(), `"state":"read"`) || strings.Contains(listRecorder.Body.String(), `"observed_at"`) {
		t.Fatalf("observer privacy exposed read state: %s", listRecorder.Body.String())
	}
}

func TestReadReceiptHTTPAppliesCurrentPublisherPrivacyToStoredRead(t *testing.T) {
	dataService, receiptService, privacy := newReceiptPrivacyHTTPFixture(t)
	publisher := newReceiptPrivacyHTTPHandler(t, dataService, receiptService, "user-2")
	observer := newReceiptPrivacyHTTPHandler(t, dataService, receiptService, "user-1")

	readBody := `{"conversation_id":"conversation-receipt-privacy","user_id":"user-2","state":"read","observed_at":"2026-08-29T14:31:00Z"}`
	readBody = strings.ReplaceAll(readBody, `\"`, `"`)
	readRequest := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-receipt-privacy/receipts", strings.NewReader(readBody))
	readRecorder := httptest.NewRecorder()
	publisher.ServeHTTP(readRecorder, readRequest)
	if readRecorder.Code != http.StatusAccepted {
		t.Fatalf("record read status = %d, body = %s", readRecorder.Code, readRecorder.Body.String())
	}

	privacy.SetPublish("conversation-receipt-privacy", "user-2", false)
	listRequest := httptest.NewRequest(http.MethodGet, "/v1/data/messages/message-receipt-privacy/receipts", nil)
	listRecorder := httptest.NewRecorder()
	observer.ServeHTTP(listRecorder, listRequest)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", listRecorder.Code, listRecorder.Body.String())
	}
	if strings.Contains(listRecorder.Body.String(), `"state":"read"`) || strings.Contains(listRecorder.Body.String(), `"observed_at"`) {
		t.Fatalf("publisher privacy exposed stored read state: %s", listRecorder.Body.String())
	}
}
