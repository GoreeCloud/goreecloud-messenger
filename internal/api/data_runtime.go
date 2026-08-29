// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"errors"
	"net/http"

	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

// DataRuntimeHandler is the application-facing HTTP composition boundary for
// GoreeCloud Data messaging, receipts, and encrypted attachments. It does not
// own credential validation, cryptographic sessions, or persistence authority;
// those remain injected through the existing service and Authenticator boundaries.
type DataRuntimeHandler struct {
	messages    *Handler
	attachments *AttachmentHTTPHandler
}

func NewDataRuntimeHandler(
	data *messagingservice.DataService,
	receipts *messagingservice.ReceiptService,
	attachments *messagingservice.AttachmentService,
	auth Authenticator,
) (*DataRuntimeHandler, error) {
	if data == nil || receipts == nil || attachments == nil || auth == nil {
		return nil, errors.New("Data, receipt, attachment, and authentication boundaries are required")
	}

	messageHandler, err := NewHandler(data, receipts, auth)
	if err != nil {
		return nil, err
	}
	attachmentHandler, err := NewAttachmentHTTPHandler(attachments, auth)
	if err != nil {
		return nil, err
	}
	return &DataRuntimeHandler{messages: messageHandler, attachments: attachmentHandler}, nil
}

// Routes returns one mux containing every currently implemented Data HTTP route.
// Shared conversation prefixes are registered directly so message and attachment
// listing endpoints remain simultaneously reachable.
func (h *DataRuntimeHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	h.messages.RegisterRoutes(mux)
	h.attachments.RegisterRoutes(mux)
	return mux
}
