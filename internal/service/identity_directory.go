// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

var (
	ErrIdentityDirectoryInvalidRequest  = errors.New("identity directory request is invalid")
	ErrIdentityDirectoryUnavailable     = errors.New("identity directory is unavailable")
	ErrIdentityDirectoryInvalidResponse = errors.New("identity directory returned an invalid projection")
)

// IdentityDirectoryResolver is the future authenticated service-to-service seam
// to GoreeCloud Identity's privacy-preserving exact-handle consumer directory.
//
// The resolver owns service authentication and must derive Messenger's verified
// service principal from trusted runtime configuration. No client-supplied
// requester-service identifier is accepted by this interface. A false resolved
// value with a nil error intentionally represents all uniform negative Identity
// outcomes, including nonexistent, private, and disclosure-unauthorized users.
type IdentityDirectoryResolver interface {
	ResolveExact(context.Context, string) (projection domain.IdentityDirectoryProjection, resolved bool, err error)
}

type IdentityDirectoryService struct {
	resolver IdentityDirectoryResolver
}

func NewIdentityDirectoryService(resolver IdentityDirectoryResolver) (*IdentityDirectoryService, error) {
	if resolver == nil {
		return nil, errors.New("Identity directory resolver is required")
	}
	return &IdentityDirectoryService{resolver: resolver}, nil
}

// ResolveExact forwards only the user-supplied exact handle to the injected
// Identity resolver. Messenger does not normalize, browse, prefix-search, or
// fuzzy-search Identity accounts and never receives Identity discovery-policy
// details. Invalid provider projections fail closed before they reach clients.
func (s *IdentityDirectoryService) ResolveExact(
	ctx context.Context,
	handle string,
) (domain.IdentityDirectoryProjection, bool, error) {
	if handle == "" {
		return domain.IdentityDirectoryProjection{}, false, ErrIdentityDirectoryInvalidRequest
	}

	projection, resolved, err := s.resolver.ResolveExact(ctx, handle)
	if err != nil {
		if errors.Is(err, ErrIdentityDirectoryInvalidRequest) {
			return domain.IdentityDirectoryProjection{}, false, ErrIdentityDirectoryInvalidRequest
		}
		return domain.IdentityDirectoryProjection{}, false, ErrIdentityDirectoryUnavailable
	}
	if !resolved {
		return domain.IdentityDirectoryProjection{}, false, nil
	}
	if err := domain.ValidateIdentityDirectoryProjection(projection); err != nil {
		return domain.IdentityDirectoryProjection{}, false, ErrIdentityDirectoryInvalidResponse
	}
	return projection, true, nil
}
