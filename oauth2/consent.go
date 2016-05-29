package oauth2

import "github.com/ory-am/fosite"

type ConsentStrategy interface {
	ValidateResponse(authorizeRequest fosite.AuthorizeRequester, token string) (claims *Session, err error)
	IssueChallenge(authorizeRequest fosite.AuthorizeRequester, redirectURL string) (token string, err error)
}
