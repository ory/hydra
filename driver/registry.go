// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"

	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/x/contextx"

	"github.com/ory/hydra/v2/oauth2/trust"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/v2/persistence"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/dbal"
)

type registry interface {
	x.HTTPClientProvider

	contextx.Provider
	config.Provider
	persistence.Provider
	x.RegistryLogger
	x.RegistryWriter
	x.RegistryCookieStore
	client.Registry
	consent.Registry
	jwk.Registry
	trust.Registry
	oauth2.Registry
	x.TracingProvider
	x.NetworkProvider

	kratos.Provider
}

func newRegistryWithoutInit(c *config.DefaultProvider, l *logrusx.Logger) (*RegistrySQL, error) {
	r := &RegistrySQL{
		l:           l,
		conf:        c,
		initialPing: defaultInitialPing,
	}

	if !r.CanHandle(c.DSN()) {
		if dbal.IsSQLite(c.DSN()) {
			return nil, dbal.ErrSQLiteSupportMissing
		}

		return nil, dbal.ErrNoResponsibleDriverFound
	}

	return r, nil
}

func callRegistry(ctx context.Context, r *RegistrySQL) {
	r.ClientValidator()
	r.ClientManager()
	r.ClientHasher()
	r.ConsentManager()
	r.ConsentStrategy()
	r.SubjectIdentifierAlgorithm(ctx)
	r.KeyManager()
	r.KeyCipher()
	r.FlowCipher()
	r.OAuth2Storage()
	r.OAuth2Provider()
	r.AccessTokenJWTStrategy()
	r.OpenIDJWTStrategy()
	r.OpenIDConnectRequestValidator()
	r.Tracer(ctx)
}
