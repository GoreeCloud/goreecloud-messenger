// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

// MessageReactionHTTPHandler exposes encrypted reaction state without decoding reaction values.
type MessageReactionHTTPHandler struct {
	service *messagingservice.MessageReactionService
	auth    Authenticator
}

func NewMessageReactionHTTPHandler(service *messagingservice.MessageReactionService, auth Authenticator) (*MessageReactionHTTPHandler, error) {
	if service == nil {
		return nil, errors.New("message reaction service is required")
	}
	if auth == nil {
		return nil, errors.New("authenticator is required")
	}
	return &MessageReactionHTTPHandler{service: service, auth: auth}, nil
}

func (h *MessageReactionHTTPHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/data/messages/{messageID}/reactions", h.recordReaction)
	mux.HandleFunc("GET /v1/data/messages/{messageID}/reactions", h.listReactions)
	return mux
}

type messageReactionRequest struct {
	ReactionID     string                   `json:"reaction_id"`
	ConversationID string                   `json:"conversation_id"`
	ReactorID      string                   `json:"reactor_id"`
	ClientNonce    string                   `json:"client_nonce"`
	Operation      domain.ReactionOperation `json:"operation"`
	Ciphertext     string                   `json:"ciphertext,omitempty"`
	Encryption     domain.EncryptionState   `json:"encryption"`
	ReactedAt      time.Time                `json:"reacted_at"`
}

type messageReactionResponse struct {
	ReactionID     string                 `json:"reaction_id"`
	MessageID      string                 `json:"message_id"`
	ConversationID string                 `json:"conversation_id"`
	ReactorID      string                 `json:"reactor_id"`
	Ciphertext     string                 `json:"ciphertext"`
	Encryption     domain.EncryptionState `json:"encryption"`
	ReactedAt      time.Time              `json:"reacted_at"`
}

func (h *MessageReactionHTTPHandler) recordReaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authenticate(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var input messageReactionRequest
	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if decoder.Decode(&struct{}{}) == nil {
		writeError(w, http.StatusBadRequest, "request body must contain one JSON object")
		return
	}

	var ciphertext []byte
	if input.Ciphertext != "" {
		decoded, err := base64.StdEncoding.DecodeString(input.Ciphertext)
		if err != nil {
			writeError(w, http.StatusBadRequest, "ciphertext must be base64 encoded")
			return
		}
		ciphertext = decoded
	}

	reaction := domain.MessageReaction{
		ReactionID:     input.ReactionID,
		MessageID:      strings.TrimSpace(r.PathValue("messageID")),
		ConversationID: input.ConversationID,
		ReactorID:      input.ReactorID,
		ClientNonce:    input.ClientNonce,
		Operation:      input.Operation,
		Ciphertext:     ciphertext,
		Encryption:     input.Encryption,
		ReactedAt:      input.ReactedAt,
	}
	if err := h.service.Record(r.Context(), userID, reaction); err != nil {
		writeMessageReactionServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *MessageReactionHTTPHandler) listReactions(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authenticate(w, r)
	if !ok {
		return
	}

	messageID := strings.TrimSpace(r.PathValue("messageID"))
	reactions, err := h.service.List(r.Context(), userID, messageID)
	if err != nil {
		writeMessageReactionServiceError(w, err)
		return
	}

	response := make([]messageReactionResponse, 0, len(reactions))
	for _, reaction := range reactions {
		response = append(response, messageReactionResponse{
			ReactionID:     reaction.ReactionID,
			MessageID:      reaction.MessageID,
			ConversationID: reaction.ConversationID,
			ReactorID:      reaction.ReactorID,
			Ciphertext:     base64.StdEncoding.EncodeToString(reaction.Ciphertext),
			Encryption:     reaction.Encryption,
			ReactedAt:      reaction.ReactedAt,
		})
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *MessageReactionHTTPHandler) authenticate(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, err := h.auth.Authenticate(r.Context(), r)
	if err != nil || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return "", false
	}
	return userID, true
}

func writeMessageReactionServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, messagingservice.ErrConversationAccess):
		writeError(w, http.StatusForbidden, "conversation access denied")
	case errors.Is(err, messagingservice.ErrReactionUserMismatch):
		writeError(w, http.StatusForbidden, "reactor does not match authenticated user")
	case errors.Is(err, messagingservice.ErrReactionTargetUnavailable):
		writeError(w, http.StatusBadRequest, "reaction target unavailable")
	case errors.Is(err, messagingservice.ErrDuplicateReaction), errors.Is(err, messagingservice.ErrReactionNonceReuse), errors.Is(err, messagingservice.ErrReactionBeforeMessage), errors.Is(err, messagingservice.ErrReactionStale):
		writeError(w, http.StatusConflict, "message reaction conflicts with existing state")
	default:
		writeError(w, http.StatusBadRequest, "message reaction rejected")
	}
}
