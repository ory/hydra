package jwk

import (
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/x"
)

type InternalRegistry interface {
	x.RegistryWriter
	x.RegistryLogger
	Registry
}

type Registry interface {
	KeyManager() Manager
	KeyGenerators() map[string]KeyGenerator
	KeyCipher() *AEAD
}

type Configuration interface {
	configuration.Provider
}
