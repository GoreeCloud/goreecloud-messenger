// SPDX-License-Identifier: AGPL-3.0-only

package domain

import (
	"errors"
	"strings"
	"time"
)

// ReceiptState records recipient-observed delivery progress for a GoreeCloud Data message.
type ReceiptState string

const (
	ReceiptDelivered ReceiptState = "delivered"
	ReceiptRead      ReceiptState = "read"
)

// DeliveryReceipt is a recipient-authenticated acknowledgement for one Data message.
// It intentionally contains no plaintext message content or carrier-transport state.
type DeliveryReceipt struct {
	MessageID      string
	ConversationID string
	UserID         string
	State          ReceiptState
	ObservedAt     time.Time
}

func (r DeliveryReceipt) Validate() error {
	if strings.TrimSpace(r.MessageID) == "" {
		return errors.New("receipt message id is required")
	}
	if strings.TrimSpace(r.ConversationID) == "" {
		return errors.New("receipt conversation id is required")
	}
	if strings.TrimSpace(r.UserID) == "" {
		return errors.New("receipt user id is required")
	}
	switch r.State {
	case ReceiptDelivered, ReceiptRead:
	default:
		return errors.New("receipt state is invalid")
	}
	if r.ObservedAt.IsZero() {
		return errors.New("receipt observation timestamp is required")
	}
	return nil
}

func (s ReceiptState) Rank() int {
	switch s {
	case ReceiptDelivered:
		return 1
	case ReceiptRead:
		return 2
	default:
		return 0
	}
}
