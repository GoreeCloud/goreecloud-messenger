// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"fmt"
	"strings"
)

// CallKind distinguishes voice and video communication.
type CallKind string

const (
	CallVoice CallKind = "voice"
	CallVideo CallKind = "video"
)

func (k CallKind) Valid() bool {
	switch k {
	case CallVoice, CallVideo:
		return true
	default:
		return false
	}
}

// CallSession records communication and security provenance for a GoreeCloud call.
type CallSession struct {
	ID             string
	ConversationID string
	Kind           CallKind
	Transport      Transport
	Encryption     EncryptionState
}

func (c CallSession) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return errors.New("call id is required")
	}
	if strings.TrimSpace(c.ConversationID) == "" {
		return errors.New("call conversation id is required")
	}
	if !c.Kind.Valid() {
		return fmt.Errorf("unsupported call kind %q", c.Kind)
	}
	if !c.Transport.Valid() {
		return fmt.Errorf("unsupported call transport %q", c.Transport)
	}
	if !c.Encryption.Valid() {
		return fmt.Errorf("unsupported call encryption state %q", c.Encryption)
	}
	if c.Transport != TransportData {
		return errors.New("GoreeCloud voice and video calls require GoreeCloud Data transport")
	}
	return nil
}

// ProvenanceLabel returns the compact user-facing call state.
func (c CallSession) ProvenanceLabel() string {
	kind := map[CallKind]string{
		CallVoice: "Voice",
		CallVideo: "Video",
	}[c.Kind]

	if c.Transport == TransportData && c.Encryption == EncryptionE2EE {
		return kind + " · E2EE · Data"
	}
	if c.Transport == TransportData {
		return kind + " · Data"
	}
	return kind
}
