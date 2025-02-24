// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"github.com/ory/hydra/v2/driver/config"

	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/rfc8628"
	enigma "github.com/ory/fosite/token/hmac"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
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
	OAuth2HMACStrategy() foauth2.CoreStrategy
	OAuth2EnigmaStrategy() *enigma.HMACStrategy
	RFC8628HMACStrategy() rfc8628.RFC8628CodeStrategy
	config.Provider
}
