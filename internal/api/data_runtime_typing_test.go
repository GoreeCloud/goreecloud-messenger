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

func TestDataRuntimeTypingIsExplicitlyEnabledAndSharesAuthentication(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-typing-runtime",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-a", "user-b"},
	}); err != nil {
		t.Fatal(err)
	}
	dataStore := messagingservice.NewMemoryDataStore()
	data, _ := messagingservice.NewDataService(dataStore, access)
	receipts, _ := messagingservice.NewReceiptService(dataStore, messagingservice.NewMemoryReceiptStore(), access)
	attachments, _ := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)
	base, err := NewDataRuntimeHandler(data, receipts, attachments, runtimeAuthenticator{userID: "user-a"})
	if err != nil {
		t.Fatal(err)
	}

	path := "/v1/data/conversations/conversation-typing-runtime/typing"
	withoutTyping := httptest.NewRecorder()
	base.Routes().ServeHTTP(withoutTyping, httptest.NewRequest(http.MethodGet, path, nil))
	if withoutTyping.Code != http.StatusNotFound {
		t.Fatalf("typing unexpectedly enabled by default: %d", withoutTyping.Code)
	}

	now := time.Date(2026, 8, 30, 22, 0, 0, 0, time.UTC)
	typingService, err := messagingservice.NewTypingService(
		messagingservice.NewMemoryTypingStore(),
		access,
		messagingservice.NewMemoryTypingPrivacyPolicy(true),
		func() time.Time { return now },
	)
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := base.WithTypingIndicators(typingService)
	if err != nil {
		t.Fatal(err)
	}

	publish := httptest.NewRecorder()
	publishRequest := httptest.NewRequest(
		http.MethodPost,
		path,
		strings.NewReader(`{"user_id":"user-a","sequence":1,"state":"typing"}`),
	)
	runtime.Routes().ServeHTTP(publish, publishRequest)
	if publish.Code != http.StatusAccepted {
		t.Fatalf("typing publish status = %d, body = %s", publish.Code, publish.Body.String())
	}
}

func TestDataRuntimeTypingUsesSameFailClosedAuthenticationBoundary(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	dataStore := messagingservice.NewMemoryDataStore()
	data, _ := messagingservice.NewDataService(dataStore, access)
	receipts, _ := messagingservice.NewReceiptService(dataStore, messagingservice.NewMemoryReceiptStore(), access)
	attachments, _ := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)
	base, err := NewDataRuntimeHandler(data, receipts, attachments, runtimeAuthenticator{err: errRuntimeAuthDenied})
	if err != nil {
		t.Fatal(err)
	}
	typingService, err := messagingservice.NewTypingService(
		messagingservice.NewMemoryTypingStore(),
		access,
		messagingservice.NewMemoryTypingPrivacyPolicy(true),
		time.Now,
	)
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := base.WithTypingIndicators(typingService)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/conversation-a/typing", nil)
	runtime.Routes().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

var errRuntimeAuthDenied = &runtimeAuthError{}

type runtimeAuthError struct{}

func (*runtimeAuthError) Error() string { return "denied" }
