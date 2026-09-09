// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

func TestWriteTypingServiceErrorMapsRateLimitTo429(t *testing.T) {
	response := httptest.NewRecorder()

	writeTypingServiceError(response, messagingservice.ErrTypingRateLimited)

	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusTooManyRequests)
	}
}
