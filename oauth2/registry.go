package oauth2

import (
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2/trust"
	"github.com/ory/hydra/x"
)

type InternalRegistry interface {
	client.Registry
	jwk.Registry
	trust.Registry
	x.RegistryWriter
	x.RegistryLogger
	consent.Registry
	Registry
}

type Registry interface {
	OAuth2Storage() x.FositeStorer
	OAuth2Provider() fosite.OAuth2Provider
	AudienceStrategy() fosite.AudienceMatchingStrategy
	ScopeStrategy() fosite.ScopeStrategy

	AccessTokenJWTStrategy() jwk.JWTStrategy
	OpenIDJWTStrategy() jwk.JWTStrategy

	OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator

	AccessRequestHooks() []AccessRequestHook
}
