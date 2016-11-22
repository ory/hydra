package oauth2

import (
	"bytes"
	"encoding/gob"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/openid"
	"github.com/ory-am/fosite/token/jwt"
)

type Session struct {
	*openid.DefaultSession `json:"idToken"`
	Extra                  map[string]interface{} `json:"extra"`
}

func NewSession(subject string) *Session {
	return &Session{
		DefaultSession: &openid.DefaultSession{
			Claims:  new(jwt.IDTokenClaims),
			Headers: new(jwt.Headers),
			Subject: subject,
		},
	}
}

func (s *Session) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	var clone Session
	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)
	_ = enc.Encode(s)
	_ = dec.Decode(&clone)
	return &clone
}
