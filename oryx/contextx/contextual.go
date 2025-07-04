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
	NoOp struct{}
)

func (d *Static) Network(ctx context.Context, network uuid.UUID) uuid.UUID {
	return d.NID
}

func (d *Static) Config(ctx context.Context, config *configx.Provider) *configx.Provider {
	return d.C
}

func (d *NoOp) Network(ctx context.Context, network uuid.UUID) uuid.UUID {
	return network
}

func (d *NoOp) Config(ctx context.Context, config *configx.Provider) *configx.Provider {
	return config
}
