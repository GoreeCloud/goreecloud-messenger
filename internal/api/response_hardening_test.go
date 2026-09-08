// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func assertAcceptedResponseHardening(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}
	if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q, want nosniff", got)
	}
	if got := recorder.Header().Get("Referrer-Policy"); got != "no-referrer" {
		t.Fatalf("Referrer-Policy = %q, want no-referrer", got)
	}
	if got := recorder.Header().Get("Content-Security-Policy"); got != "default-src 'none'; frame-ancestors 'none'" {
		t.Fatalf("Content-Security-Policy = %q", got)
	}
	if got := recorder.Header().Get("Cross-Origin-Resource-Policy"); got != "same-origin" {
		t.Fatalf("Cross-Origin-Resource-Policy = %q, want same-origin", got)
	}
	if recorder.Body.Len() != 0 {
		t.Fatalf("accepted mutation response must be bodyless, got %q", recorder.Body.String())
	}
}

func TestSubmitAcceptedResponseIsHardened(t *testing.T) {
	handler := newTestHandler(t, "user-1")
	body := `{"message_id":"message-hardening","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-hardening","ciphertext":"Y2lwaGVydGV4dA==","encryption":"e2ee","created_at":"2026-09-08T01:00:00Z"}`
	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assertAcceptedResponseHardening(t, recorder)
}

func TestReceiptAcceptedResponseIsHardened(t *testing.T) {
	handler := newTestHandlerWithSeededMessage(t, "user-2")
	body := `{"conversation_id":"conversation-1","user_id":"user-2","state":"delivered","observed_at":"2026-09-08T01:01:00Z"}`
	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages/message-1/receipts", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assertAcceptedResponseHardening(t, recorder)
}
