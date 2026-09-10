// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

func TestIdentityDirectoryHTTPRejectsTrailingJSONBeforeResolution(t *testing.T) {
	calls := 0
	handler := newIdentityDirectoryHTTPTestHandler(
		t,
		testAuthenticator{userID: "user-1"},
		apiIdentityDirectoryResolverFunc(func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			calls++
			return domain.IdentityDirectoryProjection{}, false, nil
		}),
	)

	for _, body := range []string{
		`{"handle":"alice"} {"handle":"bob"}`,
		`{"handle":"alice"} trailing-garbage`,
		`{"handle":"alice"} [`,
	} {
		request := httptest.NewRequest(http.MethodPost, "/v1/identity/resolve", strings.NewReader(body))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("body %q: status = %d, response = %s", body, recorder.Code, recorder.Body.String())
		}
	}

	if calls != 0 {
		t.Fatalf("resolver calls = %d, want 0 for ambiguous or malformed request bodies", calls)
	}
}

func TestIdentityDirectoryHTTPAllowsTrailingJSONWhitespace(t *testing.T) {
	calls := 0
	handler := newIdentityDirectoryHTTPTestHandler(
		t,
		testAuthenticator{userID: "user-1"},
		apiIdentityDirectoryResolverFunc(func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			calls++
			return domain.IdentityDirectoryProjection{}, false, nil
		}),
	)

	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/identity/resolve",
		strings.NewReader("{\"handle\":\"alice\"}\n\t "),
	)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if calls != 1 {
		t.Fatalf("resolver calls = %d, want 1", calls)
	}
}
