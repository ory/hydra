// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/negroni"

	"github.com/gorilla/sessions"
	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/ory/herodot"
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	foauth2 "github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/handler/rfc8628"
	"github.com/ory/hydra/v2/fosite/token/hmac"
	"github.com/ory/hydra/v2/fositex"
	"github.com/ory/hydra/v2/hsm"
	"github.com/ory/hydra/v2/internal/kratos"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/hydra/v2/x/oauth2cors"
	"github.com/ory/pop/v6"
	"github.com/ory/x/contextx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/healthx"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/otelx"
	"github.com/ory/x/popx"
	prometheus "github.com/ory/x/prometheusx"
	"github.com/ory/x/resilience"
	"github.com/ory/x/sqlcon"
)

type RegistrySQL struct {
	l, al           *logrusx.Logger
	conf            *config.DefaultProvider
	fh              fosite.Hasher
	cv              *client.Validator
	ctxer           contextx.Contextualizer
	hh              *healthx.Handler
	kc              *aead.AESGCM
	flowc           *aead.XChaCha20Poly1305
	cos             consent.Strategy
	writer          herodot.Writer
	hsm             hsm.Context
	forv            *openid.OpenIDConnectRequestValidator
	fop             fosite.OAuth2Provider
	sia             map[string]consent.SubjectIdentifierAlgorithm
	trc             *otelx.Tracer
	tracerWrapper   func(*otelx.Tracer) *otelx.Tracer
	arhs            []oauth2.AccessRequestHook
	basePersister   *sql.BasePersister
	oc              fosite.Configurator
	oidcs           jwk.JWTSigner
	ats             jwk.JWTSigner
	hmacs           foauth2.CoreStrategy
	enigmaHMAC      *hmac.HMACStrategy
	deviceHmac      rfc8628.RFC8628CodeStrategy
	fc              *fositex.Config
	publicCORS      *cors.Cors
	kratos          kratos.Client
	fositeFactories []fositex.Factory
	migrator        *sql.MigrationManager
	dbOptsModifier  []func(details *pop.ConnectionDetails)

	keyManager  jwk.Manager
	initialPing func(ctx context.Context, l *logrusx.Logger, p *sql.BasePersister) error
	middlewares []negroni.Handler
}

var (
	_ contextx.Provider = (*RegistrySQL)(nil)
	_ registry          = (*RegistrySQL)(nil)
)

// defaultInitialPing is the default function that will be called within RegistrySQL.Init to make sure
// the database is reachable. It can be injected for test purposes by changing the value
// of RegistrySQL.initialPing.
func defaultInitialPing(ctx context.Context, l *logrusx.Logger, p *sql.BasePersister) error {
	return errors.WithStack(resilience.Retry(l, 5*time.Second, 5*time.Minute, func() error {
		return p.Ping(ctx)
	}))
}

