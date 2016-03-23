package consent

import "github.com/ory-am/fosite"

type ConsentClaims struct {
	Subject string `json:"sub"`
}

type Validator interface {
	ValidateConsentToken(authorizeRequest fosite.AuthorizeRequester, token string) (claims *ConsentClaims, err error)
}