// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"io/fs"
	"net/http"

	enigma "github.com/ory/fosite/token/hmac"
	"github.com/ory/x/popx"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/x/contextx"
	"github.com/ory/x/httprouterx"

	"github.com/ory/hydra/v2/oauth2/trust"

	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"

	"github.com/ory/x/logrusx"

	"github.com/ory/hydra/v2/persistence"

	prometheus "github.com/ory/x/prometheusx"

	"github.com/ory/x/dbal"
	"github.com/ory/x/healthx"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
)

type Registry interface {
	dbal.Driver
	WritableRegistry

	Init(ctx context.Context, skipNetworkInit bool, migrate bool, ctxer contextx.Contextualizer, extraMigrations []fs.FS, goMigrations []popx.Migration) error

	x.HTTPClientProvider
	GetJWKSFetcherStrategy() fosite.JWKSFetcherStrategy

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
	PrometheusManager() *prometheus.MetricsManager
	x.TracingProvider
	FlowCipher() *aead.XChaCha20Poly1305

	kratos.Provider

	RegisterRoutes(ctx context.Context, admin *httprouterx.RouterAdmin, public *httprouterx.RouterPublic)
	ClientHandler() *client.Handler
	KeyHandler() *jwk.Handler
	ConsentHandler() *consent.Handler
	OAuth2Handler() *oauth2.Handler
	HealthHandler() *healthx.Handler
	OAuth2EnigmaStrategy() *enigma.HMACStrategy
	OAuth2AwareMiddleware() func(h http.Handler) http.Handler

	OAuth2HMACStrategy() foauth2.CoreStrategy
}

func NewRegistryFromDSN(ctx context.Context, c *config.DefaultProvider, l *logrusx.Logger, skipNetworkInit bool, migrate bool, ctxer contextx.Contextualizer) (Registry, error) {
	registry, err := NewRegistryWithoutInit(c, l)
	if err != nil {
		return nil, err
	}

	if err := registry.Init(ctx, skipNetworkInit, migrate, ctxer, nil, nil); err != nil {
		return nil, err
	}

	return registry, nil
}

func NewRegistryWithoutInit(c *config.DefaultProvider, l *logrusx.Logger) (Registry, error) {
	registry := NewRegistrySQL(
		c, l, config.Version, config.Commit, config.Date,
	)

	if !registry.CanHandle(c.DSN()) {
		if dbal.IsSQLite(c.DSN()) {
			return nil, dbal.ErrSQLiteSupportMissing
		}

		return nil, dbal.ErrNoResponsibleDriverFound
	}

	return registry, nil
}

func CallRegistry(ctx context.Context, r Registry) {
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
	r.AudienceStrategy()
	r.AccessTokenJWTStrategy()
	r.OpenIDJWTStrategy()
	r.OpenIDConnectRequestValidator()
	r.PrometheusManager()
	r.Tracer(ctx)
}
