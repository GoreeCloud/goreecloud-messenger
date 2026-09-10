// SPDX-License-Identifier: AGPL-3.0-only

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/GoreeCloud/goreecloud-messenger/internal/domain"
)

type identityDirectoryResolverFunc func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error)

func (f identityDirectoryResolverFunc) ResolveExact(
	ctx context.Context,
	handle string,
) (domain.IdentityDirectoryProjection, bool, error) {
	return f(ctx, handle)
}

func TestIdentityDirectoryServiceForwardsExactHandleWithoutClaimingCanonicalization(t *testing.T) {
	var observed string
	service, err := NewIdentityDirectoryService(identityDirectoryResolverFunc(
		func(_ context.Context, handle string) (domain.IdentityDirectoryProjection, bool, error) {
			observed = handle
			return domain.IdentityDirectoryProjection{
				Subject:     "identity-subject-1",
				Handle:      "alice.example",
				DisplayName: "Alice Example",
			}, true, nil
		},
	))
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	projection, resolved, err := service.ResolveExact(context.Background(), " @Alice.Example ")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !resolved {
		t.Fatal("expected resolved identity")
	}
	if observed != " @Alice.Example " {
		t.Fatalf("resolver saw %q; Messenger must forward the supplied exact handle without normalization", observed)
	}
	if projection.Handle != "alice.example" || projection.Subject != "identity-subject-1" {
		t.Fatalf("unexpected projection: %#v", projection)
	}
}

func TestIdentityDirectoryServicePreservesUniformNotResolvedState(t *testing.T) {
	service, err := NewIdentityDirectoryService(identityDirectoryResolverFunc(
		func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			return domain.IdentityDirectoryProjection{}, false, nil
		},
	))
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	projection, resolved, err := service.ResolveExact(context.Background(), "private-or-missing")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if resolved || projection != (domain.IdentityDirectoryProjection{}) {
		t.Fatalf("uniform negative result leaked projection: resolved=%v projection=%#v", resolved, projection)
	}
}

func TestIdentityDirectoryServiceMapsUnexpectedResolverFailureToUnavailable(t *testing.T) {
	backendErr := errors.New("sensitive backend detail")
	service, err := NewIdentityDirectoryService(identityDirectoryResolverFunc(
		func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			return domain.IdentityDirectoryProjection{}, false, backendErr
		},
	))
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	_, _, err = service.ResolveExact(context.Background(), "alice")
	if !errors.Is(err, ErrIdentityDirectoryUnavailable) {
		t.Fatalf("error = %v, want unavailable", err)
	}
	if errors.Is(err, backendErr) {
		t.Fatal("backend error detail escaped the service boundary")
	}
}

func TestIdentityDirectoryServicePreservesExplicitInvalidRequest(t *testing.T) {
	service, err := NewIdentityDirectoryService(identityDirectoryResolverFunc(
		func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			return domain.IdentityDirectoryProjection{}, false, ErrIdentityDirectoryInvalidRequest
		},
	))
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	_, _, err = service.ResolveExact(context.Background(), "bad handle")
	if !errors.Is(err, ErrIdentityDirectoryInvalidRequest) {
		t.Fatalf("error = %v, want invalid request", err)
	}
}

func TestIdentityDirectoryServiceRejectsInvalidProviderProjection(t *testing.T) {
	service, err := NewIdentityDirectoryService(identityDirectoryResolverFunc(
		func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			return domain.IdentityDirectoryProjection{
				Subject: "identity-subject-1",
				Handle:  "Alice",
			}, true, nil
		},
	))
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	_, _, err = service.ResolveExact(context.Background(), "alice")
	if !errors.Is(err, ErrIdentityDirectoryInvalidResponse) {
		t.Fatalf("error = %v, want invalid response", err)
	}
}

func TestIdentityDirectoryServiceRejectsEmptyHandleWithoutCallingResolver(t *testing.T) {
	calls := 0
	service, err := NewIdentityDirectoryService(identityDirectoryResolverFunc(
		func(context.Context, string) (domain.IdentityDirectoryProjection, bool, error) {
			calls++
			return domain.IdentityDirectoryProjection{}, false, nil
		},
	))
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	_, _, err = service.ResolveExact(context.Background(), "")
	if !errors.Is(err, ErrIdentityDirectoryInvalidRequest) {
		t.Fatalf("error = %v, want invalid request", err)
	}
	if calls != 0 {
		t.Fatalf("resolver calls = %d, want 0", calls)
	}
}
