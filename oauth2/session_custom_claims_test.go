package oauth2_test

import (
	"fmt"

	"testing"

	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"

	"github.com/ory/hydra/oauth2"

	"github.com/stretchr/testify/assert"
)

func createSessionWithCustomClaims(subject string, extra map[string]interface{}, issuer string) oauth2.Session {
	session := &oauth2.Session{
		DefaultSession: &openid.DefaultSession{
			Claims: &jwt.IDTokenClaims{
				Subject: subject,
				Issuer:  issuer,
			},
			Headers: new(jwt.Headers),
			Subject: subject,
		},
		Extra: extra,
	}

	return *session
}

func TestCustomClaimsInSession(t *testing.T) {
	for k, c := range []struct {
		caseName string
		assert   func(*testing.T)
	}{
		{
			caseName: "everything works while no custom claims are added",
			assert: func(t *testing.T) {
				session := createSessionWithCustomClaims("alice", nil, "hydra.localhost")
				claims := session.GetJWTClaims().ToMapClaims()

				assert.EqualValues(t, "alice", claims["sub"])
				assert.NotEqual(t, "another-alice", claims["sub"])

				assert.Contains(t, claims, "iss")
				assert.EqualValues(t, "hydra.localhost", claims["iss"])

				assert.Empty(t, claims["ext"])
			},
		},
		{
			caseName: "custom claims with no overrides get mirrored",
			assert: func(t *testing.T) {
				extra := map[string]interface{}{"foo": "bar"}
				session := createSessionWithCustomClaims("alice", extra, "hydra.localhost")
				claims := session.GetJWTClaims().ToMapClaims()

				assert.EqualValues(t, "alice", claims["sub"])
				assert.NotEqual(t, "another-alice", claims["sub"])

				assert.Contains(t, claims, "iss")
				assert.EqualValues(t, "hydra.localhost", claims["iss"])

				assert.Contains(t, claims, "foo")
				assert.EqualValues(t, claims["foo"], "bar")

				assert.Contains(t, claims, "ext")
				extClaims := claims["ext"].(map[string]interface{})

				assert.Contains(t, extClaims, "foo")
				assert.EqualValues(t, extClaims["foo"], "bar")
			},
		},
		{
			caseName: "custom claims with overrides get mirrored, but without reserved ones",
			assert: func(t *testing.T) {
				extra := map[string]interface{}{"foo": "bar", "iss": "hydra.remote", "sub": "another-alice"}
				session := createSessionWithCustomClaims("alice", extra, "hydra.localhost")
				claims := session.GetJWTClaims().ToMapClaims()

				assert.EqualValues(t, "alice", claims["sub"])
				assert.NotEqual(t, "another-alice", claims["sub"])

				assert.Contains(t, claims, "iss")
				assert.EqualValues(t, "hydra.localhost", claims["iss"])
				assert.NotEqual(t, "hydra.remote", claims["iss"])

				assert.Contains(t, claims, "foo")
				assert.EqualValues(t, claims["foo"], "bar")

				assert.Contains(t, claims, "ext")
				extClaims := claims["ext"].(map[string]interface{})

				assert.Contains(t, extClaims, "foo")
				assert.EqualValues(t, extClaims["foo"], "bar")

				assert.Contains(t, extClaims, "iss")
				assert.EqualValues(t, extClaims["iss"], "hydra.remote")

				assert.Contains(t, extClaims, "sub")
				assert.EqualValues(t, extClaims["sub"], "another-alice")
			},
		},
		{
			caseName: "no_custom_claims_in_config",
			assert: func(t *testing.T) {
				//given custom claims get just mapped under "ext
			},
		},
		{
			caseName: "more_config_claims_than_given",
			assert: func(t *testing.T) {
				//all given claims get mirrored
			},
		},
		{
			caseName: "less_config_claims_than_given",
			assert: func(t *testing.T) {
				//just custom claims of config get top-level, all other just "ext"
			},
		},
		{
			caseName: "config_claims_contain_reserved_claims",
			assert: func(t *testing.T) {
				//what to do here?
				//suggestion: config_claims override reserved claims
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d %s", k, c.caseName), func(t *testing.T) {
			if c.assert != nil {
				c.assert(t)
			}
		})
	}
}
