package consent

import (
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/x"
)

type registry interface {
	x.RegistryWriter
	x.RegistryCookieStore
	oauth2.Registry
	Registry
}

type Registry interface {
	ConsentManager() Manager
	ConsentStrategy() Strategy
	SubjectIdentifierAlgorithm() map[string]SubjectIdentifierAlgorithm
}

type Configuration interface {
	configuration.Provider
}
