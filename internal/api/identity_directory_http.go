// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

const maxIdentityDirectoryRequestBodyBytes = 4096

type IdentityDirectoryHTTPHandler struct {
	service *messagingservice.IdentityDirectoryService
	auth    Authenticator
}

func NewIdentityDirectoryHTTPHandler(
	service *messagingservice.IdentityDirectoryService,
	auth Authenticator,
) (*IdentityDirectoryHTTPHandler, error) {
	if service == nil || auth == nil {
		return nil, errors.New("Identity directory service and authenticator are required")
	}
	return &IdentityDirectoryHTTPHandler{service: service, auth: auth}, nil
}

func (h *IdentityDirectoryHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	if mux == nil {
		panic("HTTP mux is required")
	}
	mux.HandleFunc("POST /v1/identity/resolve", h.resolveExactHandle)
}

type identityDirectoryResolveRequest struct {
	Handle string `json:"handle"`
}

type identityDirectoryResolveResponse struct {
	Subject     string `json:"subject"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name,omitempty"`
}

func (h *IdentityDirectoryHTTPHandler) resolveExactHandle(w http.ResponseWriter, r *http.Request) {
	if !h.authenticateIdentityDirectoryRequest(w, r) {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxIdentityDirectoryRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input identityDirectoryResolveRequest
	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// A request is valid only when the first JSON object consumes the complete body
	// apart from JSON whitespace. A second successful value and trailing malformed
	// JSON both fail closed instead of allowing parser ambiguity before resolution.
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "request body must contain one JSON object")
		return
	}

	projection, resolved, err := h.service.ResolveExact(r.Context(), input.Handle)
	if err != nil {
		switch {
		case errors.Is(err, messagingservice.ErrIdentityDirectoryInvalidRequest):
			writeError(w, http.StatusBadRequest, "identity handle request rejected")
		case errors.Is(err, messagingservice.ErrIdentityDirectoryUnavailable),
			errors.Is(err, messagingservice.ErrIdentityDirectoryInvalidResponse):
			writeError(w, http.StatusServiceUnavailable, "identity directory unavailable")
		default:
			writeError(w, http.StatusServiceUnavailable, "identity directory unavailable")
		}
		return
	}
	if !resolved {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"error": map[string]string{
				"code":    "not_resolved",
				"message": "The requested identity could not be resolved.",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, identityDirectoryResolveResponse{
		Subject:     projection.Subject,
		Handle:      projection.Handle,
		DisplayName: projection.DisplayName,
	})
}

func (h *IdentityDirectoryHTTPHandler) authenticateIdentityDirectoryRequest(
	w http.ResponseWriter,
	r *http.Request,
) bool {
	userID, err := h.auth.Authenticate(r.Context(), r)
	if err != nil || domain.ValidateOpaqueIdentifier(userID, "authenticated user id") != nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return false
	}
	return true
}
