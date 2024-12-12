// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ory/hydra/v2/internal/testhelpers"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/x/httprouterx"

	"github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/x/contextx"

	hydra "github.com/ory/hydra-client-go/v2"

	"github.com/ory/hydra/v2/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/x"
)

func createAccessTokenSession(subject, client string, token string, expiresAt time.Time, fs x.FositeStorer, scopes fosite.Arguments) {
	createAccessTokenSessionPairwise(subject, client, token, expiresAt, fs, scopes, "")
}

func createAccessTokenSessionPairwise(subject, client string, token string, expiresAt time.Time, fs x.FositeStorer, scopes fosite.Arguments, obfuscated string) {
	ar := fosite.NewAccessRequest(oauth2.NewSession(subject))
	ar.GrantedScope = fosite.Arguments{"core"}
	if scopes != nil {
		ar.GrantedScope = scopes
	}
	ar.RequestedAt = time.Now().UTC().Round(time.Minute)
	ar.Client = &fosite.DefaultClient{ID: client}
	ar.Session.SetExpiresAt(fosite.AccessToken, expiresAt)
	ar.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	if obfuscated != "" {
		ar.Session.(*oauth2.Session).Claims.Subject = obfuscated
	}

	if err := fs.CreateAccessTokenSession(context.Background(), token, ar); err != nil {
		panic(err)
	}
}

func countAccessTokens(t *testing.T, c *pop.Connection) int {
	n, err := c.Count(&sql.OAuth2RequestSQL{Table: "access"})
	require.NoError(t, err)
	return n
}

func TestRevoke(t *testing.T) {
	conf := testhelpers.NewConfigurationWithDefaults()
	reg := testhelpers.NewRegistryMemory(t, conf, &contextx.Default{})

	testhelpers.MustEnsureRegistryKeys(context.Background(), reg, x.OpenIDConnectKeyName)
	internal.AddFositeExamples(reg)

	tokens := Tokens(reg.OAuth2ProviderConfig(), 4)
	now := time.Now().UTC().Round(time.Second)

	handler := reg.OAuth2Handler()
	router := x.NewRouterAdmin(conf.AdminURL)
	handler.SetRoutes(router, &httprouterx.RouterPublic{Router: router.Router}, func(h http.Handler) http.Handler {
		return h
	})
	server := httptest.NewServer(router)
	defer server.Close()

	createAccessTokenSession("alice", "my-client", tokens[0][0], now.Add(time.Hour), reg.OAuth2Storage(), nil)
	createAccessTokenSession("siri", "my-client", tokens[1][0], now.Add(time.Hour), reg.OAuth2Storage(), nil)
	createAccessTokenSession("siri", "my-client", tokens[2][0], now.Add(-time.Hour), reg.OAuth2Storage(), nil)
	createAccessTokenSession("siri", "encoded:client", tokens[3][0], now.Add(-time.Hour), reg.OAuth2Storage(), nil)
	require.Equal(t, 4, countAccessTokens(t, reg.Persister().Connection(context.Background())))

	client := hydra.NewAPIClient(hydra.NewConfiguration())
	client.GetConfig().Servers = hydra.ServerConfigurations{{URL: server.URL}}
	for k, c := range []struct {
		token  string
		assert func(*testing.T)
	}{
		{
			token: "invalid",
			assert: func(t *testing.T) {
				assert.Equal(t, 4, countAccessTokens(t, reg.Persister().Connection(context.Background())))
			},
		},
		{
			token: tokens[3][1],
			assert: func(t *testing.T) {
				assert.Equal(t, 4, countAccessTokens(t, reg.Persister().Connection(context.Background())))
			},
		},
		{
			token: tokens[0][1],
			assert: func(t *testing.T) {
				t.Logf("Tried to delete: %s %s", tokens[0][0], tokens[0][1])
				assert.Equal(t, 3, countAccessTokens(t, reg.Persister().Connection(context.Background())))
			},
		},
		{
			token: tokens[0][1],
		},
		{
			token: tokens[2][1],
			assert: func(t *testing.T) {
				assert.Equal(t, 2, countAccessTokens(t, reg.Persister().Connection(context.Background())))
			},
		},
		{
			token: tokens[1][1],
			assert: func(t *testing.T) {
				assert.Equal(t, 1, countAccessTokens(t, reg.Persister().Connection(context.Background())))
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			_, err := client.OAuth2API.RevokeOAuth2Token(
				context.WithValue(
					context.Background(),
					hydra.ContextBasicAuth,
					hydra.BasicAuth{UserName: "my-client", Password: "foobar"},
				)).Token(c.token).Execute()
			require.NoError(t, err)

			if c.assert != nil {
				c.assert(t)
			}
		})
	}
}
