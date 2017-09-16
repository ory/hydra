package oauth2

import (
	"net/url"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/sirupsen/logrus"
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
