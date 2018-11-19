/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package oauth2_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jose "gopkg.in/square/go-jose.v2"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/herodot"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
)

func createAccessTokenSession(subject, client string, token string, expiresAt time.Time, fs *storage.MemoryStore, scopes fosite.Arguments) {
	createAccessTokenSessionPairwise(subject, client, token, expiresAt, fs, scopes, "")
}

func createAccessTokenSessionPairwise(subject, client string, token string, expiresAt time.Time, fs *storage.MemoryStore, scopes fosite.Arguments, obfuscated string) {
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

	fs.CreateAccessTokenSession(nil, token, ar)
}

func TestRevoke(t *testing.T) {
	var (
		tokens = pkg.Tokens(4)
		store  = storage.NewExampleStore()
		now    = time.Now().UTC().Round(time.Second)
	)

	jm := &jwk.MemoryManager{Keys: map[string]*jose.JSONWebKeySet{}}
	keys, err := (&jwk.RS256Generator{}).Generate("", "sig")
	require.NoError(t, err)
	require.NoError(t, jm.AddKeySet(context.TODO(), oauth2.OpenIDConnectKeyName, keys))
	jwtStrategy, err := jwk.NewRS256JWTStrategy(jm, oauth2.OpenIDConnectKeyName)

	handler := &oauth2.Handler{
		OAuth2: compose.Compose(
			fc,
			store,
			&compose.CommonStrategy{
				CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, []byte("1234567890123456789012345678901234567890"), nil),
				OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(fc, pkg.MustINSECURELOWENTROPYRSAKEYFORTEST()),
			},
			nil,
			compose.OAuth2TokenIntrospectionFactory,
			compose.OAuth2TokenRevocationFactory,
		),
		H:                 herodot.NewJSONWriter(nil),
		OpenIDJWTStrategy: jwtStrategy,
	}

	router := httprouter.New()
	handler.SetRoutes(router, router, func(h http.Handler) http.Handler {
		return h
	})
	server := httptest.NewServer(router)

	createAccessTokenSession("alice", "my-client", tokens[0][0], now.Add(time.Hour), store, nil)
	createAccessTokenSession("siri", "my-client", tokens[1][0], now.Add(time.Hour), store, nil)
	createAccessTokenSession("siri", "my-client", tokens[2][0], now.Add(-time.Hour), store, nil)
	createAccessTokenSession("siri", "doesnt-exist", tokens[3][0], now.Add(-time.Hour), store, nil)

	client := hydra.NewPublicApiWithBasePath(server.URL)
	client.Configuration.Username = "my-client"
	client.Configuration.Password = "foobar"

	for k, c := range []struct {
		token  string
		assert func(*testing.T)
	}{
		{
			token: "invalid",
		},
		{
			token: tokens[3][1],
			assert: func(t *testing.T) {
				assert.Len(t, store.AccessTokens, 4)
			},
		},
		{
			token: tokens[0][1],
			assert: func(t *testing.T) {
				assert.Len(t, store.AccessTokens, 3)
			},
		},
		{
			token: tokens[0][1],
		},
		{
			token: tokens[2][1],
			assert: func(t *testing.T) {
				assert.Len(t, store.AccessTokens, 2)
			},
		},
		{
			token: tokens[1][1],
			assert: func(t *testing.T) {
				assert.Len(t, store.AccessTokens, 1)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			response, err := client.RevokeOAuth2Token(c.token)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode)

			if c.assert != nil {
				c.assert(t)
			}
		})
	}
}
