// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestDataRuntimeWithTypingPresenceComposesEventsAndPreferences(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-a",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-a", "user-b"},
	}); err != nil {
		t.Fatal(err)
	}

	dataStore := messagingservice.NewMemoryDataStore()
	data, _ := messagingservice.NewDataService(dataStore, access)
	receipts, _ := messagingservice.NewReceiptService(dataStore, messagingservice.NewMemoryReceiptStore(), access)
	attachments, _ := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)
	runtime, err := NewDataRuntimeHandler(data, receipts, attachments, runtimeAuthenticator{userID: "user-a"})
	if err != nil {
		t.Fatal(err)
	}

	policy := messagingservice.NewMemoryTypingPrivacyPolicy(true)
	typing, err := messagingservice.NewTypingService(
		messagingservice.NewMemoryTypingStore(),
		access,
		policy,
		func() time.Time { return time.Date(2026, 8, 31, 12, 0, 0, 0, time.UTC) },
	)
	if err != nil {
		t.Fatal(err)
	}
	preferences, err := messagingservice.NewTypingPrivacyPreferenceService(policy, access)
	if err != nil {
		t.Fatal(err)
	}

	composed, err := runtime.WithTypingPresence(typing, preferences)
	if err != nil {
		t.Fatal(err)
	}
	handler := composed.Routes()

	publish := httptest.NewRecorder()
	handler.ServeHTTP(publish, httptest.NewRequest(
		http.MethodPost,
		"/v1/data/conversations/conversation-a/typing",
		strings.NewReader(`{"sequence":1,"state":"typing"}`),
	))
	if publish.Code != http.StatusAccepted {
		t.Fatalf("typing publish expected %d, got %d: %s", http.StatusAccepted, publish.Code, publish.Body.String())
	}

	preferenceRead := httptest.NewRecorder()
	handler.ServeHTTP(preferenceRead, httptest.NewRequest(
		http.MethodGet,
		"/v1/data/conversations/conversation-a/typing/preferences",
		nil,
	))
	if preferenceRead.Code != http.StatusOK {
		t.Fatalf("typing preferences expected %d, got %d: %s", http.StatusOK, preferenceRead.Code, preferenceRead.Body.String())
	}
}

func TestDataRuntimeWithTypingPresenceRequiresBothBoundaries(t *testing.T) {
	var runtime *DataRuntimeHandler
	if _, err := runtime.WithTypingPresence(nil, nil); err == nil {
		t.Fatal("expected nil runtime and service boundaries to be rejected")
	}
}
