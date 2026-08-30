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

// RuntimeCryptographyProbe is equally narrow: it reports only whether the
// composed cryptography/session dependency can complete a bounded diagnostic
// check. Implementations must not return or project session identifiers, key
// material, algorithms, ciphertext, credentials, user identifiers, or errors.
type RuntimeCryptographyProbe interface {
	Check(context.Context) error
}

type RuntimeCryptographyProbeFunc func(context.Context) error

func (f RuntimeCryptographyProbeFunc) Check(ctx context.Context) error { return f(ctx) }

type runtimeCompositionResponse struct {
	Scope           string `json:"scope"`
	Messages        string `json:"messages"`
	Receipts        string `json:"receipts"`
	Attachments     string `json:"attachments"`
	Authentication  string `json:"authentication"`
	Persistence     string `json:"persistence"`
	Cryptography    string `json:"cryptography"`
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

	writeJSON(w, http.StatusOK, runtimeCompositionResponse{
		Scope:           "development-composition",
		Messages:        "configured",
		Receipts:        "configured",
		Attachments:     "configured",
		Authentication:  "accepted",
		Persistence:     runRuntimeDependencyProbe(r.Context(), h.persistenceProbe),
		Cryptography:    runRuntimeDependencyProbe(r.Context(), h.cryptographyProbe),
		ProductionReady: false,
	})
}

func runRuntimeDependencyProbe(ctx context.Context, probe interface{ Check(context.Context) error }) string {
	if probe == nil {
		return "not-assessed"
	}
	probeContext, cancel := context.WithTimeout(ctx, runtimeDependencyProbeTimeout)
	defer cancel()
	if err := probe.Check(probeContext); err != nil {
		return "unavailable"
	}
	return "available"
}

const runtimeDependencyProbeTimeout = 500 * time.Millisecond
