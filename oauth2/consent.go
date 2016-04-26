package oauth2

import "github.com/ory-am/fosite"

type Session struct {
	Subject string `json:"sub"`
}

type ConsentValidator interface {
	ValidateConsentToken(authorizeRequest fosite.AuthorizeRequester, token string) (claims *Session, err error)
}
