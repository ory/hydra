package oauth2

import (
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"net/url"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"time"
)

type Handler struct {
	OAuth2  fosite.OAuth2Provider
	Consent ConsentStrategy

	H herodot.Writer

	ForcedHTTP bool
	ConsentURL url.URL

	AccessTokenLifespan time.Duration
	CookieStore         sessions.Store

	L logrus.FieldLogger

	ScopeStrategy fosite.ScopeStrategy

	Issuer string
}
