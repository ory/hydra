// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package persistence

import (
	"context"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/pop/v6"
	"github.com/ory/x/networkx"
)

type (
	Persister interface {
		consent.Manager
		consent.ObfuscatedSubjectManager
		consent.LoginManager
		consent.LogoutManager
		client.Manager
		x.FositeStorer
		trust.GrantManager

		Connection(context.Context) *pop.Connection
		Transaction(context.Context, func(ctx context.Context, c *pop.Connection) error) error
		Ping(context.Context) error
		DetermineNetwork(ctx context.Context) (*networkx.Network, error)
		x.Networker
	}
	Provider interface {
		Persister() Persister
	}
)
