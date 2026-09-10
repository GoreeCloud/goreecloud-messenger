// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

type apiIdentityDirectoryResolverFunc func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error)

func (f apiIdentityDirectoryResolverFunc) ResolveExact(
	ctx context.Context,
	handle string,
) (domain.IdentityDirectoryProjection, bool, error) {
	return f(ctx, handle)
}

func newIdentityDirectoryHTTPTestHandler(
	t *testing.T,
	auth Authenticator,
	resolver messagingservice.IdentityDirectoryResolver,
) http.Handler {
	t.Helper()
	service, err := messagingservice.NewIdentityDirectoryService(resolver)
	if err != nil {
		t.Fatalf("new Identity directory service: %v", err)
	}
	handler, err := NewIdentityDirectoryHTTPHandler(service, auth)
	if err != nil {
		t.Fatalf("new Identity directory handler: %v", err)
	}
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux
}

func TestIdentityDirectoryHTTPReturnsOnlyMinimizedResolvedProjection(t *testing.T) {
	var observedHandle string
	handler := newIdentityDirectoryHTTPTestHandler(
		t,
		testAuthenticator{userID: "user-1"},
		apiIdentityDirectoryResolverFunc(func(_ context.Context, handle string) (domain.IdentityDirectoryProjection, bool, error) {
			observedHandle = handle
			return domain.IdentityDirectoryProjection{
				Subject:     "identity-subject-1",
				Handle:      "alice.example",
				DisplayName: "Alice Example",
			}, true, nil
		}),
	)

	request := httptest.NewRequest(http.MethodPost, "/v1/identity/resolve", strings.NewReader(`{"handle":"@Alice.Example"}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if observedHandle != "@Alice.Example" {
		t.Fatalf("resolver saw %q; client handle must reach Identity resolver without Messenger canonicalization", observedHandle)
	}
	body := recorder.Body.String()
	for _, required := range []string{`"subject":"identity-subject-1"`, `"handle":"alice.example"`, `"display_name":"Alice Example"`} {
		if !strings.Contains(body, required) {
			t.Fatalf("response missing %s: %s", required, body)
		}
	}
	for _, forbidden := range []string{"email", "phone", "discoverable", "allowed_services", "requester_service"} {
		if strings.Contains(body, forbidden) {
			t.Fatalf("response leaked forbidden Identity detail %q: %s", forbidden, body)
		}
	}
	if recorder.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", recorder.Header().Get("Cache-Control"))
	}
}

func TestIdentityDirectoryHTTPUsesUniformNotResolvedResponse(t *testing.T) {
	handler := newIdentityDirectoryHTTPTestHandler(
		t,
		testAuthenticator{userID: "user-1"},
		apiIdentityDirectoryResolverFunc(func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			return domain.IdentityDirectoryProjection{}, false, nil
		}),
	)

	var expectedBody string
	for _, handle := range []string{"does-not-exist", "private-account", "service-not-authorized"} {
		request := httptest.NewRequest(http.MethodPost, "/v1/identity/resolve", strings.NewReader(`{"handle":"`+handle+`"}`))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("%s: status = %d, body = %s", handle, recorder.Code, recorder.Body.String())
		}
		if expectedBody == "" {
			expectedBody = recorder.Body.String()
		} else if recorder.Body.String() != expectedBody {
			t.Fatalf("negative responses differ: %q != %q", recorder.Body.String(), expectedBody)
		}
	}
	if !strings.Contains(expectedBody, `"code":"not_resolved"`) {
		t.Fatalf("uniform response missing not_resolved code: %s", expectedBody)
	}
}

func TestIdentityDirectoryHTTPRejectsClientSuppliedRequesterService(t *testing.T) {
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
		strings.NewReader(`{"handle":"alice","requester_service":"goreecloud-messenger"}`),
	)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if calls != 0 {
		t.Fatalf("resolver calls = %d, want 0", calls)
	}
}

func TestIdentityDirectoryHTTPRequiresAuthenticatedMessengerUser(t *testing.T) {
	calls := 0
	handler := newIdentityDirectoryHTTPTestHandler(
		t,
		testAuthenticator{err: errors.New("not authenticated")},
		apiIdentityDirectoryResolverFunc(func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			calls++
			return domain.IdentityDirectoryProjection{}, false, nil
		}),
	)

	request := httptest.NewRequest(http.MethodPost, "/v1/identity/resolve", strings.NewReader(`{"handle":"alice"}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if calls != 0 {
		t.Fatalf("resolver calls = %d, want 0", calls)
	}
}

func TestIdentityDirectoryHTTPDoesNotLeakResolverFailure(t *testing.T) {
	handler := newIdentityDirectoryHTTPTestHandler(
		t,
		testAuthenticator{userID: "user-1"},
		apiIdentityDirectoryResolverFunc(func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			return domain.IdentityDirectoryProjection{}, false, errors.New("secret upstream topology detail")
		}),
	)

	request := httptest.NewRequest(http.MethodPost, "/v1/identity/resolve", strings.NewReader(`{"handle":"alice"}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), "secret upstream topology detail") {
		t.Fatalf("response leaked resolver failure: %s", recorder.Body.String())
	}
}

func TestDataRuntimeKeepsIdentityDirectoryAbsentUntilExplicitlyComposed(t *testing.T) {
	access := messagingservice.NewMemoryConversationAccess()
	dataStore := messagingservice.NewMemoryDataStore()
	data, err := messagingservice.NewDataService(dataStore, access)
	if err != nil {
		t.Fatal(err)
	}
	receipts, err := messagingservice.NewReceiptService(dataStore, messagingservice.NewMemoryReceiptStore(), access)
	if err != nil {
		t.Fatal(err)
	}
	attachments, err := messagingservice.NewAttachmentService(messagingservice.NewMemoryAttachmentStore(), access)
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := NewDataRuntimeHandler(data, receipts, attachments, testAuthenticator{userID: "user-1"})
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/v1/identity/resolve", strings.NewReader(`{"handle":"alice"}`))
	recorder := httptest.NewRecorder()
	runtime.Routes().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("uncomposed runtime status = %d, want 404", recorder.Code)
	}

	directoryService, err := messagingservice.NewIdentityDirectoryService(apiIdentityDirectoryResolverFunc(
		func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			return domain.IdentityDirectoryProjection{Subject: "identity-subject-1", Handle: "alice"}, true, nil
		},
	))
	if err != nil {
		t.Fatal(err)
	}
	withDirectory, err := runtime.WithIdentityDirectory(directoryService)
	if err != nil {
		t.Fatal(err)
	}

	request = httptest.NewRequest(http.MethodPost, "/v1/identity/resolve", strings.NewReader(`{"handle":"alice"}`))
	recorder = httptest.NewRecorder()
	withDirectory.Routes().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("composed runtime status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
