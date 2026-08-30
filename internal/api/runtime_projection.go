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

// RuntimeCryptographyProbe is the dependency-level diagnostic retained for
// compatibility with the previous development composition slice.
type RuntimeCryptographyProbe interface {
	Check(context.Context) error
}

type RuntimeCryptographyProbeFunc func(context.Context) error

func (f RuntimeCryptographyProbeFunc) Check(ctx context.Context) error { return f(ctx) }

// RuntimeCryptographySessionBoundary is a stronger composition seam: it receives
// the already-authenticated user identifier and may verify that the composed
// cryptographic runtime can resolve/check that user's session state. It must not
// return session identifiers, algorithms, key material, ciphertext, credentials,
// or other sensitive details to this HTTP boundary.
type RuntimeCryptographySessionBoundary interface {
	CheckSession(context.Context, string) error
}

type RuntimeCryptographySessionBoundaryFunc func(context.Context, string) error

func (f RuntimeCryptographySessionBoundaryFunc) CheckSession(ctx context.Context, userID string) error {
	return f(ctx, userID)
}

type runtimeCompositionResponse struct {
	Scope              string `json:"scope"`
	Messages           string `json:"messages"`
	Receipts           string `json:"receipts"`
	Attachments        string `json:"attachments"`
	Authentication     string `json:"authentication"`
	Persistence        string `json:"persistence"`
	Cryptography       string `json:"cryptography"`
	CryptographyScope  string `json:"cryptography_scope"`
	ProductionReady    bool   `json:"production_ready"`
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

	cryptography, cryptographyScope := h.runtimeCryptographyState(r.Context(), userID)
	writeJSON(w, http.StatusOK, runtimeCompositionResponse{
		Scope:             "development-composition",
		Messages:          "configured",
		Receipts:          "configured",
		Attachments:       "configured",
		Authentication:    "accepted",
		Persistence:       runRuntimeDependencyProbe(r.Context(), h.persistenceProbe),
		Cryptography:      cryptography,
		CryptographyScope: cryptographyScope,
		ProductionReady:   false,
	})
}

func (h *DataRuntimeHandler) runtimeCryptographyState(ctx context.Context, userID string) (string, string) {
	if h.cryptographySessionBoundary != nil {
		probeContext, cancel := context.WithTimeout(ctx, runtimeDependencyProbeTimeout)
		defer cancel()
		if err := h.cryptographySessionBoundary.CheckSession(probeContext, userID); err != nil {
			return "unavailable", "authenticated-session"
		}
		return "available", "authenticated-session"
	}
	if h.cryptographyProbe != nil {
		return runRuntimeDependencyProbe(ctx, h.cryptographyProbe), "dependency"
	}
	return "not-assessed", "not-assessed"
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
