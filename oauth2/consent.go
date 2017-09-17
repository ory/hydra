package oauth2

import (
	"github.com/gorilla/sessions"
	"github.com/ory/fosite"
)

type ConsentStrategy interface {
	ValidateConsentRequest(req fosite.AuthorizeRequester, session string, cookie *sessions.Session) (claims *Session, err error)
	CreateConsentRequest(req fosite.AuthorizeRequester, redirectURL string, cookie *sessions.Session) (token string, err error)
}
