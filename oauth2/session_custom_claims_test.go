package oauth2_test

import (
	"testing"

	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"

	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/oauth2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createSessionWithCustomClaims(extra map[string]interface{}, allowedTopLevelClaims []string) oauth2.Session {
	session := &oauth2.Session{
		DefaultSession: &openid.DefaultSession{
			Claims: &jwt.IDTokenClaims{
				Subject: "alice",
				Issuer:  "hydra.localhost",
			},
			Headers: new(jwt.Headers),
			Subject: "alice",
		},
		Extra:                 extra,
		AllowedTopLevelClaims: allowedTopLevelClaims,
	}
	return *session
}

func TestCustomClaimsInSession(t *testing.T) {
	c := internal.NewConfigurationWithDefaults()

	t.Run("no_custom_claims", func(t *testing.T) {
		c.MustSet(config.KeyAllowedTopLevelClaims, []string{})

		session := createSessionWithCustomClaims(nil, c.AllowedTopLevelClaims())
		claims := session.GetJWTClaims().ToMapClaims()

		assert.EqualValues(t, "alice", claims["sub"])
		assert.NotEqual(t, "another-alice", claims["sub"])

		require.Contains(t, claims, "iss")
		assert.EqualValues(t, "hydra.localhost", claims["iss"])

		assert.Empty(t, claims["ext"])
	})
	t.Run("custom_claim_gets_mirrored", func(t *testing.T) {
		c.MustSet(config.KeyAllowedTopLevelClaims, []string{"foo"})
		extra := map[string]interface{}{"foo": "bar"}

		session := createSessionWithCustomClaims(extra, c.AllowedTopLevelClaims())

		claims := session.GetJWTClaims().ToMapClaims()

		assert.EqualValues(t, "alice", claims["sub"])
		assert.NotEqual(t, "another-alice", claims["sub"])

		require.Contains(t, claims, "iss")
		assert.EqualValues(t, "hydra.localhost", claims["iss"])

		require.Contains(t, claims, "foo")
		assert.EqualValues(t, "bar", claims["foo"])

		require.Contains(t, claims, "ext")
		extClaims, ok := claims["ext"].(map[string]interface{})
		require.True(t, ok)

		require.Contains(t, extClaims, "foo")
		assert.EqualValues(t, "bar", extClaims["foo"])
	})
	t.Run("only_non_reserved_claims_get_mirrored", func(t *testing.T) {
		c.MustSet(config.KeyAllowedTopLevelClaims, []string{"foo", "iss", "sub"})
		extra := map[string]interface{}{"foo": "bar", "iss": "hydra.remote", "sub": "another-alice"}

		session := createSessionWithCustomClaims(extra, c.AllowedTopLevelClaims())

		claims := session.GetJWTClaims().ToMapClaims()

		assert.EqualValues(t, "alice", claims["sub"])
		assert.NotEqual(t, "another-alice", claims["sub"])

		require.Contains(t, claims, "iss")
		assert.EqualValues(t, "hydra.localhost", claims["iss"])
		assert.NotEqual(t, "hydra.remote", claims["iss"])

		require.Contains(t, claims, "foo")
		assert.EqualValues(t, "bar", claims["foo"])

		require.Contains(t, claims, "ext")
		extClaims, ok := claims["ext"].(map[string]interface{})
		require.True(t, ok)

		require.Contains(t, extClaims, "foo")
		assert.EqualValues(t, "bar", extClaims["foo"])

		require.Contains(t, extClaims, "iss")
		assert.EqualValues(t, "hydra.remote", extClaims["iss"])

		require.Contains(t, extClaims, "sub")
		assert.EqualValues(t, "another-alice", extClaims["sub"])
	})
	t.Run("no_custom_claims_in_config", func(t *testing.T) {
		c.MustSet(config.KeyAllowedTopLevelClaims, []string{})
		extra := map[string]interface{}{"foo": "bar", "iss": "hydra.remote", "sub": "another-alice"}

		session := createSessionWithCustomClaims(extra, c.AllowedTopLevelClaims())

		claims := session.GetJWTClaims().ToMapClaims()

		assert.EqualValues(t, "alice", claims["sub"])
		assert.NotEqual(t, "another-alice", claims["sub"])

		require.Contains(t, claims, "iss")
		assert.EqualValues(t, "hydra.localhost", claims["iss"])

		assert.NotContains(t, claims, "foo")

		require.Contains(t, claims, "ext")
		extClaims, ok := claims["ext"].(map[string]interface{})
		require.True(t, ok)

		require.Contains(t, extClaims, "foo")
		assert.EqualValues(t, "bar", extClaims["foo"])

		require.Contains(t, extClaims, "sub")
		assert.EqualValues(t, "another-alice", extClaims["sub"])

		require.Contains(t, extClaims, "iss")
		assert.EqualValues(t, "hydra.remote", extClaims["iss"])
	})
	t.Run("more_config_claims_than_given", func(t *testing.T) {
		c.MustSet(config.KeyAllowedTopLevelClaims, []string{"foo", "baz", "bar", "iss"})
		extra := map[string]interface{}{"foo": "foo_value", "sub": "another-alice"}

		session := createSessionWithCustomClaims(extra, c.AllowedTopLevelClaims())

		claims := session.GetJWTClaims().ToMapClaims()

		assert.EqualValues(t, "alice", claims["sub"])
		assert.NotEqual(t, "another-alice", claims["sub"])

		require.Contains(t, claims, "iss")
		assert.EqualValues(t, "hydra.localhost", claims["iss"])
		assert.NotEqual(t, "hydra.remote", claims["iss"])

		require.Contains(t, claims, "foo")
		assert.EqualValues(t, "foo_value", claims["foo"])

		require.Contains(t, claims, "ext")
		extClaims, ok := claims["ext"].(map[string]interface{})
		require.True(t, ok)

		require.Contains(t, extClaims, "foo")
		assert.EqualValues(t, "foo_value", extClaims["foo"])

		require.Contains(t, extClaims, "sub")
		assert.EqualValues(t, "another-alice", extClaims["sub"])
	})
	t.Run("less_config_claims_than_given", func(t *testing.T) {
		c.MustSet(config.KeyAllowedTopLevelClaims, []string{"foo", "sub"})
		extra := map[string]interface{}{"foo": "foo_value", "bar": "bar_value", "baz": "baz_value", "sub": "another-alice"}

		session := createSessionWithCustomClaims(extra, c.AllowedTopLevelClaims())

		claims := session.GetJWTClaims().ToMapClaims()

		assert.EqualValues(t, "alice", claims["sub"])
		assert.NotEqual(t, "another-alice", claims["sub"])

		require.Contains(t, claims, "iss")
		assert.EqualValues(t, "hydra.localhost", claims["iss"])

		require.Contains(t, claims, "foo")
		assert.EqualValues(t, "foo_value", claims["foo"])

		assert.NotContains(t, claims, "bar")
		assert.NotContains(t, claims, "baz")

		require.Contains(t, claims, "ext")
		extClaims, ok := claims["ext"].(map[string]interface{})
		require.True(t, ok)

		require.Contains(t, extClaims, "foo")
		assert.EqualValues(t, "foo_value", extClaims["foo"])

		require.Contains(t, extClaims, "sub")
		assert.EqualValues(t, "another-alice", extClaims["sub"])
	})
	t.Run("config_claims_contain_reserved_claims", func(t *testing.T) {
		c.MustSet(config.KeyAllowedTopLevelClaims, []string{"iss", "sub"})
		extra := map[string]interface{}{"iss": "hydra.remote", "sub": "another-alice"}

		session := createSessionWithCustomClaims(extra, c.AllowedTopLevelClaims())

		claims := session.GetJWTClaims().ToMapClaims()

		assert.EqualValues(t, "alice", claims["sub"])
		assert.NotEqual(t, "another-alice", claims["sub"])

		require.Contains(t, claims, "iss")
		assert.EqualValues(t, "hydra.localhost", claims["iss"])
		assert.NotEqualValues(t, "hydra.remote", claims["iss"])

		require.Contains(t, claims, "ext")
		extClaims, ok := claims["ext"].(map[string]interface{})
		require.True(t, ok)

		require.Contains(t, extClaims, "sub")
		assert.EqualValues(t, "another-alice", extClaims["sub"])

		require.Contains(t, extClaims, "iss")
		assert.EqualValues(t, "hydra.remote", extClaims["iss"])
	})
}
