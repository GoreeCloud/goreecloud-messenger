// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAttachmentCiphertextDownloadReturnsExactOpaqueBytes(t *testing.T) {
	handler := newTestAttachmentHandler(t, "user-1")
	ciphertext := []byte{0x00, 0x01, 0x7f, 0x80, 0xff, 'G', 'C'}
	body := `{"attachment_id":"attachment-raw-1","conversation_id":"conversation-1","sender_id":"user-1","client_nonce":"attachment-raw-nonce-1","filename":"private-photo.jpg","mime_type":"image/jpeg","ciphertext":"` + base64.StdEncoding.EncodeToString(ciphertext) + `"}`

	submit := httptest.NewRequest(http.MethodPost, "/v1/data/attachments", strings.NewReader(body))
	submitRecorder := httptest.NewRecorder()
	handler.ServeHTTP(submitRecorder, submit)
	if submitRecorder.Code != http.StatusAccepted {
		t.Fatalf("submit status = %d, body = %s", submitRecorder.Code, submitRecorder.Body.String())
	}

	download := httptest.NewRequest(http.MethodGet, "/v1/data/attachments/attachment-raw-1/ciphertext", nil)
	downloadRecorder := httptest.NewRecorder()
	handler.ServeHTTP(downloadRecorder, download)
	if downloadRecorder.Code != http.StatusOK {
		t.Fatalf("download status = %d, body = %s", downloadRecorder.Code, downloadRecorder.Body.String())
	}
	if got := downloadRecorder.Body.Bytes(); string(got) != string(ciphertext) {
		t.Fatalf("ciphertext bytes changed: got %v want %v", got, ciphertext)
	}
	if got := downloadRecorder.Header().Get("Content-Type"); got != "application/octet-stream" {
		t.Fatalf("content type = %q", got)
	}
	if got := downloadRecorder.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("cache control = %q", got)
	}
	if got := downloadRecorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("nosniff = %q", got)
	}
	if got := downloadRecorder.Header().Get("Content-Length"); got != "7" {
		t.Fatalf("content length = %q", got)
	}
}

func TestAttachmentCiphertextDownloadRequiresAuthentication(t *testing.T) {
	handler := newTestAttachmentHandler(t, "")
	request := httptest.NewRequest(http.MethodGet, "/v1/data/attachments/attachment-1/ciphertext", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
