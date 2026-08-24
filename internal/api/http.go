// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

const maxRequestBodyBytes = 1 << 20

// Authenticator resolves an authenticated GoreeCloud user from an HTTP request.
// Credential issuance and validation remain external to this transport boundary.
type Authenticator interface {
	Authenticate(context.Context, *http.Request) (string, error)
}

// Handler exposes the GoreeCloud Data service without weakening its authorization boundary.
type Handler struct {
	service *messagingservice.DataService
	auth    Authenticator
}

func NewHandler(service *messagingservice.DataService, auth Authenticator) (*Handler, error) {
	if service == nil {
		return nil, errors.New("Data service is required")
	}
	if auth == nil {
		return nil, errors.New("authenticator is required")
	}
	return &Handler{service: service, auth: auth}, nil
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/data/messages", h.submitMessage)
	mux.HandleFunc("GET /v1/data/conversations/{conversationID}/messages", h.listConversation)
	return mux
}

type submitMessageRequest struct {
	MessageID      string                 `json:"message_id"`
	ConversationID string                 `json:"conversation_id"`
	SenderID       string                 `json:"sender_id"`
	ClientNonce    string                 `json:"client_nonce"`
	Ciphertext     string                 `json:"ciphertext"`
	Encryption     domain.EncryptionState `json:"encryption"`
	CreatedAt      time.Time              `json:"created_at"`
}

type messageResponse struct {
	MessageID      string                 `json:"message_id"`
	ConversationID string                 `json:"conversation_id"`
	SenderID       string                 `json:"sender_id"`
	ClientNonce    string                 `json:"client_nonce"`
	Ciphertext     string                 `json:"ciphertext"`
	Encryption     domain.EncryptionState `json:"encryption"`
	CreatedAt      time.Time              `json:"created_at"`
}

func (h *Handler) submitMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authenticate(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input submitMessageRequest
	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if decoder.Decode(&struct{}{}) == nil {
		writeError(w, http.StatusBadRequest, "request body must contain one JSON object")
		return
	}

	ciphertext, err := base64.StdEncoding.DecodeString(input.Ciphertext)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ciphertext must be base64 encoded")
		return
	}

	envelope := domain.DataEnvelope{
		MessageID:      input.MessageID,
		ConversationID: input.ConversationID,
		SenderID:       input.SenderID,
		ClientNonce:    input.ClientNonce,
		Ciphertext:     ciphertext,
		Encryption:     input.Encryption,
		CreatedAt:      input.CreatedAt,
	}

	if err := h.service.Submit(r.Context(), userID, envelope); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) listConversation(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authenticate(w, r)
	if !ok {
		return
	}

	conversationID := strings.TrimSpace(r.PathValue("conversationID"))
	messages, err := h.service.ListConversation(r.Context(), userID, conversationID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response := make([]messageResponse, 0, len(messages))
	for _, message := range messages {
		response = append(response, messageResponse{
			MessageID:      message.MessageID,
			ConversationID: message.ConversationID,
			SenderID:       message.SenderID,
			ClientNonce:    message.ClientNonce,
			Ciphertext:     base64.StdEncoding.EncodeToString(message.Ciphertext),
			Encryption:     message.Encryption,
			CreatedAt:      message.CreatedAt,
		})
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, err := h.auth.Authenticate(r.Context(), r)
	if err != nil || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return "", false
	}
	return userID, true
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, messagingservice.ErrConversationAccess):
		writeError(w, http.StatusForbidden, "conversation access denied")
	case errors.Is(err, messagingservice.ErrSenderMismatch):
		writeError(w, http.StatusForbidden, "sender does not match authenticated user")
	case errors.Is(err, messagingservice.ErrDuplicateMessage), errors.Is(err, messagingservice.ErrNonceReuse):
		writeError(w, http.StatusConflict, "message conflicts with existing state")
	default:
		writeError(w, http.StatusBadRequest, "message request rejected")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
