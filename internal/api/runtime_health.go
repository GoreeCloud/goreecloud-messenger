// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"net/http"
)

const (
	RuntimeHealthPath    = "/healthz"
	RuntimeReadinessPath = "/readyz"
)

// RuntimeReadinessProbe is the narrow authority used by the HTTP runtime to decide whether
// this process is ready to accept application traffic. Implementations may evaluate required
// dependencies internally, but the public response remains categorical and must not disclose
// dependency names, credentials, cryptographic state, user data, or configuration details.
type RuntimeReadinessProbe interface {
	Ready(context.Context) error
}

func (h *DataRuntimeHandler) registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET "+RuntimeHealthPath, h.runtimeHealth)
	mux.HandleFunc("GET "+RuntimeReadinessPath, h.runtimeReadiness)
}

func (h *DataRuntimeHandler) runtimeHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *DataRuntimeHandler) runtimeReadiness(w http.ResponseWriter, r *http.Request) {
	if h.readinessProbe == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready"})
		return
	}
	if err := h.readinessProbe.Ready(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
