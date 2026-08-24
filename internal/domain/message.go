// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Transport identifies the actual network transport used for one message.
type Transport string

const (
	TransportData Transport = "data"
	TransportSMS  Transport = "sms"
	TransportMMS  Transport = "mms"
	TransportRCS  Transport = "rcs"
)

func (t Transport) Valid() bool {
	switch t {
	case TransportData, TransportSMS, TransportMMS, TransportRCS:
		return true
	default:
		return false
	}
}

// EncryptionState records the evidence-backed protection state for a message.
type EncryptionState string

const (
	EncryptionNone    EncryptionState = "none"
	EncryptionE2EE    EncryptionState = "e2ee"
	EncryptionUnknown EncryptionState = "unknown"
)

func (e EncryptionState) Valid() bool {
	switch e {
	case EncryptionNone, EncryptionE2EE, EncryptionUnknown:
		return true
	default:
		return false
	}
}

// DeliveryState records message lifecycle state independently from transport.
type DeliveryState string

const (
	DeliveryPending   DeliveryState = "pending"
	DeliverySent      DeliveryState = "sent"
	DeliveryDelivered DeliveryState = "delivered"
	DeliveryRead      DeliveryState = "read"
	DeliveryFailed    DeliveryState = "failed"
)

func (d DeliveryState) Valid() bool {
	switch d {
	case DeliveryPending, DeliverySent, DeliveryDelivered, DeliveryRead, DeliveryFailed:
		return true
	default:
		return false
	}
}

// Identity is a GoreeCloud user identity. Username is first-class and does not require a phone number.
type Identity struct {
	UserID        string
	Username      string
	DisplayName   string
	VerifiedPhone string
}

func (i Identity) Validate() error {
	if strings.TrimSpace(i.UserID) == "" {
		return errors.New("identity user id is required")
	}
	if strings.TrimSpace(i.Username) == "" {
		return errors.New("identity username is required")
	}
	if !strings.HasPrefix(i.Username, "@") {
		return errors.New("identity username must begin with @")
	}
	if len(i.Username) < 2 {
		return errors.New("identity username must contain characters after @")
	}
	return nil
}

// Message stores transport and protection provenance alongside content.
type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	Body           string
	Transport      Transport
	Encryption     EncryptionState
	Delivery       DeliveryState
	SentAt         time.Time
}

func (m Message) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return errors.New("message id is required")
	}
	if strings.TrimSpace(m.ConversationID) == "" {
		return errors.New("conversation id is required")
	}
	if strings.TrimSpace(m.SenderID) == "" {
		return errors.New("sender id is required")
	}
	if strings.TrimSpace(m.Body) == "" {
		return errors.New("message body is required")
	}
	if !m.Transport.Valid() {
		return fmt.Errorf("unsupported transport %q", m.Transport)
	}
	if !m.Encryption.Valid() {
		return fmt.Errorf("unsupported encryption state %q", m.Encryption)
	}
	if !m.Delivery.Valid() {
		return fmt.Errorf("unsupported delivery state %q", m.Delivery)
	}
	if m.Encryption == EncryptionE2EE && m.Transport != TransportData {
		return errors.New("GoreeCloud E2EE may only be asserted for GoreeCloud Data transport")
	}
	if m.SentAt.IsZero() {
		return errors.New("sent timestamp is required")
	}
	return nil
}

// ProvenanceLabel is the compact user-facing state derived from verified message metadata.
func (m Message) ProvenanceLabel() string {
	transport := map[Transport]string{
		TransportData: "Data",
		TransportSMS:  "SMS",
		TransportMMS:  "MMS",
		TransportRCS:  "RCS",
	}[m.Transport]

	if m.Transport == TransportData && m.Encryption == EncryptionE2EE {
		return "E2EE · Data"
	}
	return transport
}
