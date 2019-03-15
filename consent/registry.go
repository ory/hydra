package consent

import (
	"github.com/gorilla/sessions"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/herodot"
	"github.com/ory/hydra/x"
	"time"
)

type registry interface {
	x.RegistryWriter

	ConsentManager() Manager
	CookieStore() sessions.Store
	Writer() herodot.Writer

	SubjectIdentifierAlgorithm() map[string]SubjectIdentifierAlgorithm
	ScopeStrategy() fosite.ScopeStrategy
	JWTStrategy() jwt.JWTStrategy
	OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator
}

type Registry {

}

type Configuration interface {
	LogoutRedirectURL() string
	RequestMaxAge() time.Duration

	AuthenticationURL() string
	ConsentURL() string
	IssuerURL() string
	OAuth2AuthURL() string
	RunsHTTPS() bool
}
