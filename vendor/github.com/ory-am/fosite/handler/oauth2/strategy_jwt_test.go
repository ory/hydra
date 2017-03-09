package oauth2

import (
	"strings"
	"testing"
	"time"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/internal"
	"github.com/ory-am/fosite/token/jwt"
	"github.com/stretchr/testify/assert"
)

var j = &RS256JWTStrategy{
	RS256JWTStrategy: &jwt.RS256JWTStrategy{
		PrivateKey: internal.MustRSAKey(),
	},
}

// returns a valid JWT type. The JWTClaims.ExpiresAt time is intentionally
// left empty to ensure it is pulled from the session's ExpiresAt map for
// the given fosite.TokenType.
var jwtValidCase = func(tokenType fosite.TokenType) *fosite.Request {
	return &fosite.Request{
		Client: &fosite.DefaultClient{
			Secret: []byte("foobarfoobarfoobarfoobar"),
		},
		Session: &JWTSession{
			JWTClaims: &jwt.JWTClaims{
				Issuer:    "fosite",
				Subject:   "peter",
				Audience:  "group0",
				IssuedAt:  time.Now(),
				NotBefore: time.Now(),
				Extra:     make(map[string]interface{}),
			},
			JWTHeader: &jwt.Headers{
				Extra: make(map[string]interface{}),
			},
			ExpiresAt: map[fosite.TokenType]time.Time{
				tokenType: time.Now().Add(time.Hour),
			},
		},
	}
}

// returns an expired JWT type. The JWTClaims.ExpiresAt time is intentionally
// left empty to ensure it is pulled from the session's ExpiresAt map for
// the given fosite.TokenType.
var jwtExpiredCase = func(tokenType fosite.TokenType) *fosite.Request {
	return &fosite.Request{
		Client: &fosite.DefaultClient{
			Secret: []byte("foobarfoobarfoobarfoobar"),
		},
		Session: &JWTSession{
			JWTClaims: &jwt.JWTClaims{
				Issuer:    "fosite",
				Subject:   "peter",
				Audience:  "group0",
				IssuedAt:  time.Now(),
				NotBefore: time.Now(),
				Extra:     make(map[string]interface{}),
			},
			JWTHeader: &jwt.Headers{
				Extra: make(map[string]interface{}),
			},
			ExpiresAt: map[fosite.TokenType]time.Time{
				tokenType: time.Now().Add(-time.Hour),
			},
		},
	}
}

func TestAccessToken(t *testing.T) {
	for _, c := range []struct {
		r    *fosite.Request
		pass bool
	}{
		{
			r:    jwtValidCase(fosite.AccessToken),
			pass: true,
		},
		{
			r:    jwtExpiredCase(fosite.AccessToken),
			pass: false,
		},
	} {
		token, signature, err := j.GenerateAccessToken(nil, c.r)
		assert.Nil(t, err, "%s", err)
		assert.Equal(t, strings.Split(token, ".")[2], signature)

		validate := j.signature(token)
		err = j.ValidateAccessToken(nil, c.r, token)
		if c.pass {
			assert.Nil(t, err, "%s", err)
			assert.Equal(t, signature, validate)
		} else {
			assert.NotNil(t, err, "%s", err)
		}
	}
}

func TestRefreshToken(t *testing.T) {
	token, signature, err := j.GenerateRefreshToken(nil, jwtValidCase(fosite.RefreshToken))
	assert.Nil(t, err, "%s", err)
	assert.Equal(t, strings.Split(token, ".")[2], signature)

	validate := j.signature(token)
	err = j.ValidateRefreshToken(nil, jwtValidCase(fosite.RefreshToken), token)
	assert.Nil(t, err, "%s", err)
	assert.Equal(t, signature, validate)
}

func TestGenerateAuthorizeCode(t *testing.T) {
	for _, c := range []struct {
		r    *fosite.Request
		pass bool
	}{
		{
			r:    jwtValidCase(fosite.AuthorizeCode),
			pass: true,
		},
		{
			r:    jwtExpiredCase(fosite.AuthorizeCode),
			pass: false,
		},
	} {
		token, signature, err := j.GenerateAuthorizeCode(nil, c.r)
		assert.Nil(t, err, "%s", err)
		assert.Equal(t, strings.Split(token, ".")[2], signature)

		validate := j.signature(token)
		err = j.ValidateAuthorizeCode(nil, c.r, token)
		if c.pass {
			assert.Nil(t, err, "%s", err)
			assert.Equal(t, signature, validate)
		} else {
			assert.NotNil(t, err, "%s", err)
		}
	}
}
