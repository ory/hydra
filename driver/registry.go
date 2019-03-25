package driver

import (
	"github.com/go-errors/errors"

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

	WithConfig(c configuration.Provider) Registry
	WithBuildVersion(bv string) Registry

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

	return registry, nil
}
