// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestRuntimeProjectionRejectsNoncanonicalAuthenticatedIdentity(t *testing.T) {
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
	handler, err := NewDataRuntimeHandler(data, receipts, attachments, runtimeAuthenticator{userID: "user-a "})
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	handler.Routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/data/runtime", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
