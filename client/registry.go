package client

import (
	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/jwk"
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
	OpenIDJWTStrategy() jwk.JWTStrategy
	OAuth2HMACStrategy() *foauth2.HMACSHAStrategy
	Config() *config.Provider
}
