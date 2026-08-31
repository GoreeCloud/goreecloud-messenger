// SPDX-License-Identifier: AGPL-3.0-only

package api

import (
	"errors"
	"net/http"

	messagingservice "github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

// DataRuntimeHandler is the application-facing HTTP composition boundary for
// GoreeCloud Data messaging, receipts, encrypted attachments, and optional
// privacy-controlled typing presence. It does not own credential validation,
// cryptographic sessions, or persistence authority; those remain injected
// through the existing service and Authenticator boundaries.
type DataRuntimeHandler struct {
	messages          *Handler
	attachments       *AttachmentHTTPHandler
	typing            *TypingHTTPHandler
	typingPreferences *TypingPreferencesHTTPHandler
	auth              Authenticator
	persistenceProbe  RuntimePersistenceProbe
	cryptographyProbe RuntimeCryptographyProbe
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

// WithTypingIndicators returns a copy of the runtime with the already-authorized
// typing service composed under the same Authenticator boundary. Typing remains
// optional so deployments do not gain presence behavior merely by upgrading the
// message runtime.
func (h *DataRuntimeHandler) WithTypingIndicators(service *messagingservice.TypingService) (*DataRuntimeHandler, error) {
	if h == nil || service == nil {
		return nil, errors.New("runtime handler and typing service are required")
	}
	typing, err := NewTypingHTTPHandler(service, h.auth)
	if err != nil {
		return nil, err
	}
	copy := *h
	copy.typing = typing
	return &copy, nil
}

// WithTypingPrivacyPreferences returns a copy of the runtime with authenticated
// per-conversation typing privacy controls. This is independently optional so a
// deployment must deliberately compose the mutable preference surface.
func (h *DataRuntimeHandler) WithTypingPrivacyPreferences(
	service *messagingservice.TypingPrivacyPreferenceService,
) (*DataRuntimeHandler, error) {
	if h == nil || service == nil {
		return nil, errors.New("runtime handler and typing privacy preference service are required")
	}
	preferences, err := NewTypingPreferencesHTTPHandler(service, h.auth)
	if err != nil {
		return nil, err
	}
	copy := *h
	copy.typingPreferences = preferences
	return &copy, nil
}

// WithTypingPresence returns a copy of the runtime with both ephemeral typing
// publication/observation and the matching authenticated privacy controls
// composed together. Deployments can still opt into either surface separately,
// but this helper prevents an application runtime from accidentally enabling
// typing events while omitting the user-facing control surface.
func (h *DataRuntimeHandler) WithTypingPresence(
	typingService *messagingservice.TypingService,
	preferenceService *messagingservice.TypingPrivacyPreferenceService,
) (*DataRuntimeHandler, error) {
	if h == nil || typingService == nil || preferenceService == nil {
		return nil, errors.New("runtime handler, typing service, and typing privacy preference service are required")
	}
	withTyping, err := h.WithTypingIndicators(typingService)
	if err != nil {
		return nil, err
	}
	return withTyping.WithTypingPrivacyPreferences(preferenceService)
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

// WithRuntimeCryptographyProbe returns a copy of the composition handler with a
// bounded, diagnostic-only cryptography dependency probe. The probe may report
// only categorical availability through the minimized runtime projection; it
// does not expose session, algorithm, key, credential, or identity details and
// does not establish production cryptography readiness.
func (h *DataRuntimeHandler) WithRuntimeCryptographyProbe(probe RuntimeCryptographyProbe) (*DataRuntimeHandler, error) {
	if h == nil || probe == nil {
		return nil, errors.New("runtime handler and cryptography probe are required")
	}
	copy := *h
	copy.cryptographyProbe = probe
	return &copy, nil
}

// Routes returns one mux containing every enabled Data HTTP route. Shared
// conversation prefixes are registered directly so message, attachment, and
// optional typing endpoints remain simultaneously reachable.
func (h *DataRuntimeHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	h.messages.RegisterRoutes(mux)
	h.attachments.RegisterRoutes(mux)
	if h.typing != nil {
		h.typing.RegisterRoutes(mux)
	}
	if h.typingPreferences != nil {
		h.typingPreferences.RegisterRoutes(mux)
	}
	h.registerRuntimeProjection(mux)
	return mux
}
