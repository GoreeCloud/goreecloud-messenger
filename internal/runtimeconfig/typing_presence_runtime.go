// SPDX-License-Identifier: AGPL-3.0-only

package runtimeconfig

import (
	"fmt"
	"time"

	"github.com/GoreeCloud/goreecloud-messenger/internal/service"
)

// TypingPresenceRuntime is the Development composition of ephemeral typing
// presence and its authenticated privacy-preference control surface. Both
// services are deliberately backed by the same TypingPrivacyRuntimeStore.
type TypingPresenceRuntime struct {
	Typing          *service.TypingService
	Preferences     *service.TypingPrivacyPreferenceService
	PersistenceMode service.TypingPrivacyPersistenceMode
}

// TypingPresenceRuntimeFromEnvironment resolves the explicit process-level
// typing privacy persistence selection, constructs exactly one shared privacy
// store, and injects it into both typing authorization and preference mutation.
func TypingPresenceRuntimeFromEnvironment(
	typingStore service.TypingStore,
	access service.ConversationAccess,
	now func() time.Time,
) (TypingPresenceRuntime, error) {
	config, err := TypingPrivacyPersistenceFromEnvironment()
	if err != nil {
		return TypingPresenceRuntime{}, err
	}
	return typingPresenceRuntimeFromConfig(config, typingStore, access, now)
}

func typingPresenceRuntimeFromConfig(
	config service.TypingPrivacyPersistenceConfig,
	typingStore service.TypingStore,
	access service.ConversationAccess,
	now func() time.Time,
) (TypingPresenceRuntime, error) {
	privacyStore, err := service.NewTypingPrivacyRuntimeStore(config)
	if err != nil {
		return TypingPresenceRuntime{}, fmt.Errorf("initialize typing privacy runtime store: %w", err)
	}
	typing, err := service.NewTypingService(typingStore, access, privacyStore, now)
	if err != nil {
		return TypingPresenceRuntime{}, fmt.Errorf("initialize typing service: %w", err)
	}
	preferences, err := service.NewTypingPrivacyPreferenceService(privacyStore, access)
	if err != nil {
		return TypingPresenceRuntime{}, fmt.Errorf("initialize typing privacy preference service: %w", err)
	}
	return TypingPresenceRuntime{
		Typing:          typing,
		Preferences:     preferences,
		PersistenceMode: config.Mode,
	}, nil
}
