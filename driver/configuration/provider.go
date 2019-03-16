package configuration

import (
	"net/url"
	"time"
)

type Provider interface {
	MustValidate()

	DSN() string
	AdminListenOn() string
	PublicListenOn() string

	DefaultClientScope() []string
	SupportedSubjectTypes() []string


	ConsentRequestMaxAge() time.Duration

	LogoutRedirectURL() string
	LoginURL() string
	ConsentURL() string
	ErrorURL() *url.URL

	IssuerURL() string
	OAuth2AuthURL() string

	ServesHTTPS() bool


	HashSignature() bool
	IsUsingJWTAsAccessTokens() bool

	AccessTokenLifespan() time.Duration
	AccessTokenStrategy() string
	ClientRegistrationURL() string
	ClaimsSupported() string
	ScopesSupported() string
	SubjectTypes() []string
	UserinfoEndpoint() string
	ShareOAuth2Debug() bool
}

func MustValidate(p Provider) {
}
