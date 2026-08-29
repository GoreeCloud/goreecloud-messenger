// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

type runtimeAuthenticator struct {
	userID string
	err    error
}

func (a runtimeAuthenticator) Authenticate(context.Context, *http.Request) (string, error) {
	return a.userID, a.err
}

func TestDataRuntimeRequiresAllBoundaries(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	dataStore := messagingservice.NewMemoryDataStore()
	data, err := messagingservice.NewDataService(dataStore, access)
	if err != nil {
		t.Fatal(err)
	}
	receipts, err := messagingservice.NewReceiptService(dataStore, messagingservice.NewMemoryReceiptStore(), access)
	if err != nil {
		t.Fatal(err)
	}
	attachments, err := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name        string
		data        *messagingservice.DataService
		receipts    *messagingservice.ReceiptService
		attachments *messagingservice.AttachmentService
		auth        Authenticator
	}{
		{name: "data", receipts: receipts, attachments: attachments, auth: runtimeAuthenticator{userID: "user-a"}},
		{name: "receipts", data: data, attachments: attachments, auth: runtimeAuthenticator{userID: "user-a"}},
		{name: "attachments", data: data, receipts: receipts, auth: runtimeAuthenticator{userID: "user-a"}},
		{name: "auth", data: data, receipts: receipts, attachments: attachments},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := NewDataRuntimeHandler(tc.data, tc.receipts, tc.attachments, tc.auth); err == nil {
				t.Fatal("expected missing boundary to fail closed")
			}
		})
	}
}

func TestDataRuntimeComposesMessageAndAttachmentConversationRoutes(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	if err := access.SetConversation(domain.Conversation{
		ID:             "conversation-a",
		Kind:           domain.ConversationDirect,
		ParticipantIDs: []string{"user-a", "user-b"},
	}); err != nil {
		t.Fatal(err)
	}

	dataStore := messagingservice.NewMemoryDataStore()
	data, err := messagingservice.NewDataService(dataStore, access)
	if err != nil {
		t.Fatal(err)
	}
	receipts, err := messagingservice.NewReceiptService(dataStore, messagingservice.NewMemoryReceiptStore(), access)
	if err != nil {
		t.Fatal(err)
	}
	attachments, err := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)
	if err != nil {
		t.Fatal(err)
	}

	runtime, err := NewDataRuntimeHandler(data, receipts, attachments, runtimeAuthenticator{userID: "user-a"})
	if err != nil {
		t.Fatal(err)
	}
	handler := runtime.Routes()

	for _, path := range []string{
		"/v1/data/conversations/conversation-a/messages",
		"/v1/data/conversations/conversation-a/attachments",
	} {
		t.Run(path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, path, nil)
			handler.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusOK {
				t.Fatalf("expected %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestDataRuntimeUsesSameFailClosedAuthenticationBoundaryForAttachmentsAndMessages(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	dataStore := messagingservice.NewMemoryDataStore()
	data, _ := messagingservice.NewDataService(dataStore, access)
	receipts, _ := messagingservice.NewReceiptService(dataStore, messagingservice.NewMemoryReceiptStore(), access)
	attachments, _ := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)
	runtime, err := NewDataRuntimeHandler(data, receipts, attachments, runtimeAuthenticator{err: errors.New("denied")})
	if err != nil {
		t.Fatal(err)
	}
	handler := runtime.Routes()

	for _, path := range []string{
		"/v1/data/conversations/conversation-a/messages",
		"/v1/data/conversations/conversation-a/attachments",
		"/v1/data/attachments/attachment-a/ciphertext",
	} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, path, nil)
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("%s expected %d, got %d", path, http.StatusUnauthorized, recorder.Code)
		}
	}
}