func (m *RegistrySQL) Init(
	ctx context.Context,
	skipNetworkInit bool,
	migrate bool,
	extraMigrations []fs.FS,
	goMigrations []popx.Migration,
) error {
	if m.basePersister == nil {
		if m.Config().CGroupsV1AutoMaxProcsEnabled() {
			_, err := maxprocs.Set(maxprocs.Logger(m.Logger().Infof))
			if err != nil {
				return fmt.Errorf("could not set GOMAXPROCS: %w", err)
			}
		}

		// new db connection
		pool, idlePool, connMaxLifetime, connMaxIdleTime, cleanedDSN := sqlcon.ParseConnectionOptions(
			m.l, m.Config().DSN(),
		)

		opts := &pop.ConnectionDetails{
			URL:             sqlcon.FinalizeDSN(m.l, cleanedDSN),
			IdlePool:        idlePool,
			ConnMaxLifetime: connMaxLifetime,
			ConnMaxIdleTime: connMaxIdleTime,
			Pool:            pool,
			TracerProvider:  m.Tracer(ctx).Provider(),
			Unsafe:          m.Config().DbIgnoreUnknownTableColumns(),
		}

		for _, f := range m.dbOptsModifier {
			f(opts)
		}

		c, err := pop.NewConnection(opts)
		if err != nil {
			return errors.WithStack(err)
		}
		if err := resilience.Retry(m.l, 5*time.Second, 5*time.Minute, c.Open); err != nil {
			return errors.WithStack(err)
		}

		m.basePersister = sql.NewBasePersister(c, m)
		if err := m.initialPing(ctx, m.Logger(), m.basePersister); err != nil {
			m.Logger().Print("Could not ping database: ", err)
			return err
		}

		m.migrator = sql.NewMigrationManager(c, m, extraMigrations, goMigrations)

		// if dsn is memory we have to run the migrations on every start
		// use case - such as
		// - just in memory
		// - shared connection
		// - shared but unique in the same process
		// see: https://sqlite.org/inmemorydb.html
		switch {
		case dbal.IsMemorySQLite(m.Config().DSN()):
			m.Logger().Println("Hydra is running migrations on every startup as DSN is memory.")
			m.Logger().Println("This means your data is lost when Hydra terminates.")
			fallthrough
		case migrate:
			if err := m.migrator.MigrateUp(ctx); err != nil {
				return err
			}
		}

		if !skipNetworkInit {
			net, err := m.basePersister.DetermineNetwork(ctx)
			if err != nil {
				m.Logger().WithError(err).Warnf("Unable to determine network, retrying.")
				return err
			}

			m.basePersister = m.basePersister.WithFallbackNetworkID(net.ID)
		}
	}

	return nil
}

func (m *RegistrySQL) alwaysCanHandle(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	s := dbal.Canonicalize(scheme)
	return s == dbal.DriverMySQL || s == dbal.DriverPostgreSQL || s == dbal.DriverCockroachDB
}

func (m *RegistrySQL) PingContext(ctx context.Context) error {
	return m.basePersister.Ping(ctx)
}

func (m *RegistrySQL) ClientManager() client.Manager {
	return m.Persister()
}

func (m *RegistrySQL) ConsentManager() consent.Manager {
	return m.Persister()
}

func (m *RegistrySQL) ObfuscatedSubjectManager() consent.ObfuscatedSubjectManager {
	return m.Persister()
}
func (m *RegistrySQL) LoginManager() consent.LoginManager   { return m.Persister() }
func (m *RegistrySQL) LogoutManager() consent.LogoutManager { return m.Persister() }

func (m *RegistrySQL) OAuth2Storage() x.FositeStorer {
	return m.Persister()
}

func (m *RegistrySQL) KeyManager() jwk.Manager {
	if m.keyManager == nil {
		softwareKeyManager := &sql.JWKPersister{BasePersister: m.basePersister}
		if m.Config().HSMEnabled() {
			hardwareKeyManager := hsm.NewKeyManager(m.HSMContext(), m.Config())
			m.keyManager = jwk.NewManagerStrategy(hardwareKeyManager, softwareKeyManager)
		} else {
			m.keyManager = softwareKeyManager
		}
	}
	return m.keyManager
}

func (m *RegistrySQL) GrantManager() trust.GrantManager {
	return m.Persister()
}

func (m *RegistrySQL) Contextualizer() contextx.Contextualizer {
	if m.ctxer == nil {
		panic("registry Contextualizer not set")
	}
	return m.ctxer
}

func (m *RegistrySQL) addPublicCORSOnHandler(ctx context.Context) func(http.Handler) http.Handler {
	corsConfig, corsEnabled := m.Config().CORSPublic(ctx)
	if !corsEnabled {
		return func(h http.Handler) http.Handler {
			return h
		}
	}
	if m.publicCORS == nil {
		m.publicCORS = cors.New(corsConfig)
	}
	return func(h http.Handler) http.Handler {
		return m.publicCORS.Handler(h)
	}
}

