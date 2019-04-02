package client

import (
	"github.com/ory/fosite"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/x"
)

type InternalRegistry interface {
	x.RegistryWriter
	Registry
}

type Registry interface {
	ClientValidator() *Validator
	ClientManager() Manager
	ClientHasher() fosite.Hasher
}

type Configuration interface {
	configuration.Provider
}
