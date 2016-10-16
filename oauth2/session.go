package oauth2

import (
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/fosite/handler/openid"
	"github.com/ory-am/fosite/token/jwt"
)

type Session struct {
	Subject string `json:"sub"`
	*openid.DefaultSession `json:"idToken"`
	*oauth2.HMACSession    `json:"session"`
	Extra   map[string]interface{} `json:"extra"`
}

func NewSession(subject string) *Session {
	return &Session{
		Subject: subject,
		DefaultSession: &openid.DefaultSession{
			Claims:  new(jwt.IDTokenClaims),
			Headers: new(jwt.Headers),
		},
		HMACSession: new(oauth2.HMACSession),
	}
}
