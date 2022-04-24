package jwk

import (
	"context"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/x"
)

type InternalRegistry interface {
	x.RegistryWriter
	x.RegistryLogger
	Registry
}

type Registry interface {
	Config(ctx context.Context) *config.Provider
	KeyManager() Manager
	SoftwareKeyManager() Manager
	KeyGenerators() map[string]KeyGenerator
	KeyCipher() *AEAD
}