func (m *RegistrySQL) RegisterPublicRoutes(ctx context.Context, public *httprouterx.RouterPublic) {
	m.HealthHandler().SetHealthRoutes(public, false, healthx.WithMiddleware(m.addPublicCORSOnHandler(ctx)))

	corsMW := oauth2cors.Middleware(m)
	jwk.NewHandler(m).SetPublicRoutes(public, corsMW)
	client.NewHandler(m).SetPublicRoutes(public)
	oauth2.NewHandler(m).SetPublicRoutes(public, corsMW)
}

func (m *RegistrySQL) RegisterAdminRoutes(admin *httprouterx.RouterAdmin) {
	m.HealthHandler().SetHealthRoutes(admin, true)
	m.HealthHandler().SetVersionRoutes(admin)
	admin.Handler("GET", prometheus.MetricsPrometheusPath, promhttp.Handler())

	consent.NewHandler(m).SetRoutes(admin)
	jwk.NewHandler(m).SetAdminRoutes(admin)
	client.NewHandler(m).SetAdminRoutes(admin)
	oauth2.NewHandler(m).SetAdminRoutes(admin)
	trust.NewHandler(m).SetRoutes(admin)
}

func (m *RegistrySQL) Writer() herodot.Writer {
	if m.writer == nil {
		h := herodot.NewJSONWriter(m.Logger())
		h.ErrorEnhancer = x.ErrorEnhancer
		m.writer = h
	}
	return m.writer
}

func (m *RegistrySQL) Logger() *logrusx.Logger {
	if m.l == nil {
		m.l = logrusx.New("Ory Hydra", config.Version)
	}
	return m.l
}

func (m *RegistrySQL) AuditLogger() *logrusx.Logger {
	if m.al == nil {
		m.al = logrusx.NewAudit("Ory Hydra", config.Version)
		m.al.UseConfig(m.Config().Source(contextx.RootContext))
	}
	return m.al
}

func (m *RegistrySQL) ClientHasher() fosite.Hasher {
	if m.fh == nil {
		m.fh = x.NewHasher(m, m.Config())
	}
	return m.fh
}

func (m *RegistrySQL) ClientValidator() *client.Validator {
	if m.cv == nil {
		m.cv = client.NewValidator(m)
	}
	return m.cv
}

func (m *RegistrySQL) HealthHandler() *healthx.Handler {
	if m.hh == nil {
		m.hh = healthx.NewHandler(m.Writer(), config.Version, healthx.ReadyCheckers{
			"database": func(r *http.Request) error {
				return m.PingContext(r.Context())
			},
			"migrations": func(r *http.Request) error {
				status, err := m.migrator.MigrationStatus(r.Context())
				if err != nil {
					return err
				}

				if status.HasPending() {
					err := errors.Errorf("migrations have not yet been fully applied: %+v", status)
					m.Logger().WithField("status", fmt.Sprintf("%+v", status)).WithError(err).Warn("Instance is not yet ready because migrations have not yet been fully applied.")
					return err
				}
				return nil
			},
		})
	}

	return m.hh
}

func (m *RegistrySQL) ConsentStrategy() consent.Strategy {
	if m.cos == nil {
		m.cos = consent.NewStrategy(m)
	}
	return m.cos
}

func (m *RegistrySQL) KeyCipher() *aead.AESGCM {
	if m.kc == nil {
		m.kc = aead.NewAESGCM(m.Config())
	}
	return m.kc
}

func (m *RegistrySQL) FlowCipher() *aead.XChaCha20Poly1305 {
	if m.flowc == nil {
		m.flowc = aead.NewXChaCha20Poly1305(m.Config())
	}
	return m.flowc
}

