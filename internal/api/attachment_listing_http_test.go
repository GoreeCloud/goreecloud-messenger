// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAttachmentListReturnsMetadataWithoutCiphertextOrNonce(t *testing.T) {
	handler := newTestAttachmentHandler(t, "user-1")
	ciphertext := []byte("opaque-listing-ciphertext")
	body := `{"attachment_id":"attachment-1","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"private-nonce-1","filename":"photo.jpg","mime_type":"image/jpeg","ciphertext":"` + base64.StdEncoding.EncodeToString(ciphertext) + `"}`

	submit := httptest.NewRequest(http.MethodPost, "/v1/data/attachments", strings.NewReader(body))
	submitRecorder := httptest.NewRecorder()
	handler.ServeHTTP(submitRecorder, submit)
	if submitRecorder.Code != http.StatusAccepted {
		t.Fatalf("submit status = %d, body = %s", submitRecorder.Code, submitRecorder.Body.String())
	}

	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/conversation-1/attachments?limit=10", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	response := recorder.Body.String()
	for _, expected := range []string{`"attachment_id":"attachment-1"`, `"filename":"photo.jpg"`, `"mime_type":"image/jpeg"`, `"ciphertext_bytes":25`} {
		if !strings.Contains(response, expected) {
			t.Fatalf("listing missing %s: %s", expected, response)
		}
	}
	for _, forbidden := range []string{"ciphertext\"", "private-nonce-1", "client_nonce"} {
		if strings.Contains(response, forbidden) {
			t.Fatalf("listing exposed forbidden attachment material %q: %s", forbidden, response)
		}
	}
}

func TestAttachmentListRejectsInvalidLimit(t *testing.T) {
	handler := newTestAttachmentHandler(t, "user-1")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/conversations/conversation-1/attachments?limit=101", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
