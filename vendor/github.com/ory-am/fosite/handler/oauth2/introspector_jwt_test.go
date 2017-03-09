package oauth2

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/internal"
	"github.com/ory-am/fosite/token/jwt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestIntrospectJWT(t *testing.T) {
	strat := &RS256JWTStrategy{
		RS256JWTStrategy: &jwt.RS256JWTStrategy{
			PrivateKey: internal.MustRSAKey(),
		},
	}

	v := &StatelessJWTValidator{
		JWTAccessTokenStrategy: strat,
		ScopeStrategy:          fosite.HierarchicScopeStrategy,
	}

	for k, c := range []struct {
		description string
		token       func() string
		expectErr   error
		scopes      []string
	}{
		{
			description: "should fail because jwt is expired",
			token: func() string {
				jwt := jwtExpiredCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(nil, jwt)
				assert.NoError(t, err)
				return token
			},
			expectErr: fosite.ErrTokenExpired,
		},
		{
			description: "should pass because scope was granted",
			token: func() string {
				jwt := jwtValidCase(fosite.AccessToken)
				jwt.GrantedScopes = []string{"foo", "bar"}
				token, _, err := strat.GenerateAccessToken(nil, jwt)
				assert.NoError(t, err)
				return token
			},
			scopes: []string{"foo"},
		},
		{
			description: "should fail because scope was not granted",
			token: func() string {
				jwt := jwtValidCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(nil, jwt)
				assert.NoError(t, err)
				return token
			},
			scopes:    []string{"foo"},
			expectErr: fosite.ErrInvalidScope,
		},
		{
			description: "should fail because signature is invalid",
			token: func() string {
				jwt := jwtValidCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(nil, jwt)
				assert.NoError(t, err)
				parts := strings.Split(token, ".")
				dec, err := base64.RawURLEncoding.DecodeString(parts[1])
				assert.NoError(t, err)
				s := strings.Replace(string(dec), "peter", "piper", -1)
				parts[1] = base64.RawURLEncoding.EncodeToString([]byte(s))
				return strings.Join(parts, ".")
			},
			expectErr: fosite.ErrTokenSignatureMismatch,
		},
		{
			description: "should pass",
			token: func() string {
				jwt := jwtValidCase(fosite.AccessToken)
				token, _, err := strat.GenerateAccessToken(nil, jwt)
				assert.NoError(t, err)
				return token
			},
		},
	} {
		if c.scopes == nil {
			c.scopes = []string{}
		}
		areq := fosite.NewAccessRequest(nil)
		err := v.IntrospectToken(nil, c.token(), fosite.AccessToken, areq, c.scopes)

		assert.True(t, errors.Cause(err) == c.expectErr, "(%d) %s\n%s\n%s", k, c.description, err, c.expectErr)

		if err == nil {
			assert.Equal(t, "peter", areq.Session.GetSubject())
		}

		t.Logf("Passed test case %d", k)
	}
}

func BenchmarkIntrospectJWT(b *testing.B) {
	strat := &RS256JWTStrategy{
		RS256JWTStrategy: &jwt.RS256JWTStrategy{
			PrivateKey: internal.MustRSAKey(),
		},
	}

	v := &StatelessJWTValidator{
		JWTAccessTokenStrategy: strat,
	}

	jwt := jwtValidCase(fosite.AccessToken)
	token, _, err := strat.GenerateAccessToken(nil, jwt)
	assert.NoError(b, err)
	areq := fosite.NewAccessRequest(nil)

	for n := 0; n < b.N; n++ {
		err = v.IntrospectToken(nil, token, fosite.AccessToken, areq, []string{})
	}

	assert.NoError(b, err)
}
