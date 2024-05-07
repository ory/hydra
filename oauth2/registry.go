// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/hydra/v2/x"
)

type InternalRegistry interface {
	client.Registry
	jwk.Registry
	trust.Registry
	x.RegistryWriter
	x.RegistryLogger
	consent.Registry
	persistence.Provider
	Registry
	FlowCipher() *aead.XChaCha20Poly1305
}

type Registry interface {
	OAuth2Storage() x.FositeStorer
	OAuth2Provider() fosite.OAuth2Provider
	AudienceStrategy() fosite.AudienceMatchingStrategy
	AccessTokenJWTStrategy() jwk.JWTSigner
	OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator
	AccessRequestHooks() []AccessRequestHook
	OAuth2ProviderConfig() fosite.Configurator
}
