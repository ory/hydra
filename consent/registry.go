package consent

import (
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/metrics"
	"github.com/ory/hydra/x"
)

type InternalRegistry interface {
	x.RegistryWriter
	x.RegistryCookieStore
	Registry
	client.Registry

	OAuth2Storage() x.FositeStorer
	OpenIDJWTStrategy() jwk.JWTStrategy
	PrometheusManager() *metrics.Prometheus
	OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator
	ScopeStrategy() fosite.ScopeStrategy
}

type Registry interface {
	ConsentManager() Manager
	ConsentStrategy() Strategy
	SubjectIdentifierAlgorithm() map[string]SubjectIdentifierAlgorithm
}

type Configuration interface {
	configuration.Provider
}
