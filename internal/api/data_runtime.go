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
	messages                    *Handler
	attachments                 *AttachmentHTTPHandler
	auth                        Authenticator
	persistenceProbe            RuntimePersistenceProbe
	cryptographyProbe           RuntimeCryptographyProbe
	cryptographySessionBoundary RuntimeCryptographySessionBoundary
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
	return &DataRuntimeHandler{
		messages:    messageHandler,
		attachments: attachmentHandler,
		auth:        auth,
	}, nil
}

// WithRuntimePersistenceProbe returns a copy of the composition handler with a
// bounded, diagnostic-only persistence probe. The probe does not gain authority
// over message, receipt, or attachment operations and is never treated as a
// production-readiness decision.
func (h *DataRuntimeHandler) WithRuntimePersistenceProbe(probe RuntimePersistenceProbe) (*DataRuntimeHandler, error) {
	if h == nil || probe == nil {
		return nil, errors.New("runtime handler and persistence probe are required")
	}
	copy := *h
	copy.persistenceProbe = probe
	return &copy, nil
}

// WithRuntimeCryptographyProbe preserves the dependency-level diagnostic from
// the previous development slice. Prefer WithRuntimeCryptographySessionBoundary
// when the composed cryptography implementation can validate an authenticated
// user's session state without exporting session details.
func (h *DataRuntimeHandler) WithRuntimeCryptographyProbe(probe RuntimeCryptographyProbe) (*DataRuntimeHandler, error) {
	if h == nil || probe == nil {
		return nil, errors.New("runtime handler and cryptography probe are required")
	}
	copy := *h
	copy.cryptographyProbe = probe
	return &copy, nil
}

// WithRuntimeCryptographySessionBoundary injects a bounded, authenticated-user
// session check into the minimized runtime projection. The user identifier is
// passed only to the injected boundary and is never serialized in the response.
// This composition check is diagnostic development evidence, not E2EE or
// production-readiness acceptance.
func (h *DataRuntimeHandler) WithRuntimeCryptographySessionBoundary(
	boundary RuntimeCryptographySessionBoundary,
) (*DataRuntimeHandler, error) {
	if h == nil || boundary == nil {
		return nil, errors.New("runtime handler and cryptography session boundary are required")
	}
	copy := *h
	copy.cryptographySessionBoundary = boundary
	return &copy, nil
}

// Routes returns one mux containing every currently implemented Data HTTP route.
// Shared conversation prefixes are registered directly so message and attachment
// listing endpoints remain simultaneously reachable.
func (h *DataRuntimeHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	h.messages.RegisterRoutes(mux)
	h.attachments.RegisterRoutes(mux)
	h.registerRuntimeProjection(mux)
	return mux
}
