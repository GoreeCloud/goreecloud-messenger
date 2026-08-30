// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"net/http"
	"strings"
)

type runtimeCompositionResponse struct {
	Scope           string `json:"scope"`
	Messages        string `json:"messages"`
	Receipts        string `json:"receipts"`
	Attachments     string `json:"attachments"`
	Authentication  string `json:"authentication"`
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
		ProductionReady: false,
	})
}
