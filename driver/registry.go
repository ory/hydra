// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/ory/x/contextx"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/pop/v6"
	"github.com/ory/x/contextx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/logrusx"
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
	scheme, _, _ := strings.Cut(c.DSN(), "://")
	if !pop.DialectSupported(pop.CanonicalDialect(scheme)) {
		if dbal.IsSQLite(c.DSN()) {
			return nil, errors.New("The DSN connection string looks like a SQLite connection, but SQLite support was not built into the binary. Please check if you have downloaded the correct binary or are using the correct Docker Image. Binary archives and Docker Images indicate SQLite support by appending the -sqlite suffix.")
		}
		return nil, errors.New("unsupported DSN type")
	}

	return &RegistrySQL{
		l:           l,
		conf:        c,
		initialPing: defaultInitialPing,
	}, nil
}

func callRegistry(ctx context.Context, r *RegistrySQL) {
	r.ClientValidator()
	r.ClientManager()
	r.ClientHasher()
	r.ConsentManager()
	r.ConsentStrategy()
	r.KeyManager()
	r.KeyCipher()
	r.FlowCipher()
	r.OAuth2Storage()
	r.OAuth2Provider()
	r.AccessTokenJWTSigner()
	r.OpenIDJWTSigner()
	r.OpenIDConnectRequestValidator()
	r.Tracer(ctx)
}
