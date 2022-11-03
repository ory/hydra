// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"github.com/ory/hydra/driver/config"

	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
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
	OpenIDJWTStrategy() jwk.JWTSigner
	OAuth2HMACStrategy() *foauth2.HMACSHAStrategy
	config.Provider
}
