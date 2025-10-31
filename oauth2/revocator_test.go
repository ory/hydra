// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/internal"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/oauth2"
	"github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/pop/v6"
	"github.com/ory/x/httprouterx"
	"github.com/ory/x/prometheusx"
)

func createAccessTokenSession(t testing.TB, subject, client, token string, expiresAt time.Time, fs x.FositeStorer, scopes fosite.Arguments) {
	createAccessTokenSessionPairwise(t, subject, client, token, expiresAt, fs, scopes, "")
}

func createAccessTokenSessionPairwise(t testing.TB, subject, client, token string, expiresAt time.Time, fs x.FositeStorer, scopes fosite.Arguments, obfuscated string) {
	ar := fosite.NewAccessRequest(oauth2.NewTestSession(t, subject))
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
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t)

	testhelpers.MustEnsureRegistryKeys(t, reg, x.OpenIDConnectKeyName)
	internal.AddFositeExamples(t, reg)

	tokens := Tokens(reg.OAuth2ProviderConfig(), 4)
	now := time.Now().UTC().Round(time.Second)

	metrics := prometheusx.NewMetricsManagerWithPrefix("hydra", prometheusx.HTTPMetrics, config.Version, config.Commit, config.Date)
	handler := oauth2.NewHandler(reg)
	router := httprouterx.NewRouterAdminWithPrefix(metrics)
	handler.SetPublicRoutes(router.ToPublic(), func(h http.Handler) http.Handler { return h })
	handler.SetAdminRoutes(router)
	server := httptest.NewServer(router)
	defer server.Close()

	createAccessTokenSession(t, "alice", "my-client", tokens[0].sig, now.Add(time.Hour), reg.OAuth2Storage(), nil)
	createAccessTokenSession(t, "siri", "my-client", tokens[1].sig, now.Add(time.Hour), reg.OAuth2Storage(), nil)
	createAccessTokenSession(t, "siri", "my-client", tokens[2].sig, now.Add(-time.Hour), reg.OAuth2Storage(), nil)
	createAccessTokenSession(t, "siri", "encoded:client", tokens[3].sig, now.Add(-time.Hour), reg.OAuth2Storage(), nil)
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
			token: tokens[3].tok,
			assert: func(t *testing.T) {
				assert.Equal(t, 4, countAccessTokens(t, reg.Persister().Connection(context.Background())))
			},
		},
		{
			token: tokens[0].tok,
			assert: func(t *testing.T) {
				t.Logf("Tried to delete: %s %s", tokens[0].sig, tokens[0].tok)
				assert.Equal(t, 3, countAccessTokens(t, reg.Persister().Connection(context.Background())))
			},
		},
		{
			token: tokens[0].tok,
		},
		{
			token: tokens[2].tok,
			assert: func(t *testing.T) {
				assert.Equal(t, 2, countAccessTokens(t, reg.Persister().Connection(context.Background())))
			},
		},
		{
			token: tokens[1].tok,
			assert: func(t *testing.T) {
				assert.Equal(t, 1, countAccessTokens(t, reg.Persister().Connection(context.Background())))
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			resp, err := client.OAuth2API.RevokeOAuth2Token(
				context.WithValue(
					context.Background(),
					hydra.ContextBasicAuth,
					hydra.BasicAuth{UserName: "my-client", Password: "foobar"},
				)).Token(c.token).Execute()
			body, _ := io.ReadAll(resp.Body)
			require.NoErrorf(t, err, "body: %s", body)

			if c.assert != nil {
				c.assert(t)
			}
		})
	}
}
