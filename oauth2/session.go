package oauth2

import (
	csh "github.com/ory-am/fosite/handler/core/strategy"
	"github.com/ory-am/fosite/handler/oidc/strategy"
	"github.com/ory-am/fosite/token/jwt"
)

type Session struct {
	Subject                  string `json:"sub"`
	*strategy.DefaultSession `json:"idToken"`
}

func NewSession(subject string) *Session {
	return &Session{
		Subject: subject,
		DefaultSession: &strategy.DefaultSession{
			Claims:      new(jwt.IDTokenClaims),
			Headers:     new(jwt.Headers),
			HMACSession: new(csh.HMACSession),
		},
	}
}
