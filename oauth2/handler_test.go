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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/oauth2"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var lifespan = time.Hour
var flushRequests = []*fosite.Request{
	{
		ID:            "flush-1",
		RequestedAt:   time.Now().Round(time.Second),
		Client:        &client.Client{ID: "foobar"},
		Scopes:        fosite.Arguments{"fa", "ba"},
		GrantedScopes: fosite.Arguments{"fa", "ba"},
		Form:          url.Values{"foo": []string{"bar", "baz"}},
		Session:       &fosite.DefaultSession{Subject: "bar"},
	},
	{
		ID:            "flush-2",
		RequestedAt:   time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
		Client:        &client.Client{ID: "foobar"},
		Scopes:        fosite.Arguments{"fa", "ba"},
		GrantedScopes: fosite.Arguments{"fa", "ba"},
		Form:          url.Values{"foo": []string{"bar", "baz"}},
		Session:       &fosite.DefaultSession{Subject: "bar"},
	},
	{
		ID:            "flush-3",
		RequestedAt:   time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
		Client:        &client.Client{ID: "foobar"},
		Scopes:        fosite.Arguments{"fa", "ba"},
		GrantedScopes: fosite.Arguments{"fa", "ba"},
		Form:          url.Values{"foo": []string{"bar", "baz"}},
		Session:       &fosite.DefaultSession{Subject: "bar"},
	},
}

func TestHandlerFlushHandler(t *testing.T) {
	store := oauth2.NewFositeMemoryStore(nil, lifespan)
	h := &oauth2.Handler{
		H:             herodot.NewJSONWriter(nil),
		ScopeStrategy: fosite.HierarchicScopeStrategy,
		IssuerURL:     "http://hydra.localhost",
		Storage:       store,
	}

	for _, r := range flushRequests {
		require.NoError(t, store.CreateAccessTokenSession(nil, r.ID, r))
	}

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	c := hydra.NewOAuth2ApiWithBasePath(ts.URL)

	ds := new(fosite.DefaultSession)
	ctx := context.Background()

	resp, err := c.FlushInactiveOAuth2Tokens(hydra.FlushInactiveOAuth2TokensRequest{NotAfter: time.Now().Add(-time.Hour * 24)})
	require.NoError(t, err)
	assert.EqualValues(t, http.StatusNoContent, resp.StatusCode)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.NoError(t, err)

	resp, err = c.FlushInactiveOAuth2Tokens(hydra.FlushInactiveOAuth2TokensRequest{NotAfter: time.Now().Add(-(lifespan + time.Hour/2))})
	require.NoError(t, err)
	assert.EqualValues(t, http.StatusNoContent, resp.StatusCode)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.Error(t, err)

	resp, err = c.FlushInactiveOAuth2Tokens(hydra.FlushInactiveOAuth2TokensRequest{NotAfter: time.Now()})
	require.NoError(t, err)
	assert.EqualValues(t, http.StatusNoContent, resp.StatusCode)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.Error(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.Error(t, err)
}

func TestHandlerWellKnown(t *testing.T) {
	h := &oauth2.Handler{
		H:             herodot.NewJSONWriter(nil),
		ScopeStrategy: fosite.HierarchicScopeStrategy,
		IssuerURL:     "http://hydra.localhost",
	}

	AuthPathT := "/oauth2/auth"
	TokenPathT := "/oauth2/token"
	JWKPathT := "/.well-known/jwks.json"

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)

	res, err := http.Get(ts.URL + "/.well-known/openid-configuration")
	require.NoError(t, err)
	defer res.Body.Close()

	trueConfig := oauth2.WellKnown{
		Issuer:                            strings.TrimRight(h.IssuerURL, "/") + "/",
		AuthURL:                           strings.TrimRight(h.IssuerURL, "/") + AuthPathT,
		TokenURL:                          strings.TrimRight(h.IssuerURL, "/") + TokenPathT,
		JWKsURI:                           strings.TrimRight(h.IssuerURL, "/") + JWKPathT,
		RegistrationEndpoint:              strings.TrimRight(h.IssuerURL, "/") + client.ClientsHandlerPath,
		SubjectTypes:                      []string{"pairwise", "public"},
		ResponseTypes:                     []string{"code", "code id_token", "id_token", "token id_token", "token", "token id_token code"},
		ClaimsSupported:                   []string{"sub"},
		ScopesSupported:                   []string{"offline", "openid"},
		UserinfoEndpoint:                  strings.TrimRight(h.IssuerURL, "/") + oauth2.UserinfoPath,
		TokenEndpointAuthMethodsSupported: []string{"client_secret_post", "client_secret_basic", "private_key_jwt", "none"},
		GrantTypesSupported:               []string{"authorization_code", "implicit", "client_credentials"},
		ResponseModesSupported:            []string{"query", "fragment"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
		RequestParameterSupported:         true,
		RequestURIParameterSupported:      true,
		RequireRequestURIRegistration:     true,
	}
	var wellKnownResp oauth2.WellKnown
	err = json.NewDecoder(res.Body).Decode(&wellKnownResp)
	require.NoError(t, err, "problem decoding wellknown json response: %+v", err)
	assert.EqualValues(t, trueConfig, wellKnownResp)

	h.ScopesSupported = "foo,bar"
	h.ClaimsSupported = "baz,oof"
	h.UserinfoEndpoint = "bar"

	res, err = http.Get(ts.URL + "/.well-known/openid-configuration")
	require.NoError(t, err)
	defer res.Body.Close()
	require.NoError(t, json.NewDecoder(res.Body).Decode(&wellKnownResp))

	assert.EqualValues(t, wellKnownResp.ClaimsSupported, []string{"sub", "baz", "oof"})
	assert.EqualValues(t, wellKnownResp.ScopesSupported, []string{"offline", "openid", "foo", "bar"})
	assert.Equal(t, wellKnownResp.UserinfoEndpoint, "bar")
}
