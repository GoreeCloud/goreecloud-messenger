// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestRuntimeProjectionRequiresAuthenticationAndStaysMinimized(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	dataStore := messagingservice.NewMemoryDataStore()
	data, _ := messagingservice.NewDataService(dataStore, access)
	receipts, _ := messagingservice.NewReceiptService(dataStore, messagingservice.NewMemoryReceiptStore(), access)
	attachments, _ := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)

	denied, err := NewDataRuntimeHandler(data, receipts, attachments, runtimeAuthenticator{err: errors.New("denied")})
	if err != nil {
		t.Fatal(err)
	}
	deniedRecorder := httptest.NewRecorder()
	denied.Routes().ServeHTTP(deniedRecorder, httptest.NewRequest(http.MethodGet, "/v1/data/runtime", nil))
	if deniedRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected authentication failure, got %d", deniedRecorder.Code)
	}

	allowed, err := NewDataRuntimeHandler(data, receipts, attachments, runtimeAuthenticator{userID: "user-a"})
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	allowed.Routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/v1/data/runtime", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	body := recorder.Body.String()
	for _, expected := range []string{
		`"scope":"development-composition"`,
		`"messages":"configured"`,
		`"receipts":"configured"`,
		`"attachments":"configured"`,
		`"authentication":"accepted"`,
		`"persistence":"not-assessed"`,
		`"production_ready":false`,
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("projection missing %s in %s", expected, body)
		}
	}
	assertRuntimeProjectionMinimized(t, body)
}

func TestRuntimeProjectionReportsOnlyBoundedPersistenceState(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	dataStore := messagingservice.NewMemoryDataStore()
	data, _ := messagingservice.NewDataService(dataStore, access)
	receipts, _ := messagingservice.NewReceiptService(dataStore, messagingservice.NewMemoryReceiptStore(), access)
	attachments, _ := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)

	base, err := NewDataRuntimeHandler(data, receipts, attachments, runtimeAuthenticator{userID: "user-a"})
	if err != nil {
		t.Fatal(err)
	}
	available, err := base.WithRuntimePersistenceProbe(RuntimePersistenceProbeFunc(func(context.Context) error { return nil }))
	if err != nil {
		t.Fatal(err)
	}
	availableRecorder := httptest.NewRecorder()
	available.Routes().ServeHTTP(availableRecorder, httptest.NewRequest(http.MethodGet, "/v1/data/runtime", nil))
	if !strings.Contains(availableRecorder.Body.String(), `"persistence":"available"`) {
		t.Fatalf("expected available persistence state: %s", availableRecorder.Body.String())
	}
	assertRuntimeProjectionMinimized(t, availableRecorder.Body.String())

	unavailable, err := base.WithRuntimePersistenceProbe(RuntimePersistenceProbeFunc(func(context.Context) error {
		return errors.New("/private/store credential=secret key=opaque")
	}))
	if err != nil {
		t.Fatal(err)
	}
	unavailableRecorder := httptest.NewRecorder()
	unavailable.Routes().ServeHTTP(unavailableRecorder, httptest.NewRequest(http.MethodGet, "/v1/data/runtime", nil))
	body := unavailableRecorder.Body.String()
	if !strings.Contains(body, `"persistence":"unavailable"`) {
		t.Fatalf("expected unavailable persistence state: %s", body)
	}
	assertRuntimeProjectionMinimized(t, body)
}

func assertRuntimeProjectionMinimized(t *testing.T, body string) {
	t.Helper()
	for _, forbidden := range []string{"user-a", "root", "path", "ciphertext", "credential", "secret", "private", "key"} {
		if strings.Contains(strings.ToLower(body), forbidden) {
			t.Fatalf("projection leaked forbidden detail %q in %s", forbidden, body)
		}
	}
}
