// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestDataRuntimeComposesTypingPrivacyPreferencesOnlyWhenExplicitlyEnabled(t *testing.T) {
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

	path := "/v1/data/conversations/conversation-a/typing/preferences"
	baseRecorder := httptest.NewRecorder()
	runtime.Routes().ServeHTTP(baseRecorder, httptest.NewRequest(http.MethodGet, path, nil))
	if baseRecorder.Code != http.StatusNotFound {
		t.Fatalf("base runtime expected %d, got %d", http.StatusNotFound, baseRecorder.Code)
	}

	preferences, err := messagingservice.NewTypingPrivacyPreferenceService(
		messagingservice.NewMemoryTypingPrivacyPolicy(true),
		access,
	)
	if err != nil {
		t.Fatal(err)
	}
	composed, err := runtime.WithTypingPrivacyPreferences(preferences)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	composed.Routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("composed runtime expected %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
}
