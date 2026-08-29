// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDirectReplyRoundTripKeepsPayloadOpaque(t *testing.T) {
	handler := newTestHandler(t, "user-1")
	originalCiphertext := base64.StdEncoding.EncodeToString([]byte("opaque original ciphertext"))
	replyCiphertext := base64.StdEncoding.EncodeToString([]byte("opaque reply ciphertext"))

	originalBody := `{"message_id":"message-1","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-1","ciphertext":"` + originalCiphertext + `","encryption":"e2ee","created_at":"2026-08-29T13:30:00Z"}`
	original := httptest.NewRequest(http.MethodPost, "/v1/data/messages", strings.NewReader(originalBody))
	originalRecorder := httptest.NewRecorder()
	handler.ServeHTTP(originalRecorder, original)
	if originalRecorder.Code != http.StatusAccepted {
		t.Fatalf("original status = %d, body = %s", originalRecorder.Code, originalRecorder.Body.String())
	}

	replyBody := `{"message_id":"message-2","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-2","reply_to_message_id":"message-1","ciphertext":"` + replyCiphertext + `","encryption":"e2ee","created_at":"2026-08-29T13:31:00Z"}`
	reply := httptest.NewRequest(http.MethodPost, "/v1/data/messages", strings.NewReader(replyBody))
	replyRecorder := httptest.NewRecorder()
	handler.ServeHTTP(replyRecorder, reply)
	if replyRecorder.Code != http.StatusAccepted {
		t.Fatalf("reply status = %d, body = %s", replyRecorder.Code, replyRecorder.Body.String())
	}

	list := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/conversation-1/messages", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, list)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", listRecorder.Code, listRecorder.Body.String())
	}
	response := listRecorder.Body.String()
	if !strings.Contains(response, `"reply_to_message_id":"message-1"`) {
		t.Fatalf("response missing reply target: %s", response)
	}
	if !strings.Contains(response, `"ciphertext":"`+replyCiphertext+`"`) {
		t.Fatalf("response missing opaque reply ciphertext: %s", response)
	}
	if strings.Contains(response, "opaque reply ciphertext") || strings.Contains(response, "opaque original ciphertext") {
		t.Fatalf("response exposed decoded message content: %s", response)
	}
}

func TestDirectReplyRejectsUnavailableTargetWithBoundedError(t *testing.T) {
	handler := newTestHandler(t, "user-1")
	body := `{"message_id":"message-2","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-2","reply_to_message_id":"missing-message","ciphertext":"Y2lwaGVydGV4dA==","encryption":"e2ee","created_at":"2026-08-29T13:31:00Z"}`

	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"error":"reply target unavailable"`) {
		t.Fatalf("unexpected bounded error: %s", recorder.Body.String())
	}
}
