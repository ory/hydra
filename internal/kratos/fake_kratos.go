// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package kratos

import "context"

type (
	FakeKratos struct {
		DisableSessionWasCalled bool
		LastDisabledSession     string
	}
)

const (
	FakeSessionID = "fake-kratos-session-id"
)

var _ Client = new(FakeKratos)

func NewFake() *FakeKratos {
	return &FakeKratos{}
}

func (f *FakeKratos) DisableSession(ctx context.Context, identityProviderSessionID string) error {
	f.DisableSessionWasCalled = true
	f.LastDisabledSession = identityProviderSessionID

	return nil
}

func (f *FakeKratos) Reset() {
	f.DisableSessionWasCalled = false
	f.LastDisabledSession = ""
}
