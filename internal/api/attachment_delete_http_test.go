// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAttachmentDeleteIsIdempotentAndMakesCiphertextUnavailable(t *testing.T) {
	handler := newTestAttachmentHandler(t, "user-1")
	ciphertext := base64.StdEncoding.EncodeToString([]byte("opaque attachment ciphertext"))
	body := `{"attachment_id":"attachment-1","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"attachment-nonce-1","filename":"photo.jpg","mime_type":"image/jpeg","ciphertext":"` + ciphertext + `"}`

	submit := httptest.NewRequest(http.MethodPost, "/v1/data/attachments", strings.NewReader(body))
	submitRecorder := httptest.NewRecorder()
	handler.ServeHTTP(submitRecorder, submit)
	if submitRecorder.Code != http.StatusAccepted {
		t.Fatalf("submit status = %d, body = %s", submitRecorder.Code, submitRecorder.Body.String())
	}

	for attempt := 0; attempt < 2; attempt++ {
		request := httptest.NewRequest(http.MethodDelete, "/v1/data/attachments/attachment-1", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusNoContent {
			t.Fatalf("delete attempt %d status = %d, body = %s", attempt+1, recorder.Code, recorder.Body.String())
		}
	}

	get := httptest.NewRequest(http.MethodGet, "/v1/data/attachments/attachment-1", nil)
	getRecorder := httptest.NewRecorder()
	handler.ServeHTTP(getRecorder, get)
	if getRecorder.Code != http.StatusNotFound {
		t.Fatalf("get after delete status = %d, body = %s", getRecorder.Code, getRecorder.Body.String())
	}
}