func (m *RegistrySQL) CookieStore(ctx context.Context) (sessions.Store, error) {
	var keys [][]byte
	secrets, err := m.conf.GetCookieSecrets(ctx)
	if err != nil {
		return nil, err
	}

	for _, k := range secrets {
		encrypt := sha256.Sum256(k)
		keys = append(keys, k, encrypt[:])
	}

	cs := sessions.NewCookieStore(keys...)
	cs.Options.Secure = m.Config().CookieSecure(ctx)
	cs.Options.HttpOnly = true

	// CookieStore MaxAge is set to 86400 * 30 by default. This prevents secure cookies retrieval with expiration > 30 days.
	// MaxAge(0) disables internal MaxAge check by SecureCookie, see:
	//
	// https://github.com/ory/hydra/pull/2488#discussion_r618992698
	cs.MaxAge(0)

	if domain := m.Config().CookieDomain(ctx); domain != "" {
		cs.Options.Domain = domain
	}

	cs.Options.Path = "/"
	if sameSite := m.Config().CookieSameSiteMode(ctx); sameSite != 0 {
		cs.Options.SameSite = sameSite
	}

	return cs, nil
}

func (m *RegistrySQL) HTTPClient(ctx context.Context, opts ...httpx.ResilientOptions) *retryablehttp.Client {
	opts = append(opts,
		httpx.ResilientClientWithLogger(m.Logger()),
		httpx.ResilientClientWithMaxRetry(2),
		httpx.ResilientClientWithConnectionTimeout(30*time.Second))

	tracer := m.Tracer(ctx)
	if tracer.IsLoaded() {
		opts = append(opts, httpx.ResilientClientWithTracer(tracer.Tracer()))
	}

	if m.Config().ClientHTTPNoPrivateIPRanges() {
		opts = append(
			opts,
			httpx.ResilientClientDisallowInternalIPs(),
			httpx.ResilientClientAllowInternalIPRequestsTo(m.Config().ClientHTTPPrivateIPExceptionURLs()...),
		)
	}
	return httpx.NewResilientClient(opts...)
}

func (m *RegistrySQL) OAuth2Provider() fosite.OAuth2Provider {
	if m.fop == nil {
		m.fop = fosite.NewOAuth2Provider(m.OAuth2Storage(), m.OAuth2ProviderConfig())
	}
	return m.fop
}

func (m *RegistrySQL) OpenIDJWTStrategy() jwk.JWTSigner {
	if m.oidcs == nil {
		m.oidcs = jwk.NewDefaultJWTSigner(m, x.OpenIDConnectKeyName)
	}
	return m.oidcs
}

func (m *RegistrySQL) AccessTokenJWTStrategy() jwk.JWTSigner {
	if m.ats == nil {
		m.ats = jwk.NewDefaultJWTSigner(m, x.OAuth2JWTKeyName)
	}
	return m.ats
}

func (m *RegistrySQL) OAuth2EnigmaStrategy() *hmac.HMACStrategy {
	if m.enigmaHMAC == nil {
		m.enigmaHMAC = &hmac.HMACStrategy{Config: m.OAuth2Config()}
	}
	return m.enigmaHMAC
}

func (m *RegistrySQL) OAuth2HMACStrategy() foauth2.CoreStrategy {
	if m.hmacs == nil {
		m.hmacs = foauth2.NewHMACSHAStrategy(m.OAuth2EnigmaStrategy(), m.OAuth2Config())
	}
	return m.hmacs
}

// RFC8628HMACStrategy returns the rfc8628 strategy
func (m *RegistrySQL) RFC8628HMACStrategy() rfc8628.RFC8628CodeStrategy {
	if m.deviceHmac == nil {
		m.deviceHmac = compose.NewDeviceStrategy(m.OAuth2Config())
	}
	return m.deviceHmac
}

func (m *RegistrySQL) OAuth2Config() *fositex.Config {
	if m.fc == nil {
		m.fc = fositex.NewConfig(m)
	}
	return m.fc
}

func (m *RegistrySQL) ExtraFositeFactories() []fositex.Factory {
	return m.fositeFactories
}

