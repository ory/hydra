// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package contextx

import (
	"context"

	"github.com/ory/x/configx"

	"github.com/gofrs/uuid"
)

type (
	Contextualizer interface {
		// Network returns the network id for the given context.
		Network(ctx context.Context, network uuid.UUID) uuid.UUID

		// Config returns the config for the given context.
		Config(ctx context.Context, config *configx.Provider) *configx.Provider
	}
	Provider interface {
		Contextualizer() Contextualizer
	}
	Static struct {
		NID uuid.UUID
		C   *configx.Provider
	}
)

func (d *Static) Network(_ context.Context, nid uuid.UUID) uuid.UUID {
	if d.NID == uuid.Nil {
		return nid
	}
	return d.NID
}

func (d *Static) Config(_ context.Context, c *configx.Provider) *configx.Provider {
	if d.C == nil {
		return c
	}
	return d.C
}
