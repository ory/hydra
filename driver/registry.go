package driver

import (
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"

	"github.com/ory/hydra/metrics/prometheus"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/tracing"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
	"github.com/ory/x/dbal"
	"github.com/ory/x/healthx"
)

type Registry interface {
	dbal.Driver

	Init() error

	WithConfig(c configuration.Provider) Registry
	WithLogger(l logrus.FieldLogger) Registry

	WithBuildInfo(version, hash, date string) Registry
	BuildVersion() string
	BuildDate() string
	BuildHash() string

	x.RegistryLogger
	x.RegistryWriter
	x.RegistryCookieStore
	client.Registry
	consent.Registry
	jwk.Registry
	oauth2.Registry

	ClientHandler() *client.Handler
	KeyHandler() *jwk.Handler
	ConsentHandler() *consent.Handler
	OAuth2Handler() *oauth2.Handler
	HealthHandler() *healthx.Handler

	RegisterRoutes(admin *x.RouterAdmin, public *x.RouterPublic)

	PrometheusManager() *prometheus.MetricsManager

	Tracer() *tracing.Tracer
}

func MustNewRegistry(c configuration.Provider) Registry {
	r, err := NewRegistry(c)
	cmdx.Must(err, "unable to initialize services: %s", err)
	return r
}

func NewRegistry(c configuration.Provider) (Registry, error) {
	driver, err := dbal.GetDriverFor(c.DSN())
	if err != nil {
		return nil, err
	}

	registry, ok := driver.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}

	registry = registry.WithConfig(c)

	if err := registry.Init(); err != nil {
		return nil, err
	}

	_ = registry.ClientHandler()
	_ = registry.KeyHandler()
	_ = registry.ConsentHandler()
	_ = registry.OAuth2Handler()
	_ = registry.HealthHandler()
	_ = registry.PrometheusManager()
	_ = registry.Tracer()
	_ = registry.OAuth2Storage()
	_ = registry.OAuth2Provider()
	_ = registry.AudienceStrategy()
	_ = registry.ScopeStrategy()
	_ = registry.AccessTokenJWTStrategy()
	_ = registry.OpenIDJWTStrategy()
	_ = registry.OpenIDConnectRequestValidator()
	_ = registry.KeyManager()
	_ = registry.KeyGenerators()
	_ = registry.KeyCipher()
	_ = registry.ConsentManager()
	_ = registry.ConsentStrategy()
	_ = registry.SubjectIdentifierAlgorithm()
	_ = registry.ClientValidator()
	_ = registry.ClientManager()
	_ = registry.ClientHasher()
	_ = registry.CookieStore()

	return registry, nil
}
