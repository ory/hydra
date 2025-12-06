// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/token/hmac"
)

var hmacshaStrategy = oauth2.NewHMACSHAStrategy(
	&hmac.HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar")}},
	&fosite.Config{
		AccessTokenLifespan:   time.Hour * 24,
		AuthorizeCodeLifespan: time.Hour * 24,
	},
)

var hmacshaStrategyUnprefixed = oauth2.NewHMACSHAStrategyUnPrefixed(
	&hmac.HMACStrategy{Config: &fosite.Config{GlobalSecret: []byte("foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobar")}},
	&fosite.Config{
		AccessTokenLifespan:   time.Hour * 24,
		AuthorizeCodeLifespan: time.Hour * 24,
	},
)

var hmacExpiredCase = fosite.Request{
	Client: &fosite.DefaultClient{
		Secret: []byte("foobarfoobarfoobarfoobar"),
	},
	Session: &fosite.DefaultSession{
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.AccessToken:   time.Now().UTC().Add(-time.Hour),
			fosite.AuthorizeCode: time.Now().UTC().Add(-time.Hour),
			fosite.RefreshToken:  time.Now().UTC().Add(-time.Hour),
		},
	},
}

var hmacValidCase = fosite.Request{
	Client: &fosite.DefaultClient{
		Secret: []byte("foobarfoobarfoobarfoobar"),
	},
	Session: &fosite.DefaultSession{
		ExpiresAt: map[fosite.TokenType]time.Time{
			fosite.AccessToken:   time.Now().UTC().Add(time.Hour),
			fosite.AuthorizeCode: time.Now().UTC().Add(time.Hour),
			fosite.RefreshToken:  time.Now().UTC().Add(time.Hour),
		},
	},
}

func TestHMACAccessToken(t *testing.T) {
	for k, c := range []struct {
		r      fosite.Request
		pass   bool
		strat  any
		prefix string
	}{
		{
			r:      hmacValidCase,
			pass:   true,
			strat:  hmacshaStrategy,
			prefix: "ory_at_",
		},
		{
			r:      hmacExpiredCase,
			pass:   false,
			strat:  hmacshaStrategy,
			prefix: "ory_at_",
		},
		{
			r:     hmacValidCase,
			pass:  true,
			strat: hmacshaStrategyUnprefixed,
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			token, signature, err := hmacshaStrategy.GenerateAccessToken(context.Background(), &c.r)
			assert.NoError(t, err)
			assert.Equal(t, strings.Split(token, ".")[1], signature)
			assert.Contains(t, token, c.prefix)

			cases := []string{
				token,
			}
			if c.prefix != "" {
				cases = append(cases, strings.TrimPrefix(token, c.prefix))
			}

			for k, token := range cases {
				t.Run(fmt.Sprintf("prefix=%v", k == 0), func(t *testing.T) {
					err = hmacshaStrategy.ValidateAccessToken(context.Background(), &c.r, token)
					if c.pass {
						assert.NoError(t, err)
						validate := hmacshaStrategy.Enigma.Signature(token)
						assert.Equal(t, signature, validate)
					} else {
						assert.Error(t, err)
					}
				})
			}
		})
	}
}

func TestHMACRefreshToken(t *testing.T) {
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
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			token, signature, err := hmacshaStrategy.GenerateRefreshToken(context.Background(), &c.r)
			assert.NoError(t, err)
			assert.Equal(t, strings.Split(token, ".")[1], signature)
			assert.Contains(t, token, "ory_rt_")

			for k, token := range []string{
				token,
				strings.TrimPrefix(token, "ory_rt_"),
			} {
				t.Run(fmt.Sprintf("prefix=%v", k == 0), func(t *testing.T) {
					err = hmacshaStrategy.ValidateRefreshToken(context.Background(), &c.r, token)
					if c.pass {
						assert.NoError(t, err)
						validate := hmacshaStrategy.Enigma.Signature(token)
						assert.Equal(t, signature, validate)
					} else {
						assert.Error(t, err)
					}
				})
			}
		})
	}
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
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			token, signature, err := hmacshaStrategy.GenerateAuthorizeCode(context.Background(), &c.r)
			assert.NoError(t, err)
			assert.Equal(t, strings.Split(token, ".")[1], signature)
			assert.Contains(t, token, "ory_ac_")

			for k, token := range []string{
				token,
				strings.TrimPrefix(token, "ory_ac_"),
			} {
				t.Run(fmt.Sprintf("prefix=%v", k == 0), func(t *testing.T) {
					err = hmacshaStrategy.ValidateAuthorizeCode(context.Background(), &c.r, token)
					if c.pass {
						assert.NoError(t, err)
						validate := hmacshaStrategy.Enigma.Signature(token)
						assert.Equal(t, signature, validate)
					} else {
						assert.Error(t, err)
					}
				})
			}
		})
	}
}
