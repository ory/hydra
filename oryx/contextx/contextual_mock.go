// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package contextx

import (
	"context"

	"github.com/ory/x/configx"

	"github.com/gofrs/uuid"
)

// TestContextualizer is a mock implementation of the Contextualizer interface.
type TestContextualizer struct{}

type contextKeyFake int

// fakeNIDContext is a test key for NID.
const fakeNIDContext contextKeyFake = 1

// SetNIDContext sets the nid for the given context.
func SetNIDContext(ctx context.Context, nid uuid.UUID) context.Context {
	return context.WithValue(ctx, fakeNIDContext, nid) //nolint:staticcheck
}

// Network returns the network id for the given context.
func (d *TestContextualizer) Network(ctx context.Context, network uuid.UUID) uuid.UUID {
	nid, ok := ctx.Value(fakeNIDContext).(uuid.UUID)
	if !ok {
		return network
	}
	return nid
}

// Config returns the config for the given context.
func (d *TestContextualizer) Config(ctx context.Context, config *configx.Provider) *configx.Provider {
	return config
}
