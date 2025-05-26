// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package kratos

import (
	"context"

	"github.com/ory/fosite"
	client "github.com/ory/kratos-client-go"
)

type (
	FakeKratos struct {
		DisableSessionWasCalled bool
		DisableSessionCB        func()
		LastDisabledSession     string
	}
)

const (
	FakeSessionID  = "fake-kratos-session-id"
	FakeUsername   = "fake-kratos-username"
	FakePassword   = "fake-kratos-password" // nolint: gosec
	FakeIdentityID = "fake-kratos-identity-id"
)

var _ Client = new(FakeKratos)

func NewFake() *FakeKratos {
	return &FakeKratos{}
}

func (f *FakeKratos) DisableSession(_ context.Context, identityProviderSessionID string) error {
	f.DisableSessionWasCalled = true
	f.LastDisabledSession = identityProviderSessionID
	if f.DisableSessionCB != nil {
		f.DisableSessionCB()
	}

	return nil
}

func (f *FakeKratos) Authenticate(_ context.Context, username, password string) (*client.Session, error) {
	if username == FakeUsername && password == FakePassword {
		return &client.Session{Identity: &client.Identity{Id: FakeIdentityID}}, nil
	}
	return nil, fosite.ErrNotFound
}

func (f *FakeKratos) Reset() {
	(*f) = *NewFake()
}
