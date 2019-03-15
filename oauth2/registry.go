package oauth2

import (
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

type Registry interface {
	ClientManager() client.Manager
	Writer() herodot.Writer
	Logger() logrus.FieldLogger

	OAuth2Provider() fosite.OAuth2Provider

	ConsentStrategy() consent.Strategy
	OAuth2Storage() pkg.FositeStorer

	OpenIDJWTStrategy() jwk.JWTStrategy
	AccessTokenJWTStrategy() jwk.JWTStrategy
	ScopeStrategy() fosite.ScopeStrategy
	AudienceStrategy() fosite.AudienceMatchingStrategy
}

type Configuration interface {
	HashSignature() bool
	IsUsingJWTAsAccessTokens() bool
	ForcedHTTP() bool
	ErrorURL() *url.URL
	AccessTokenLifespan() time.Duration
	AccessTokenStrategy() string
	IssuerURL() string
	ClientRegistrationURL() string
	ClaimsSupported() string
	ScopesSupported() string
	SubjectTypes() []string
	UserinfoEndpoint() string
	ShareOAuth2Debug() bool
}
