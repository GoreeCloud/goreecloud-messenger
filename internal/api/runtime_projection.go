// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// RuntimePersistenceProbe is a deliberately narrow diagnostic boundary. It may
// answer whether the composed persistence dependency can complete a bounded
// health check, but must not return paths, credentials, keys, user identifiers,
// record counts, or other sensitive implementation details.
type RuntimePersistenceProbe interface {
	Check(context.Context) error
}

type RuntimePersistenceProbeFunc func(context.Context) error

func (f RuntimePersistenceProbeFunc) Check(ctx context.Context) error { return f(ctx) }

type runtimeCompositionResponse struct {
	Scope           string `json:"scope"`
	Messages        string `json:"messages"`
	Receipts        string `json:"receipts"`
	Attachments     string `json:"attachments"`
	Authentication  string `json:"authentication"`
	Persistence     string `json:"persistence"`
	ProductionReady bool   `json:"production_ready"`
}

func (h *DataRuntimeHandler) registerRuntimeProjection(mux *http.ServeMux) {
	if mux == nil {
		panic("HTTP mux is required")
	}
	mux.HandleFunc("GET /v1/data/runtime", h.runtimeProjection)
}

func (h *DataRuntimeHandler) runtimeProjection(w http.ResponseWriter, r *http.Request) {
	userID, err := h.auth.Authenticate(r.Context(), r)
	if err != nil || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	persistence := "not-assessed"
	if h.persistenceProbe != nil {
		ctx, cancel := context.WithTimeout(r.Context(), runtimePersistenceProbeTimeout)
		defer cancel()
		if err := h.persistenceProbe.Check(ctx); err != nil {
			persistence = "unavailable"
		} else {
			persistence = "available"
		}
	}

	writeJSON(w, http.StatusOK, runtimeCompositionResponse{
		Scope:           "development-composition",
		Messages:        "configured",
		Receipts:        "configured",
		Attachments:     "configured",
		Authentication:  "accepted",
		Persistence:     persistence,
		ProductionReady: false,
	})
}

const runtimePersistenceProbeTimeout = 500 * time.Millisecond
