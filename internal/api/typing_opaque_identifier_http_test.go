// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTypingHTTPRejectsBoundaryWhitespaceConversationPathInsteadOfNormalizing(t *testing.T) {
	now := time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)
	service, _ := newTypingHTTPTestService(t, &now)
	handler := newTypingHTTPHandler(t, service, "user-1")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/%20conversation-typing-http/typing", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestTypingHTTPRejectsNoncanonicalAuthenticatedIdentity(t *testing.T) {
	now := time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)
	service, _ := newTypingHTTPTestService(t, &now)
	handler := newTypingHTTPHandler(t, service, "user-1 ")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/conversation-typing-http/typing", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestTypingPreferencesHTTPRejectsBoundaryWhitespaceConversationPathInsteadOfNormalizing(t *testing.T) {
	service, _ := typingPreferenceTestService(t)
	handler, err := NewTypingPreferencesHTTPHandler(service, typingPreferencesAuthenticator{userID: "user-a"})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/%20conversation-a/typing/preferences", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}

func TestTypingPreferencesHTTPRejectsNoncanonicalAuthenticatedIdentity(t *testing.T) {
	service, _ := typingPreferenceTestService(t)
	handler, err := NewTypingPreferencesHTTPHandler(service, typingPreferencesAuthenticator{userID: "user-a "})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/conversation-a/typing/preferences", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
