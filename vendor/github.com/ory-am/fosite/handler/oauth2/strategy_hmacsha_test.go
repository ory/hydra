package oauth2

import (
	"strings"
	"testing"
	"time"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/token/hmac"
	"github.com/stretchr/testify/assert"
)

var s = HMACSHAStrategy{
	Enigma: &hmac.HMACStrategy{GlobalSecret: []byte("foobarfoobarfoobarfoobar")},
}

var hmacExpiredCase = fosite.Request{
	Client: &fosite.DefaultClient{
		Secret: []byte("foobarfoobarfoobarfoobar"),
	},
	Session: &fosite.DefaultSession{
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.AccessToken:   time.Now().Add(-time.Hour),
			fosite.AuthorizeCode: time.Now().Add(-time.Hour),
		},
	},
}

var hmacValidCase = fosite.Request{
	Client: &fosite.DefaultClient{
		Secret: []byte("foobarfoobarfoobarfoobar"),
	},
	Session: &fosite.DefaultSession{
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.AccessToken:   time.Now().Add(time.Hour),
			fosite.AuthorizeCode: time.Now().Add(time.Hour),
		},
	},
}

func TestHMACAccessToken(t *testing.T) {
	for _, c := range []struct {
		r    fosite.Request
		pass bool
	}{
		{
			r:    hmacValidCase,
			pass: true,
		},
		{
			r:    hmacExpiredCase,
			pass: false,
		},
	} {
		token, signature, err := s.GenerateAccessToken(nil, &c.r)
		assert.Nil(t, err, "%s", err)
		assert.Equal(t, strings.Split(token, ".")[1], signature)

		err = s.ValidateAccessToken(nil, &c.r, token)
		if c.pass {
			assert.Nil(t, err, "%s", err)
			validate := s.Enigma.Signature(token)
			assert.Equal(t, signature, validate)
		} else {
			assert.NotNil(t, err, "%s", err)
		}
	}
}

func TestHMACRefreshToken(t *testing.T) {
	token, signature, err := s.GenerateRefreshToken(nil, &hmacValidCase)
	assert.Nil(t, err, "%s", err)
	assert.Equal(t, strings.Split(token, ".")[1], signature)

	validate := s.Enigma.Signature(token)
	err = s.ValidateRefreshToken(nil, &hmacValidCase, token)
	assert.Nil(t, err, "%s", err)
	assert.Equal(t, signature, validate)
}

func TestHMACAuthorizeCode(t *testing.T) {
	for k, c := range []struct {
		r    fosite.Request
		pass bool
	}{
		{
			r:    hmacValidCase,
			pass: true,
		},
		{
			r:    hmacExpiredCase,
			pass: false,
		},
	} {
		token, signature, err := s.GenerateAuthorizeCode(nil, &c.r)
		assert.Nil(t, err, "%s", err)
		assert.Equal(t, strings.Split(token, ".")[1], signature)

		err = s.ValidateAuthorizeCode(nil, &c.r, token)
		if c.pass {
			assert.Nil(t, err, "%d: %s", k, err)
			validate := s.Enigma.Signature(token)
			assert.Equal(t, signature, validate)
		} else {
			assert.NotNil(t, err, "%d: %s", k, err)
		}
	}
}
