package oauth2

import (
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/oidc"
	"time"
	"strings"
	"github.com/ory-am/fosite/enigma/jwt"
)

type Session struct {
	Issuer string `json:"iss"`
	Subject string `json:"sub"`
	Audience string `json:"aud"`
	ExpireAt time.Time `json:"exp"`
	NotBefore time.Time `json:"nbf"`
	IssuedAt time.Time `json:"iat"`
	JTI string `json:"jti"`
}

func (s *Session) String() string {
	return strings.Join([]string{
		s.Issuer,
		s.Audience,
		s.ExpireAt,
		s.NotBefore,
		s.IssuedAt,
		s.JTI,
		// Subject is not included because it will change between request and response
	}, ";")
}

type ConsentStrategy interface {
	ValidateResponseToken(authorizeRequest fosite.AuthorizeRequester, token string) (claims *Session, err error)
	IssueRequestToken(authorizeRequest fosite.AuthorizeRequester) (token string, err error)
}

type DefaultConsentStrategy struct {
	JWT jwt.Enigma
}

func (s *DefaultConsentStrategy) ValidateResponseToken(authorizeRequest fosite.AuthorizeRequester, token string) (claims *Session, err error) {

}

func (s *DefaultConsentStrategy) IssueRequestToken(authorizeRequest fosite.AuthorizeRequester) (token string, err error){

}
