package oauth2

import (
	"github.com/ory-am/fosite"
	"github.com/gorilla/sessions"
)

type ConsentStrategy interface {
	ValidateResponse(authorizeRequest fosite.AuthorizeRequester, token string, session *sessions.Session) (claims *Session, err error)
	IssueChallenge(authorizeRequest fosite.AuthorizeRequester, redirectURL string, session *sessions.Session) (token string, err error)
}
