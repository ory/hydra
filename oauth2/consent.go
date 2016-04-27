package oauth2

import "github.com/ory-am/fosite"

type Session struct {
	Subject string `json:"sub"`
}

type ConsentStrategy interface {
	ValidateResponseToken(authorizeRequest fosite.AuthorizeRequester, token string) (claims *Session, err error)
	IssueRequestToken(authorizeRequest fosite.AuthorizeRequester) (token string, err error)
}
