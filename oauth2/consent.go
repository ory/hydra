package oauth2

import (
	"github.com/gorilla/sessions"
	"github.com/ory-am/fosite"
)

type ConsentStrategy interface {
	ValidateResponse(authorizeRequest fosite.AuthorizeRequester, token string, session *sessions.Session) (claims *Session, err error)
	IssueChallenge(authorizeRequest fosite.AuthorizeRequester, redirectURL string, session *sessions.Session) (token string, err error)
}
