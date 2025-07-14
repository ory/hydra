// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package contextx

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ory/x/configx"
)

type Default struct{}

var _ Contextualizer = (*Default)(nil)

func (d *Default) Network(_ context.Context, network uuid.UUID) uuid.UUID {
	if network == uuid.Nil {
		panic("nid must be not nil")
	}
	return network
}

func (d *Default) Config(_ context.Context, config *configx.Provider) *configx.Provider {
	return config
}
