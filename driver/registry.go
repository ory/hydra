package driver

import (
	"fmt"

	"github.com/pkg/errors"
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
	PrometheusManager() *prometheus.MetricsManager
	Tracer() *tracing.Tracer

	RegisterRoutes(admin *x.RouterAdmin, public *x.RouterPublic)
	ClientHandler() *client.Handler
	KeyHandler() *jwk.Handler
	ConsentHandler() *consent.Handler
	OAuth2Handler() *oauth2.Handler
	HealthHandler() *healthx.Handler
}

func MustNewRegistry(c configuration.Provider) Registry {
	r, err := NewRegistry(c)
	cmdx.Must(err, "unable to initialize services: %s", err)
	return r
}

func NewRegistry(c configuration.Provider) (Registry, error) {
	fmt.Print("Got dsn: ", c.DSN())
	driver, err := dbal.GetDriverFor(c.DSN())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	registry, ok := driver.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}

	registry = registry.WithConfig(c)

	if err := registry.Init(); err != nil {
		return nil, err
	}

	return registry, nil
}

func CallRegistry(r Registry) {
	r.ClientValidator()
	r.ClientManager()
	r.ClientHasher()
	r.ConsentManager()
	r.ConsentStrategy()
	r.SubjectIdentifierAlgorithm()
	r.KeyManager()
	r.KeyGenerators()
	r.KeyCipher()
	r.OAuth2Storage()
	r.OAuth2Provider()
	r.AudienceStrategy()
	r.ScopeStrategy()
	r.AccessTokenJWTStrategy()
	r.OpenIDJWTStrategy()
	r.OpenIDConnectRequestValidator()
	r.PrometheusManager()
	r.Tracer()
}
