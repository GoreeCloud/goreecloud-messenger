// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type runtimeReadinessProbeFunc func(context.Context) error

func (f runtimeReadinessProbeFunc) Ready(ctx context.Context) error { return f(ctx) }

func TestRuntimeHealthIsPublicCategoricalAndNonCacheable(t *testing.T) {
	h := &DataRuntimeHandler{}
	mux := http.NewServeMux()
	h.registerHealthRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, RuntimeHealthPath, nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("health Cache-Control = %q, want no-store", response.Header().Get("Cache-Control"))
	}
	if body := response.Body.String(); !strings.Contains(body, `"status":"ok"`) {
		t.Fatalf("health body = %q, want categorical ok status", body)
	}
}

func TestRuntimeReadinessFailsClosedWithoutAuthority(t *testing.T) {
	h := &DataRuntimeHandler{}
	mux := http.NewServeMux()
	h.registerHealthRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, RuntimeReadinessPath, nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("readiness status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
	if body := response.Body.String(); !strings.Contains(body, `"status":"not_ready"`) {
		t.Fatalf("readiness body = %q, want categorical not_ready status", body)
	}
}

func TestRuntimeReadinessDoesNotLeakProbeFailureDetails(t *testing.T) {
	h := &DataRuntimeHandler{
		readinessProbe: runtimeReadinessProbeFunc(func(context.Context) error {
			return errors.New("secret dependency host and credential detail")
		}),
	}
	mux := http.NewServeMux()
	h.registerHealthRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, RuntimeReadinessPath, nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("readiness status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
	body := response.Body.String()
	if strings.Contains(body, "secret") || strings.Contains(body, "credential") || strings.Contains(body, "host") {
		t.Fatalf("readiness response leaked probe detail: %q", body)
	}
	if !strings.Contains(body, `"status":"not_ready"`) {
		t.Fatalf("readiness body = %q, want categorical not_ready status", body)
	}
}

func TestRuntimeReadinessSucceedsOnlyAfterExplicitReadyDecision(t *testing.T) {
	called := false
	h := &DataRuntimeHandler{
		readinessProbe: runtimeReadinessProbeFunc(func(ctx context.Context) error {
			called = true
			if ctx == nil {
				t.Fatal("readiness probe received nil context")
			}
			return nil
		}),
	}
	mux := http.NewServeMux()
	h.registerHealthRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, RuntimeReadinessPath, nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	if !called {
		t.Fatal("readiness authority was not consulted")
	}
	if response.Code != http.StatusOK {
		t.Fatalf("readiness status = %d, want %d", response.Code, http.StatusOK)
	}
	if body := response.Body.String(); !strings.Contains(body, `"status":"ready"`) {
		t.Fatalf("readiness body = %q, want categorical ready status", body)
	}
}

func TestWithRuntimeReadinessProbeRequiresExplicitAuthority(t *testing.T) {
	h := &DataRuntimeHandler{}
	if _, err := h.WithRuntimeReadinessProbe(nil); err == nil {
		t.Fatal("nil readiness probe was accepted")
	}
	var nilRuntime *DataRuntimeHandler
	if _, err := nilRuntime.WithRuntimeReadinessProbe(runtimeReadinessProbeFunc(func(context.Context) error { return nil })); err == nil {
		t.Fatal("nil runtime accepted readiness probe")
	}
}
