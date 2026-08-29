// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

// MessageDeletionHTTPHandler exposes delete-for-everyone tombstones without exposing message content.
type MessageDeletionHTTPHandler struct {
	service *messagingservice.MessageDeletionService
	auth    Authenticator
}

func NewMessageDeletionHTTPHandler(service *messagingservice.MessageDeletionService, auth Authenticator) (*MessageDeletionHTTPHandler, error) {
	if service == nil {
		return nil, errors.New("message deletion service is required")
	}
	if auth == nil {
		return nil, errors.New("authenticator is required")
	}
	return &MessageDeletionHTTPHandler{service: service, auth: auth}, nil
}

func (h *MessageDeletionHTTPHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/data/messages/{messageID}/deletions", h.recordDeletion)
	mux.HandleFunc("GET /v1/data/messages/{messageID}/deletions", h.listDeletions)
	return mux
}

type messageDeletionRequest struct {
	DeletionID     string                      `json:"deletion_id"`
	ConversationID string                      `json:"conversation_id"`
	DeleterID      string                      `json:"deleter_id"`
	ClientNonce    string                      `json:"client_nonce"`
	Scope          domain.MessageDeletionScope `json:"scope"`
	DeletedAt      time.Time                   `json:"deleted_at"`
}

type messageDeletionResponse struct {
	DeletionID     string                      `json:"deletion_id"`
	MessageID      string                      `json:"message_id"`
	ConversationID string                      `json:"conversation_id"`
	DeleterID      string                      `json:"deleter_id"`
	Scope          domain.MessageDeletionScope `json:"scope"`
	DeletedAt      time.Time                   `json:"deleted_at"`
}

func (h *MessageDeletionHTTPHandler) recordDeletion(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authenticate(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input messageDeletionRequest
	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if decoder.Decode(&struct{}{}) == nil {
		writeError(w, http.StatusBadRequest, "request body must contain one JSON object")
		return
	}

	deletion := domain.MessageDeletion{
		DeletionID:     input.DeletionID,
		MessageID:      strings.TrimSpace(r.PathValue("messageID")),
		ConversationID: input.ConversationID,
		DeleterID:      input.DeleterID,
		ClientNonce:    input.ClientNonce,
		Scope:          input.Scope,
		DeletedAt:      input.DeletedAt,
	}
	if err := h.service.Record(r.Context(), userID, deletion); err != nil {
		writeMessageDeletionServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *MessageDeletionHTTPHandler) listDeletions(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authenticate(w, r)
	if !ok {
		return
	}

	messageID := strings.TrimSpace(r.PathValue("messageID"))
	deletions, err := h.service.List(r.Context(), userID, messageID)
	if err != nil {
		writeMessageDeletionServiceError(w, err)
		return
	}

	response := make([]messageDeletionResponse, 0, len(deletions))
	for _, deletion := range deletions {
		response = append(response, messageDeletionResponse{
			DeletionID:     deletion.DeletionID,
			MessageID:      deletion.MessageID,
			ConversationID: deletion.ConversationID,
			DeleterID:      deletion.DeleterID,
			Scope:          deletion.Scope,
			DeletedAt:      deletion.DeletedAt,
		})
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *MessageDeletionHTTPHandler) authenticate(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, err := h.auth.Authenticate(r.Context(), r)
	if err != nil || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return "", false
	}
	return userID, true
}

func writeMessageDeletionServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, messagingservice.ErrConversationAccess):
		writeError(w, http.StatusForbidden, "conversation access denied")
	case errors.Is(err, messagingservice.ErrDeletionUserMismatch):
		writeError(w, http.StatusForbidden, "deleter does not match authenticated user")
	case errors.Is(err, messagingservice.ErrDeletionNotSender):
		writeError(w, http.StatusForbidden, "only the original sender may delete this message for everyone")
	case errors.Is(err, messagingservice.ErrMessageNotFound):
		writeError(w, http.StatusNotFound, "message not found")
	case errors.Is(err, messagingservice.ErrDuplicateDeletion),
		errors.Is(err, messagingservice.ErrDeletionNonceReuse),
		errors.Is(err, messagingservice.ErrDeletionBeforeMessage),
		errors.Is(err, messagingservice.ErrMessageAlreadyDeleted):
		writeError(w, http.StatusConflict, "message deletion conflicts with existing state")
	default:
		writeError(w, http.StatusBadRequest, "message deletion rejected")
	}
}