func (m *RegistrySQL) OAuth2ProviderConfig() fosite.Configurator {
	if m.oc != nil {
		return m.oc
	}

	conf := m.OAuth2Config()
	hmacAtStrategy := m.OAuth2HMACStrategy()
	deviceHmacAtStrategy := m.RFC8628HMACStrategy()
	oidcSigner := m.OpenIDJWTStrategy()
	atSigner := m.AccessTokenJWTStrategy()
	jwtAtStrategy := &foauth2.DefaultJWTStrategy{
		Signer:   atSigner,
		Strategy: hmacAtStrategy,
		Config:   conf,
	}

	conf.LoadDefaultHandlers(&compose.CommonStrategy{
		CoreStrategy: fositex.NewTokenStrategy(m.Config(), hmacAtStrategy, &foauth2.DefaultJWTStrategy{
			Signer:   jwtAtStrategy,
			Strategy: hmacAtStrategy,
			Config:   conf,
		}),
		RFC8628CodeStrategy: deviceHmacAtStrategy,
		OIDCTokenStrategy: &openid.DefaultStrategy{
			Config: conf,
			Signer: oidcSigner,
		},
		Signer: oidcSigner,
	})

	m.oc = conf
	return m.oc
}

func (m *RegistrySQL) OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator {
	if m.forv == nil {
		m.forv = openid.NewOpenIDConnectRequestValidator(&openid.DefaultStrategy{
			Config: m.OAuth2ProviderConfig(),
			Signer: m.OpenIDJWTStrategy(),
		}, m.OAuth2ProviderConfig())
	}
	return m.forv
}

func (m *RegistrySQL) Networker() x.Networker { return m.basePersister }

func (m *RegistrySQL) SubjectIdentifierAlgorithm(ctx context.Context) map[string]consent.SubjectIdentifierAlgorithm {
	if m.sia == nil {
		m.sia = map[string]consent.SubjectIdentifierAlgorithm{}
		for _, t := range m.Config().SubjectTypesSupported(ctx) {
			switch t {
			case "public":
				m.sia["public"] = consent.NewSubjectIdentifierAlgorithmPublic()
			case "pairwise":
				m.sia["pairwise"] = consent.NewSubjectIdentifierAlgorithmPairwise([]byte(m.Config().SubjectIdentifierAlgorithmSalt(ctx)))
			}
		}
	}
	return m.sia
}

func (m *RegistrySQL) Tracer(_ context.Context) *otelx.Tracer {
	if m.trc == nil {
		t, err := otelx.New("Ory Hydra", m.l, m.conf.Tracing())
		if err != nil {
			m.Logger().WithError(err).Error("Unable to initialize Tracer.")
		} else {
			// Wrap the tracer if required
			if m.tracerWrapper != nil {
				t = m.tracerWrapper(t)
			}

			m.trc = t
		}
	}
	if m.trc == nil || m.trc.Tracer() == nil {
		m.trc = otelx.NewNoop(m.l, m.Config().Tracing())
	}

	return m.trc
}

func (m *RegistrySQL) Persister() persistence.Persister {
	return sql.NewPersister(m.basePersister, m)
}

// Config returns the configuration for the given context. It may or may not be the same as the global configuration.
func (m *RegistrySQL) Config() *config.DefaultProvider {
	return m.conf
}

// WithConsentStrategy forces a consent strategy which is only used for testing.
func (m *RegistrySQL) WithConsentStrategy(c consent.Strategy) {
	m.cos = c
}

func (m *RegistrySQL) AccessRequestHooks() []oauth2.AccessRequestHook {
	if m.arhs == nil {
		m.arhs = []oauth2.AccessRequestHook{
			oauth2.RefreshTokenHook(m),
			oauth2.TokenHook(m),
		}
	}
	return m.arhs
}

func (m *RegistrySQL) HSMContext() hsm.Context {
	if m.hsm == nil {
		m.hsm = hsm.NewContext(m.Config(), m.l)
	}
	return m.hsm
}

func (m *RegistrySQL) Kratos() kratos.Client {
	if m.kratos == nil {
		m.kratos = kratos.New(m)
	}
	return m.kratos
}

func (m *RegistrySQL) HTTPMiddlewares() []negroni.Handler {
	return m.middlewares
}

func (m *RegistrySQL) Migrator() *sql.MigrationManager {
	return m.migrator
}

func (m *RegistrySQL) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.basePersister.Transaction(ctx, func(ctx context.Context, _ *pop.Connection) error { return fn(ctx) })
}
