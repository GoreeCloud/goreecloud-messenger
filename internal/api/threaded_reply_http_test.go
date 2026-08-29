// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestThreadedReplyRoundTripKeepsPayloadOpaque(t *testing.T) {
	handler := newTestHandler(t, "user-1")
	rootCiphertext := base64.StdEncoding.EncodeToString([]byte("opaque thread root ciphertext"))
	firstCiphertext := base64.StdEncoding.EncodeToString([]byte("opaque first thread ciphertext"))
	nestedCiphertext := base64.StdEncoding.EncodeToString([]byte("opaque nested thread ciphertext"))

	requests := []string{
		`{"message_id":"message-1","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-1","ciphertext":"` + rootCiphertext + `","encryption":"e2ee","created_at":"2026-08-29T13:42:00Z"}`,
		`{"message_id":"message-2","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-2","reply_to_message_id":"message-1","thread_root_message_id":"message-1","ciphertext":"` + firstCiphertext + `","encryption":"e2ee","created_at":"2026-08-29T13:43:00Z"}`,
		`{"message_id":"message-3","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-3","reply_to_message_id":"message-2","thread_root_message_id":"message-1","ciphertext":"` + nestedCiphertext + `","encryption":"e2ee","created_at":"2026-08-29T13:44:00Z"}`,
	}
	for _, body := range requests {
		request := httptest.NewRequest(http.MethodPost, "/v1/data/messages", strings.NewReader(body))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusAccepted {
			t.Fatalf("submit status = %d, body = %s", recorder.Code, recorder.Body.String())
		}
	}

	list := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/conversation-1/threads/message-1/messages", nil)
	listRecorder := httptest.NewRecorder()
	handler.ServeHTTP(listRecorder, list)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("thread list status = %d, body = %s", listRecorder.Code, listRecorder.Body.String())
	}
	response := listRecorder.Body.String()
	if !strings.Contains(response, `"thread_root_message_id":"message-1"`) {
		t.Fatalf("response missing thread root metadata: %s", response)
	}
	if !strings.Contains(response, `"reply_to_message_id":"message-2"`) {
		t.Fatalf("response missing nested reply metadata: %s", response)
	}
	if !strings.Contains(response, `"ciphertext":"`+nestedCiphertext+`"`) {
		t.Fatalf("response missing opaque nested ciphertext: %s", response)
	}
	for _, plaintext := range []string{"opaque thread root ciphertext", "opaque first thread ciphertext", "opaque nested thread ciphertext"} {
		if strings.Contains(response, plaintext) {
			t.Fatalf("response exposed decoded thread content %q: %s", plaintext, response)
		}
	}
}

func TestThreadedReplyRejectsUnrelatedParentWithBoundedError(t *testing.T) {
	handler := newTestHandler(t, "user-1")
	for _, body := range []string{
		`{"message_id":"message-1","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-1","ciphertext":"Y2lwaGVydGV4dC0x","encryption":"e2ee","created_at":"2026-08-29T13:42:00Z"}`,
		`{"message_id":"message-2","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-2","ciphertext":"Y2lwaGVydGV4dC0y","encryption":"e2ee","created_at":"2026-08-29T13:43:00Z"}`,
	} {
		request := httptest.NewRequest(http.MethodPost, "/v1/data/messages", strings.NewReader(body))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusAccepted {
			t.Fatalf("seed status = %d, body = %s", recorder.Code, recorder.Body.String())
		}
	}

	body := `{"message_id":"message-3","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"nonce-3","reply_to_message_id":"message-2","thread_root_message_id":"message-1","ciphertext":"Y2lwaGVydGV4dC0z","encryption":"e2ee","created_at":"2026-08-29T13:44:00Z"}`
	request := httptest.NewRequest(http.MethodPost, "/v1/data/messages", strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"error":"thread target unavailable"`) {
		t.Fatalf("unexpected bounded error: %s", recorder.Body.String())
	}
}

func TestThreadListRejectsNonParticipant(t *testing.T) {
	handler := newTestHandler(t, "user-3")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/conversation-1/threads/message-1/messages", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
