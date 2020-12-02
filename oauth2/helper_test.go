package oauth2_test

import (
	"time"

	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/hmac"
	"github.com/ory/hydra/driver/config"
)

func Tokens(c *config.Provider, length int) (res [][]string) {
	s := &oauth2.HMACSHAStrategy{
		Enigma: &hmac.HMACStrategy{
			GlobalSecret: c.GetSystemSecret(),
		},
		AccessTokenLifespan:   time.Hour,
		AuthorizeCodeLifespan: time.Hour,
	}

	for i := 0; i < length; i++ {
		tok, sig, _ := s.Enigma.Generate()
		res = append(res, []string{sig, tok})
	}
	return res
}
