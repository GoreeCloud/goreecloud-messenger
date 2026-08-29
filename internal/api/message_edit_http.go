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

// MessageEditHTTPHandler exposes encrypted message-revision transport without decrypting content.
type MessageEditHTTPHandler struct {
	service *messagingservice.MessageEditService
	auth    Authenticator
}

func NewMessageEditHTTPHandler(service *messagingservice.MessageEditService, auth Authenticator) (*MessageEditHTTPHandler, error) {
	if service == nil {
		return nil, errors.New("message edit service is required")
	}
	if auth == nil {
		return nil, errors.New("authenticator is required")
	}
	return &MessageEditHTTPHandler{service: service, auth: auth}, nil
}

func (h *MessageEditHTTPHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/data/messages/{messageID}/edits", h.recordEdit)
	mux.HandleFunc("GET /v1/data/messages/{messageID}/edits", h.listEdits)
	return mux
}

type messageEditRequest struct {
	EditID         string                 `json:"edit_id"`
	ConversationID string                 `json:"conversation_id"`
	EditorID       string                 `json:"editor_id"`
	ClientNonce    string                 `json:"client_nonce"`
	Ciphertext     string                 `json:"ciphertext"`
	Encryption     domain.EncryptionState `json:"encryption"`
	EditedAt       time.Time              `json:"edited_at"`
}

type messageEditResponse struct {
	EditID         string                 `json:"edit_id"`
	MessageID      string                 `json:"message_id"`
	ConversationID string                 `json:"conversation_id"`
	EditorID       string                 `json:"editor_id"`
	ClientNonce    string                 `json:"client_nonce"`
	Ciphertext     string                 `json:"ciphertext"`
	Encryption     domain.EncryptionState `json:"encryption"`
	EditedAt       time.Time              `json:"edited_at"`
}

func (h *MessageEditHTTPHandler) recordEdit(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authenticate(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input messageEditRequest
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

	edit := domain.MessageEdit{
		EditID:         input.EditID,
		MessageID:      strings.TrimSpace(r.PathValue("messageID")),
		ConversationID: input.ConversationID,
		EditorID:       input.EditorID,
		ClientNonce:    input.ClientNonce,
		Ciphertext:     ciphertext,
		Encryption:     input.Encryption,
		EditedAt:       input.EditedAt,
	}
	if err := h.service.Record(r.Context(), userID, edit); err != nil {
		writeMessageEditServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *MessageEditHTTPHandler) listEdits(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.authenticate(w, r)
	if !ok {
		return
	}

	messageID := strings.TrimSpace(r.PathValue("messageID"))
	edits, err := h.service.List(r.Context(), userID, messageID)
	if err != nil {
		writeMessageEditServiceError(w, err)
		return
	}

	response := make([]messageEditResponse, 0, len(edits))
	for _, edit := range edits {
		response = append(response, messageEditResponse{
			EditID:         edit.EditID,
			MessageID:      edit.MessageID,
			ConversationID: edit.ConversationID,
			EditorID:       edit.EditorID,
			ClientNonce:    edit.ClientNonce,
			Ciphertext:     base64.StdEncoding.EncodeToString(edit.Ciphertext),
			Encryption:     edit.Encryption,
			EditedAt:       edit.EditedAt,
		})
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *MessageEditHTTPHandler) authenticate(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, err := h.auth.Authenticate(r.Context(), r)
	if err != nil || strings.TrimSpace(userID) == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return "", false
	}
	return userID, true
}

func writeMessageEditServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, messagingservice.ErrConversationAccess):
		writeError(w, http.StatusForbidden, "conversation access denied")
	case errors.Is(err, messagingservice.ErrEditUserMismatch):
		writeError(w, http.StatusForbidden, "editor does not match authenticated user")
	case errors.Is(err, messagingservice.ErrEditNotSender):
		writeError(w, http.StatusForbidden, "only the original sender may edit this message")
	case errors.Is(err, messagingservice.ErrMessageNotFound):
		writeError(w, http.StatusNotFound, "message not found")
	case errors.Is(err, messagingservice.ErrDuplicateEdit), errors.Is(err, messagingservice.ErrEditNonceReuse), errors.Is(err, messagingservice.ErrEditBeforeMessage):
		writeError(w, http.StatusConflict, "message edit conflicts with existing state")
	default:
		writeError(w, http.StatusBadRequest, "message edit rejected")
	}
}
